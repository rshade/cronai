package processor

import (
	"fmt"
	"os"
	"path/filepath"
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
		Metadata:    make(map[string]string), // Initialize metadata map
	}

	// Add additional metadata if needed
	tmplData.Metadata["processor"] = processorName
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
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

	// Handle webhook formats with type
	if strings.HasPrefix(processorName, "webhook-") {
		webhookType := strings.TrimPrefix(processorName, "webhook-")
		return processWebhookWithTemplate(webhookType, tmplData, templateName)
	}

	// Handle standard processors
	switch processorName {
	case "log-to-file", "file":
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

	// Add to metadata for logging
	data.Metadata["slack_channel"] = channel
	data.Metadata["template_used"] = templateName

	// TODO: Implement actual Slack API call
	fmt.Printf("Would send to Slack channel %s with template %s\n", channel, templateName)
	fmt.Printf("Payload size: %d bytes\n", len(payload))

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

	// Validate that required templates exist
	subjectTemplateName := templateName + "_subject"
	htmlTemplateName := templateName + "_html"
	textTemplateName := templateName + "_text"

	// Execute subject template
	subject := manager.SafeExecute(subjectTemplateName, data)

	// Execute HTML body template
	htmlBody := manager.SafeExecute(htmlTemplateName, data)

	// Execute text body template (fallback)
	textBody := manager.SafeExecute(textTemplateName, data)

	// Add to metadata for logging
	data.Metadata["email_recipient"] = email
	data.Metadata["subject_template"] = subjectTemplateName
	data.Metadata["html_template"] = htmlTemplateName
	data.Metadata["text_template"] = textTemplateName

	// TODO: Implement actual email sending
	fmt.Printf("Would send email to %s with subject: %s\n", email, subject)
	fmt.Printf("HTML body length: %d bytes, Text body length: %d bytes\n",
		len(htmlBody), len(textBody))

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
		if !manager.TemplateExists(templateName) {
			templateName = "default_webhook"
		}
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)

	// Add to metadata for logging
	data.Metadata["webhook_type"] = webhookType
	data.Metadata["webhook_url_env"] = fmt.Sprintf("WEBHOOK_URL_%s", strings.ToUpper(webhookType))
	data.Metadata["template_used"] = templateName

	// TODO: Implement actual webhook call
	fmt.Printf("Would send to webhook type %s with template %s\n",
		webhookType, templateName)
	fmt.Printf("Payload size: %d bytes\n", len(payload))

	return nil
}

// processFileWithTemplate saves response to file using template
func processFileWithTemplate(data template.TemplateData, templateName string) error {
	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_file"
	}

	// Execute filename template
	filenameTemplateName := templateName + "_filename"
	filename := manager.SafeExecute(filenameTemplateName, data)

	// Execute content template
	contentTemplateName := templateName + "_content"
	content := manager.SafeExecute(contentTemplateName, data)

	// Add to metadata for logging
	data.Metadata["filename_template"] = filenameTemplateName
	data.Metadata["content_template"] = contentTemplateName
	data.Metadata["output_file"] = filename

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for output file: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	fmt.Printf("Response saved to file: %s (%d bytes)\n", filename, len(content))
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
