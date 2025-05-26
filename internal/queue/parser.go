// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file implements message parsing and validation logic, supporting multiple message
// formats and ensuring proper task configuration before processing.
package queue

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DefaultMessageParser implements the MessageParser interface
type DefaultMessageParser struct{}

// NewMessageParser creates a new message parser instance
func NewMessageParser() MessageParser {
	return &DefaultMessageParser{}
}

// MinimalMessage represents the minimal message format (variables only)
type MinimalMessage struct {
	Variables map[string]string `json:"variables"`
}

// ComprehensiveMessage represents the comprehensive message format
type ComprehensiveMessage struct {
	Model     string            `json:"model"`
	Prompt    string            `json:"prompt"`
	Processor string            `json:"processor"`
	Variables map[string]string `json:"variables,omitempty"`
	IsInline  bool              `json:"is_inline,omitempty"`
}

// Parse converts a raw message into a TaskMessage
func (p *DefaultMessageParser) Parse(message *Message) (*TaskMessage, error) {
	if message == nil {
		return nil, fmt.Errorf("message cannot be nil")
	}

	if len(message.Body) == 0 {
		return nil, fmt.Errorf("message body is empty")
	}

	// Check what type of message we have by examining the structure
	var rawMessage map[string]interface{}
	if err := json.Unmarshal(message.Body, &rawMessage); err != nil {
		return nil, fmt.Errorf("unable to parse message: invalid format")
	}

	// Check if it's a comprehensive format (has model, prompt, or processor fields)
	_, hasModel := rawMessage["model"]
	_, hasPrompt := rawMessage["prompt"]
	_, hasProcessor := rawMessage["processor"]
	_, hasIsInline := rawMessage["is_inline"]

	isComprehensiveFormat := hasModel || hasPrompt || hasProcessor || hasIsInline

	if isComprehensiveFormat {
		// Try comprehensive format
		var comprehensive ComprehensiveMessage
		if err := json.Unmarshal(message.Body, &comprehensive); err != nil {
			return nil, fmt.Errorf("unable to parse message: invalid format")
		}

		if comprehensive.Model == "" || comprehensive.Prompt == "" || comprehensive.Processor == "" {
			return nil, fmt.Errorf("comprehensive message format missing required fields: model, prompt, processor")
		}

		return &TaskMessage{
			Model:     comprehensive.Model,
			Prompt:    comprehensive.Prompt,
			Processor: comprehensive.Processor,
			Variables: comprehensive.Variables,
			IsInline:  comprehensive.IsInline,
		}, nil
	}

	// If not comprehensive, try minimal format
	var minimal MinimalMessage
	if err := json.Unmarshal(message.Body, &minimal); err == nil {
		// For minimal format, extract configuration from message attributes
		model := message.Attributes["model"]
		prompt := message.Attributes["prompt"]
		processor := message.Attributes["processor"]

		if model == "" || prompt == "" || processor == "" {
			return nil, fmt.Errorf("minimal message format requires model, prompt, and processor in message attributes")
		}

		return &TaskMessage{
			Model:     model,
			Prompt:    prompt,
			Processor: processor,
			Variables: minimal.Variables,
			IsInline:  false,
		}, nil
	}

	return nil, fmt.Errorf("unable to parse message: invalid format")
}

// Validate checks if a TaskMessage is valid
func (p *DefaultMessageParser) Validate(task *TaskMessage) error {
	if task == nil {
		return fmt.Errorf("task message cannot be nil")
	}

	if strings.TrimSpace(task.Model) == "" {
		return fmt.Errorf("model cannot be empty")
	}

	if strings.TrimSpace(task.Prompt) == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	if strings.TrimSpace(task.Processor) == "" {
		return fmt.Errorf("processor cannot be empty")
	}

	// Validate model is one of the supported types
	validModels := map[string]bool{
		"openai": true,
		"claude": true,
		"gemini": true,
	}

	if !validModels[strings.ToLower(task.Model)] {
		return fmt.Errorf("unsupported model: %s", task.Model)
	}

	// If it's not an inline prompt, validate it doesn't contain newlines
	if !task.IsInline && strings.Contains(task.Prompt, "\n") {
		return fmt.Errorf("prompt file reference cannot contain newlines")
	}

	// Validate variables don't contain invalid characters
	for key, value := range task.Variables {
		if strings.ContainsAny(key, " \t\n\r") {
			return fmt.Errorf("variable key '%s' contains whitespace", key)
		}
		if key == "" {
			return fmt.Errorf("variable key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("variable value for key '%s' cannot be empty", key)
		}
	}

	return nil
}
