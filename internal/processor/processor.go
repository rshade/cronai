package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// Default logger
var log = logger.DefaultLogger()

// SetLogger sets the logger for the processor package
func SetLogger(l *logger.Logger) {
	log = l
}

// ProcessorOptions contains options for processors
type ProcessorOptions struct {
	TemplateDir string // Directory containing custom templates
}

// ProcessResponse processes a model response using the specified processor
func ProcessResponse(processorName string, response *models.ModelResponse, templateName string) error {
	log.Info("Processing response", logger.Fields{
		"processor":   processorName,
		"model":       response.Model,
		"prompt":      response.PromptName,
		"template":    templateName,
		"execution":   response.ExecutionID,
		"timestamp":   response.Timestamp,
		"content_len": len(response.Content),
	})

	// Create template data
	tmplData := template.TemplateData{
		Content:     response.Content,
		Model:       response.Model,
		Timestamp:   response.Timestamp,
		PromptName:  response.PromptName,
		Variables:   response.Variables,
		ExecutionID: response.ExecutionID,
		Metadata:    make(map[string]string), // Initialize metadata map
	}

	// Add standard metadata fields
	tmplData.Metadata["timestamp"] = response.Timestamp.Format(time.RFC3339)
	tmplData.Metadata["date"] = response.Timestamp.Format("2006-01-02")
	tmplData.Metadata["time"] = response.Timestamp.Format("15:04:05")
	tmplData.Metadata["execution_id"] = response.ExecutionID
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
	case "console":
		return processConsoleOutput(tmplData, templateName)
	default:
		log.Error("Unsupported processor", logger.Fields{
			"processor": processorName,
		})
		return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("unsupported processor: %s", processorName),
			"processor type not recognized")
	}
}

// processSlackWithTemplate sends formatted messages to Slack
func processSlackWithTemplate(channel string, data template.TemplateData, templateName string) error {
	// Check for Slack token
	slackToken := os.Getenv("SLACK_TOKEN")
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

// processEmailWithTemplate with multipart support
func processEmailWithTemplate(email string, data template.TemplateData, templateName string) error {
	// Check for SMTP settings
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		log.Error("SMTP server not set", logger.Fields{
			"email": email,
		})
		return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("SMTP_SERVER environment variable not set"),
			"missing email configuration")
	}

	// Get additional SMTP settings
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // Default SMTP port
	}

	smtpUser := os.Getenv("SMTP_USER")
	// We'll need SMTP_PASSWORD for actual sending, but just check it exists for now
	if os.Getenv("SMTP_PASSWORD") == "" {
		log.Warn("SMTP_PASSWORD not set", nil)
	}

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		log.Warn("SMTP_FROM not set, using SMTP_USER", nil)
		smtpFrom = smtpUser
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
	if subject == "" {
		log.Warn("Empty email subject, using default", logger.Fields{
			"template": templateName,
		})
		subject = fmt.Sprintf("AI Response: %s", data.PromptName)
	}

	// Execute HTML body template
	htmlBody := manager.SafeExecute(htmlTemplateName, data)
	if htmlBody == "" {
		log.Warn("Empty HTML body, falling back to text", logger.Fields{
			"template": templateName,
		})
	}

	// Execute text body template (fallback)
	textBody := manager.SafeExecute(textTemplateName, data)
	if textBody == "" && htmlBody == "" {
		log.Error("Both HTML and text bodies empty", logger.Fields{
			"template": templateName,
		})
		return errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("both HTML and text templates for %s produced empty output", templateName),
			"email body generation failed")
	}

	// Add to metadata for logging
	data.Metadata["email_recipient"] = email
	data.Metadata["subject_template"] = subjectTemplateName
	data.Metadata["html_template"] = htmlTemplateName
	data.Metadata["text_template"] = textTemplateName

	// In MVP, just log rather than actually sending
	log.Info("Would send email", logger.Fields{
		"to":      email,
		"subject": subject,
		"from":    smtpFrom,
		"server":  smtpServer,
		"port":    smtpPort,
	})

	// For production implementation, we would use a proper email library:
	/*
		m := gomail.NewMessage()
		m.SetHeader("From", smtpFrom)
		m.SetHeader("To", email)
		m.SetHeader("Subject", subject)

		// Add text part
		if textBody != "" {
			m.SetBody("text/plain", textBody)
		}

		// Add HTML part if available
		if htmlBody != "" {
			if textBody != "" {
				m.AddAlternative("text/html", htmlBody)
			} else {
				m.SetBody("text/html", htmlBody)
			}
		}

		// Create dialer
		port, _ := strconv.Atoi(smtpPort)
		d := gomail.NewDialer(smtpServer, port, smtpUser, smtpPass)

		// Send email
		if err := d.DialAndSend(m); err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "failed to send email")
		}
	*/

	return nil
}

