package processor

import (
	"encoding/json"
	"fmt"
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
	config ProcessorConfig
}

// NewSlackProcessor creates a new Slack processor
func NewSlackProcessor(config ProcessorConfig) (Processor, error) {
	return &SlackProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (s *SlackProcessor) Process(response *models.ModelResponse, templateName string) error {
	// Create template data
	tmplData := template.TemplateData{
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

	// Check for required environment variables
	slackToken := os.Getenv(EnvSlackToken)
	if slackToken == "" {
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("SLACK_TOKEN environment variable not set"),
			"missing Slack configuration")
	}

	return nil
}

// GetType returns the processor type identifier
func (s *SlackProcessor) GetType() string {
	return "slack"
}

// GetConfig returns the processor configuration
func (s *SlackProcessor) GetConfig() ProcessorConfig {
	return s.config
}

// processSlackWithTemplate sends formatted messages to Slack
func (s *SlackProcessor) processSlackWithTemplate(channel string, data template.TemplateData, templateName string) error {
	// Check for Slack token
	slackToken := os.Getenv(EnvSlackToken)
	if slackToken == "" {
		log.Error("Slack token not set", logger.Fields{
			"channel": channel,
		})
		return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("SLACK_TOKEN environment variable not set"),
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

	// In MVP, just log rather than actually sending
	log.Info("Would send to Slack", logger.Fields{
		"channel": channel,
		"payload": string(payloadBytes),
	})

	// For production implementation, we would POST to Slack API
	// slackURL := "https://slack.com/api/chat.postMessage"
	/*
		req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return errors.Wrap(errors.CategoryApplication, err, "failed to create Slack request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+slackToken)

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "Slack API request failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			return errors.Wrap(errors.CategoryExternal,
				fmt.Errorf("Slack API error: %d - %s", resp.StatusCode, string(body)),
				"failed to send message to Slack")
		}
	*/

	return nil
}
