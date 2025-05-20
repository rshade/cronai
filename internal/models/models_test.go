package models

import (
	"errors"
	"fmt"
	"os"
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
	Variables    map[string]string
}

func (m *MockModelClient) Execute(_ string) (*ModelResponse, error) {
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
		assert.Contains(t, err.Error(), "temperature must be between 0 and 1")
	})

	t.Run("should execute successfully with valid config", func(t *testing.T) {
		// Set up env vars for the test
		oldOpenAIKey := os.Getenv("OPENAI_API_KEY")
		if err := os.Setenv("OPENAI_API_KEY", "test-key"); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("OPENAI_API_KEY", oldOpenAIKey); err != nil {
				t.Fatal(err)
			}
		}()

		// Mock the createModelClient function
		originalCreateModelClient := createModelClient
		createModelClient = func(_ string, _ *config.ModelConfig) (ModelClient, error) {
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
		if err := os.Setenv("OPENAI_API_KEY", "test-key"); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("OPENAI_API_KEY", oldOpenAIKey); err != nil {
				t.Fatal(err)
			}
		}()

		// Mock the createModelClient function
		originalCreateModelClient := createModelClient
		createModelClient = func(_ string, _ *config.ModelConfig) (ModelClient, error) {
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
	// No need to store original function since we'll use defaultCreateModelClient directly

	t.Run("should fail for unknown model", func(t *testing.T) {
		client, err := defaultCreateModelClient("unknown", nil)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "unsupported model")
	})

	t.Run("should create OpenAI client", func(t *testing.T) {
		// Set up env vars for the test
		oldOpenAIKey := os.Getenv("OPENAI_API_KEY")
		if err := os.Setenv("OPENAI_API_KEY", "test-key"); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("OPENAI_API_KEY", oldOpenAIKey); err != nil {
				t.Fatal(err)
			}
		}()

		client, err := defaultCreateModelClient("openai", &config.ModelConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.IsType(t, &OpenAIClient{}, client)
	})

	t.Run("should create Claude client", func(t *testing.T) {
		// Set up env vars for the test
		oldClaudeKey := os.Getenv("ANTHROPIC_API_KEY")
		if err := os.Setenv("ANTHROPIC_API_KEY", "test-key"); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("ANTHROPIC_API_KEY", oldClaudeKey); err != nil {
				t.Fatal(err)
			}
		}()

		client, err := defaultCreateModelClient("claude", &config.ModelConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.IsType(t, &ClaudeClient{}, client)
	})

	t.Run("should create Gemini client", func(t *testing.T) {
		// Set up env vars for the test
		oldGeminiKey := os.Getenv("GOOGLE_API_KEY")
		if err := os.Setenv("GOOGLE_API_KEY", "test-key"); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("GOOGLE_API_KEY", oldGeminiKey); err != nil {
				t.Fatal(err)
			}
		}()

		client, err := defaultCreateModelClient("gemini", &config.ModelConfig{})
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
		createModelClient = func(modelName string, _ *config.ModelConfig) (ModelClient, error) {
			if modelName == "openai" {
				return &MockModelClient{
					Content: "OpenAI success",
					Model:   "openai",
				}, nil
			}
			return nil, errors.New("should not be called")
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig, "test")

		assert.NotNil(t, result.Response)
		assert.Equal(t, "OpenAI success", result.Response.Content)
		assert.Equal(t, "openai", result.FinalModel)
		assert.Empty(t, result.Errors)
	})

	t.Run("should fallback to claude when openai fails", func(t *testing.T) {
		// Mock createModelClient
		createModelClient = func(modelName string, _ *config.ModelConfig) (ModelClient, error) {
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
				return nil, fmt.Errorf("unexpected model: %s", modelName)
			}
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig, "test")

		assert.NotNil(t, result.Response)
		assert.Equal(t, "Claude success", result.Response.Content)
		assert.Equal(t, "claude", result.FinalModel)
		assert.Len(t, result.Errors, 2) // Two attempts for openai
		assert.Equal(t, "openai", result.Errors[0].Model)
		assert.Equal(t, "openai", result.Errors[1].Model)
	})

	t.Run("should try all fallbacks and fail if all fail", func(t *testing.T) {
		// Mock createModelClient - all models fail
		createModelClient = func(modelName string, _ *config.ModelConfig) (ModelClient, error) {
			return &MockModelClient{
				ShouldFail:   true,
				ErrorMessage: fmt.Sprintf("%s error", modelName),
				Model:        modelName,
			}, nil
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig, "test")

		assert.Nil(t, result.Response)
		assert.Len(t, result.Errors, 6) // 2 attempts per model * 3 models
		assert.Equal(t, "openai", result.Errors[0].Model)
		assert.Equal(t, "openai", result.Errors[1].Model)
		assert.Equal(t, "claude", result.Errors[2].Model)
		assert.Equal(t, "claude", result.Errors[3].Model)
		assert.Equal(t, "gemini", result.Errors[4].Model)
		assert.Equal(t, "gemini", result.Errors[5].Model)
	})

	t.Run("should retry within the same model before trying fallbacks", func(t *testing.T) {
		// Counter to track retry attempts for openai
		openaiAttempts := 0

		// Mock createModelClient - openai fails on first attempt
		createModelClient = func(modelName string, _ *config.ModelConfig) (ModelClient, error) {
			if modelName == "openai" {
				openaiAttempts++
				if openaiAttempts == 1 { // Fail on first attempt only
					return &MockModelClient{
						ShouldFail:   true,
						ErrorMessage: "openai error",
						Model:        "openai",
					}, nil
				}
				return &MockModelClient{
					Content: "OpenAI retry success",
					Model:   "openai",
				}, nil
			}
			return nil, fmt.Errorf("unexpected model: %s", modelName)
		}

		result := executeWithFallback("openai", "test prompt", nil, modelConfig, "test")

		assert.NotNil(t, result.Response)
		assert.Equal(t, "OpenAI retry success", result.Response.Content)
		assert.Equal(t, "openai", result.FinalModel)
		assert.Len(t, result.Errors, 1) // One failed attempt
		assert.Equal(t, 2, openaiAttempts)
	})
}

