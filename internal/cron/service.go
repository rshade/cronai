// Package cron provides functionality for scheduling and executing AI tasks.
package cron

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
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

// executeModel is a variable function for mocking in tests
var executeModel = models.ExecuteModel

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

// ScheduledTask represents a task with its schedule
type ScheduledTask struct {
	Schedule    string
	Model       string
	Prompt      string
	Processor   string
	Variables   map[string]string
	ModelParams string
	Template    string
	Task        Task
}

// EntryMetadata contains metadata about a scheduled task.
type EntryMetadata struct {
	Model     string
	Prompt    string
	Processor string
	Schedule  string
	Variables map[string]string
}

// Service manages the scheduling and execution of AI tasks.
type Service struct {
	configFile string
	scheduler  *cron.Cron
	entries    map[string]EntryMetadata
	mu         sync.Mutex
}

// Default logger for the cron package
var log = logger.DefaultLogger()

// SetLogger sets the logger for the cron package
func SetLogger(l *logger.Logger) {
	log = l
}

// Initialize the prompt manager
func init() {
	// Ensure the prompt manager is initialized
	_ = prompt.GetPromptManager()
}

// NewCronService creates a new cron service
func NewCronService(configFile string) *Service {
	return &Service{
		configFile: configFile,
		entries:    make(map[string]EntryMetadata),
	}
}

// StartService starts the CronAI service
func (s *Service) StartService(ctx context.Context) error {
	log.Info("Starting CronAI service", logger.Fields{"config_path": s.configFile})

	// Parse config file
	tasks, err := parseConfigFile(s.configFile)
	if err != nil {
		log.Error("Failed to parse config file", logger.Fields{"error": err.Error()})
		return errors.Wrap(errors.CategoryConfiguration, err, "failed to parse configuration file")
	}

	log.Info("Parsed configuration file", logger.Fields{"task_count": len(tasks)})

	// Create a new cron scheduler
	s.scheduler = cron.New()

	// Add each task to the scheduler
	for i, task := range tasks {
		task := task // Create a copy of the task for the closure

		// Convert Task to ScheduledTask
		scheduledTask := &ScheduledTask{
			Schedule:    task.Schedule,
			Model:       task.Model,
			Prompt:      task.Prompt,
			Processor:   task.Processor,
			Variables:   task.Variables,
			ModelParams: task.ModelParams,
			Task: Task{
				Model:     task.Model,
				Prompt:    task.Prompt,
				Processor: task.Processor,
				Variables: task.Variables,
			},
		}

		err = s.scheduleTask(scheduledTask)
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
	s.scheduler.Start()
	log.Info("Cron scheduler started")

	// Run until context is cancelled
	<-ctx.Done()

	// Stop the scheduler
	s.scheduler.Stop()
	log.Info("Cron scheduler stopped")

	return nil
}

// scheduleTask adds a task to the scheduler
func (s *Service) scheduleTask(task *ScheduledTask) error {
	if s.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	entryID, err := s.scheduler.AddFunc(task.Schedule, func() {
		s.executeTask(task.Task)
	})
	if err != nil {
		return err
	}

	// Store metadata
	s.mu.Lock()
	s.entries[fmt.Sprintf("%d", entryID)] = EntryMetadata{
		Model:     task.Model,
		Prompt:    task.Prompt,
		Processor: task.Processor,
		Schedule:  task.Schedule,
		Variables: task.Variables,
	}
	s.mu.Unlock()

	return nil
}

// ListTasks returns a list of scheduled tasks
func (s *Service) ListTasks() []ScheduledTask {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]ScheduledTask, 0, len(s.entries))
	for _, entry := range s.entries {
		tasks = append(tasks, ScheduledTask{
			Schedule:  entry.Schedule,
			Model:     entry.Model,
			Prompt:    entry.Prompt,
			Processor: entry.Processor,
			Variables: entry.Variables,
			Task: Task{
				Model:     entry.Model,
				Prompt:    entry.Prompt,
				Processor: entry.Processor,
				Variables: entry.Variables,
			},
		})
	}
	return tasks
}

