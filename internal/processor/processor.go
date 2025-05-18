package processor

import (
	"fmt"
	"strings"

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

	// Parse processor name to determine type and target
	var processorType, target string

	// Handle special processor formats
	if strings.HasPrefix(processorName, "slack-") {
		processorType = "slack"
		target = strings.TrimPrefix(processorName, "slack-")
	} else if strings.HasPrefix(processorName, "email-") {
		processorType = "email"
		target = strings.TrimPrefix(processorName, "email-")
	} else if strings.HasPrefix(processorName, "webhook-") {
		processorType = "webhook"
		target = strings.TrimPrefix(processorName, "webhook-")
	} else {
		// Handle standard processors
		switch processorName {
		case "log-to-file", "file":
			processorType = "file"
		case "console":
			processorType = "console"
		default:
			log.Error("Unsupported processor", logger.Fields{
				"processor": processorName,
			})
			return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("unsupported processor: %s", processorName),
				"processor type not recognized")
		}
	}

	// Create processor configuration
	config := ProcessorConfig{
		Type:         processorType,
		Target:       target,
		TemplateName: templateName,
	}

	// Create processor using registry
	registry := GetRegistry()
	processor, err := registry.CreateProcessor(config)
	if err != nil {
		return errors.Wrap(errors.CategoryApplication, err, "failed to create processor")
	}

	// Process the response
	return processor.Process(response, templateName)
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
