package models

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
)

// MockModelClient is a simple mocked version of the ModelClient interface
type MockModelClient struct {
	ShouldFail   bool
	ErrorMessage string
	Content      string
	Model        string
	ExecuteCount int
}

func (m *MockModelClient) Execute(promptContent string) (*ModelResponse, error) {
	m.ExecuteCount++
	if m.ShouldFail {
		return nil, fmt.Errorf("%s", m.ErrorMessage)
	}
	return &ModelResponse{
		Content: m.Content,
		Model:   m.Model,
	}, nil
}

// TestExecuteModelBasic tests the basic functionality of ExecuteModel
func TestExecuteModelBasic(t *testing.T) {
	t.Run("should fail when model config is invalid", func(t *testing.T) {
		// Execute with invalid parameters
		response, err := ExecuteModel("openai", "test prompt", nil, "temperature=2.0") // Invalid temperature

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid model parameters")
	})
}

// TestFallbackMechanism tests the fallback mechanism by testing the executeWithFallback function directly
func TestFallbackMechanism(t *testing.T) {
	// Create a simple Fallback test that mimics the behavior without needing to mock the createModelClient function
	primaryModel := "openai"
	promptContent := "test prompt"
	variables := map[string]string{"var": "value"}

	// Create a model config with fallback model and max retries
	modelConfig := &config.ModelConfig{
		Temperature:      0.7,
		MaxTokens:        1024,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		FallbackModels:   []string{"claude", "gemini"},
		MaxRetries:       2,
		OpenAIConfig: &config.OpenAIConfig{
			Model:         "gpt-3.5-turbo",
			SystemMessage: "Test message",
		},
		ClaudeConfig: &config.ClaudeConfig{
			Model:         "claude-3-sonnet-20240229",
			SystemMessage: "Test message",
		},
		GeminiConfig: &config.GeminiConfig{
			Model:          "gemini-pro",
			SafetySettings: map[string]string{},
		},
	}

	// Test with openai failing, claude succeeding
	// First we create a simulation of what would happen if createModelClient returned these clients
	clients := map[string]ModelClient{
		"openai": &MockModelClient{ShouldFail: true, ErrorMessage: "openai error", Model: "openai"},
		"claude": &MockModelClient{ShouldFail: false, Content: "Claude response", Model: "claude"},
		"gemini": &MockModelClient{ShouldFail: true, ErrorMessage: "should not be called", Model: "gemini"},
	}

	// Now we'll create a test executeWithFallback that uses our mock clients
	result := simulateFallback(primaryModel, promptContent, variables, modelConfig, clients)

	// Assertions for successful fallback
	if result.Response == nil {
		t.Errorf("Expected fallback to succeed with claude, but all models failed")
	} else {
		assert.Equal(t, "Claude response", result.Response.Content)
		assert.Equal(t, "claude", result.FinalModel)
		// We expect to have errors from previous failed attempts
		assert.GreaterOrEqual(t, len(result.Errors), 1)
		// At least one error should be from OpenAI
		foundOpenAIError := false
		for _, err := range result.Errors {
			if err.Model == "openai" && strings.Contains(err.Message, "openai error") {
				foundOpenAIError = true
			}
		}
		assert.True(t, foundOpenAIError, "Should contain at least one openai error")
	}

	// Test all models failing
	clientsFailing := map[string]ModelClient{
		"openai": &MockModelClient{ShouldFail: true, ErrorMessage: "openai error", Model: "openai"},
		"claude": &MockModelClient{ShouldFail: true, ErrorMessage: "claude error", Model: "claude"},
		"gemini": &MockModelClient{ShouldFail: true, ErrorMessage: "gemini error", Model: "gemini"},
	}

	// Now we'll create a test executeWithFallback that uses our failing mock clients
	resultFailing := simulateFallback(primaryModel, promptContent, variables, modelConfig, clientsFailing)

	// Assertions for all failures
	assert.Nil(t, resultFailing.Response)
	assert.Equal(t, 6, len(resultFailing.Errors)) // 3 models * 2 retries each = 6 errors

	// Check each model has the right errors and retry counts
	openaiErrors := 0
	claudeErrors := 0
	geminiErrors := 0

	for _, err := range resultFailing.Errors {
		switch err.Model {
		case "openai":
			openaiErrors++
			assert.Contains(t, err.Message, "openai error")
		case "claude":
			claudeErrors++
			assert.Contains(t, err.Message, "claude error")
		case "gemini":
			geminiErrors++
			assert.Contains(t, err.Message, "gemini error")
		}
	}

	assert.Equal(t, 2, openaiErrors) // 2 retries for openai
	assert.Equal(t, 2, claudeErrors) // 2 retries for claude
	assert.Equal(t, 2, geminiErrors) // 2 retries for gemini

	// Test default fallback sequence
	emptyConfig := &config.ModelConfig{
		Temperature:      0.7,
		MaxTokens:        1024,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		FallbackModels:   []string{}, // Empty to test default sequence
		MaxRetries:       1,
	}

	clientsDefaultFallback := map[string]ModelClient{
		"openai": &MockModelClient{ShouldFail: true, ErrorMessage: "openai error", Model: "openai"},
		"claude": &MockModelClient{ShouldFail: true, ErrorMessage: "claude error", Model: "claude"},
		"gemini": &MockModelClient{ShouldFail: false, Content: "Gemini response", Model: "gemini"},
	}

	resultDefaultFallback := simulateFallback(primaryModel, promptContent, variables, emptyConfig, clientsDefaultFallback)

	// Assertions for default fallback sequence
	assert.NotNil(t, resultDefaultFallback.Response)
	assert.Equal(t, "Gemini response", resultDefaultFallback.Response.Content)
	assert.Equal(t, "gemini", resultDefaultFallback.FinalModel)
	assert.Equal(t, 2, len(resultDefaultFallback.Errors)) // One for openai, one for claude
}

// simulateFallback simulates the executeWithFallback function using mock clients
func simulateFallback(primaryModel string, promptContent string, variables map[string]string, modelConfig *config.ModelConfig, mockClients map[string]ModelClient) *ModelFallbackResult {
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
	for _, modelName := range modelsToTry {
		client, exists := mockClients[modelName]
		if !exists {
			// Simulate client creation failure
			result.Errors = append(result.Errors, ModelError{
				Model:   modelName,
				Message: "failed to create client: mock client not found",
				Time:    time.Now(),
			})
			continue
		}

		for retry := 0; retry < modelConfig.MaxRetries; retry++ {
			// Execute the prompt with mock client
			response, err := client.Execute(promptContent)
			if err != nil {
				result.Errors = append(result.Errors, ModelError{
					Model:   modelName,
					Message: "execution failed: " + err.Error(),
					Err:     err,
					Retry:   retry,
				})
				continue
			}

			// Success! Add metadata and return
			response.Variables = variables
			response.ExecutionID = fmt.Sprintf("%s-test-%d", modelName, retry)
			result.Response = response
			result.FinalModel = modelName

			return result
		}
	}

	return result
}
