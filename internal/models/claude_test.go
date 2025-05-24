package models

import (
	"os"
	"testing"

	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewClaudeClient(t *testing.T) {
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
				if err := os.Unsetenv("ANTHROPIC_API_KEY"); err != nil {
					t.Fatal(err)
				}
			},
			config:  &config.ModelConfig{},
			wantErr: true,
			errMsg:  "ANTHROPIC_API_KEY environment variable not set",
		},
		{
			name: "valid API key",
			setupEnv: func() {
				if err := os.Setenv("ANTHROPIC_API_KEY", "test-key"); err != nil {
					t.Fatal(err)
				}
			},
			config:  &config.ModelConfig{},
			wantErr: false,
		},
		{
			name: "with base URL",
			setupEnv: func() {
				if err := os.Setenv("ANTHROPIC_API_KEY", "test-key"); err != nil {
					t.Fatal(err)
				}
				if err := os.Setenv("ANTHROPIC_BASE_URL", "https://custom.anthropic.com"); err != nil {
					t.Fatal(err)
				}
			},
			config:  &config.ModelConfig{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore environment
			oldKey := os.Getenv("ANTHROPIC_API_KEY")
			oldURL := os.Getenv("ANTHROPIC_BASE_URL")
			defer func() {
				if err := os.Setenv("ANTHROPIC_API_KEY", oldKey); err != nil {
					t.Fatal(err)
				}
				if err := os.Setenv("ANTHROPIC_BASE_URL", oldURL); err != nil {
					t.Fatal(err)
				}
			}()

			tt.setupEnv()

			client, err := NewClaudeClient(tt.config)

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

func TestClaudeClient_GetModelName(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ModelConfig
		expected string
	}{
		{
			name:     "default model",
			config:   &config.ModelConfig{},
			expected: DefaultClaudeModel,
		},
		{
			name: "claude 3 opus model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "claude-3-opus-20240229",
				},
			},
			expected: "claude-3-opus-20240229",
		},
		{
			name: "claude 3.5 sonnet model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "claude-3-5-sonnet-20241022",
				},
			},
			expected: "claude-3-5-sonnet-20241022",
		},
		{
			name: "claude 4 opus model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: Claude4OpusLatest,
				},
			},
			expected: Claude4OpusLatest,
		},
		{
			name: "claude 4 sonnet model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: Claude4SonnetLatest,
				},
			},
			expected: Claude4SonnetLatest,
		},
		{
			name: "claude 4 haiku model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: Claude4HaikuLatest,
				},
			},
			expected: Claude4HaikuLatest,
		},
		{
			name: "opus alias",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "opus",
				},
			},
			expected: Claude4OpusLatest,
		},
		{
			name: "sonnet alias",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "sonnet",
				},
			},
			expected: Claude4SonnetLatest,
		},
		{
			name: "haiku alias",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "haiku",
				},
			},
			expected: Claude4HaikuLatest,
		},
		{
			name: "3.5-opus alias",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "3.5-opus",
				},
			},
			expected: Claude35OpusLatest,
		},
		{
			name: "3-sonnet alias",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "3-sonnet",
				},
			},
			expected: Claude3SonnetLatest,
		},
		{
			name: "unsupported model falls back to default",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "unsupported-model",
				},
			},
			expected: DefaultClaudeModel,
		},
		{
			name:     "nil config",
			config:   nil,
			expected: DefaultClaudeModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &ClaudeClient{
				config: tt.config,
			}
			result := client.getModelName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClaudeClient_GetSystemMessage(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ModelConfig
		expected string
	}{
		{
			name:     "default message",
			config:   &config.ModelConfig{},
			expected: "You are a helpful assistant.",
		},
		{
			name: "custom message",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					SystemMessage: "You are a programming assistant.",
				},
			},
			expected: "You are a programming assistant.",
		},
		{
			name:     "nil config",
			config:   nil,
			expected: "You are a helpful assistant.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &ClaudeClient{
				config: tt.config,
			}
			result := client.getSystemMessage()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Note: Testing the Execute method would require mocking the anthropic client
// Since the anthropic SDK doesn't provide an interface that's easily mockable,
// we focus on testing the configuration and helper methods.

func TestGetAvailableClaudeModels(t *testing.T) {
	models := GetAvailableClaudeModels()

	// Test that we have all expected model categories
	expectedCategories := []string{
		"claude-4", "claude-3-5", "claude-3",
	}

	for _, category := range expectedCategories {
		found := false
		for model := range models {
			if contains(model, category) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected to find models for category: %s", category)
	}

	// Test specific models exist
	expectedModels := []string{
		Claude4OpusLatest,
		Claude4SonnetLatest,
		Claude4HaikuLatest,
		Claude35OpusLatest,
		Claude35SonnetLatest,
		Claude35HaikuLatest,
		Claude3OpusLatest,
		Claude3SonnetLatest,
		Claude3HaikuLatest,
	}

	for _, model := range expectedModels {
		_, exists := models[model]
		assert.True(t, exists, "Expected model %s to exist in available models", model)
	}
}

func TestGetClaudeModelAliases(t *testing.T) {
	aliases := GetClaudeModelAliases()

	// Test key aliases exist
	expectedAliases := []string{
		"opus", "sonnet", "haiku",
		"4-opus", "4-sonnet", "4-haiku",
		"3.5-opus", "3.5-sonnet", "3.5-haiku",
		"3-opus", "3-sonnet", "3-haiku",
	}

	for _, alias := range expectedAliases {
		_, exists := aliases[alias]
		assert.True(t, exists, "Expected alias %s to exist", alias)
	}

	// Test that default aliases point to Claude 4 models
	assert.Equal(t, Claude4OpusLatest, aliases["opus"])
	assert.Equal(t, Claude4SonnetLatest, aliases["sonnet"])
	assert.Equal(t, Claude4HaikuLatest, aliases["haiku"])
}

func TestSupportedClaudeModels(t *testing.T) {
	// Test all Claude 4 models are supported
	claude4Models := []string{
		Claude4OpusLatest,
		Claude4Opus20250514,
		Claude4SonnetLatest,
		Claude4HaikuLatest,
	}

	for _, model := range claude4Models {
		assert.True(t, SupportedClaudeModels[model], "Claude 4 model %s should be supported", model)
	}

	// Test all Claude 3.5 models are supported
	claude35Models := []string{
		Claude35OpusLatest,
		Claude35Opus20250120,
		Claude35SonnetLatest,
		Claude35Sonnet20241022,
		Claude35Sonnet20240620,
		Claude35HaikuLatest,
		Claude35Haiku20241022,
	}

	for _, model := range claude35Models {
		assert.True(t, SupportedClaudeModels[model], "Claude 3.5 model %s should be supported", model)
	}

	// Test all Claude 3 models are supported
	claude3Models := []string{
		Claude3OpusLatest,
		Claude3Opus20240229,
		Claude3SonnetLatest,
		Claude3Sonnet20240229,
		Claude3HaikuLatest,
		Claude3Haiku20240307,
	}

	for _, model := range claude3Models {
		assert.True(t, SupportedClaudeModels[model], "Claude 3 model %s should be supported", model)
	}
}

// Helper function for string contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || (len(substr) < len(s) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
