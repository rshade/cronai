package cron

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron/v3"
	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/rshade/cronai/pkg/config"
)

// Task represents a scheduled task
type Task struct {
	Schedule    string
	Model       string
	Prompt      string
	Processor   string
	Template    string            // Optional template name
	Variables   map[string]string // Variables for the prompt
	ModelParams string            // Model-specific parameters (temperature, tokens, etc.)
}

// Default logger for the cron package
var log = logger.DefaultLogger()

// SetLogger sets the logger for the cron package
func SetLogger(l *logger.Logger) {
	log = l
}

// StartService starts the CronAI service with the given configuration file
// validateTask validates a task configuration
func validateTask(task Task, lineNum int) error {
	var validateErrors *multierror.Error

	// Validate cron schedule format
	_, err := cron.ParseStandard(task.Schedule)
	if err != nil {
		validateErrors = multierror.Append(validateErrors,
			fmt.Errorf("line %d: invalid cron schedule '%s': %w", lineNum, task.Schedule, err))
	}

	// Validate model
	switch task.Model {
	case "openai", "claude", "gemini":
		// Valid models
	default:
		validateErrors = multierror.Append(validateErrors,
			fmt.Errorf("line %d: unsupported model '%s' (supported: openai, claude, gemini)", lineNum, task.Model))
	}

	// Validate prompt file exists
	promptPath := task.Prompt
	if !strings.HasSuffix(promptPath, ".md") {
		promptPath += ".md"
	}

	_, err = os.Stat(fmt.Sprintf("cron_prompts/%s", promptPath))
	if err != nil {
		validateErrors = multierror.Append(validateErrors,
			fmt.Errorf("line %d: prompt file 'cron_prompts/%s' not found: %w", lineNum, promptPath, err))
	}

	// Validate processor
	// Check if processor starts with known prefixes
	validProcessor := false
	for _, prefix := range []string{"slack-", "email-", "webhook-", "log-to-"} {
		if strings.HasPrefix(task.Processor, prefix) {
			validProcessor = true
			break
		}
	}
	if !validProcessor {
		validateErrors = multierror.Append(validateErrors,
			fmt.Errorf("line %d: invalid processor '%s' (should start with slack-, email-, webhook-, or log-to-)",
				lineNum, task.Processor))
	}

	// Validate template if specified
	if task.Template != "" {
		// Check if template file exists (assuming they're in templates/ directory)
		_, err = os.Stat(fmt.Sprintf("templates/%s.tmpl", task.Template))
		if err != nil {
			// Try checking library/ subdirectory
			_, err = os.Stat(fmt.Sprintf("templates/library/%s.tmpl", task.Template))
			if err != nil {
				validateErrors = multierror.Append(validateErrors,
					fmt.Errorf("line %d: template '%s.tmpl' not found in templates/ or templates/library/",
						lineNum, task.Template))
			}
		}
	}

	// Validate model parameters if specified
	if task.ModelParams != "" {
		params, err := config.ParseModelParams(task.ModelParams)
		if err != nil {
			validateErrors = multierror.Append(validateErrors,
				fmt.Errorf("line %d: invalid model parameters: %w", lineNum, err))
		} else {
			// Apply parameters to verify they're valid
			cfg := config.DefaultModelConfig()
			if err := cfg.UpdateFromParams(params); err != nil {
				validateErrors = multierror.Append(validateErrors,
					fmt.Errorf("line %d: invalid model parameters: %w", lineNum, err))
			}
			// Validate the resulting configuration
			if err := cfg.Validate(); err != nil {
				validateErrors = multierror.Append(validateErrors,
					fmt.Errorf("line %d: model configuration validation failed: %w", lineNum, err))
			}
		}
	}

	return validateErrors.ErrorOrNil()
}

func StartService(configPath string) error {
	log.Info("Starting CronAI service", logger.Fields{"config_path": configPath})

	// Parse config file
	tasks, err := parseConfigFile(configPath)
	if err != nil {
		log.Error("Failed to parse config file", logger.Fields{"error": err.Error()})
		return errors.Wrap(errors.CategoryConfiguration, err, "failed to parse configuration file")
	}

	log.Info("Parsed configuration file", logger.Fields{"task_count": len(tasks)})

	// Create a new cron scheduler
	c := cron.New()

	// Add each task to the scheduler
	for i, task := range tasks {
		task := task // Create a copy of the task for the closure
		_, err = c.AddFunc(task.Schedule, func() {
			executeTask(task)
		})
		if err != nil {
			log.Error("Error scheduling task", logger.Fields{
				"task_index": i,
				"schedule":   task.Schedule,
				"model":      task.Model,
				"prompt":     task.Prompt,
				"error":      err.Error(),
			})
			continue
		}
		log.Info("Scheduled task", logger.Fields{
			"task_index": i,
			"schedule":   task.Schedule,
			"model":      task.Model,
			"prompt":     task.Prompt,
			"processor":  task.Processor,
		})
	}

	// Start the scheduler
	c.Start()
	log.Info("Cron scheduler started")

	// Keep running until terminated
	select {}
}

