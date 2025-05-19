package models

import (
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewGeminiClient(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		config   *config.ModelConfig
		wantErr  bool
		errMsg   string
	}{
		{
			name: "missing API key",
			setupEnv: func() {
				os.Unsetenv("GOOGLE_API_KEY")
			},
			config:  &config.ModelConfig{},
			wantErr: true,
			errMsg:  "GOOGLE_API_KEY environment variable not set",
		},
		{
			name: "valid API key",
			setupEnv: func() {
				os.Setenv("GOOGLE_API_KEY", "test-key")
			},
			config:  &config.ModelConfig{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore environment
			oldKey := os.Getenv("GOOGLE_API_KEY")
			defer os.Setenv("GOOGLE_API_KEY", oldKey)

			tt.setupEnv()

			client, err := NewGeminiClient(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestGeminiClient_GetModelName(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ModelConfig
		expected string
	}{
		{
			name:     "default model",
			config:   &config.ModelConfig{},
			expected: "gemini-pro",
		},
		{
			name: "custom model",
			config: &config.ModelConfig{
				GeminiConfig: &config.GeminiConfig{
					Model: "gemini-pro-vision",
				},
			},
			expected: "gemini-pro-vision",
		},
		{
			name:     "nil config",
			config:   nil,
			expected: "gemini-pro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &GeminiClient{
				config: tt.config,
			}
			result := client.getModelName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarmCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected genai.HarmCategory
		wantErr  bool
	}{
		{"harassment", "harassment", genai.HarmCategoryHarassment, false},
		{"hate_speech", "hate_speech", genai.HarmCategoryHateSpeech, false},
		{"hate", "hate", genai.HarmCategoryHateSpeech, false},
		{"sexually_explicit", "sexually_explicit", genai.HarmCategorySexuallyExplicit, false},
		{"sexual", "sexual", genai.HarmCategorySexuallyExplicit, false},
		{"dangerous_content", "dangerous_content", genai.HarmCategoryDangerousContent, false},
		{"dangerous", "dangerous", genai.HarmCategoryDangerousContent, false},
		{"derogatory", "derogatory", genai.HarmCategoryDerogatory, false},
		{"toxicity", "toxicity", genai.HarmCategoryToxicity, false},
		{"toxic", "toxic", genai.HarmCategoryToxicity, false},
		{"violence", "violence", genai.HarmCategoryViolence, false},
		{"medical", "medical", genai.HarmCategoryMedical, false},
		{"unknown", "unknown", genai.HarmCategoryUnspecified, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHarmCategory(tt.category)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHarmLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected genai.HarmBlockThreshold
		wantErr  bool
	}{
		{"block_none", "block_none", genai.HarmBlockNone, false},
		{"none", "none", genai.HarmBlockNone, false},
		{"block_low", "block_low", genai.HarmBlockLowAndAbove, false},
		{"low", "low", genai.HarmBlockLowAndAbove, false},
		{"block_medium", "block_medium", genai.HarmBlockMediumAndAbove, false},
		{"medium", "medium", genai.HarmBlockMediumAndAbove, false},
		{"block_high", "block_high", genai.HarmBlockOnlyHigh, false},
		{"high", "high", genai.HarmBlockOnlyHigh, false},
		{"block", "block", genai.HarmBlockUnspecified, false},
		{"block_all", "block_all", genai.HarmBlockUnspecified, false},
		{"unknown", "unknown", genai.HarmBlockUnspecified, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHarmLevel(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Note: Testing the Execute method would require mocking the genai client
// The Google AI SDK uses complex internal structures that are difficult to mock.
// We focus on testing the configuration and helper methods.
