// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file implements the task processor that executes queue tasks by loading prompts,
// executing AI models, and processing responses through configured processors.
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
)

// executeModel is a variable function for mocking in tests
var executeModel = models.ExecuteModel

// DefaultTaskProcessor implements the TaskProcessor interface
type DefaultTaskProcessor struct {
	promptManager prompt.Manager
}

// NewTaskProcessor creates a new task processor
func NewTaskProcessor() TaskProcessor {
	return &DefaultTaskProcessor{
		promptManager: prompt.GetPromptManager(),
	}
}

// Process processes a task message
func (p *DefaultTaskProcessor) Process(_ context.Context, task *TaskMessage) error {
	startTime := time.Now()

	log.Info("Processing queue task", logger.Fields{
		"model":     task.Model,
		"prompt":    task.Prompt,
		"processor": task.Processor,
		"isInline":  task.IsInline,
	})

	// Add the prompt name to variables for tracking before loading
	if task.Variables == nil {
		task.Variables = make(map[string]string)
	}
	task.Variables["promptName"] = task.Prompt

	// Load or use the prompt content
	var promptContent string
	var err error

	if task.IsInline {
		// Use the prompt field directly as content
		log.Debug("Using inline prompt content", logger.Fields{"content_length": len(task.Prompt)})
		promptContent = task.Prompt
	} else {
		// Load prompt from file
		if len(task.Variables) > 0 {
			log.Debug("Loading prompt with variables", logger.Fields{
				"prompt":    task.Prompt,
				"var_count": len(task.Variables),
			})
			promptContent, err = p.promptManager.LoadPromptWithVariables(task.Prompt, task.Variables)
		} else {
			log.Debug("Loading prompt without variables", logger.Fields{"prompt": task.Prompt})
			promptContent, err = p.promptManager.LoadPrompt(task.Prompt)
		}

		if err != nil {
			return fmt.Errorf("failed to load prompt: %w", err)
		}
	}

	log.Debug("Executing model", logger.Fields{
		"model":     task.Model,
		"var_count": len(task.Variables),
	})

	response, err := executeModel(task.Model, promptContent, task.Variables, "")
	if err != nil {
		return fmt.Errorf("failed to execute model: %w", err)
	}

	// Process the response
	log.Debug("Processing response", logger.Fields{"processor": task.Processor})

	// Parse processor configuration
	procType, target := parseProcessorName(task.Processor)

	procConfig := processor.Config{
		Type:   procType,
		Target: target,
	}

	// Get the registry and create processor
	registry := processor.GetRegistry()
	proc, err := registry.CreateProcessor(procType, procConfig)
	if err != nil {
		return fmt.Errorf("failed to create processor: %w", err)
	}

	// Create a model response
	currentTime := time.Now()
	modelResponse := &models.ModelResponse{
		Model:       task.Model,
		PromptName:  task.Prompt,
		Content:     response.Content,
		Timestamp:   currentTime,
		ExecutionID: fmt.Sprintf("queue-%d", currentTime.UnixNano()),
	}

	// Process the response
	if err := proc.Process(modelResponse, ""); err != nil {
		return fmt.Errorf("failed to process response: %w", err)
	}

	duration := time.Since(startTime)
	log.Info("Queue task completed successfully", logger.Fields{
		"duration":  duration.String(),
		"model":     task.Model,
		"processor": task.Processor,
	})

	return nil
}

// parseProcessorName parses the processor name to determine type and target
func parseProcessorName(processorName string) (string, string) {
	// This is similar to the function in cron/service.go
	// but kept separate to avoid circular dependencies

	prefixes := map[string]string{
		"slack-":   "slack",
		"email-":   "email",
		"webhook-": "webhook",
		"github-":  "github",
		"file-":    "file",
	}

	for prefix, procType := range prefixes {
		if len(processorName) > len(prefix) && processorName[:len(prefix)] == prefix {
			return procType, processorName[len(prefix):]
		}
	}

	// Special cases
	if processorName == "log-to-file" {
		return "file", "to-file"
	}

	if processorName == "console" {
		return "console", ""
	}

	// Default
	return "console", ""
}