// ListTasks returns a list of tasks from the configuration file
func ListTasks(configPath string) ([]Task, error) {
	return parseConfigFile(configPath)
}

// parseConfigFile parses the configuration file and returns a list of tasks
func parseConfigFile(configPath string) (tasks []Task, err error) {
	log.Info("Parsing configuration file", logger.Fields{"path": configPath})

	file, err := os.Open(configPath)
	if err != nil {
		log.Error("Failed to open config file", logger.Fields{"path": configPath, "error": err.Error()})
		return nil, errors.Wrap(errors.CategoryConfiguration, err, "failed to open config file")
	}

	// Using named return values and multierror to preserve both errors
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Warn("Failed to close config file", logger.Fields{"error": closeErr.Error()})
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
	var parseErrors *multierror.Error

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
			fieldErr := fmt.Errorf("line %d: insufficient fields (need at least 8, got %d)", lineNum, len(parts))
			log.Warn("Invalid format in config file", logger.Fields{
				"line":         lineNum,
				"line_content": line,
				"reason":       "insufficient fields",
				"field_count":  len(parts),
			})
			parseErrors = multierror.Append(parseErrors, fieldErr)
			continue
		}

		// Extract the cron schedule (first 5 parts)
		schedule := strings.Join(parts[0:5], " ")

		// Extract model, prompt, and processor
		model := parts[5]
		prompt := parts[6]
		processor := parts[7]

		log.Debug("Parsed task parameters", logger.Fields{
			"line":      lineNum,
			"schedule":  schedule,
			"model":     model,
			"prompt":    prompt,
			"processor": processor,
		})

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

		// Create the task
		task := Task{
			Schedule:    schedule,
			Model:       model,
			Prompt:      prompt,
			Processor:   processor,
			Template:    template,
			Variables:   variables,
			ModelParams: modelParams,
		}

		// Validate the task
		if validateErr := validateTask(task, lineNum); validateErr != nil {
			parseErrors = multierror.Append(parseErrors, validateErr)
			log.Warn("Task validation failed", logger.Fields{
				"line":  lineNum,
				"error": validateErr.Error(),
			})
			continue
		}

		// Add the validated task
		tasks = append(tasks, task)
	}

	if err := scanner.Err(); err != nil {
		log.Error("Error reading config file", logger.Fields{"path": configPath, "error": err.Error()})
		return nil, errors.Wrap(errors.CategoryConfiguration, err, "error reading config file")
	}

	// Check if we had any parse errors
	if parseErrors != nil && parseErrors.ErrorOrNil() != nil {
		log.Error("Configuration file contains errors", logger.Fields{
			"path":        configPath,
			"error_count": parseErrors.Len(),
		})
		return tasks, errors.Wrap(errors.CategoryConfiguration, parseErrors.ErrorOrNil(),
			fmt.Sprintf("configuration file contains %d validation errors", parseErrors.Len()))
	}

	log.Info("Successfully parsed configuration file", logger.Fields{"path": configPath, "task_count": len(tasks)})
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
	startTime := time.Now()
	log.Info("Executing task", logger.Fields{
		"time":      startTime.Format(time.RFC3339),
		"model":     task.Model,
		"prompt":    task.Prompt,
		"processor": task.Processor,
	})

	// Load the prompt with variables
	var promptContent string
	var err error

	if len(task.Variables) > 0 {
		log.Debug("Loading prompt with variables", logger.Fields{"prompt": task.Prompt, "var_count": len(task.Variables)})
		promptContent, err = prompt.LoadPromptWithVariables(task.Prompt, task.Variables)
	} else {
		log.Debug("Loading prompt without variables", logger.Fields{"prompt": task.Prompt})
		promptContent, err = prompt.LoadPrompt(task.Prompt)
	}

	if err != nil {
		log.Error("Error loading prompt", logger.Fields{"prompt": task.Prompt, "error": err.Error()})
		return
	}

	// Add the prompt name to the variables map for tracking execution
	if task.Variables == nil {
		task.Variables = make(map[string]string)
	}
	task.Variables["promptName"] = task.Prompt

	// Execute the model with model parameters
	log.Debug("Executing model", logger.Fields{"model": task.Model, "prompt_length": len(promptContent)})
	response, err := models.ExecuteModel(task.Model, promptContent, task.Variables, task.ModelParams)
	if err != nil {
		log.Error("Error executing model", logger.Fields{"model": task.Model, "error": err.Error()})
		return
	}

	// Process the response with template
	log.Debug("Processing response", logger.Fields{"processor": task.Processor, "template": task.Template})
	err = processor.ProcessResponse(task.Processor, response, task.Template)
	if err != nil {
		log.Error("Error processing response", logger.Fields{"processor": task.Processor, "error": err.Error()})
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Info("Task completed successfully", logger.Fields{
		"time":      endTime.Format(time.RFC3339),
		"duration":  duration.String(),
		"model":     task.Model,
		"prompt":    task.Prompt,
		"processor": task.Processor,
	})
}