// TestModelParameterParsing tests parameter parsing
func TestModelParameterParsing(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid empty params",
			params:    "",
			expectErr: false,
		},
		{
			name:      "valid single param",
			params:    "temperature=0.7",
			expectErr: false,
		},
		{
			name:      "valid multiple params",
			params:    "temperature=0.7,max_tokens=500",
			expectErr: false,
		},
		{
			name:      "valid nested model-specific params",
			params:    "temperature=0.7,openai.model=gpt-4",
			expectErr: false,
		},
		{
			name:      "invalid param format (no value)",
			params:    "temperature",
			expectErr: true,
			errMsg:    "invalid parameter format",
		},
		{
			name:      "invalid param format (no equals)",
			params:    "temperature 0.7",
			expectErr: true,
			errMsg:    "invalid parameter format",
		},
		{
			name:      "invalid temperature value",
			params:    "temperature=abc",
			expectErr: true,
			errMsg:    "invalid temperature value",
		},
		{
			name:      "temperature too high",
			params:    "temperature=2.0",
			expectErr: true,
			errMsg:    "temperature must be between 0 and 1",
		},
		{
			name:      "temperature negative",
			params:    "temperature=-0.5",
			expectErr: true,
			errMsg:    "temperature must be between 0 and 1",
		},
		{
			name:      "invalid max_tokens",
			params:    "max_tokens=abc",
			expectErr: true,
			errMsg:    "invalid max_tokens value",
		},
		{
			name:      "invalid top_p",
			params:    "top_p=2.0",
			expectErr: true,
			errMsg:    "top_p must be between 0 and 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse parameters like ExecuteModel does
			modelConfig := config.DefaultModelConfig()
			modelConfig.LoadFromEnvironment()
			var err error
			if tt.params != "" {
				params, parseErr := config.ParseModelParams(tt.params)
				if parseErr == nil {
					err = modelConfig.UpdateFromParams(params)
				} else {
					err = parseErr
				}
			}

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

// TestValidateConfiguration tests configuration validation
func TestValidateConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(config *config.ModelConfig)
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid basic config",
			setup:     func(_ *config.ModelConfig) {},
			expectErr: false,
		},
		{
			name: "invalid temperature (too high)",
			setup: func(config *config.ModelConfig) {
				config.Temperature = 2.0
			},
			expectErr: true,
			errMsg:    "temperature must be between 0 and 1",
		},
		{
			name: "invalid temperature (too low)",
			setup: func(config *config.ModelConfig) {
				config.Temperature = -0.5
			},
			expectErr: true,
			errMsg:    "temperature must be between 0 and 1",
		},
		{
			name: "invalid top_p (too high)",
			setup: func(config *config.ModelConfig) {
				config.TopP = 2.0
			},
			expectErr: true,
			errMsg:    "top_p must be between 0 and 1",
		},
		{
			name: "invalid top_p (too low)",
			setup: func(config *config.ModelConfig) {
				config.TopP = -0.5
			},
			expectErr: true,
			errMsg:    "top_p must be between 0 and 1",
		},
		{
			name: "invalid frequency penalty (too high)",
			setup: func(config *config.ModelConfig) {
				config.FrequencyPenalty = 3.0
			},
			expectErr: true,
			errMsg:    "frequency_penalty must be between -2.0 and 2.0",
		},
		{
			name: "invalid presence penalty (too low)",
			setup: func(config *config.ModelConfig) {
				config.PresencePenalty = -3.0
			},
			expectErr: true,
			errMsg:    "presence_penalty must be between -2.0 and 2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.DefaultModelConfig()
			tt.setup(config)

			err := config.Validate()

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
