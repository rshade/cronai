package processor

import (
	"github.com/rshade/cronai/internal/models"
)

// Processor defines the interface for all response processors
type Processor interface {
	// Process handles the model response with optional template
	Process(response *models.ModelResponse, templateName string) error

	// Validate checks if the processor is properly configured
	Validate() error

	// GetType returns the processor type identifier
	GetType() string

	// GetConfig returns the processor configuration
	GetConfig() ProcessorConfig
}

// ProcessorConfig represents standardized configuration for processors
type ProcessorConfig struct {
	// Type identifies the processor type (email, slack, webhook, file, etc.)
	Type string `json:"type"`

	// Target specifies the destination (email address, webhook URL, file path, etc.)
	Target string `json:"target"`

	// Options contains processor-specific configuration
	Options map[string]interface{} `json:"options,omitempty"`

	// TemplateName specifies the template to use (optional)
	TemplateName string `json:"template_name,omitempty"`

	// Environment contains environment variable names used by the processor
	Environment map[string]string `json:"environment,omitempty"`
}

// ProcessorFactory is a function that creates a new processor instance
type ProcessorFactory func(config ProcessorConfig) (Processor, error)
