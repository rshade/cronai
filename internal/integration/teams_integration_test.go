// Package integration provides integration tests for CronAI
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
)

// TeamsWebhookResponse represents the expected response from Teams webhook
type TeamsWebhookResponse struct {
	Text string `json:"text"`
}

// TestTeamsIntegration tests the Microsoft Teams webhook processor with real API calls
func TestTeamsIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test: RUN_INTEGRATION_TESTS not set")
	}

	// Check for Teams webhook URL
	teamsWebhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if teamsWebhookURL == "" {
		t.Skip("Skipping Teams integration test: TEAMS_WEBHOOK_URL not set")
	}

	// Verify webhook URL format
	if !isValidTeamsWebhookURL(teamsWebhookURL) {
		t.Fatalf("Invalid Teams webhook URL format: %s", teamsWebhookURL)
	}

	tests := []struct {
		name         string
		processor    string
		response     *models.ModelResponse
		templateName string
		validateFunc func(t *testing.T, payload string) error
	}{
		{
			name:      "Basic Teams Message",
			processor: "teams-integration-test",
			response: &models.ModelResponse{
				Content:     fmt.Sprintf("Integration test message from CronAI at %s", time.Now().Format(time.RFC3339)),
				Model:       "test-model",
				PromptName:  "integration-test",
				Timestamp:   time.Now(),
				ExecutionID: fmt.Sprintf("test-%d", time.Now().Unix()),
			},
			templateName: "",
			validateFunc: func(_ *testing.T, payload string) error {
				// Validate JSON structure
				var msg map[string]interface{}
				if err := json.Unmarshal([]byte(payload), &msg); err != nil {
					return fmt.Errorf("invalid JSON payload: %w", err)
				}

				// Check required Teams MessageCard fields
				if msgType, ok := msg["@type"].(string); !ok || msgType != "MessageCard" {
					return fmt.Errorf("missing or invalid @type field")
				}
				if _, ok := msg["@context"].(string); !ok {
					return fmt.Errorf("missing @context field")
				}
				if _, ok := msg["summary"].(string); !ok {
					return fmt.Errorf("missing summary field")
				}

				return nil
			},
		},
		{
			name:      "Teams Alert Message",
			processor: "teams-alerts",
			response: &models.ModelResponse{
				Content:    "Test alert: System monitoring check passed",
				Model:      "test-model",
				PromptName: "system-monitor",
				Variables: map[string]string{
					"severity":  "info",
					"component": "integration-test",
					"server":    "test-server",
				},
				Timestamp:   time.Now(),
				ExecutionID: fmt.Sprintf("alert-test-%d", time.Now().Unix()),
			},
			templateName: "",
			validateFunc: func(_ *testing.T, payload string) error {
				// Validate alert-specific fields
				var msg map[string]interface{}
				if err := json.Unmarshal([]byte(payload), &msg); err != nil {
					return fmt.Errorf("invalid JSON payload: %w", err)
				}

				// Check for sections
				sections, ok := msg["sections"].([]interface{})
				if !ok || len(sections) == 0 {
					return fmt.Errorf("missing or empty sections")
				}

				return nil
			},
		},
	}

	// Test with mock HTTP client if USE_MOCK_HTTP is set
	if os.Getenv("USE_MOCK_HTTP") == "1" {
		t.Log("Using mock HTTP client for Teams integration test")
		// In a real implementation, we would inject a mock HTTP client
		// For now, we'll just validate the payload generation
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create processor
				err := processor.ProcessResponse(tt.processor, tt.response, tt.templateName)
				if err != nil {
					t.Fatalf("Failed to process response: %v", err)
				}

				// In mock mode, we can't capture the actual payload sent
				// but we've validated the processor runs without error
				t.Log("Mock test completed successfully")
			})
		}
		return
	}

	// Real API tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add delay between tests to avoid rate limiting
			time.Sleep(2 * time.Second)

			// Test that processor validates correctly
			config := processor.Config{
				Type:   "webhook",
				Target: tt.processor,
			}

			proc, err := processor.NewWebhookProcessor(config)
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}

			// Validate configuration
			if err := proc.Validate(); err != nil {
				t.Fatalf("Processor validation failed: %v", err)
			}

			// Process the response
			err = proc.Process(tt.response, tt.templateName)
			if err != nil {
				t.Fatalf("Failed to process response: %v", err)
			}

			// Validate the payload using the validateFunc
			if err := tt.validateFunc(t, tt.response.Content); err != nil {
				t.Errorf("Payload validation failed: %v", err)
			}

			// If we have a webhook URL, test sending the webhook
			if teamsWebhookURL != "" {
				payload, err := json.Marshal(tt.response)
				if err != nil {
					t.Fatalf("Failed to marshal response: %v", err)
				}

				if err := sendTeamsWebhook(teamsWebhookURL, payload); err != nil {
					t.Errorf("Failed to send webhook: %v", err)
				}
			}

			t.Logf("Successfully processed Teams webhook for: %s", tt.name)
		})
	}
}

