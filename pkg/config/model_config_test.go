package config

import (
	"os"
	"reflect"
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
