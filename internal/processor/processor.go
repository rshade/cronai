package processor

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// ProcessResponse processes a model response using the specified processor
func ProcessResponse(processorName string, response *models.ModelResponse, templateName string) error {
	// Create template data
	tmplData := template.TemplateData{
		Content:     response.Content,
		Model:       response.Model,
		Timestamp:   time.Now(),
		PromptName:  response.PromptName,
		Variables:   response.Variables,
		ExecutionID: response.ExecutionID,
	}

	// Handle special processor formats
	if strings.HasPrefix(processorName, "slack-") {
		slackChannel := strings.TrimPrefix(processorName, "slack-")
		return processSlackWithTemplate(slackChannel, tmplData, templateName)
	}

	if strings.HasPrefix(processorName, "email-") {
		emailAddress := strings.TrimPrefix(processorName, "email-")
		return processEmailWithTemplate(emailAddress, tmplData, templateName)
	}

	// Handle standard processors
	switch processorName {
	case "webhook-monitoring":
		return processWebhookWithTemplate("monitoring", tmplData, templateName)
	case "log-to-file":
		return processFileWithTemplate(tmplData, templateName)
	default:
		return fmt.Errorf("unsupported processor: %s", processorName)
	}
}

// processSlackWithTemplate sends formatted messages to Slack
func processSlackWithTemplate(channel string, data template.TemplateData, templateName string) error {
	// Check for Slack token
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("SLACK_TOKEN environment variable not set")
	}

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_slack"
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)

	// TODO: Implement actual Slack API call
	fmt.Printf("Would send to Slack channel %s: %s\n", channel, payload)

	return nil
}

// processEmailWithTemplate with multipart support
func processEmailWithTemplate(email string, data template.TemplateData, templateName string) error {
	// Check for SMTP settings
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		return fmt.Errorf("SMTP_SERVER environment variable not set")
	}

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_email"
	}

	// Execute subject template
	subject := manager.SafeExecute(templateName+"_subject", data)

	// Execute HTML body template - currently just logging but would use in actual implementation
	_ = manager.SafeExecute(templateName+"_html", data)

	// Execute text body template (fallback) - currently just logging but would use in actual implementation
	_ = manager.SafeExecute(templateName+"_text", data)

	// TODO: Implement actual email sending
	fmt.Printf("Would send email to %s: %s\n", email, subject)

	return nil
}

// processWebhookWithTemplate sends webhook payload using template
func processWebhookWithTemplate(webhookType string, data template.TemplateData, templateName string) error {
	// Check for webhook URL
	webhookURL := os.Getenv(fmt.Sprintf("WEBHOOK_URL_%s", strings.ToUpper(webhookType)))
	if webhookURL == "" {
		webhookURL = os.Getenv("WEBHOOK_URL") // fallback
		if webhookURL == "" {
			return fmt.Errorf("webhook URL environment variable not set")
		}
	}

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = fmt.Sprintf("default_webhook_%s", webhookType)
		// If specific webhook type template doesn't exist, use generic
		if _, err := manager.GetTemplate(templateName); err != nil {
			templateName = "default_webhook"
		}
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)

	// TODO: Implement actual webhook call
	fmt.Printf("Would send to webhook %s: %s\n", webhookURL, payload)

	return nil
}

// processFileWithTemplate saves response to file using template
func processFileWithTemplate(data template.TemplateData, templateName string) error {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_file"
	}

	// Execute filename template
	filename := manager.SafeExecute(templateName+"_filename", data)

	// Execute content template
	content := manager.SafeExecute(templateName+"_content", data)

	// Write to file
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	fmt.Printf("Response saved to file: %s\n", filename)
	return nil
}

// InitTemplates initializes the template system
func InitTemplates(templateDir string) error {
	manager := template.GetManager()

	// Load templates from directory if specified
	if templateDir != "" {
		if err := manager.LoadTemplatesFromDir(templateDir); err != nil {
			return fmt.Errorf("failed to load templates from directory: %w", err)
		}
	}

	return nil
}
