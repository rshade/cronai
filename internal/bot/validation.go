// Package bot provides the main service for running CronAI in bot mode.
package bot

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// SupportedModels defines the allowed AI models
var SupportedModels = map[string]bool{
	"openai": true,
	"claude": true,
	"gemini": true,
}

// SupportedProcessors defines the allowed processor types
var SupportedProcessors = map[string]bool{
	"console": true,
	"file":    true,
	"email":   true,
	"slack":   true,
	"webhook": true,
	"github":  true,
	"teams":   true,
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s='%s': %s", e.Field, e.Value, e.Message)
}

// ValidateModel validates that the model name is supported
func ValidateModel(model string) error {
	if model == "" {
		return &ValidationError{
			Field:   "model",
			Value:   model,
			Message: "model cannot be empty",
		}
	}

	if !SupportedModels[model] {
		return &ValidationError{
			Field:   "model",
			Value:   model,
			Message: fmt.Sprintf("unsupported model, must be one of: %s", getSupportedKeys(SupportedModels)),
		}
	}

	return nil
}

// ValidateProcessor validates that the processor type is supported
func ValidateProcessor(processor string) error {
	if processor == "" {
		return nil // Processor is optional
	}

	// Extract base processor type (remove prefixes like "file-", "slack-", etc.)
	baseType := strings.SplitN(processor, "-", 2)[0]

	if !SupportedProcessors[baseType] {
		return &ValidationError{
			Field:   "processor",
			Value:   processor,
			Message: fmt.Sprintf("unsupported processor type '%s', must be one of: %s", baseType, getSupportedKeys(SupportedProcessors)),
		}
	}

	return nil
}

// ValidatePort validates that the port is valid and available
func ValidatePort(portStr string) error {
	if portStr == "" {
		return &ValidationError{
			Field:   "port",
			Value:   portStr,
			Message: "port cannot be empty",
		}
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return &ValidationError{
			Field:   "port",
			Value:   portStr,
			Message: "port must be a valid integer",
		}
	}

	// Skip availability check for port 0 (used in tests)
	if port == 0 {
		return nil
	}

	if port < 1 || port > 65535 {
		return &ValidationError{
			Field:   "port",
			Value:   portStr,
			Message: "port must be between 1 and 65535",
		}
	}

	// Check if port is available (optional check)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &ValidationError{
			Field:   "port",
			Value:   portStr,
			Message: fmt.Sprintf("port %d is not available: %v", port, err),
		}
	}
	_ = listener.Close() //nolint:errcheck // Port validation, close error not critical

	return nil
}

// ValidateWebhookSecret validates the webhook secret
func ValidateWebhookSecret(secret string) error {
	if secret == "" {
		return nil // Secret is optional but recommended
	}

	if len(secret) < 8 {
		return &ValidationError{
			Field:   "webhook_secret",
			Value:   "[REDACTED]",
			Message: "webhook secret should be at least 8 characters long for security",
		}
	}

	return nil
}

// getSupportedKeys returns a comma-separated list of supported keys
func getSupportedKeys(m map[string]bool) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
