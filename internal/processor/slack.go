package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// SlackProcessor handles Slack messaging
type SlackProcessor struct {
	config Config
}

// NewSlackProcessor creates a new Slack processor
func NewSlackProcessor(config Config) (Processor, error) {
	return &SlackProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (s *SlackProcessor) Process(response *models.ModelResponse, templateName string) error {
	// Create template data
	tmplData := template.Data{
		Content:     response.Content,
		Model:       response.Model,
		Timestamp:   response.Timestamp,
		PromptName:  response.PromptName,
		Variables:   response.Variables,
		ExecutionID: response.ExecutionID,
		Metadata:    make(map[string]string),
	}

	// Add standard metadata fields
	tmplData.Metadata["timestamp"] = response.Timestamp.Format(time.RFC3339)
	tmplData.Metadata["date"] = response.Timestamp.Format("2006-01-02")
	tmplData.Metadata["time"] = response.Timestamp.Format("15:04:05")
	tmplData.Metadata["execution_id"] = response.ExecutionID
	tmplData.Metadata["processor"] = s.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	return s.processSlackWithTemplate(s.config.Target, tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (s *SlackProcessor) Validate() error {
	if s.config.Target == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("slack channel cannot be empty"),
			"invalid slack processor configuration")
	}

	// Check for required environment variables - either token or webhook URL
	slackToken := os.Getenv(EnvSlackToken)
	slackWebhookURL := os.Getenv(EnvSlackWebhookURL)

	if slackToken == "" && slackWebhookURL == "" {
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("either SLACK_TOKEN or SLACK_WEBHOOK_URL environment variable must be set"),
			"missing Slack configuration")
	}

	return nil
}

// GetType returns the processor type identifier
func (s *SlackProcessor) GetType() string {
	return "slack"
}

// GetConfig returns the processor configuration
func (s *SlackProcessor) GetConfig() Config {
	return s.config
}

// processSlackWithTemplate sends formatted messages to Slack
func (s *SlackProcessor) processSlackWithTemplate(channel string, data template.Data, templateName string) error {
	// Check for Slack configuration
	slackToken := os.Getenv(EnvSlackToken)
	slackWebhookURL := os.Getenv(EnvSlackWebhookURL)

	if slackToken == "" && slackWebhookURL == "" {
		log.Error("Neither Slack token nor webhook URL set", logger.Fields{
			"channel": channel,
		})
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("either SLACK_TOKEN or SLACK_WEBHOOK_URL environment variable must be set"),
			"missing Slack configuration")
	}

	// Get template manager
	manager := template.GetManager()

	// If monitoring-related prompt, use monitoring template as default
	isMonitoring := strings.Contains(strings.ToLower(data.PromptName), "monitor") ||
		strings.Contains(strings.ToLower(data.PromptName), "alert") ||
		strings.Contains(strings.ToLower(data.PromptName), "health")

	// Use default template if none specified
	if templateName == "" {
		if isMonitoring {
			templateName = "default_slack_monitoring"
		} else {
			templateName = "default_slack"
		}
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)
	if payload == "" {
		log.Error("Failed to generate Slack payload", logger.Fields{
			"template": templateName,
			"channel":  channel,
		})
		return errors.Wrap(errors.CategoryApplication, fmt.Errorf("empty payload generated from template %s", templateName),
			"Slack message generation failed")
	}

	// Add to metadata for logging
	data.Metadata["slack_channel"] = channel
	data.Metadata["template_used"] = templateName

	// Validate JSON payload
	var jsonPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &jsonPayload); err != nil {
		log.Error("Invalid Slack JSON payload", logger.Fields{
			"template": templateName,
			"error":    err.Error(),
		})
		return errors.Wrap(errors.CategoryApplication, err, "Slack payload is not valid JSON")
	}

	// Add channel to payload if not present
	if _, ok := jsonPayload["channel"]; !ok {
		jsonPayload["channel"] = channel
	}

	// Convert back to JSON
	payloadBytes, err := json.Marshal(jsonPayload)
	if err != nil {
		return errors.Wrap(errors.CategoryApplication, err, "failed to marshal Slack payload")
	}

	// Send to Slack using appropriate method
	if slackWebhookURL != "" {
		// Use webhook method
		return s.sendViaWebhook(slackWebhookURL, payloadBytes)
	}

	// Use OAuth token method
	return s.sendViaOAuth(slackToken, payloadBytes)
}

// sendViaWebhook sends the message using a webhook URL
func (s *SlackProcessor) sendViaWebhook(webhookURL string, payload []byte) error {
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(errors.CategoryApplication, err, "failed to create webhook request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(errors.CategoryExternal, err, "Slack webhook request failed")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error("Failed to close response body", logger.Fields{
				"error": closeErr.Error(),
			})
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.CategoryExternal, err, "failed to read webhook response")
	}

	if resp.StatusCode >= 400 {
		return errors.Wrap(errors.CategoryExternal,
			fmt.Errorf("slack webhook error: %d - %s", resp.StatusCode, string(body)),
			"failed to send message via webhook")
	}

	// Slack webhooks return "ok" on success
	if string(body) != "ok" {
		return errors.Wrap(errors.CategoryExternal,
			fmt.Errorf("unexpected webhook response: %s", string(body)),
			"webhook response indicates failure")
	}

	log.Info("Successfully sent message to Slack via webhook", logger.Fields{
		"response": string(body),
	})

	return nil
}

// sendViaOAuth sends the message using OAuth token and the Slack Web API
func (s *SlackProcessor) sendViaOAuth(token string, payload []byte) error {
	return s.sendViaOAuthWithURL(token, payload, "https://slack.com/api/chat.postMessage")
}

// sendViaOAuthWithURL allows testing with custom API endpoint
func (s *SlackProcessor) sendViaOAuthWithURL(token string, payload []byte, apiURL string) error {
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(errors.CategoryApplication, err, "failed to create Slack API request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(errors.CategoryExternal, err, "Slack API request failed")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error("Failed to close response body", logger.Fields{
				"error": closeErr.Error(),
			})
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(errors.CategoryExternal, err, "failed to read API response")
	}

	if resp.StatusCode >= 400 {
		return errors.Wrap(errors.CategoryExternal,
			fmt.Errorf("slack API error: %d - %s", resp.StatusCode, string(body)),
			"failed to send message via API")
	}

	// Parse the API response
	var apiResp struct {
		OK    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return errors.Wrap(errors.CategoryApplication, err, "failed to parse Slack API response")
	}

	if !apiResp.OK {
		return errors.Wrap(errors.CategoryExternal,
			fmt.Errorf("slack API returned error: %s", apiResp.Error),
			"Slack API request failed")
	}

	log.Info("Successfully sent message to Slack via API", logger.Fields{
		"ok": apiResp.OK,
	})

	return nil
}