// processWebhookWithTemplate sends webhook payload using template
func processWebhookWithTemplate(webhookType string, data template.TemplateData, templateName string) error {
	// Check for webhook URL
	webhookURL := os.Getenv(fmt.Sprintf("WEBHOOK_URL_%s", strings.ToUpper(webhookType)))
	if webhookURL == "" {
		webhookURL = os.Getenv("WEBHOOK_URL") // fallback
		if webhookURL == "" {
			log.Error("Webhook URL not set", logger.Fields{
				"webhook_type": webhookType,
			})
			return errors.Wrap(errors.CategoryConfiguration,
				fmt.Errorf("webhook URL environment variable not set"),
				"missing webhook configuration")
		}
	}

	// Get webhook method
	webhookMethod := os.Getenv(fmt.Sprintf("WEBHOOK_METHOD_%s", strings.ToUpper(webhookType)))
	if webhookMethod == "" {
		webhookMethod = os.Getenv("WEBHOOK_METHOD") // fallback
		if webhookMethod == "" {
			webhookMethod = "POST" // default
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
		customHeaders := os.Getenv(fmt.Sprintf("WEBHOOK_HEADERS_%s", strings.ToUpper(webhookType)))
		if customHeaders == "" {
			customHeaders = os.Getenv("WEBHOOK_HEADERS") // fallback
		}

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

// processFileWithTemplate saves response to file using template
func processFileWithTemplate(data template.TemplateData, templateName string) error {
	// Create logs directory if it doesn't exist
	logsDir := os.Getenv("LOGS_DIRECTORY")
	if logsDir == "" {
		logsDir = "logs"
	}

	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		log.Error("Failed to create logs directory", logger.Fields{
			"directory": logsDir,
			"error":     err.Error(),
		})
		return errors.Wrap(errors.CategorySystem, err, "failed to create logs directory")
	}

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_file"
	}

	// Execute filename template
	filenameTemplateName := templateName + "_filename"
	filename := manager.SafeExecute(filenameTemplateName, data)
	if filename == "" {
		log.Error("Failed to generate filename", logger.Fields{
			"template": templateName,
		})
		return errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("empty filename generated from template %s", templateName),
			"filename generation failed")
	}

	// Make sure filename is within logs directory
	if !strings.HasPrefix(filename, logsDir) {
		filename = filepath.Join(logsDir, filepath.Base(filename))
	}

	// Execute content template
	contentTemplateName := templateName + "_content"
	content := manager.SafeExecute(contentTemplateName, data)
	if content == "" {
		log.Warn("Empty content generated from template", logger.Fields{
			"template": templateName,
			"filename": filename,
		})
		// Use raw content as fallback
		content = data.Content
	}

	// Add to metadata for logging
	data.Metadata["filename_template"] = filenameTemplateName
	data.Metadata["content_template"] = contentTemplateName
	data.Metadata["output_file"] = filename

	// Create parent directory if needed
	parentDir := filepath.Dir(filename)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		log.Error("Failed to create parent directory", logger.Fields{
			"directory": parentDir,
			"error":     err.Error(),
		})
		return errors.Wrap(errors.CategorySystem, err, "failed to create parent directory for output file")
	}

	// Write to file
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Error("Failed to write response to file", logger.Fields{
			"filename": filename,
			"error":    err.Error(),
		})
		return errors.Wrap(errors.CategorySystem, err, "failed to write response to file")
	}

	log.Info("Response saved to file", logger.Fields{
		"filename":    filename,
		"content_len": len(content),
	})
	return nil
}

// processConsoleOutput prints formatted response to console
func processConsoleOutput(data template.TemplateData, templateName string) error {
	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = "default_console"

		// Register default console template if it doesn't exist
		if _, err := manager.GetTemplate(templateName); err != nil {
			err = manager.RegisterTemplate(templateName, `
==========================================================
AI Response: {{.PromptName}}
==========================================================
Model: {{.Model}}
Time: {{.Timestamp.Format "2006-01-02 15:04:05"}}
==========================================================

{{.Content}}

==========================================================
`)
			if err != nil {
				log.Error("Failed to register default console template", logger.Fields{
					"error": err.Error(),
				})
			}
		}
	}

	// Execute template to get output
	output := manager.SafeExecute(templateName, data)
	if output == "" {
		// Fallback to basic output
		output = fmt.Sprintf("AI Response (%s): %s\n\n%s",
			data.Model,
			data.PromptName,
			data.Content)
	}

	// Print to console
	fmt.Println(output)
	return nil
}

// InitTemplates initializes the template system
func InitTemplates(templateDir string) error {
	log.Info("Initializing template system", logger.Fields{
		"template_dir": templateDir,
	})

	manager := template.GetManager()

	// Register default templates
	log.Debug("Registering default templates", nil)

	// First, try to load library templates
	if err := manager.LoadLibraryTemplates(); err != nil {
		log.Warn("Failed to load library templates", logger.Fields{
			"error": err.Error(),
		})
		// Continue anyway as we have fallbacks
	}

	// Load templates from directory if specified
	if templateDir != "" {
		log.Debug("Loading templates from directory", logger.Fields{
			"directory": templateDir,
		})

		if err := manager.LoadTemplatesFromDir(templateDir); err != nil {
			log.Error("Failed to load templates from directory", logger.Fields{
				"directory": templateDir,
				"error":     err.Error(),
			})
			return errors.Wrap(errors.CategorySystem, err, "failed to load templates from directory")
		}
	}

	return nil
}
