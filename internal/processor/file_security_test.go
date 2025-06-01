// Package processor provides response processing functionality.
package processor

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFileProcessor_sanitizeFilename(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		baseDir     string
		expectError bool
		expectPath  string
	}{
		{
			name:        "normal filename",
			filename:    "test.txt",
			baseDir:     "/logs",
			expectError: false,
			expectPath:  "/logs/test.txt",
		},
		{
			name:        "filename with path in base directory",
			filename:    "/logs/subdir/test.txt",
			baseDir:     "/logs",
			expectError: false,
			expectPath:  "/logs/subdir/test.txt",
		},
		{
			name:        "path traversal attempt with ..",
			filename:    "../../../etc/passwd",
			baseDir:     "/logs",
			expectError: true,
		},
		{
			name:        "path traversal attempt with absolute path",
			filename:    "/etc/passwd",
			baseDir:     "/logs",
			expectError: false,
			expectPath:  "/logs/passwd", // Only basename is used
		},
		{
			name:        "complex path traversal attempt",
			filename:    "/logs/../../../etc/passwd",
			baseDir:     "/logs",
			expectError: false,
			expectPath:  "/logs/passwd", // Only basename is used for absolute paths outside base
		},
		{
			name:        "relative path outside base",
			filename:    "../../outside.txt",
			baseDir:     "/logs",
			expectError: true, // Path traversal detected
		},
		{
			name:        "subdirectory creation allowed",
			filename:    "subdir/test.txt",
			baseDir:     "/logs",
			expectError: false,
			expectPath:  "/logs/subdir/test.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &FileProcessor{}
			result, err := processor.sanitizeFilename(tt.filename, tt.baseDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for filename %s, but got none", tt.filename)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for filename %s: %v", tt.filename, err)
				}
				if result != tt.expectPath {
					t.Errorf("Expected path %s, got %s", tt.expectPath, result)
				}
			}
		})
	}
}

func TestFileProcessor_sanitizeFilename_pathTraversalPrevention(t *testing.T) {
	processor := &FileProcessor{}
	baseDir := "/safe/logs"

	// Test various path traversal attempts
	maliciousInputs := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"/etc/passwd",
		"./../../secret.txt",
		"logs/../../../etc/hosts",
		"subdir/../../etc/shadow",
	}

	for _, input := range maliciousInputs {
		t.Run(input, func(t *testing.T) {
			result, err := processor.sanitizeFilename(input, baseDir)

			if err != nil {
				// Error is acceptable for obvious traversal attempts
				return
			}

			// If no error, ensure result is safe
			if result != "" {
				// Check that the result is within the base directory
				absResult, err := filepath.Abs(result)
				if err != nil {
					t.Errorf("Failed to get absolute path of result: %v", err)
					return
				}

				absBase, err := filepath.Abs(baseDir)
				if err != nil {
					t.Errorf("Failed to get absolute path of base: %v", err)
					return
				}

				// Check if the result path is within the base directory
				relPath, err := filepath.Rel(absBase, absResult)
				if err != nil || strings.HasPrefix(relPath, "..") {
					t.Errorf("Sanitized path %s is not within base directory %s", absResult, absBase)
				}
			}
		})
	}
}
