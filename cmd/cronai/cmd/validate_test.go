package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateCommand(t *testing.T) {
	// Test that validate command is properly configured
	if validateCmd.Use != "validate" {
		t.Errorf("Expected validate command Use to be 'validate', got %s", validateCmd.Use)
	}

	if validateCmd.Short != "Validate response template files" {
		t.Errorf("Unexpected short description: %s", validateCmd.Short)
	}

	// Verify flags exist
	fileFlag := validateCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("Expected 'file' flag to exist")
	}

	dirFlag := validateCmd.Flags().Lookup("dir")
	if dirFlag == nil {
		t.Error("Expected 'dir' flag to exist")
	}
}

func TestValidateTemplateFile(t *testing.T) {
	// Create test template files
	tmpDir := t.TempDir()

	// Valid template
	validTemplate := filepath.Join(tmpDir, "valid.tmpl")
	validContent := `{{define "test"}}
Hello {{.Name}}!
{{end}}`
	if err := os.WriteFile(validTemplate, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to create valid template: %v", err)
	}

	// Invalid template
	invalidTemplate := filepath.Join(tmpDir, "invalid.tmpl")
	invalidContent := `{{define "test"}}
Hello {{.Name  <!-- Missing closing brace -->
{{end}}`
	if err := os.WriteFile(invalidTemplate, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create invalid template: %v", err)
	}

	// Non-existent template
	nonExistentTemplate := filepath.Join(tmpDir, "nonexistent.tmpl")

	tests := []struct {
		name           string
		filePath       string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:     "valid template",
			filePath: validTemplate,
			expectedOutput: []string{
				"✅ Template",
				"is valid",
			},
		},
		{
			name:     "invalid template",
			filePath: invalidTemplate,
			expectedOutput: []string{
				"❌ Invalid template",
			},
		},
		{
			name:     "non-existent file",
			filePath: nonExistentTemplate,
			expectedOutput: []string{
				"Error reading file",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Run the validation
			validateTemplateFile(tt.filePath)

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := buf.String()

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestValidateTemplateDir(t *testing.T) {
	// Create test template directory
	tmpDir := t.TempDir()

	// Create valid templates
	validTemplate1 := filepath.Join(tmpDir, "valid1.tmpl")
	if err := os.WriteFile(validTemplate1, []byte(`{{define "test"}}Valid{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to create valid template 1: %v", err)
	}

	validTemplate2 := filepath.Join(tmpDir, "valid2.tmpl")
	if err := os.WriteFile(validTemplate2, []byte(`{{define "test2"}}Also valid{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to create valid template 2: %v", err)
	}

	// Create invalid template
	invalidTemplate := filepath.Join(tmpDir, "invalid.tmpl")
	if err := os.WriteFile(invalidTemplate, []byte(`{{define "bad"}}{{.Missing{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to create invalid template: %v", err)
	}

	// Create non-template file (should be ignored)
	nonTemplate := filepath.Join(tmpDir, "readme.txt")
	if err := os.WriteFile(nonTemplate, []byte("Not a template"), 0644); err != nil {
		t.Fatalf("Failed to create non-template file: %v", err)
	}

	// Empty directory
	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	tests := []struct {
		name           string
		dirPath        string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:    "directory with mixed templates",
			dirPath: tmpDir,
			expectedOutput: []string{
				"✅ Template",
				"valid1.tmpl",
				"✅ Template",
				"valid2.tmpl",
				"❌ Invalid template",
				"invalid.tmpl",
				"Some templates have errors",
			},
		},
		{
			name:    "empty directory",
			dirPath: emptyDir,
			expectedOutput: []string{
				"No template files found",
			},
		},
		{
			name:    "non-existent directory",
			dirPath: filepath.Join(tmpDir, "nonexistent"),
			expectedOutput: []string{
				"No template files found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = w

			// Run the validation
			validateTemplateDir(tt.dirPath)

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			os.Stdout = old
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatal(err)
			}
			output := buf.String()

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestValidateCommandExecution(t *testing.T) {
	// Create test templates
	tmpDir := t.TempDir()
	templateFile := filepath.Join(tmpDir, "test.tmpl")
	if err := os.WriteFile(templateFile, []byte(`{{define "test"}}Valid{{end}}`), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	tests := []struct {
		name           string
		fileArg        string
		dirArg         string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:    "validate file",
			fileArg: templateFile,
			dirArg:  "",
			expectedOutput: []string{
				"✅ Template",
				"is valid",
			},
		},
		{
			name:    "validate directory",
			fileArg: "",
			dirArg:  tmpDir,
			expectedOutput: []string{
				"✅ Template",
				"is valid",
				"All templates are valid",
			},
		},
		{
			name:    "no arguments",
			fileArg: "",
			dirArg:  "",
			expectedOutput: []string{
				"Please specify either --file or --dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Mock the Run function by calling it with simulated flag values
			cmd := &cobra.Command{}
			flags := cmd.Flags()
			flags.String("file", tt.fileArg, "")
			flags.String("dir", tt.dirArg, "")

			// Call the Run function directly
			validateCmd.Run(cmd, []string{})

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := buf.String()

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestValidateCommandLongDescription(t *testing.T) {
	expectedStrings := []string{
		"Validate response template files",
		"syntax errors",
		"template syntax errors early",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(validateCmd.Long, expected) {
			t.Errorf("Expected Long description to contain '%s'", expected)
		}
	}
}

func TestValidateCommandExamples(t *testing.T) {
	expectedExamples := []string{
		"cronai validate --file=templates/email_report.tmpl",
		"cronai validate --dir=templates/",
	}

	for _, expected := range expectedExamples {
		if !strings.Contains(validateCmd.Example, expected) {
			t.Errorf("Expected Example to contain '%s'", expected)
		}
	}
}
