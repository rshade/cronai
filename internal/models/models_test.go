package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockModelClient is a simple mocked version of the ModelClient interface
type MockModelClient struct {
	ShouldFail   bool
	ErrorMessage string
	Content      string
	Model        string
	ExecuteCount int
	Variables    map[string]string
}

func (m *MockModelClient) Execute(promptContent string) (*ModelResponse, error) {
	m.ExecuteCount++
	if m.ShouldFail {
		return nil, fmt.Errorf("%s", m.ErrorMessage)
	}
	return &ModelResponse{
		Content:   m.Content,
		Model:     m.Model,
		Variables: m.Variables,
		Timestamp: time.Now(),
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

	t.Run("should execute successfully with valid config", func(t *testing.T) {
		// Set up env vars for the test
		oldOpenAIKey := os.Getenv("OPENAI_API_KEY")
		os.Setenv("OPENAI_API_KEY", "test-key")
		defer os.Setenv("OPENAI_API_KEY", oldOpenAIKey)

		// Mock the createModelClient function
		originalCreateModelClient := createModelClient
		createModelClient = func(modelName string, modelConfig *config.ModelConfig) (ModelClient, error) {
			return &MockModelClient{
				Content: "Test response",
				Model:   "test-model",
			}, nil
		}
		defer func() { createModelClient = originalCreateModelClient }()

		// Execute
		response, err := ExecuteModel("openai", "test prompt", nil, "temperature=0.5")

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Test response", response.Content)
	})

	t.Run("should handle variables correctly", func(t *testing.T) {
		// Set up env vars for the test
		oldOpenAIKey := os.Getenv("OPENAI_API_KEY")
		os.Setenv("OPENAI_API_KEY", "test-key")
		defer os.Setenv("OPENAI_API_KEY", oldOpenAIKey)

		// Mock the createModelClient function
		originalCreateModelClient := createModelClient
		createModelClient = func(modelName string, modelConfig *config.ModelConfig) (ModelClient, error) {
			return &MockModelClient{
				Content: "Test response",
				Model:   "test-model",
			}, nil
		}
		defer func() { createModelClient = originalCreateModelClient }()

		// Execute with variables
		variables := map[string]string{"key": "value"}
		response, err := ExecuteModel("openai", "test prompt", variables, "")

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, variables, response.Variables)
	})
}

