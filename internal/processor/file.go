package processor

import (
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

// FileProcessor handles file output
type FileProcessor struct {
	config Config
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(config Config) (Processor, error) {
	return &FileProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (f *FileProcessor) Process(response *models.ModelResponse, templateName string) error {
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
	tmplData.Metadata["processor"] = f.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	return f.processFileWithTemplate(tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (f *FileProcessor) Validate() error {
	// File processor doesn't require specific validation
	// The directory will be created if it doesn't exist
	return nil
}

// GetType returns the processor type identifier
func (f *FileProcessor) GetType() string {
	return "file"
}

// GetConfig returns the processor configuration
func (f *FileProcessor) GetConfig() Config {
	return f.config
}

// processFileWithTemplate saves response to file using template
func (f *FileProcessor) processFileWithTemplate(data template.Data, templateName string) error {
	// Create logs directory if it doesn't exist
	logsDir := GetEnvWithDefault(EnvLogsDirectory, DefaultLogsDirectory)

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
