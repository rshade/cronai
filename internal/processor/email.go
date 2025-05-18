package processor

import (
	"fmt"
	"os"
	"time"

	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// EmailProcessor handles email processing
type EmailProcessor struct {
	config Config
}

// NewEmailProcessor creates a new email processor
func NewEmailProcessor(config Config) (Processor, error) {
	return &EmailProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (e *EmailProcessor) Process(response *models.ModelResponse, templateName string) error {
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
	tmplData.Metadata["processor"] = e.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	return e.processEmailWithTemplate(e.config.Target, tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (e *EmailProcessor) Validate() error {
	if e.config.Target == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("email target cannot be empty"),
			"invalid email processor configuration")
	}

	// Check for required environment variables
	smtpServer := os.Getenv(EnvSMTPServer)
	if smtpServer == "" {
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("SMTP_SERVER environment variable not set"),
			"missing email configuration")
	}

	return nil
}

// GetType returns the processor type identifier
func (e *EmailProcessor) GetType() string {
	return "email"
}

// GetConfig returns the processor configuration
func (e *EmailProcessor) GetConfig() Config {
	return e.config
}

// processEmailWithTemplate with multipart support
func (e *EmailProcessor) processEmailWithTemplate(email string, data template.Data, templateName string) error {
	// Check for SMTP settings
	smtpServer := os.Getenv(EnvSMTPServer)
	if smtpServer == "" {
		log.Error("SMTP server not set", logger.Fields{
			"email": email,
		})
		return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("SMTP_SERVER environment variable not set"),
			"missing email configuration")
	}

	// Get additional SMTP settings
	smtpPort := GetEnvWithDefault(EnvSMTPPort, DefaultSMTPPort)
	smtpUser := os.Getenv(EnvSMTPUser)

	// We'll need SMTP_PASSWORD for actual sending, but just check it exists for now
	if os.Getenv(EnvSMTPPassword) == "" {
		log.Warn("SMTP_PASSWORD not set", nil)
	}

	smtpFrom := os.Getenv(EnvSMTPFrom)
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
