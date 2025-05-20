// Package processor provides implementations for different response processors
package processor

import (
	"fmt"
	"time"

	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

// ConsoleProcessor handles console output
type ConsoleProcessor struct {
	config Config
}

// NewConsoleProcessor creates a new console processor
func NewConsoleProcessor(config Config) (Processor, error) {
	return &ConsoleProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (c *ConsoleProcessor) Process(response *models.ModelResponse, templateName string) error {
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
	tmplData.Metadata["processor"] = c.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	return c.processConsoleOutput(tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (c *ConsoleProcessor) Validate() error {
	// Console processor doesn't require validation
	return nil
}

// GetType returns the processor type identifier
func (c *ConsoleProcessor) GetType() string {
	return "console"
}

// GetConfig returns the processor configuration
func (c *ConsoleProcessor) GetConfig() Config {
	return c.config
}

// processConsoleOutput prints formatted response to console
func (c *ConsoleProcessor) processConsoleOutput(data template.Data, templateName string) error {
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
