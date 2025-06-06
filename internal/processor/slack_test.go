// Package processor implements the response processors for CronAI.
package processor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

func TestNewSlackProcessor(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			config: Config{
				Type:   "slack",
				Target: "",
			},
			wantErr: false, // Constructor doesn't validate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSlackProcessor(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSlackProcessor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlackProcessor_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		envVars     map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid with token",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			envVars: map[string]string{
				"SLACK_TOKEN": "xoxb-test-token",
			},
			wantErr: false,
		},
		{
			name: "valid with webhook",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			envVars: map[string]string{
				"SLACK_WEBHOOK_URL": "https://hooks.slack.com/services/TEST/TEST/TEST",
			},
			wantErr: false,
		},
		{
			name: "valid with both token and webhook",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			envVars: map[string]string{
				"SLACK_TOKEN":       "xoxb-test-token",
				"SLACK_WEBHOOK_URL": "https://hooks.slack.com/services/TEST/TEST/TEST",
			},
			wantErr: false,
		},
		{
			name: "missing target",
			config: Config{
				Type:   "slack",
				Target: "",
			},
			envVars: map[string]string{
				"SLACK_TOKEN": "xoxb-test-token",
			},
			wantErr:     true,
			errContains: "slack channel cannot be empty",
		},
		{
			name: "missing credentials",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			envVars:     map[string]string{},
			wantErr:     true,
			errContains: "either SLACK_TOKEN or SLACK_WEBHOOK_URL environment variable must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("Failed to set env var %s: %v", k, err)
				}
			}
			defer func() {
				for k := range tt.envVars {
					if err := os.Unsetenv(k); err != nil {
						t.Errorf("Failed to unset env var %s: %v", k, err)
					}
				}
			}()

			processor, err := NewSlackProcessor(tt.config)
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}
			err = processor.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

