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
			expected: "claude-3-sonnet-20240229",
		},
		{
			name: "custom model",
			config: &config.ModelConfig{
				ClaudeConfig: &config.ClaudeConfig{
					Model: "claude-3-opus-20240229",
				},
			},
			expected: "claude-3-opus-20240229",
		},
		{
			name:     "nil config",
			config:   nil,
			expected: "claude-3-sonnet-20240229",
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
