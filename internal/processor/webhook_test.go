// Package processor provides functionality for processing model responses through various output channels
package processor

import (
	"os"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

func TestTeamsWebhookProcessor(t *testing.T) {
	// Ensure default processors are registered
	registry := GetRegistry()
	registry.RegisterDefaults()

	tests := []struct {
		name          string
		processor     string
		response      *models.ModelResponse
		templateName  string
		expectedError bool
		setup         func() error
		cleanup       func() error
	}{
		{
			name:      "Teams Webhook Without URL - Should Fail",
			processor: "teams-general",
			response: &models.ModelResponse{
				Content:     "Test content for Teams",
				Model:       "claude-3-haiku",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution-001",
			},
			templateName:  "",
			expectedError: true,
			setup: func() error {
				// Ensure Teams webhook URL env vars are not set
				if err := os.Unsetenv("TEAMS_WEBHOOK_URL"); err != nil {
					return err
				}
				if err := os.Unsetenv("WEBHOOK_URL_TEAMS"); err != nil {
					return err
				}
				if err := os.Unsetenv("WEBHOOK_URL"); err != nil {
					return err
				}
				return nil
			},
			cleanup: func() error { return nil },
		},
		{
			name:      "Teams Webhook With TEAMS_WEBHOOK_URL - Should Pass",
			processor: "teams-general",
			response: &models.ModelResponse{
				Content:     "Test content for Teams",
				Model:       "claude-3-haiku",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution-002",
			},
			templateName:  "",
			expectedError: false,
			setup: func() error {
				// Set Teams webhook URL
				return os.Setenv("TEAMS_WEBHOOK_URL", "https://outlook.office.com/webhook/test")
			},
			cleanup: func() error {
				return os.Unsetenv("TEAMS_WEBHOOK_URL")
			},
		},
		{
			name:      "Teams Webhook With WEBHOOK_URL_TEAMS - Should Pass",
			processor: "teams-monitoring",
			response: &models.ModelResponse{
				Content:     "Test monitoring alert for Teams",
				Model:       "claude-3-haiku",
				PromptName:  "monitoring-check",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution-003",
			},
			templateName:  "",
			expectedError: false,
			setup: func() error {
				// Clear TEAMS_WEBHOOK_URL to test fallback
				if err := os.Unsetenv("TEAMS_WEBHOOK_URL"); err != nil {
					return err
				}
				// Set type-specific Teams webhook URL
				return os.Setenv("WEBHOOK_URL_TEAMS", "https://outlook.office.com/webhook/teams")
			},
			cleanup: func() error {
				return os.Unsetenv("WEBHOOK_URL_TEAMS")
			},
		},
		{
			name:      "Teams Webhook with Custom Template - Should Pass",
			processor: "teams-alerts",
			response: &models.ModelResponse{
				Content:     "System alert: High CPU usage detected on server-01",
				Model:       "claude-3-haiku",
				PromptName:  "system-monitor",
				Variables:   map[string]string{"severity": "high", "component": "cpu", "server": "server-01"},
				Timestamp:   time.Now(),
				ExecutionID: "alert-123",
			},
			templateName:  "custom_teams_alert",
			expectedError: false,
			setup: func() error {
				// Set Teams webhook URL
				if err := os.Setenv("TEAMS_WEBHOOK_URL", "https://outlook.office.com/webhook/alerts"); err != nil {
					return err
				}

				// Register custom Teams template
				manager := template.GetManager()
				customTemplate := `{
	"@type": "MessageCard",
	"@context": "https://schema.org/extensions",
	"themeColor": "FF6D00",
	"summary": "🚨 {{.Variables.severity}} Alert: {{.PromptName}}",
	"sections": [{
		"activityTitle": "🚨 {{.Variables.severity}} Alert: {{.PromptName}}",
		"activitySubtitle": "Server: {{.Variables.server}}",
		"facts": [
			{
				"name": "Severity",
				"value": "{{.Variables.severity}}"
			},
			{
				"name": "Component", 
				"value": "{{.Variables.component}}"
			},
			{
				"name": "Server",
				"value": "{{.Variables.server}}"
			},
			{
				"name": "Time",
				"value": "{{.Timestamp.Format "Jan 02, 2006 15:04:05"}}"
			}
		],
		"markdown": true,
		"text": "{{.Content}}"
	}]
}`
				return manager.RegisterTemplate("custom_teams_alert", customTemplate)
			},
			cleanup: func() error {
				return os.Unsetenv("TEAMS_WEBHOOK_URL")
			},
		},
		{
			name:      "Teams Webhook with Default Teams Template - Should Pass",
			processor: "teams-reports",
			response: &models.ModelResponse{
				Content:     "Weekly report generated successfully",
				Model:       "claude-3-sonnet",
				PromptName:  "weekly-report",
				Variables:   map[string]string{"report_type": "weekly", "period": "2024-W01"},
				Timestamp:   time.Now(),
				ExecutionID: "report-456",
			},
			templateName:  "",
			expectedError: false,
			setup: func() error {
				// Set Teams webhook URL
				return os.Setenv("TEAMS_WEBHOOK_URL", "https://outlook.office.com/webhook/reports")
			},
			cleanup: func() error {
				return os.Unsetenv("TEAMS_WEBHOOK_URL")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if err := tt.setup(); err != nil {
				t.Fatalf("Test setup failed: %v", err)
			}

			// Cleanup after test
			defer func() {
				if err := tt.cleanup(); err != nil {
					t.Logf("Test cleanup failed: %v", err)
				}
			}()

			// Execute the processor
			err := ProcessResponse(tt.processor, tt.response, tt.templateName)

			// Validate results
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTeamsWebhookEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name           string
		webhookType    string
		envVars        map[string]string
		expectedURL    string
		expectedHasURL bool
	}{
		{
			name:        "TEAMS_WEBHOOK_URL takes precedence",
			webhookType: "teams",
			envVars: map[string]string{
				"TEAMS_WEBHOOK_URL": "https://teams.url/primary",
				"WEBHOOK_URL_TEAMS": "https://teams.url/secondary",
				"WEBHOOK_URL":       "https://generic.url",
			},
			expectedURL:    "https://teams.url/primary",
			expectedHasURL: true,
		},
		{
			name:        "WEBHOOK_URL_TEAMS as fallback",
			webhookType: "teams",
			envVars: map[string]string{
				"WEBHOOK_URL_TEAMS": "https://teams.url/secondary",
				"WEBHOOK_URL":       "https://generic.url",
			},
			expectedURL:    "https://teams.url/secondary",
			expectedHasURL: true,
		},
		{
			name:        "WEBHOOK_URL as final fallback",
			webhookType: "teams",
			envVars: map[string]string{
				"WEBHOOK_URL": "https://generic.url",
			},
			expectedURL:    "https://generic.url",
			expectedHasURL: true,
		},
		{
			name:           "No URL configured",
			webhookType:    "teams",
			envVars:        map[string]string{},
			expectedURL:    "",
			expectedHasURL: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all webhook environment variables
			if err := os.Unsetenv("TEAMS_WEBHOOK_URL"); err != nil {
				t.Logf("Failed to unset TEAMS_WEBHOOK_URL: %v", err)
			}
			if err := os.Unsetenv("WEBHOOK_URL_TEAMS"); err != nil {
				t.Logf("Failed to unset WEBHOOK_URL_TEAMS: %v", err)
			}
			if err := os.Unsetenv("WEBHOOK_URL"); err != nil {
				t.Logf("Failed to unset WEBHOOK_URL: %v", err)
			}

			// Set up test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set environment variable %s: %v", key, err)
				}
			}

			// Cleanup after test
			defer func() {
				for key := range tt.envVars {
					if err := os.Unsetenv(key); err != nil {
						t.Logf("Failed to unset %s: %v", key, err)
					}
				}
			}()

			// Test GetWebhookURL function
			actualURL := GetWebhookURL(tt.webhookType)

			if tt.expectedHasURL {
				if actualURL == "" {
					t.Errorf("Expected URL %q but got empty string", tt.expectedURL)
				} else if actualURL != tt.expectedURL {
					t.Errorf("Expected URL %q but got %q", tt.expectedURL, actualURL)
				}
			} else {
				if actualURL != "" {
					t.Errorf("Expected empty URL but got %q", actualURL)
				}
			}
		})
	}
}
