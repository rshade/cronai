package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultModelConfig(t *testing.T) {
	config := DefaultModelConfig()

	// Test common parameters
	if config.Temperature != 0.7 {
		t.Errorf("Expected default Temperature to be 0.7, got %f", config.Temperature)
	}
	if config.MaxTokens != 1024 {
		t.Errorf("Expected default MaxTokens to be 1024, got %d", config.MaxTokens)
	}
	if config.TopP != 1.0 {
		t.Errorf("Expected default TopP to be 1.0, got %f", config.TopP)
	}

	// Test model-specific parameters
	if config.OpenAIConfig.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected default OpenAI model to be gpt-3.5-turbo, got %s", config.OpenAIConfig.Model)
	}
	if config.ClaudeConfig.Model != "claude-3-sonnet-20240229" {
		t.Errorf("Expected default Claude model to be claude-3-sonnet-20240229, got %s", config.ClaudeConfig.Model)
	}
	if config.GeminiConfig.Model != "gemini-pro" {
		t.Errorf("Expected default Gemini model to be gemini-pro, got %s", config.GeminiConfig.Model)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Save original environment variables to restore after test
	origTemperature := os.Getenv("MODEL_TEMPERATURE")
	origMaxTokens := os.Getenv("MODEL_MAX_TOKENS")
	origOpenAIModel := os.Getenv("OPENAI_MODEL")
	origClaudeModel := os.Getenv("CLAUDE_MODEL")
	origGeminiModel := os.Getenv("GEMINI_MODEL")
	origSafetySettings := os.Getenv("GEMINI_SAFETY_SETTINGS")

	// Set environment variables for testing
	if err := os.Setenv("MODEL_TEMPERATURE", "0.5"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("MODEL_MAX_TOKENS", "2048"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("OPENAI_MODEL", "gpt-4"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("CLAUDE_MODEL", "claude-3-opus-20240229"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("GEMINI_MODEL", "gemini-1.5-pro"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("GEMINI_SAFETY_SETTINGS", "harmful=block,harassment=warn"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Restore environment variables after test
	defer func() {
		// Ignoring errors in defer as there's not much we can do about them
		_ = os.Setenv("MODEL_TEMPERATURE", origTemperature)
		_ = os.Setenv("MODEL_MAX_TOKENS", origMaxTokens)
		_ = os.Setenv("OPENAI_MODEL", origOpenAIModel)
		_ = os.Setenv("CLAUDE_MODEL", origClaudeModel)
		_ = os.Setenv("GEMINI_MODEL", origGeminiModel)
		_ = os.Setenv("GEMINI_SAFETY_SETTINGS", origSafetySettings)
	}()

	// Create and load config from environment
	config := DefaultModelConfig()
	config.LoadFromEnvironment()

	// Test loaded values
	if config.Temperature != 0.5 {
		t.Errorf("Expected Temperature to be 0.5, got %f", config.Temperature)
	}
	if config.MaxTokens != 2048 {
		t.Errorf("Expected MaxTokens to be 2048, got %d", config.MaxTokens)
	}
	if config.OpenAIConfig.Model != "gpt-4" {
		t.Errorf("Expected OpenAI model to be gpt-4, got %s", config.OpenAIConfig.Model)
	}
	if config.ClaudeConfig.Model != "claude-3-opus-20240229" {
		t.Errorf("Expected Claude model to be claude-3-opus-20240229, got %s", config.ClaudeConfig.Model)
	}
	if config.GeminiConfig.Model != "gemini-1.5-pro" {
		t.Errorf("Expected Gemini model to be gemini-1.5-pro, got %s", config.GeminiConfig.Model)
	}

	// Test safety settings
	expectedSettings := map[string]string{
		"harmful":    "block",
		"harassment": "warn",
	}
	if !reflect.DeepEqual(config.GeminiConfig.SafetySettings, expectedSettings) {
		t.Errorf("Expected safety settings %v, got %v", expectedSettings, config.GeminiConfig.SafetySettings)
	}
}

func TestLoadFromEnvironmentWithErrors(t *testing.T) {
	// Test invalid environment variable values
	testCases := []struct {
		name      string
		envVar    string
		value     string
		setupFunc func() func()
	}{
		{
			name:   "invalid temperature",
			envVar: "MODEL_TEMPERATURE",
			value:  "invalid",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_TEMPERATURE")
				_ = os.Setenv("MODEL_TEMPERATURE", "invalid")
				return func() { _ = os.Setenv("MODEL_TEMPERATURE", orig) }
			},
		},
		{
			name:   "invalid max tokens",
			envVar: "MODEL_MAX_TOKENS",
			value:  "not_a_number",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_MAX_TOKENS")
				_ = os.Setenv("MODEL_MAX_TOKENS", "not_a_number")
				return func() { _ = os.Setenv("MODEL_MAX_TOKENS", orig) }
			},
		},
		{
			name:   "invalid top_p",
			envVar: "MODEL_TOP_P",
			value:  "abc",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_TOP_P")
				_ = os.Setenv("MODEL_TOP_P", "abc")
				return func() { _ = os.Setenv("MODEL_TOP_P", orig) }
			},
		},
		{
			name:   "invalid frequency penalty",
			envVar: "MODEL_FREQUENCY_PENALTY",
			value:  "invalid",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_FREQUENCY_PENALTY")
				_ = os.Setenv("MODEL_FREQUENCY_PENALTY", "invalid")
				return func() { _ = os.Setenv("MODEL_FREQUENCY_PENALTY", orig) }
			},
		},
		{
			name:   "invalid presence penalty",
			envVar: "MODEL_PRESENCE_PENALTY",
			value:  "invalid",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_PRESENCE_PENALTY")
				_ = os.Setenv("MODEL_PRESENCE_PENALTY", "invalid")
				return func() { _ = os.Setenv("MODEL_PRESENCE_PENALTY", orig) }
			},
		},
		{
			name:   "invalid max retries",
			envVar: "MODEL_MAX_RETRIES",
			value:  "invalid",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_MAX_RETRIES")
				_ = os.Setenv("MODEL_MAX_RETRIES", "invalid")
				return func() { _ = os.Setenv("MODEL_MAX_RETRIES", orig) }
			},
		},
		{
			name:   "fallback models",
			envVar: "MODEL_FALLBACK_MODELS",
			value:  "claude|openai|gemini",
			setupFunc: func() func() {
				orig := os.Getenv("MODEL_FALLBACK_MODELS")
				_ = os.Setenv("MODEL_FALLBACK_MODELS", "claude|openai|gemini")
				return func() { _ = os.Setenv("MODEL_FALLBACK_MODELS", orig) }
			},
		},
		{
			name:   "openai system message",
			envVar: "OPENAI_SYSTEM_MESSAGE",
			value:  "OpenAI system message",
			setupFunc: func() func() {
				orig := os.Getenv("OPENAI_SYSTEM_MESSAGE")
				_ = os.Setenv("OPENAI_SYSTEM_MESSAGE", "OpenAI system message")
				return func() { _ = os.Setenv("OPENAI_SYSTEM_MESSAGE", orig) }
			},
		},
		{
			name:   "claude system message",
			envVar: "CLAUDE_SYSTEM_MESSAGE",
			value:  "Claude system message",
			setupFunc: func() func() {
				orig := os.Getenv("CLAUDE_SYSTEM_MESSAGE")
				_ = os.Setenv("CLAUDE_SYSTEM_MESSAGE", "Claude system message")
				return func() { _ = os.Setenv("CLAUDE_SYSTEM_MESSAGE", orig) }
			},
		},
		{
			name:   "all environment variables",
			envVar: "ALL",
			value:  "various",
			setupFunc: func() func() {
				// Save all originals
				originals := map[string]string{
					"MODEL_TOP_P":             os.Getenv("MODEL_TOP_P"),
					"MODEL_FREQUENCY_PENALTY": os.Getenv("MODEL_FREQUENCY_PENALTY"),
					"MODEL_PRESENCE_PENALTY":  os.Getenv("MODEL_PRESENCE_PENALTY"),
					"MODEL_MAX_RETRIES":       os.Getenv("MODEL_MAX_RETRIES"),
				}

				// Set new values
				_ = os.Setenv("MODEL_TOP_P", "0.9")
				_ = os.Setenv("MODEL_FREQUENCY_PENALTY", "0.5")
				_ = os.Setenv("MODEL_PRESENCE_PENALTY", "0.5")
				_ = os.Setenv("MODEL_MAX_RETRIES", "3")

				return func() {
					for k, v := range originals {
						_ = os.Setenv(k, v)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cleanup := tc.setupFunc()
			defer cleanup()

			config := DefaultModelConfig()
			// LoadFromEnvironment should handle errors without panicking
			config.LoadFromEnvironment()

			// Check valid values were set correctly
			if tc.name == "fallback models" {
				expected := []string{"claude|openai|gemini"}
				if !reflect.DeepEqual(config.FallbackModels, expected) {
					t.Errorf("Expected fallback models %v, got %v", expected, config.FallbackModels)
				}
			}

			if tc.name == "openai system message" {
				if config.OpenAIConfig.SystemMessage != "OpenAI system message" {
					t.Errorf("Expected OpenAI system message, got %s", config.OpenAIConfig.SystemMessage)
				}
			}

			if tc.name == "claude system message" {
				if config.ClaudeConfig.SystemMessage != "Claude system message" {
					t.Errorf("Expected Claude system message, got %s", config.ClaudeConfig.SystemMessage)
				}
			}

			if tc.name == "all environment variables" {
				if config.TopP != 0.9 {
					t.Errorf("Expected TopP 0.9, got %f", config.TopP)
				}
				if config.FrequencyPenalty != 0.5 {
					t.Errorf("Expected FrequencyPenalty 0.5, got %f", config.FrequencyPenalty)
				}
				if config.PresencePenalty != 0.5 {
					t.Errorf("Expected PresencePenalty 0.5, got %f", config.PresencePenalty)
				}
				if config.MaxRetries != 3 {
					t.Errorf("Expected MaxRetries 3, got %d", config.MaxRetries)
				}
			}
		})
	}
}

func TestParseModelParams(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]string{},
			wantErr:  false,
		},
		{
			name:  "valid parameters",
			input: "temperature=0.8,max_tokens=2048,model=gpt-4",
			expected: map[string]string{
				"temperature": "0.8",
				"max_tokens":  "2048",
				"model":       "gpt-4",
			},
			wantErr: false,
		},
		{
			name:  "model-specific parameters",
			input: "openai.model=gpt-4,claude.system_message=Custom message",
			expected: map[string]string{
				"openai.model":          "gpt-4",
				"claude.system_message": "Custom message",
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "temperature=0.8,invalid-format",
			wantErr: true,
		},
		{
			name:  "parameters with spaces",
			input: "temperature = 0.8 , max_tokens = 2048",
			expected: map[string]string{
				"temperature": "0.8",
				"max_tokens":  "2048",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseModelParams(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseModelParams() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("ParseModelParams() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestUpdateFromParams(t *testing.T) {
	testCases := []struct {
		name       string
		params     map[string]string
		checkFunc  func(*ModelConfig) bool
		errMessage string
	}{
		{
			name: "common parameters",
			params: map[string]string{
				"temperature": "0.5",
				"max_tokens":  "2048",
				"top_p":       "0.9",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.Temperature == 0.5 && c.MaxTokens == 2048 && c.TopP == 0.9
			},
			errMessage: "Common parameters not updated correctly",
		},
		{
			name: "invalid temperature",
			params: map[string]string{
				"temperature": "invalid",
			},
			checkFunc:  nil, // We expect an error
			errMessage: "Should reject invalid temperature",
		},
		{
			name: "out of range temperature",
			params: map[string]string{
				"temperature": "2.0",
			},
			checkFunc:  nil, // We expect an error
			errMessage: "Should reject out of range temperature",
		},
		{
			name: "model parameters",
			params: map[string]string{
				"model": "gpt-4",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.OpenAIConfig.Model == "gpt-4" &&
					c.ClaudeConfig.Model == "gpt-4" &&
					c.GeminiConfig.Model == "gpt-4"
			},
			errMessage: "Model parameter not applied to all models",
		},
		{
			name: "OpenAI specific parameters",
			params: map[string]string{
				"openai.model":          "gpt-4",
				"openai.system_message": "OpenAI specific message",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.OpenAIConfig.Model == "gpt-4" &&
					c.OpenAIConfig.SystemMessage == "OpenAI specific message"
			},
			errMessage: "OpenAI specific parameters not updated correctly",
		},
		{
			name: "Claude specific parameters",
			params: map[string]string{
				"claude.model":          "claude-3-opus-20240229",
				"claude.system_message": "Claude specific message",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.ClaudeConfig.Model == "claude-3-opus-20240229" &&
					c.ClaudeConfig.SystemMessage == "Claude specific message"
			},
			errMessage: "Claude specific parameters not updated correctly",
		},
		{
			name: "Gemini specific parameters",
			params: map[string]string{
				"gemini.model":          "gemini-1.5-pro",
				"gemini.safety_setting": "harmful=block",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.GeminiConfig.Model == "gemini-1.5-pro" &&
					c.GeminiConfig.SafetySettings["harmful"] == "block"
			},
			errMessage: "Gemini specific parameters not updated correctly",
		},
		{
			name: "Multiple safety settings",
			params: map[string]string{
				"gemini.safety_setting": "harmful=block",
			},
			checkFunc: func(c *ModelConfig) bool {
				if err := c.UpdateFromParams(map[string]string{"gemini.safety_setting": "harassment=warn"}); err != nil {
					return false
				}
				return c.GeminiConfig.SafetySettings["harmful"] == "block" &&
					c.GeminiConfig.SafetySettings["harassment"] == "warn"
			},
			errMessage: "Multiple safety settings not handled correctly",
		},
		{
			name: "Mixed common and model-specific parameters",
			params: map[string]string{
				"temperature":           "0.6",
				"max_tokens":            "1500",
				"openai.model":          "gpt-4-turbo",
				"claude.model":          "claude-3-haiku-20240307",
				"gemini.model":          "gemini-1.5-pro",
				"system_message":        "Generic system message",
				"claude.system_message": "Claude-specific message",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.Temperature == 0.6 &&
					c.MaxTokens == 1500 &&
					c.OpenAIConfig.Model == "gpt-4-turbo" &&
					c.ClaudeConfig.Model == "claude-3-haiku-20240307" &&
					c.GeminiConfig.Model == "gemini-1.5-pro" &&
					c.OpenAIConfig.SystemMessage == "Generic system message" &&
					c.ClaudeConfig.SystemMessage == "Claude-specific message"
			},
			errMessage: "Mixed common and model-specific parameters not handled correctly",
		},
		{
			name: "Model-specific parameters override common parameters",
			params: map[string]string{
				"model":                 "gpt-3.5-turbo",
				"openai.model":          "gpt-4",
				"system_message":        "Generic system message",
				"openai.system_message": "OpenAI specific message",
			},
			checkFunc: func(c *ModelConfig) bool {
				return c.OpenAIConfig.Model == "gpt-4" &&
					c.ClaudeConfig.Model == "gpt-3.5-turbo" &&
					c.GeminiConfig.Model == "gpt-3.5-turbo" &&
					c.OpenAIConfig.SystemMessage == "OpenAI specific message" &&
					c.ClaudeConfig.SystemMessage == "Generic system message"
			},
			errMessage: "Model-specific parameters should override common parameters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultModelConfig()
			err := config.UpdateFromParams(tc.params)

			if tc.checkFunc == nil {
				// We expect an error
				if err == nil {
					t.Errorf("UpdateFromParams() expected error but got nil")
				}
			} else {
				// We don't expect an error
				if err != nil {
					t.Errorf("UpdateFromParams() error = %v", err)
					return
				}

				// Check if the parameters were updated correctly
				if !tc.checkFunc(config) {
					t.Errorf("%s", tc.errMessage)
				}
			}
		})
	}
}

func TestUpdateFromParamsErrorHandling(t *testing.T) {
	testCases := []struct {
		name   string
		params map[string]string
	}{
		{
			name: "invalid max tokens",
			params: map[string]string{
				"max_tokens": "invalid",
			},
		},
		{
			name: "negative max tokens",
			params: map[string]string{
				"max_tokens": "-10",
			},
		},
		{
			name: "invalid top_p",
			params: map[string]string{
				"top_p": "invalid",
			},
		},
		{
			name: "out of range top_p low",
			params: map[string]string{
				"top_p": "-0.5",
			},
		},
		{
			name: "out of range top_p high",
			params: map[string]string{
				"top_p": "2.0",
			},
		},
		{
			name: "invalid frequency penalty",
			params: map[string]string{
				"frequency_penalty": "invalid",
			},
		},
		{
			name: "out of range frequency penalty high",
			params: map[string]string{
				"frequency_penalty": "3.0",
			},
		},
		{
			name: "out of range frequency penalty low",
			params: map[string]string{
				"frequency_penalty": "-3.0",
			},
		},
		{
			name: "invalid presence penalty",
			params: map[string]string{
				"presence_penalty": "invalid",
			},
		},
		{
			name: "out of range presence penalty high",
			params: map[string]string{
				"presence_penalty": "3.0",
			},
		},
		{
			name: "out of range presence penalty low",
			params: map[string]string{
				"presence_penalty": "-2.5",
			},
		},
		{
			name: "invalid max retries",
			params: map[string]string{
				"max_retries": "invalid",
			},
		},
		{
			name: "negative max retries",
			params: map[string]string{
				"max_retries": "-5",
			},
		},
		{
			name: "negative temperature",
			params: map[string]string{
				"temperature": "-0.5",
			},
		},
		{
			name: "fallback models with pipe separator",
			params: map[string]string{
				"fallback_models": "claude|openai|gemini",
			},
		},
		{
			name: "unhandled parameter ignored",
			params: map[string]string{
				"unknown_param": "value",
			},
		},
		{
			name: "model-specific unknown parameter ignored",
			params: map[string]string{
				"unknown.model": "value",
			},
		},
		{
			name: "invalid model prefix",
			params: map[string]string{
				"invalid.model": "value",
			},
		},
		{
			name: "unknown openai parameter",
			params: map[string]string{
				"openai.unknown": "value",
			},
		},
		{
			name: "unknown claude parameter",
			params: map[string]string{
				"claude.unknown": "value",
			},
		},
		{
			name: "unknown gemini parameter",
			params: map[string]string{
				"gemini.unknown": "value",
			},
		},
		{
			name: "invalid gemini safety setting format",
			params: map[string]string{
				"gemini.safety_setting": "invalid_format",
			},
		},
		{
			name: "nil model configs with model-specific params",
			params: map[string]string{
				"openai.model": "gpt-4",
				"claude.model": "claude-3",
				"gemini.model": "gemini-pro",
			},
		},
		{
			name: "alternative parameter names",
			params: map[string]string{
				"maxtokens":        "1000",
				"topp":             "0.9",
				"frequencypenalty": "0.5",
				"presencepenalty":  "0.5",
				"fallbackmodels":   "claude|openai",
				"maxretries":       "3",
				"systemmessage":    "Test message",
			},
		},
		{
			name: "openai system message alternative name",
			params: map[string]string{
				"openai.systemmessage": "OpenAI message",
			},
		},
		{
			name: "claude system message alternative name",
			params: map[string]string{
				"claude.systemmessage": "Claude message",
			},
		},
		{
			name: "gemini safety setting alternative name",
			params: map[string]string{
				"gemini.safetysetting": "harmful=block",
			},
		},
		{
			name: "model-specific with invalid format",
			params: map[string]string{
				"openai.model": "gpt-4",
				"claude":       "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultModelConfig()

			// For nil config test, set configs to nil
			if tc.name == "nil model configs with model-specific params" {
				config.OpenAIConfig = nil
				config.ClaudeConfig = nil
				config.GeminiConfig = nil
			}

			err := config.UpdateFromParams(tc.params)

			// Most error cases should return error
			if strings.Contains(tc.name, "invalid") || strings.Contains(tc.name, "negative") ||
				strings.Contains(tc.name, "out of range") {
				if err == nil && !strings.Contains(tc.name, "ignored") && !strings.Contains(tc.name, "unknown") && !strings.Contains(tc.name, "model-specific with") && !strings.Contains(tc.name, "model prefix") && !strings.Contains(tc.name, "safety setting format") {
					t.Errorf("Expected error for %s, got nil", tc.name)
				}
			}

			// Check fallback models were set correctly
			if tc.name == "fallback models with pipe separator" && err == nil {
				expected := []string{"claude", "openai", "gemini"}
				if !reflect.DeepEqual(config.FallbackModels, expected) {
					t.Errorf("Expected fallback models %v, got %v", expected, config.FallbackModels)
				}
			}

			// Check alternative parameter names
			if tc.name == "alternative parameter names" && err == nil {
				if config.MaxTokens != 1000 {
					t.Errorf("Expected MaxTokens 1000, got %d", config.MaxTokens)
				}
				if config.TopP != 0.9 {
					t.Errorf("Expected TopP 0.9, got %f", config.TopP)
				}
				if config.FrequencyPenalty != 0.5 {
					t.Errorf("Expected FrequencyPenalty 0.5, got %f", config.FrequencyPenalty)
				}
				if config.PresencePenalty != 0.5 {
					t.Errorf("Expected PresencePenalty 0.5, got %f", config.PresencePenalty)
				}
				if config.MaxRetries != 3 {
					t.Errorf("Expected MaxRetries 3, got %d", config.MaxRetries)
				}
			}

			// Check nil configs are created when needed
			if tc.name == "nil model configs with model-specific params" && err == nil {
				if config.OpenAIConfig == nil || config.ClaudeConfig == nil || config.GeminiConfig == nil {
					t.Errorf("Expected model configs to be created, but some were nil")
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		configure func(*ModelConfig)
		wantErr   bool
	}{
		{
			name:      "valid config",
			configure: func(c *ModelConfig) {},
			wantErr:   false,
		},
		{
			name: "invalid temperature - too high",
			configure: func(c *ModelConfig) {
				c.Temperature = 1.5
			},
			wantErr: true,
		},
		{
			name: "invalid temperature - negative",
			configure: func(c *ModelConfig) {
				c.Temperature = -0.5
			},
			wantErr: true,
		},
		{
			name: "invalid maxTokens - zero",
			configure: func(c *ModelConfig) {
				c.MaxTokens = 0
			},
			wantErr: true,
		},
		{
			name: "invalid maxTokens - negative",
			configure: func(c *ModelConfig) {
				c.MaxTokens = -100
			},
			wantErr: true,
		},
		{
			name: "invalid topP - too high",
			configure: func(c *ModelConfig) {
				c.TopP = 1.5
			},
			wantErr: true,
		},
		{
			name: "invalid topP - negative",
			configure: func(c *ModelConfig) {
				c.TopP = -0.5
			},
			wantErr: true,
		},
		{
			name: "invalid frequencyPenalty - too high",
			configure: func(c *ModelConfig) {
				c.FrequencyPenalty = 3.0
			},
			wantErr: true,
		},
		{
			name: "invalid frequencyPenalty - too low",
			configure: func(c *ModelConfig) {
				c.FrequencyPenalty = -3.0
			},
			wantErr: true,
		},
		{
			name: "invalid presencePenalty - too high",
			configure: func(c *ModelConfig) {
				c.PresencePenalty = 3.0
			},
			wantErr: true,
		},
		{
			name: "invalid presencePenalty - too low",
			configure: func(c *ModelConfig) {
				c.PresencePenalty = -3.0
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultModelConfig()
			tc.configure(config)

			err := config.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
