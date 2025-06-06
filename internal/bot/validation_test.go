// Package bot provides the main service for running CronAI in bot mode.
package bot

import (
	"testing"
)

func TestValidateModel(t *testing.T) {
	tests := []struct {
		name    string
		model   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid openai model",
			model:   "openai",
			wantErr: false,
		},
		{
			name:    "valid claude model",
			model:   "claude",
			wantErr: false,
		},
		{
			name:    "valid gemini model",
			model:   "gemini",
			wantErr: false,
		},
		{
			name:    "empty model",
			model:   "",
			wantErr: true,
			errMsg:  "model cannot be empty",
		},
		{
			name:    "invalid model",
			model:   "invalid-model",
			wantErr: true,
			errMsg:  "unsupported model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModel(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !validationContains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateModel() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateProcessor(t *testing.T) {
	tests := []struct {
		name      string
		processor string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty processor (optional)",
			processor: "",
			wantErr:   false,
		},
		{
			name:      "valid console processor",
			processor: "console",
			wantErr:   false,
		},
		{
			name:      "valid file processor with target",
			processor: "file-output.txt",
			wantErr:   false,
		},
		{
			name:      "valid slack processor with channel",
			processor: "slack-general",
			wantErr:   false,
		},
		{
			name:      "invalid processor",
			processor: "invalid-processor",
			wantErr:   true,
			errMsg:    "unsupported processor type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProcessor(tt.processor)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProcessor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !validationContains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateProcessor() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid port",
			port:    "8080",
			wantErr: false,
		},
		{
			name:    "port zero (test mode)",
			port:    "0",
			wantErr: false,
		},
		{
			name:    "empty port",
			port:    "",
			wantErr: true,
			errMsg:  "port cannot be empty",
		},
		{
			name:    "invalid port format",
			port:    "abc",
			wantErr: true,
			errMsg:  "port must be a valid integer",
		},
		{
			name:    "port too low",
			port:    "-1",
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name:    "port too high",
			port:    "70000",
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !validationContains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePort() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateWebhookSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty secret (optional)",
			secret:  "",
			wantErr: false,
		},
		{
			name:    "valid secret",
			secret:  "mysecret123",
			wantErr: false,
		},
		{
			name:    "short secret",
			secret:  "short",
			wantErr: true,
			errMsg:  "webhook secret should be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWebhookSecret(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWebhookSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !validationContains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateWebhookSecret() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Value:   "test_value",
		Message: "test message",
	}

	expected := "validation error for test_field='test_value': test message"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}

// Helper function to check if string contains substring
func validationContains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 ||
			(len(s) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
