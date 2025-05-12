package cron

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron/v3"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
)

// Task represents a scheduled task
type Task struct {
	Schedule      string
	Model         string
	Prompt        string
	Processor     string
	Template      string            // Optional template name
	Variables     map[string]string // Variables for the prompt
	ModelParams   string            // Model-specific parameters (temperature, tokens, etc.)
}

// StartService starts the CronAI service with the given configuration file
func StartService(configPath string) error {
	// Parse config file
	tasks, err := parseConfigFile(configPath)
	if err != nil {
		return err
	}

	// Create a new cron scheduler
	c := cron.New()

	// Add each task to the scheduler
	for _, task := range tasks {
		task := task // Create a copy of the task for the closure
		_, err = c.AddFunc(task.Schedule, func() {
			executeTask(task)
		})
		if err != nil {
			fmt.Printf("Error scheduling task: %v\n", err)
			continue
		}
		fmt.Printf("Scheduled task: %s %s %s %s\n", task.Schedule, task.Model, task.Prompt, task.Processor)
	}

	// Start the scheduler
	c.Start()

	// Keep running until terminated
	select {}
}

// ListTasks returns a list of tasks from the configuration file
func ListTasks(configPath string) ([]Task, error) {
	return parseConfigFile(configPath)
}

// parseConfigFile parses the configuration file and returns a list of tasks
func parseConfigFile(configPath string) (tasks []Task, err error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	// Using named return values and multierror to preserve both errors
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			// If we have both processing and close errors, combine them
			if err != nil {
				err = multierror.Append(
					err,
					fmt.Errorf("failed to close config file: %w", closeErr),
				)
			} else {
				// If we only have close error, just use that
				err = fmt.Errorf("failed to close config file: %w", closeErr)
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line
		parts := strings.Fields(line)
		if len(parts) < 8 { // Need at least 8 fields: 5 for cron schedule + model + prompt + processor
			fmt.Printf("Line %d: Invalid format (need at least 8 fields)\n", lineNum)
			continue
		}

		// Extract the cron schedule (first 5 parts)
		schedule := strings.Join(parts[0:5], " ")

		// Extract model, prompt, and processor
		model := parts[5]
		prompt := parts[6]
		processor := parts[7]

		// Parse optional template and variables
		var template string
		variables := make(map[string]string)
		var modelParams string
		
		// Parse optional template if present
		if len(parts) > 8 {
			// Check if the next part is a template or variables
			if !strings.Contains(parts[8], "=") {
				template = parts[8]
				// Process variables if they exist after the template
				if len(parts) > 9 {
					varString := strings.Join(parts[9:], " ")
					parseVariables(varString, variables)
				}
			} else {
				// No template specified, just variables
				varString := strings.Join(parts[8:], " ")
				parseVariables(varString, variables)
			}
		}
		
		// Check for model parameters as field after variables
		if len(parts) > 9 && strings.HasPrefix(parts[9], "model_params:") {
			modelParams = strings.TrimPrefix(parts[9], "model_params:")
		} else if len(parts) > 8 && strings.HasPrefix(parts[8], "model_params:") {
			modelParams = strings.TrimPrefix(parts[8], "model_params:")
		}

		// Add the task
		tasks = append(tasks, Task{
			Schedule:    schedule,
			Model:       model,
			Prompt:      prompt,
			Processor:   processor,
			Template:    template,
			Variables:   variables,
			ModelParams: modelParams,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return tasks, nil
}

// parseVariables parses a variable string and adds values to the variables map
func parseVariables(varString string, variables map[string]string) {
	for _, varPair := range strings.Split(varString, ",") {
		keyValue := strings.SplitN(varPair, "=", 2)
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])

			// Handle special variables
			switch value {
			case "{{CURRENT_DATE}}":
				value = time.Now().Format("2006-01-02")
			case "{{CURRENT_TIME}}":
				value = time.Now().Format("15:04:05")
			case "{{CURRENT_DATETIME}}":
				value = time.Now().Format("2006-01-02 15:04:05")
			}

			variables[key] = value
		}
	}
}

// executeTask executes a single task
func executeTask(task Task) {
	fmt.Printf("Executing task at %s: %s %s %s\n", time.Now().Format(time.RFC3339), task.Model, task.Prompt, task.Processor)

	// Load the prompt with variables
	var promptContent string
	var err error

	if len(task.Variables) > 0 {
		promptContent, err = prompt.LoadPromptWithVariables(task.Prompt, task.Variables)
	} else {
		promptContent, err = prompt.LoadPrompt(task.Prompt)
	}

	if err != nil {
		fmt.Printf("Error loading prompt: %v\n", err)
		return
	}

	// Execute the model with model parameters
	response, err := models.ExecuteModel(task.Model, promptContent, task.Variables, task.ModelParams)
	if err != nil {
		fmt.Printf("Error executing model: %v\n", err)
		return
	}

	// Process the response with template
	err = processor.ProcessResponse(task.Processor, response, task.Template)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return
	}

	fmt.Printf("Task completed successfully at %s\n", time.Now().Format(time.RFC3339))
}