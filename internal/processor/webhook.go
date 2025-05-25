package processor

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// WebhookProcessor handles webhook requests
type WebhookProcessor struct {
	config Config
}

// NewWebhookProcessor creates a new webhook processor
func NewWebhookProcessor(config Config) (Processor, error) {
	return &WebhookProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (w *WebhookProcessor) Process(response *models.ModelResponse, templateName string) error {
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
	tmplData.Metadata["processor"] = w.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	// Extract webhook type from target or use default
	webhookType := w.config.Target
	if webhookType == "" {
		webhookType = "default"
	}

	return w.processWebhookWithTemplate(webhookType, tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (w *WebhookProcessor) Validate() error {
	// Check for webhook URL
	webhookURL := GetWebhookURL(w.config.Target)
	if webhookURL == "" {
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("webhook URL environment variable not set for type: %s", w.config.Target),
			"missing webhook configuration")
	}

	return nil
}

// GetType returns the processor type identifier
func (w *WebhookProcessor) GetType() string {
	return "webhook"
}

// GetConfig returns the processor configuration
func (w *WebhookProcessor) GetConfig() Config {
	return w.config
}

// truncatePayload truncates a JSON payload to fit within the specified size limit
// while preserving valid JSON structure
func truncatePayload(payload string, maxSize int) string {
	if len(payload) <= maxSize {
		return payload
	}

	// Parse the JSON to handle truncation properly
	var jsonData interface{}
	if err := json.Unmarshal([]byte(payload), &jsonData); err != nil {
		// If we can't parse as JSON, do a simple string truncation
		return payload[:maxSize]
	}

	// Convert back to JSON with indentation for better truncation
	truncated, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return payload[:maxSize]
	}

	// If still too large, truncate and try to close any open JSON structures
	if len(truncated) > maxSize {
		truncated = truncated[:maxSize]
		// Find the last complete JSON structure
		lastBrace := strings.LastIndex(string(truncated), "}")
		if lastBrace > 0 {
			truncated = truncated[:lastBrace+1]
		}
	}

	return string(truncated)
}

// processWebhookWithTemplate sends webhook payload using template
func (w *WebhookProcessor) processWebhookWithTemplate(webhookType string, data template.Data, templateName string) error {
	// Check for webhook URL
	webhookURL := GetWebhookURL(webhookType)
	if webhookURL == "" {
		log.Error("Webhook URL not set", logger.Fields{
			"webhook_type": webhookType,
		})
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("webhook URL environment variable not set"),
			"missing webhook configuration")
	}

	// Get webhook method
	webhookMethod := GetWebhookMethod(webhookType)

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = fmt.Sprintf("default_webhook_%s", webhookType)

		// If specific webhook type template doesn't exist, use generic
		if !manager.TemplateExists(templateName) {
			templateName = "default_webhook"
		}
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)
	if payload == "" {
		log.Error("Failed to generate webhook payload", logger.Fields{
			"template": templateName,
			"type":     webhookType,
		})
		return errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("empty payload generated from template %s", templateName),
			"webhook payload generation failed")
	}

	// Add to metadata for logging
	data.Metadata["webhook_type"] = webhookType
	data.Metadata["webhook_url_env"] = fmt.Sprintf("WEBHOOK_URL_%s", strings.ToUpper(webhookType))
	data.Metadata["template_used"] = templateName

	// Validate JSON payload
	var jsonPayload interface{}
	if err := json.Unmarshal([]byte(payload), &jsonPayload); err != nil {
		log.Error("Invalid webhook JSON payload", logger.Fields{
			"template": templateName,
			"error":    err.Error(),
		})
		return errors.Wrap(errors.CategoryApplication, err, "webhook payload is not valid JSON")
	}

	// Check message size limit for Teams (25KB)
	if webhookType == "teams" && len(payload) > 25*1024 {
		log.Warn("Teams webhook payload exceeds 25KB limit, truncating", logger.Fields{
			"original_size": len(payload),
		})
		payload = truncatePayload(payload, 25*1024)
	}

	// In MVP, just log rather than actually sending
	log.Info("Would send webhook", logger.Fields{
		"url":     webhookURL,
		"method":  webhookMethod,
		"type":    webhookType,
		"payload": payload,
	})

	// For production implementation:
	/*
		req, err := http.NewRequest(webhookMethod, webhookURL, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			return errors.Wrap(errors.CategoryApplication, err, "failed to create webhook request")
		}

		req.Header.Set("Content-Type", "application/json")

		// Add any custom headers
		customHeaders := GetWebhookHeaders(webhookType)

		if customHeaders != "" {
			headers := strings.Split(customHeaders, ",")
			for _, header := range headers {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "webhook request failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			return errors.Wrap(errors.CategoryExternal,
				fmt.Errorf("webhook error: %d - %s", resp.StatusCode, string(body)),
				"webhook request returned error status")
		}
	*/

	return nil
}
