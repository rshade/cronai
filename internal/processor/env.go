package processor

import (
	"os"
	"strings"
)

// Environment variable names for processors
const (
	// Slack processor environment variables
	EnvSlackToken = "SLACK_TOKEN"

	// Email processor environment variables
	EnvSMTPServer   = "SMTP_SERVER"
	EnvSMTPPort     = "SMTP_PORT"
	EnvSMTPUser     = "SMTP_USER"
	EnvSMTPPassword = "SMTP_PASSWORD"
	EnvSMTPFrom     = "SMTP_FROM"

	// Webhook processor environment variables
	EnvWebhookURL     = "WEBHOOK_URL"
	EnvWebhookMethod  = "WEBHOOK_METHOD"
	EnvWebhookHeaders = "WEBHOOK_HEADERS"

	// Dynamic webhook environment variables (use with strings.ToUpper)
	EnvWebhookURLPrefix     = "WEBHOOK_URL_"
	EnvWebhookMethodPrefix  = "WEBHOOK_METHOD_"
	EnvWebhookHeadersPrefix = "WEBHOOK_HEADERS_"

	// File processor environment variables
	EnvLogsDirectory = "LOGS_DIRECTORY"

	// GitHub processor environment variables
	EnvGitHubToken = "GITHUB_TOKEN"

	// Default values
	DefaultSMTPPort      = "587"
	DefaultWebhookMethod = "POST"
	DefaultLogsDirectory = "logs"
)

// GetEnvWithDefault returns the value of the environment variable or a default value
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetWebhookURL returns the webhook URL for a specific type or the default
func GetWebhookURL(webhookType string) string {
	// Try type-specific URL first
	if url := os.Getenv(EnvWebhookURLPrefix + strings.ToUpper(webhookType)); url != "" {
		return url
	}
	// Fall back to default URL
	return os.Getenv(EnvWebhookURL)
}

// GetWebhookMethod returns the webhook method for a specific type or the default
func GetWebhookMethod(webhookType string) string {
	// Try type-specific method first
	if method := os.Getenv(EnvWebhookMethodPrefix + strings.ToUpper(webhookType)); method != "" {
		return method
	}
	// Fall back to default method
	return GetEnvWithDefault(EnvWebhookMethod, DefaultWebhookMethod)
}

// GetWebhookHeaders returns the webhook headers for a specific type or the default
func GetWebhookHeaders(webhookType string) string {
	// Try type-specific headers first
	if headers := os.Getenv(EnvWebhookHeadersPrefix + strings.ToUpper(webhookType)); headers != "" {
		return headers
	}
	// Fall back to default headers
	return os.Getenv(EnvWebhookHeaders)
}
