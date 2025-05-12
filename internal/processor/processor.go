package processor

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/models"
)

// ProcessResponse processes a model response using the specified processor
func ProcessResponse(processorName string, response *models.ModelResponse) error {
	// Handle special processor formats
	if strings.HasPrefix(processorName, "slack-") {
		slackChannel := strings.TrimPrefix(processorName, "slack-")
		return processSlack(slackChannel, response)
	}

	if strings.HasPrefix(processorName, "email-") {
		emailAddress := strings.TrimPrefix(processorName, "email-")
		return processEmail(emailAddress, response)
	}

	// Handle standard processors
	switch processorName {
	case "webhook-monitoring":
		return processWebhook("monitoring", response)
	case "log-to-file":
		return processFile(response)
	default:
		return fmt.Errorf("unsupported processor: %s", processorName)
	}
}

// processSlack sends the response to a Slack channel
func processSlack(channel string, response *models.ModelResponse) error {
	// Check for Slack token
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("SLACK_TOKEN environment variable not set")
	}

	// TODO: Implement actual Slack API call
	fmt.Printf("Would send to Slack channel %s: %s\n", channel, response.Content)
	return nil
}

// processEmail sends the response via email
func processEmail(email string, response *models.ModelResponse) error {
	// Check for SMTP settings
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		return fmt.Errorf("SMTP_SERVER environment variable not set")
	}

	// TODO: Implement actual email sending
	fmt.Printf("Would send email to %s: %s\n", email, response.Content)
	return nil
}

// processWebhook sends the response to a webhook
func processWebhook(webhookType string, response *models.ModelResponse) error {
	// Check for webhook URL
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL environment variable not set")
	}

	// TODO: Implement actual webhook call
	fmt.Printf("Would send to webhook (%s): %s\n", webhookType, response.Content)
	return nil
}

// processFile saves the response to a file
func processFile(response *models.ModelResponse) error {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create filename with timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	filename := fmt.Sprintf("logs/%s-%s.txt", response.Model, timestamp)

	// Write response to file
	err = os.WriteFile(filename, []byte(response.Content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	fmt.Printf("Response saved to file: %s\n", filename)
	return nil
}