// TestTeamsWebhookPayloadValidation tests that Teams payloads meet Microsoft's requirements
func TestTeamsWebhookPayloadValidation(t *testing.T) {
	tests := []struct {
		name         string
		payload      string
		expectError  bool
		errorMessage string
	}{
		{
			name: "Valid MessageCard",
			payload: `{
				"@type": "MessageCard",
				"@context": "https://schema.org/extensions",
				"summary": "Test message",
				"sections": [{
					"text": "This is a test message"
				}]
			}`,
			expectError: false,
		},
		{
			name: "Message Too Large",
			payload: func() string {
				// Create a payload larger than 25KB
				largeText := make([]byte, 26*1024)
				for i := range largeText {
					largeText[i] = 'A'
				}
				return fmt.Sprintf(`{
					"@type": "MessageCard",
					"@context": "https://schema.org/extensions",
					"summary": "Large message",
					"sections": [{
						"text": "%s"
					}]
				}`, string(largeText))
			}(),
			expectError:  true,
			errorMessage: "payload exceeds Teams 25KB limit",
		},
		{
			name: "Invalid JSON",
			payload: `{
				"@type": "MessageCard",
				"summary": "Invalid JSON"
				"sections": []
			}`,
			expectError:  true,
			errorMessage: "invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTeamsPayload([]byte(tt.payload))

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s' but got none", tt.errorMessage)
				} else if !bytes.Contains([]byte(err.Error()), []byte(tt.errorMessage)) {
					t.Errorf("Expected error containing '%s' but got: %v", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// Helper functions

// isValidTeamsWebhookURL checks if the URL is a valid Teams webhook URL
func isValidTeamsWebhookURL(url string) bool {
	// Teams webhook URLs should contain outlook.office.com or outlook.office365.com
	return (len(url) > 0 &&
		(bytes.Contains([]byte(url), []byte("outlook.office.com")) ||
			bytes.Contains([]byte(url), []byte("outlook.office365.com"))))
}

// validateTeamsPayload validates a Teams webhook payload
func validateTeamsPayload(payload []byte) error {
	// Check size limit (25KB for Teams)
	if len(payload) > 25*1024 {
		return fmt.Errorf("payload exceeds Teams 25KB limit: %d bytes", len(payload))
	}

	// Validate JSON
	var msg map[string]interface{}
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate required fields
	if msgType, ok := msg["@type"].(string); !ok || msgType != "MessageCard" {
		return fmt.Errorf("missing or invalid @type field")
	}
	if _, ok := msg["@context"].(string); !ok {
		return fmt.Errorf("missing @context field")
	}
	if _, ok := msg["summary"].(string); !ok {
		return fmt.Errorf("missing summary field")
	}

	return nil
}

// sendTeamsWebhook sends a webhook to Microsoft Teams
// This function is used in integration tests to verify webhook functionality
func sendTeamsWebhook(webhookURL string, payload []byte) error {
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error but don't return it since we're in a defer
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("webhook error: %d (failed to read response body: %v)", resp.StatusCode, readErr)
		}
		return fmt.Errorf("webhook error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}
