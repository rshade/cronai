package models

import (
	"fmt"
	"log"
	"time"

	"github.com/rshade/cronai/pkg/config"
)

// ModelError represents an error from a model execution with full context
type ModelError struct {
	Model   string
	Message string
	Err     error
	Time    time.Time
	Retry   int
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("[%s] retry %d: %s", e.Model, e.Retry, e.Message)
}

func (e *ModelError) Unwrap() error {
	return e.Err
}

// ModelFallbackResult represents the result of a model execution with fallback attempts
type ModelFallbackResult struct {
	Response   *ModelResponse
	Errors     []ModelError
	FinalModel string // The model that ultimately succeeded (if any)
}

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

	// Get the prompt name from variables if available
	promptName := ""
	if promptNameVal, exists := variables["promptName"]; exists {
		promptName = promptNameVal
	}

	// Execute with fallback support
	result := executeWithFallback(modelName, promptContent, variables, modelConfig, promptName)

	// If we got a successful response, return it
	if result.Response != nil {
		// Log fallback usage if applicable
		if len(result.Errors) > 0 {
			log.Printf("Model '%s' succeeded after %d failed attempts", result.FinalModel, len(result.Errors))
			for i, err := range result.Errors {
				log.Printf("  Attempt %d: %v", i+1, err)
			}
		}
		return result.Response, nil
	}

	// All models failed, return comprehensive error
	errorMsg := fmt.Sprintf("all models failed after %d attempts", len(result.Errors))
	for i, err := range result.Errors {
		errorMsg += fmt.Sprintf("\n  Attempt %d - %s", i+1, err.Error())
	}
	return nil, fmt.Errorf("%s", errorMsg)
}

// executeWithFallback executes a model with fallback support
func executeWithFallback(primaryModel string, promptContent string, variables map[string]string, modelConfig *config.ModelConfig, promptName string) *ModelFallbackResult {
	result := &ModelFallbackResult{
		Errors: []ModelError{},
	}

	// Build the list of models to try (primary + fallbacks)
	modelsToTry := []string{primaryModel}

	// Add configured fallback models
	if len(modelConfig.FallbackModels) > 0 {
		modelsToTry = append(modelsToTry, modelConfig.FallbackModels...)
	} else {
		// Default fallback sequence if not configured
		modelsToTry = append(modelsToTry, getDefaultFallbackSequence(primaryModel)...)
	}

	// Try each model with retries
	for modelIndex, modelName := range modelsToTry {
		for retry := 0; retry < modelConfig.MaxRetries; retry++ {
			// Create the client for this model
			client, err := createModelClient(modelName, modelConfig)
			if err != nil {
				result.Errors = append(result.Errors, ModelError{
					Model:   modelName,
					Message: fmt.Sprintf("failed to create client: %v", err),
					Err:     err,
					Time:    time.Now(),
					Retry:   retry,
				})
				log.Printf("Failed to create %s client: %v", modelName, err)
				continue
			}

			// Execute the prompt
			response, err := client.Execute(promptContent)
			if err != nil {
				result.Errors = append(result.Errors, ModelError{
					Model:   modelName,
					Message: fmt.Sprintf("execution failed: %v", err),
					Err:     err,
					Time:    time.Now(),
					Retry:   retry,
				})
				log.Printf("Model %s (attempt %d/%d) failed: %v", modelName, retry+1, modelConfig.MaxRetries, err)
				continue
			}

			// Success! Add metadata and return
			response.Variables = variables
			response.Timestamp = time.Now()
			response.PromptName = promptName
			response.ExecutionID = generateExecutionID(modelName, promptName)
			result.Response = response
			result.FinalModel = modelName

			// Log success with fallback info if we're not on the primary model
			if modelIndex > 0 || retry > 0 {
				log.Printf("Model %s succeeded after %d total attempts across %d models",
					modelName, len(result.Errors)+1, modelIndex+1)
			}

			return result
		}
	}

	return result
}

// createModelClient creates a client for the specified model
func createModelClient(modelName string, modelConfig *config.ModelConfig) (ModelClient, error) {
	switch modelName {
	case "openai":
		return NewOpenAIClient(modelConfig)
	case "claude":
		return NewClaudeClient(modelConfig)
	case "gemini":
		return NewGeminiClient(modelConfig)
	default:
		return nil, fmt.Errorf("unsupported model: %s", modelName)
	}
}

// getDefaultFallbackSequence returns the default fallback sequence for a model
func getDefaultFallbackSequence(primaryModel string) []string {
	// Define sensible defaults for each primary model
	switch primaryModel {
	case "openai":
		return []string{"claude", "gemini"}
	case "claude":
		return []string{"openai", "gemini"}
	case "gemini":
		return []string{"openai", "claude"}
	default:
		// If unknown model, try all known models
		return []string{"openai", "claude", "gemini"}
	}
}

// generateExecutionID creates a unique ID for the execution
func generateExecutionID(modelName, promptName string) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s-%s-%s", modelName, promptName, timestamp)
}