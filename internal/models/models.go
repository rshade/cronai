package models

import (
	"fmt"
	"time"

	"github.com/rshade/cronai/pkg/config"
)

// ModelResponse represents a response from an AI model
type ModelResponse struct {
	Content     string
	Model       string
	PromptName  string            // Name of the prompt used
	Variables   map[string]string // Variables used in the prompt
	Timestamp   time.Time         // When the response was generated
	ExecutionID string            // Unique execution identifier
}

// ModelClient defines the interface for AI model clients
type ModelClient interface {
	Execute(promptContent string) (*ModelResponse, error)
}

// ExecuteModel executes a prompt using the specified model and returns the response
func ExecuteModel(modelName string, promptContent string, variables map[string]string, modelParams string) (*ModelResponse, error) {
	// Parse model parameters if provided
	params, err := config.ParseModelParams(modelParams)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model parameters: %w", err)
	}

	// Create a model configuration with default values
	modelConfig := config.DefaultModelConfig()

	// Load configuration from environment variables
	modelConfig.LoadFromEnvironment()

	// Update configuration with any provided parameters
	if err := modelConfig.UpdateFromParams(params); err != nil {
		return nil, fmt.Errorf("invalid model parameters: %w", err)
	}

	// Validate the configuration
	if err := modelConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid model configuration: %w", err)
	}

	// Create the appropriate client based on the model name
	var client ModelClient

	switch modelName {
	case "openai":
		client, err = NewOpenAIClient(modelConfig)
	case "claude":
		client, err = NewClaudeClient(modelConfig)
	case "gemini":
		client, err = NewGeminiClient(modelConfig)
	default:
		return nil, fmt.Errorf("unsupported model: %s", modelName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s client: %w", modelName, err)
	}

	// Execute the prompt using the selected client
	response, err := client.Execute(promptContent)
	if err != nil {
		return nil, err
	}

	// Add additional metadata
	response.Variables = variables
	response.Timestamp = time.Now()
	response.ExecutionID = generateExecutionID(modelName, "")

	return response, nil
}

// generateExecutionID creates a unique ID for the execution
func generateExecutionID(modelName, promptName string) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s-%s-%s", modelName, promptName, timestamp)
}