func TestSlackProcessor_Process(t *testing.T) {
	// Initialize template manager
	template.GetManager()

	tests := []struct {
		name         string
		config       Config
		response     *models.ModelResponse
		templateName string
		envVars      map[string]string
		wantErr      bool
	}{
		{
			name: "process with default template",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			response: &models.ModelResponse{
				Content:     "Test message content",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "test-prompt",
				ExecutionID: "test-123",
			},
			envVars: map[string]string{
				"SLACK_WEBHOOK_URL": "mock-webhook-url",
			},
			wantErr: false,
		},
		{
			name: "process with monitoring template",
			config: Config{
				Type:   "slack",
				Target: "#alerts",
			},
			response: &models.ModelResponse{
				Content:     "System health check failed",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "system-health-monitoring",
				ExecutionID: "test-456",
			},
			envVars: map[string]string{
				"SLACK_WEBHOOK_URL": "mock-webhook-url",
			},
			wantErr: false,
		},
		{
			name: "process with default template and token",
			config: Config{
				Type:   "slack",
				Target: "#general",
			},
			response: &models.ModelResponse{
				Content:     "Test message for token",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "test-prompt",
				ExecutionID: "test-789",
			},
			templateName: "",
			envVars: map[string]string{
				"SLACK_TOKEN": "xoxb-test-token",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("Failed to set env var %s: %v", k, err)
				}
			}
			defer func() {
				for k := range tt.envVars {
					if err := os.Unsetenv(k); err != nil {
						t.Errorf("Failed to unset env var %s: %v", k, err)
					}
				}
			}()

			// Create mock server for webhook calls
			if _, hasWebhook := tt.envVars["SLACK_WEBHOOK_URL"]; hasWebhook {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
					if _, err := w.Write([]byte("ok")); err != nil {
						t.Errorf("Failed to write response: %v", err)
					}
				}))
				defer mockServer.Close()
				if err := os.Setenv("SLACK_WEBHOOK_URL", mockServer.URL); err != nil {
					t.Fatalf("Failed to set webhook URL: %v", err)
				}
			}

			processor, err := NewSlackProcessor(tt.config)
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}

			// Skip OAuth tests since we can't easily mock the Slack API
			if _, hasToken := tt.envVars["SLACK_TOKEN"]; hasToken {
				t.Skip("OAuth tests require Slack API mocking - skipping")
			}

			err = processor.Process(tt.response, tt.templateName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlackProcessor_sendViaWebhook(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		serverResponse string
		serverStatus   int
		wantErr        bool
		errContains    string
	}{
		{
			name:           "successful webhook",
			payload:        []byte(`{"text":"test message"}`),
			serverResponse: "ok",
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "webhook error response",
			payload:        []byte(`{"text":"test message"}`),
			serverResponse: "invalid_payload",
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
			errContains:    "slack webhook error: 400",
		},
		{
			name:           "unexpected webhook response",
			payload:        []byte(`{"text":"test message"}`),
			serverResponse: "not ok",
			serverStatus:   http.StatusOK,
			wantErr:        true,
			errContains:    "unexpected webhook response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.serverStatus)
				if _, err := w.Write([]byte(tt.serverResponse)); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer mockServer.Close()

			processor := &SlackProcessor{
				config: Config{
					Type:   "slack",
					Target: "#general",
				},
			}

			err := processor.sendViaWebhook(mockServer.URL, tt.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("sendViaWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("sendViaWebhook() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}

// TestSlackProcessor_sendViaOAuth is disabled as it requires modification to accept custom API endpoint for testing
// TODO: Implement proper OAuth testing with dependency injection
/*
func TestSlackProcessor_sendViaOAuth(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      []byte
		serverResp   map[string]interface{}
		serverStatus int
		wantErr      bool
		errContains  string
	}{
		{
			name:    "successful API call",
			token:   "xoxb-test-token",
			payload: []byte(`{"channel":"#general","text":"test message"}`),
			serverResp: map[string]interface{}{
				"ok": true,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:    "API error response",
			token:   "xoxb-test-token",
			payload: []byte(`{"channel":"#general","text":"test message"}`),
			serverResp: map[string]interface{}{
				"ok":    false,
				"error": "channel_not_found",
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
			errContains:  "channel_not_found",
		},
		{
			name:         "HTTP error",
			token:        "xoxb-test-token",
			payload:      []byte(`{"channel":"#general","text":"test message"}`),
			serverResp:   map[string]interface{}{},
			serverStatus: http.StatusUnauthorized,
			wantErr:      true,
			errContains:  "Slack API error: 401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify authorization header
				authHeader := r.Header.Get("Authorization")
				expectedAuth := "Bearer " + tt.token
				if authHeader != expectedAuth {
					t.Errorf("Expected Authorization header %v, got %v", expectedAuth, authHeader)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if err := json.NewEncoder(w).Encode(tt.serverResp); err != nil {
					t.Errorf("Failed to encode response: %v", err)
				}
			}))
			defer mockServer.Close()

			// We need to modify the sendViaOAuth method to accept a custom URL for testing
			// For now, we'll skip this test as it requires code changes
			t.Skip("Requires modification to accept custom API endpoint for testing")
		})
	}
}
*/

func TestSlackProcessor_GetType(t *testing.T) {
	processor := &SlackProcessor{
		config: Config{
			Type: "slack",
		},
	}

	if got := processor.GetType(); got != "slack" {
		t.Errorf("GetType() = %v, want %v", got, "slack")
	}
}

func TestSlackProcessor_GetConfig(t *testing.T) {
	config := Config{
		Type:   "slack",
		Target: "#general",
	}

	processor := &SlackProcessor{
		config: config,
	}

	got := processor.GetConfig()
	if got.Type != config.Type || got.Target != config.Target {
		t.Errorf("GetConfig() = %v, want %v", got, config)
	}
}

func TestSlackProcessor_JSONValidation(t *testing.T) {
	// Initialize template manager
	template.GetManager()

	tests := []struct {
		name         string
		templateName string
		response     *models.ModelResponse
		wantValid    bool
	}{
		{
			name:         "valid default template JSON",
			templateName: "",
			response: &models.ModelResponse{
				Content:     "Test message",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "test-prompt",
				ExecutionID: "test-123",
			},
			wantValid: true,
		},
		{
			name:         "valid monitoring template JSON",
			templateName: "",
			response: &models.ModelResponse{
				Content:     "Alert message",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "monitoring-check",
				ExecutionID: "test-456",
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that validates JSON
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var payload map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Errorf("Invalid JSON received: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Check for required fields
				if _, ok := payload["blocks"]; !ok {
					t.Error("Missing 'blocks' field in payload")
				}
				if _, ok := payload["channel"]; !ok {
					t.Error("Missing 'channel' field in payload")
				}

				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("ok")); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer mockServer.Close()

			// Use webhook URL so we can mock the call
			if err := os.Setenv("SLACK_WEBHOOK_URL", mockServer.URL); err != nil {
				t.Fatalf("Failed to set webhook URL: %v", err)
			}
			defer func() {
				if err := os.Unsetenv("SLACK_WEBHOOK_URL"); err != nil {
					t.Errorf("Failed to unset webhook URL: %v", err)
				}
			}()

			slackProcessor, err := NewSlackProcessor(Config{
				Type:   "slack",
				Target: "#general",
			})
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}

			err = slackProcessor.Process(tt.response, tt.templateName)
			if err != nil && tt.wantValid {
				t.Errorf("Process() returned error for valid JSON: %v", err)
			}
		})
	}
}
