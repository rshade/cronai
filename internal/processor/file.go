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

	// Sanitize and validate filename to prevent path traversal attacks
	safeFilename, err := f.sanitizeFilename(filename, logsDir)
	if err != nil {
		log.Error("Failed to sanitize filename", logger.Fields{
			"original_filename": filename,
			"error":             err.Error(),
		})
		return errors.Wrap(errors.CategorySecurity, err, "filename validation failed")
	}

	// Execute content template
	contentTemplateName := templateName + "_content"
	content := manager.SafeExecute(contentTemplateName, data)
	if content == "" {
		log.Warn("Empty content generated from template", logger.Fields{
			"template": templateName,
			"filename": safeFilename,
		})
		// Use raw content as fallback
		content = data.Content
	}

	// Add to metadata for logging
	data.Metadata["filename_template"] = filenameTemplateName
	data.Metadata["content_template"] = contentTemplateName
	data.Metadata["output_file"] = safeFilename

	// Create parent directory if needed - ensure we only work with the sanitized, validated path
	// The safeFilename has been validated to be within the base directory by sanitizeFilename
	parentDir := filepath.Dir(safeFilename)
	// Additional safety check: ensure parent directory is not attempting path traversal
	if parentDir != "." && parentDir != "/" && !strings.Contains(parentDir, "..") {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			log.Error("Failed to create parent directory", logger.Fields{
				"directory": parentDir,
				"error":     err.Error(),
			})
			return errors.Wrap(errors.CategorySystem, err, "failed to create parent directory for output file")
		}
	}

	// Write to file - using the validated and sanitized filename
	// The safeFilename has been validated to be within the base directory by sanitizeFilename
	// Additional safety check: ensure filename doesn't contain path traversal patterns
	if strings.Contains(safeFilename, "..") {
		return errors.Wrap(errors.CategorySecurity,
			fmt.Errorf("sanitized filename still contains path traversal patterns: %s", safeFilename),
			"invalid sanitized filename")
	}
	err = os.WriteFile(safeFilename, []byte(content), 0644)
	if err != nil {
		log.Error("Failed to write response to file", logger.Fields{
			"filename": safeFilename,
			"error":    err.Error(),
		})
		return errors.Wrap(errors.CategorySystem, err, "failed to write response to file")
	}

	log.Info("Response saved to file", logger.Fields{
		"filename":    safeFilename,
		"content_len": len(content),
	})
	return nil
}

// sanitizeFilename validates and sanitizes a filename to prevent path traversal attacks
func (f *FileProcessor) sanitizeFilename(filename, baseDir string) (string, error) {
	// First normalize path separators to handle Windows-style paths on any OS
	normalizedFilename := strings.ReplaceAll(filename, "\\", "/")

	// Reject relative filenames containing path traversal patterns
	// Absolute paths with .. are handled later by extracting basename if outside base dir
	if !filepath.IsAbs(normalizedFilename) && strings.Contains(normalizedFilename, "..") {
		return "", fmt.Errorf("invalid filename: path traversal patterns detected in %s", filename)
	}

	// Clean the base directory path
	cleanBaseDir := filepath.Clean(baseDir)

	// Handle different cases based on whether filename is absolute or relative
	var targetPath string
	if filepath.IsAbs(normalizedFilename) {
		// For absolute paths, check if they're within the base directory
		cleanFilename := filepath.Clean(normalizedFilename)
		absBaseDir, err := filepath.Abs(cleanBaseDir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path for base directory: %w", err)
		}

		// If absolute path is within base directory, use it
		if strings.HasPrefix(cleanFilename, absBaseDir+string(filepath.Separator)) || cleanFilename == absBaseDir {
			targetPath = cleanFilename
		} else {
			// If outside base directory, use only the basename
			targetPath = filepath.Join(cleanBaseDir, filepath.Base(cleanFilename))
		}
	} else {
		// For relative paths, join with base directory first, then clean
		targetPath = filepath.Clean(filepath.Join(cleanBaseDir, normalizedFilename))
	}

	// Resolve to absolute paths for final validation
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for target: %w", err)
	}

	absBaseDir, err := filepath.Abs(cleanBaseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for base directory: %w", err)
	}

	// Final check: ensure the resolved path is within the base directory
	if !strings.HasPrefix(absTargetPath, absBaseDir+string(filepath.Separator)) &&
		absTargetPath != absBaseDir {
		return "", fmt.Errorf("path traversal detected: %s resolves to %s which is outside base directory %s", normalizedFilename, absTargetPath, absBaseDir)
	}

	return targetPath, nil
}
