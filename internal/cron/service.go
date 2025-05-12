package cron

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/robfig/cron/v3"
)

// Task represents a scheduled task
type Task struct {
	Schedule  string
	Model     string
	Prompt    string
	Processor string
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
func parseConfigFile(configPath string) ([]Task, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var tasks []Task
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
		if len(parts) < 4 {
			fmt.Printf("Line %d: Invalid format (need at least 4 fields)\n", lineNum)
			continue
		}

		// Extract the cron schedule (first 5 parts)
		schedule := strings.Join(parts[:5], " ")
		parts = parts[5:]

		// Extract model, prompt, and processor
		if len(parts) < 3 {
			fmt.Printf("Line %d: Missing fields after schedule\n", lineNum)
			continue
		}

		model := parts[0]
		prompt := parts[1]
		processor := parts[2]

		// Add the task
		tasks = append(tasks, Task{
			Schedule:  schedule,
			Model:     model,
			Prompt:    prompt,
			Processor: processor,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return tasks, nil
}

// executeTask executes a single task
func executeTask(task Task) {
	fmt.Printf("Executing task at %s: %s %s %s\n", time.Now().Format(time.RFC3339), task.Model, task.Prompt, task.Processor)

	// Load the prompt
	promptContent, err := prompt.LoadPrompt(task.Prompt)
	if err != nil {
		fmt.Printf("Error loading prompt: %v\n", err)
		return
	}

	// Execute the model
	response, err := models.ExecuteModel(task.Model, promptContent)
	if err != nil {
		fmt.Printf("Error executing model: %v\n", err)
		return
	}

	// Process the response
	err = processor.ProcessResponse(task.Processor, response)
	if err != nil {
		fmt.Printf("Error processing response: %v\n", err)
		return
	}

	fmt.Printf("Task completed successfully at %s\n", time.Now().Format(time.RFC3339))
}