// executeTask executes a single task
func (s *Service) executeTask(task Task) {
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

	// Get the prompt manager
	promptManager := prompt.GetPromptManager()

	if len(task.Variables) > 0 {
		log.Debug("Loading prompt with variables", logger.Fields{"prompt": task.Prompt, "var_count": len(task.Variables)})
		promptContent, err = promptManager.LoadPromptWithVariables(task.Prompt, task.Variables)
	} else {
		log.Debug("Loading prompt without variables", logger.Fields{"prompt": task.Prompt})
		promptContent, err = promptManager.LoadPrompt(task.Prompt)
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
	response, err := executeModel(task.Model, promptContent, task.Variables, task.ModelParams)
	if err != nil {
		log.Error("Error executing model", logger.Fields{"model": task.Model, "error": err.Error()})
		return
	}

	// Process the response
	log.Debug("Processing response", logger.Fields{"processor": task.Processor})

	// Create processor config
	procConfig := processor.Config{
		Type:   "unknown", // Will be set by ParseProcessor
		Target: "",        // Will be set by ParseProcessor
	}

	// Parse processor name to determine type and target
	procType, target := parseProcessor(task.Processor)
	procConfig.Type = procType
	procConfig.Target = target

	// Get the registry
	registry := processor.GetRegistry()

	// Create the processor
	proc, err := registry.CreateProcessor(procType, procConfig)
	if err != nil {
		log.Error("Error getting processor", logger.Fields{"error": err.Error()})
		return
	}

	// Create a models.ModelResponse
	modelResponse := &models.ModelResponse{
		Model:       task.Model,
		PromptName:  task.Prompt,
		Content:     response.Content,
		Timestamp:   time.Now(),
		ExecutionID: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	err = proc.Process(modelResponse, "")
	if err != nil {
		log.Error("Error processing response", logger.Fields{"error": err.Error()})
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

// RunTask executes a single task immediately
func (s *Service) RunTask(task Task) error {
	// Load the prompt content
	var promptContent string
	var err error

	// Get the prompt manager
	promptManager := prompt.GetPromptManager()

	if len(task.Variables) > 0 {
		promptContent, err = promptManager.LoadPromptWithVariables(task.Prompt, task.Variables)
	} else {
		promptContent, err = promptManager.LoadPrompt(task.Prompt)
	}

	if err != nil {
		return fmt.Errorf("error loading prompt: %w", err)
	}

	// Execute the model with model parameters
	response, err := executeModel(task.Model, promptContent, task.Variables, task.ModelParams)
	if err != nil {
		return fmt.Errorf("error executing model: %w", err)
	}

	// Process the response
	procConfig := processor.Config{
		Type:   "unknown", // Will be set by ParseProcessor
		Target: "",        // Will be set by ParseProcessor
	}

	// Parse processor name to determine type and target
	procType, target := parseProcessor(task.Processor)
	procConfig.Type = procType
	procConfig.Target = target

	// Get the registry
	registry := processor.GetRegistry()

	// Create the processor
	proc, err := registry.CreateProcessor(procType, procConfig)
	if err != nil {
		return fmt.Errorf("error getting processor: %w", err)
	}

	// Create a models.ModelResponse
	modelResponse := &models.ModelResponse{
		Model:       task.Model,
		PromptName:  task.Prompt,
		Content:     response.Content,
		Timestamp:   time.Now(),
		ExecutionID: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	err = proc.Process(modelResponse, "")
	if err != nil {
		return fmt.Errorf("error processing response: %w", err)
	}

	return nil
}

// Stop stops the cron service
func (s *Service) Stop() error {
	if s.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}
	s.scheduler.Stop()
	return nil
}

// Package-level functions for backward compatibility

// StartService starts the CronAI service with the given configuration file
func StartService(configPath string) error {
	service := NewCronService(configPath)
	ctx := context.Background()
	return service.StartService(ctx)
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

		task, err := parseConfigLine(line)
		// Handle empty or comment lines
		if task == nil && err == nil {
			continue
		}
		// Handle parse errors
		if err != nil {
			parseErrors = multierror.Append(parseErrors, fmt.Errorf("line %d: %v", lineNum, err))
			continue
		}

		if task != nil {
			// Validate the task
			if validateErr := validateTask(task.Task, lineNum); validateErr != nil {
				parseErrors = multierror.Append(parseErrors, validateErr)
			}
			tasks = append(tasks, task.Task)
		}
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

// parseConfigLine parses a single line from the configuration file
func parseConfigLine(line string) (*ScheduledTask, error) {
	// Skip empty lines and comments
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil // Skip empty lines
	}
	if strings.HasPrefix(line, "#") {
		return nil, nil // Skip comment lines
	}

	// Parse the line
	parts := strings.Fields(line)
	if len(parts) < 8 { // Need at least 8 fields: 5 for cron schedule + model + prompt + processor
		return nil, fmt.Errorf("invalid format: insufficient fields (need at least 8, got %d)", len(parts))
	}

	// Extract the cron schedule (first 5 parts)
	schedule := strings.Join(parts[0:5], " ")

	// Extract model, prompt, and processor
	modelPart := parts[5]
	prompt := parts[6]
	processor := parts[7]

	// Parse model and model parameters
	var model string
	var modelParams string
	if strings.Contains(modelPart, ":") {
		// Split on first colon to separate model from parameters
		modelParts := strings.SplitN(modelPart, ":", 2)
		model = modelParts[0]
		if len(modelParts) > 1 {
			modelParams = modelParts[1]
		}
	} else {
		model = modelPart
	}

	// Validate model
	if !isValidModel(model) {
		return nil, fmt.Errorf("invalid model '%s'", model)
	}

	// Validate processor format
	if !isValidProcessor(processor) {
		return nil, fmt.Errorf("invalid processor format '%s'", processor)
	}

	// Parse optional variables
	var variables map[string]string
	var template string
	if len(parts) > 8 {
		// Check for variables
		varString := strings.Join(parts[8:], " ")
		if strings.Contains(varString, "=") {
			variables = parseVariables(varString)
			if variables == nil {
				variables = make(map[string]string)
			}

			// Check if any variable has an invalid format
			for _, part := range parts[8:] {
				if strings.Contains(part, "=") {
					keyValue := strings.SplitN(part, "=", 2)
					if len(keyValue) != 2 {
						return nil, fmt.Errorf("invalid variable format '%s'", part)
					}
				}
			}

			// Extract template from variables if present
			if templateVar, ok := variables["template"]; ok {
				template = templateVar
			}
		} else {
			// If there's a part without '=' and it's not a valid format, it's an error
			return nil, fmt.Errorf("invalid variable format '%s'", varString)
		}
	}

	// Create the task
	task := &ScheduledTask{
		Schedule:    schedule,
		Model:       model,
		Prompt:      prompt,
		Processor:   processor,
		Variables:   variables,
		ModelParams: modelParams,
		Template:    template,
		Task: Task{
			Model:       model,
			Prompt:      prompt,
			Processor:   processor,
			Variables:   variables,
			ModelParams: modelParams,
			Template:    template,
		},
	}

	return task, nil
}

// parseVariables parses a variable string and returns a map
func parseVariables(varString string) map[string]string {
	if varString == "" {
		return nil
	}

	variables := make(map[string]string)
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

	if len(variables) == 0 {
		return make(map[string]string)
	}
	return variables
}

// Helper functions for validation

// isValidModel checks if the model is supported
func isValidModel(model string) bool {
	switch model {
	case "openai", "claude", "gemini":
		return true
	default:
		return false
	}
}

// isValidProcessor checks if the processor format is valid
func isValidProcessor(processor string) bool {
	// Check common processor formats
	if processor == "console" {
		return true
	}

	// Check if processor starts with known prefixes
	validPrefixes := []string{"slack-", "email-", "webhook-", "file-", "log-", "github-"}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(processor, prefix) {
			return true
		}
	}
	return false
}

// parseProcessor parses the processor name to determine type and target
func parseProcessor(processorName string) (string, string) {
	// Handle special processor formats
	if strings.HasPrefix(processorName, "slack-") {
		return "slack", strings.TrimPrefix(processorName, "slack-")
	} else if strings.HasPrefix(processorName, "email-") {
		return "email", strings.TrimPrefix(processorName, "email-")
	} else if strings.HasPrefix(processorName, "webhook-") {
		return "webhook", strings.TrimPrefix(processorName, "webhook-")
	} else if strings.HasPrefix(processorName, "github-") {
		return "github", strings.TrimPrefix(processorName, "github-")
	} else if strings.HasPrefix(processorName, "file-") {
		return "file", strings.TrimPrefix(processorName, "file-")
	} else if processorName == "log-to-file" {
		return "file", "to-file"
	} else if processorName == "console" {
		return "console", ""
	}

	// Default case
	return "console", ""
}

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
	if !isValidModel(task.Model) {
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
	if !isValidProcessor(task.Processor) {
		validateErrors = multierror.Append(validateErrors,
			fmt.Errorf("line %d: invalid processor '%s' (should start with slack-, email-, webhook-, or log-)",
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
		modelConfig := config.NewModelConfig()
		params, err := config.ParseModelParams(task.ModelParams)
		if err != nil {
			validateErrors = multierror.Append(validateErrors,
				fmt.Errorf("line %d: invalid model parameters: %w", lineNum, err))
		} else if err := modelConfig.UpdateFromParams(params); err != nil {
			validateErrors = multierror.Append(validateErrors,
				fmt.Errorf("line %d: invalid model parameters: %w", lineNum, err))
		}
	}

	return validateErrors.ErrorOrNil()
}

// ProcessResponse processes a model response using the specified processor
func (s *Service) ProcessResponse(processor processor.Processor, response *models.ModelResponse, templateName string) error {
	return processor.Process(response, templateName)
}

// CreateProcessor creates a new processor instance
func (s *Service) CreateProcessor(processorType string, config processor.Config) (processor.Processor, error) {
	return processor.CreateProcessor(processorType, config)
}

// GetProcessor creates a new processor instance
func (s *Service) GetProcessor(processorType string, config processor.Config) (processor.Processor, error) {
	return processor.CreateProcessor(processorType, config)
}