// TestCreateModelClient tests the createModelClient function
func TestCreateModelClient(t *testing.T) {
	// Store original function to restore later
	originalCreateModelClient := createModelClient
	defer func() { createModelClient = originalCreateModelClient }()

	t.Run("should fail for unknown model", func(t *testing.T) {
		client, err := originalCreateModelClient("unknown", nil)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "unsupported model")
	})

	t.Run("should create OpenAI client", func(t *testing.T) {
		// Set up env vars for the test
		oldOpenAIKey := os.Getenv("OPENAI_API_KEY")
		os.Setenv("OPENAI_API_KEY", "test-key")
		defer os.Setenv("OPENAI_API_KEY", oldOpenAIKey)

		client, err := originalCreateModelClient("openai", &config.ModelConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.IsType(t, &OpenAIClient{}, client)
	})

	t.Run("should create Claude client", func(t *testing.T) {
		// Set up env vars for the test
		oldClaudeKey := os.Getenv("ANTHROPIC_API_KEY")
		os.Setenv("ANTHROPIC_API_KEY", "test-key")
		defer os.Setenv("ANTHROPIC_API_KEY", oldClaudeKey)

		client, err := originalCreateModelClient("claude", &config.ModelConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.IsType(t, &ClaudeClient{}, client)
	})

	t.Run("should create Gemini client", func(t *testing.T) {
		// Set up env vars for the test
		oldGeminiKey := os.Getenv("GOOGLE_API_KEY")
		os.Setenv("GOOGLE_API_KEY", "test-key")
		defer os.Setenv("GOOGLE_API_KEY", oldGeminiKey)

		client, err := originalCreateModelClient("gemini", &config.ModelConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.IsType(t, &GeminiClient{}, client)
	})
}

// TestFallbackMechanism tests the fallback mechanism
func TestFallbackMechanism(t *testing.T) {
	// Store original function to restore later
	originalCreateModelClient := createModelClient
	defer func() { createModelClient = originalCreateModelClient }()

	// Create a model config with fallback models and max retries
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

	t.Run("should succeed with primary model", func(t *testing.T) {
		// Mock createModelClient to return successful primary model
		createModelClient = func(modelName string, config *config.ModelConfig) (ModelClient, error) {
			if modelName == "openai" {
				return &MockModelClient{
					Content: "OpenAI success",
					Model:   "openai",
				}, nil
			}
			return nil, errors.New("should not be called")
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig)

		assert.NotNil(t, result.Response)
		assert.Equal(t, "OpenAI success", result.Response.Content)
		assert.Equal(t, "openai", result.FinalModel)
		assert.Empty(t, result.Errors)
	})

	t.Run("should fallback to claude when openai fails", func(t *testing.T) {
		// Mock createModelClient
		createModelClient = func(modelName string, config *config.ModelConfig) (ModelClient, error) {
			switch modelName {
			case "openai":
				return &MockModelClient{
					ShouldFail:   true,
					ErrorMessage: "openai error",
					Model:        "openai",
				}, nil
			case "claude":
				return &MockModelClient{
					Content: "Claude success",
					Model:   "claude",
				}, nil
			default:
				return nil, errors.New("should not be called")
			}
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig)

		assert.NotNil(t, result.Response)
		assert.Equal(t, "Claude success", result.Response.Content)
		assert.Equal(t, "claude", result.FinalModel)
		assert.GreaterOrEqual(t, len(result.Errors), 2) // At least 2 errors from OpenAI retries
	})

	t.Run("should use all fallbacks when all fail", func(t *testing.T) {
		// Mock createModelClient to return failing clients
		createModelClient = func(modelName string, config *config.ModelConfig) (ModelClient, error) {
			return &MockModelClient{
				ShouldFail:   true,
				ErrorMessage: fmt.Sprintf("%s error", modelName),
				Model:        modelName,
			}, nil
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig)

		assert.Nil(t, result.Response)
		assert.Equal(t, "", result.FinalModel)
		assert.Equal(t, 6, len(result.Errors)) // 3 models * 2 retries each
	})

	t.Run("should use default fallback sequence when not configured", func(t *testing.T) {
		// Model config without fallback models
		emptyConfig := &config.ModelConfig{
			Temperature:      0.7,
			MaxTokens:        1024,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			FallbackModels:   []string{}, // Empty to test default sequence
			MaxRetries:       1,
		}

		// Mock createModelClient
		createModelClient = func(modelName string, config *config.ModelConfig) (ModelClient, error) {
			switch modelName {
			case "openai":
				return &MockModelClient{
					ShouldFail:   true,
					ErrorMessage: "openai error",
					Model:        "openai",
				}, nil
			case "claude":
				return &MockModelClient{
					ShouldFail:   true,
					ErrorMessage: "claude error",
					Model:        "claude",
				}, nil
			case "gemini":
				return &MockModelClient{
					Content: "Gemini success",
					Model:   "gemini",
				}, nil
			default:
				return nil, errors.New("unexpected model")
			}
		}

		result := executeWithFallback("openai", "test prompt", nil, emptyConfig)

		assert.NotNil(t, result.Response)
		assert.Equal(t, "Gemini success", result.Response.Content)
		assert.Equal(t, "gemini", result.FinalModel)
		assert.Equal(t, 2, len(result.Errors)) // One for openai, one for claude
	})

	t.Run("should handle client creation failure", func(t *testing.T) {
		// Mock createModelClient to return error
		createModelClient = func(modelName string, config *config.ModelConfig) (ModelClient, error) {
			if modelName == "openai" {
				return nil, errors.New("failed to create client")
			}
			return &MockModelClient{
				Content: "Success",
				Model:   modelName,
			}, nil
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig)

		assert.NotNil(t, result.Response)
		assert.Equal(t, "Success", result.Response.Content)
		assert.Equal(t, "claude", result.FinalModel)
		assert.GreaterOrEqual(t, len(result.Errors), 1)
	})
}

// TestModelError tests the ModelError type
func TestModelError(t *testing.T) {
	modelErr := ModelError{
		Model:   "test-model",
		Message: "test error",
		Err:     errors.New("underlying error"),
		Retry:   1,
		Time:    time.Now(),
	}

	str := modelErr.Error()
	assert.Contains(t, str, "test-model")
	assert.Contains(t, str, "test error")
	assert.Contains(t, str, "attempt 2")
	assert.Contains(t, str, "underlying error")
}

// TestGenerateExecutionID tests the generateExecutionID function
func TestGenerateExecutionID(t *testing.T) {
	id1 := generateExecutionID("openai", "test-prompt")
	id2 := generateExecutionID("openai", "test-prompt")

	// IDs should be different
	assert.NotEqual(t, id1, id2)

	// IDs should contain model and prompt
	assert.Contains(t, id1, "openai")
	assert.Contains(t, id1, "test-prompt")
}

// TestGetDefaultFallbackSequence tests the getDefaultFallbackSequence function
func TestGetDefaultFallbackSequence(t *testing.T) {
	tests := []struct {
		name     string
		primary  string
		expected []string
	}{
		{
			name:     "openai primary",
			primary:  "openai",
			expected: []string{"claude", "gemini"},
		},
		{
			name:     "claude primary",
			primary:  "claude",
			expected: []string{"openai", "gemini"},
		},
		{
			name:     "gemini primary",
			primary:  "gemini",
			expected: []string{"openai", "claude"},
		},
		{
			name:     "unknown primary",
			primary:  "unknown",
			expected: []string{"openai", "claude", "gemini"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDefaultFallbackSequence(tt.primary)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestModelValidateParameters tests the parameter validation logic
func TestModelValidateParameters(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid temperature",
			params:    "temperature=0.5",
			expectErr: false,
		},
		{
			name:      "invalid temperature too high",
			params:    "temperature=2.5",
			expectErr: true,
			errMsg:    "temperature must be between 0 and 2",
		},
		{
			name:      "invalid temperature negative",
			params:    "temperature=-0.5",
			expectErr: true,
			errMsg:    "temperature must be between 0 and 2",
		},
		{
			name:      "valid max_tokens",
			params:    "max_tokens=1000",
			expectErr: false,
		},
		{
			name:      "invalid max_tokens negative",
			params:    "max_tokens=-100",
			expectErr: true,
			errMsg:    "max_tokens must be positive",
		},
		{
			name:      "valid top_p",
			params:    "top_p=0.9",
			expectErr: false,
		},
		{
			name:      "invalid top_p too high",
			params:    "top_p=1.5",
			expectErr: true,
			errMsg:    "top_p must be between 0 and 1",
		},
		{
			name:      "valid frequency_penalty",
			params:    "frequency_penalty=0.5",
			expectErr: false,
		},
		{
			name:      "invalid frequency_penalty too high",
			params:    "frequency_penalty=2.5",
			expectErr: true,
			errMsg:    "frequency_penalty must be between -2 and 2",
		},
		{
			name:      "valid presence_penalty",
			params:    "presence_penalty=-0.5",
			expectErr: false,
		},
		{
			name:      "invalid presence_penalty too low",
			params:    "presence_penalty=-2.5",
			expectErr: true,
			errMsg:    "presence_penalty must be between -2 and 2",
		},
		{
			name:      "multiple parameters",
			params:    "temperature=0.7,max_tokens=500,top_p=0.9",
			expectErr: false,
		},
		{
			name:      "invalid format",
			params:    "invalid_format",
			expectErr: true,
			errMsg:    "parameter must be in key=value format",
		},
		{
			name:      "unknown parameter",
			params:    "unknown_param=123",
			expectErr: true,
			errMsg:    "unknown parameter",
		},
		{
			name:      "model specific parameter",
			params:    "openai.system_message=Hello",
			expectErr: false,
		},
		{
			name:      "empty parameters",
			params:    "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse parameters like ExecuteModel does
			modelConfig := config.NewModelConfig()
			err := modelConfig.ParseParameters(tt.params)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
