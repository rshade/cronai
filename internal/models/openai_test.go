package models

import (
	"os"
	"testing"

	"github.com/rshade/cronai/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAIClient(t *testing.T) {
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
				if err := os.Unsetenv("OPENAI_API_KEY"); err != nil {
					t.Fatal(err)
				}
			},
			config:  &config.ModelConfig{},
			wantErr: true,
			errMsg:  "OPENAI_API_KEY environment variable not set",
		},
		{
			name: "valid API key",
			setupEnv: func() {
				if err := os.Setenv("OPENAI_API_KEY", "test-key"); err != nil {
					t.Fatal(err)
				}
			},
			config:  &config.ModelConfig{},
			wantErr: false,
		},
		{
			name: "with base URL",
			setupEnv: func() {
				if err := os.Setenv("OPENAI_API_KEY", "test-key"); err != nil {
					t.Fatal(err)
				}
				if err := os.Setenv("OPENAI_BASE_URL", "https://custom.openai.com"); err != nil {
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
			oldKey := os.Getenv("OPENAI_API_KEY")
			oldURL := os.Getenv("OPENAI_BASE_URL")
			defer func() {
				if err := os.Setenv("OPENAI_API_KEY", oldKey); err != nil {
					t.Fatal(err)
				}
				if err := os.Setenv("OPENAI_BASE_URL", oldURL); err != nil {
					t.Fatal(err)
				}
			}()

			tt.setupEnv()

			client, err := NewOpenAIClient(tt.config)

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

func TestOpenAIClient_HelperMethods(t *testing.T) {
	tests := []struct {
		name              string
		config            *config.ModelConfig
		expectedModel     string
		expectedSystemMsg string
	}{
		{
			name:              "default values",
			config:            &config.ModelConfig{},
			expectedModel:     "gpt-3.5-turbo",
			expectedSystemMsg: "You are a helpful assistant.",
		},
		{
			name: "custom values",
			config: &config.ModelConfig{
				OpenAIConfig: &config.OpenAIConfig{
					Model:         "gpt-4",
					SystemMessage: "You are a code assistant.",
				},
			},
			expectedModel:     "gpt-4",
			expectedSystemMsg: "You are a code assistant.",
		},
		{
			name:              "nil config",
			config:            nil,
			expectedModel:     "gpt-3.5-turbo",
			expectedSystemMsg: "You are a helpful assistant.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OpenAIClient{
				config: tt.config,
			}

			// Test getModelName
			modelName := client.getModelName()
			assert.Equal(t, tt.expectedModel, modelName)

			// Test getSystemMessage
			systemMsg := client.getSystemMessage()
			assert.Equal(t, tt.expectedSystemMsg, systemMsg)
		})
	}
}

func TestOpenAIClient_GetModelName(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ModelConfig
		expected string
	}{
		{
			name:     "default model",
			config:   &config.ModelConfig{},
			expected: "gpt-3.5-turbo",
		},
		{
			name: "custom model",
			config: &config.ModelConfig{
				OpenAIConfig: &config.OpenAIConfig{
					Model: "gpt-4",
				},
			},
			expected: "gpt-4",
		},
		{
			name:     "nil config",
			config:   nil,
			expected: "gpt-3.5-turbo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OpenAIClient{
				config: tt.config,
			}
			result := client.getModelName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOpenAIClient_GetSystemMessage(t *testing.T) {
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
				OpenAIConfig: &config.OpenAIConfig{
					SystemMessage: "You are a code assistant.",
				},
			},
			expected: "You are a code assistant.",
		},
		{
			name:     "nil config",
			config:   nil,
			expected: "You are a helpful assistant.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OpenAIClient{
				config: tt.config,
			}
			result := client.getSystemMessage()
			assert.Equal(t, tt.expected, result)
		})
	}
}
