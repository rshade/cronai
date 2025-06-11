package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

func TestProcessResponse(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "processor_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Ensure default processors are registered
	registry := GetRegistry()
	registry.RegisterDefaults()

	// Setup test cases
	testCases := []struct {
		name          string
		processor     string
		response      *models.ModelResponse
		templateName  string
		expectedError bool
		setup         func() error
		validate      func() error
	}{
		{
			name:      "File Processor",
			processor: "log-to-file",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Variables:   map[string]string{"key": "value"},
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "test_file",
			expectedError: false,
			setup: func() error {
				// Set logs directory to temp dir
				if err := os.Setenv("LOGS_DIRECTORY", tempDir); err != nil {
					return err
				}

				// Register test templates
				manager := template.GetManager()
				err := manager.RegisterTemplate("test_file_filename", filepath.Join(tempDir, "test-output.txt"))
				if err != nil {
					return err
				}
				return manager.RegisterTemplate("test_file_content", "Content: {{.Content}}\nModel: {{.Model}}")
			},
			validate: func() error {
				// Check if file was created
				content, err := os.ReadFile(filepath.Join(tempDir, "test-output.txt"))
				if err != nil {
					return err
				}
				// Verify content
				expected := "Content: Test content\nModel: test-model"
				if string(content) != expected {
					return &testError{msg: "File content mismatch, got: " + string(content) + ", expected: " + expected}
				}
				return nil
			},
		},
		{
			name:      "Invalid Processor",
			processor: "invalid-processor",
			response: &models.ModelResponse{
				Content:    "Test content",
				Model:      "test-model",
				PromptName: "test-prompt",
			},
			templateName:  "",
			expectedError: true,
			setup:         func() error { return nil },
			validate:      func() error { return nil },
		},
		{
			name:      "Slack Processor Without Token",
			processor: "slack-test-channel",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: true,
			setup: func() error {
				// Ensure SLACK_TOKEN is not set
				if err := os.Unsetenv("SLACK_TOKEN"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
		{
			name:      "Email Processor Without SMTP",
			processor: "email-test@example.com",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: true,
			setup: func() error {
				// Ensure SMTP_SERVER is not set
				if err := os.Unsetenv("SMTP_SERVER"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
		{
			name:      "Webhook Processor Without URL",
			processor: "webhook-monitoring",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: true,
			setup: func() error {
				// Ensure webhook URL env vars are not set
				if err := os.Unsetenv("WEBHOOK_URL"); err != nil {
					return err
				}
				if err := os.Unsetenv("WEBHOOK_URL_MONITORING"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
		{
			name:      "Console Processor",
			processor: "console",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: false,
			setup:         func() error { return nil },
			validate:      func() error { return nil },
		},
		{
			name:      "File Processor with alias",
			processor: "file",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Variables:   map[string]string{"key": "value"},
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "test_file_alias",
			expectedError: false,
			setup: func() error {
				// Set logs directory to temp dir
				if err := os.Setenv("LOGS_DIRECTORY", tempDir); err != nil {
					return err
				}

				// Register test templates
				manager := template.GetManager()
				err := manager.RegisterTemplate("test_file_alias_filename", filepath.Join(tempDir, "test-output-alias.txt"))
				if err != nil {
					return err
				}
				return manager.RegisterTemplate("test_file_alias_content", "Content: {{.Content}}")
			},
			validate: func() error {
				// Check if file was created
				content, err := os.ReadFile(filepath.Join(tempDir, "test-output-alias.txt"))
				if err != nil {
					return err
				}
				// Verify content
				expected := "Content: Test content"
				if string(content) != expected {
					return &testError{msg: "File content mismatch, got: " + string(content) + ", expected: " + expected}
				}
				return nil
			},
		},
		{
			name:      "Slack Processor With Fake Webhook URL",
			processor: "slack-test-channel",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: true,
			setup: func() error {
				// Set fake SLACK_WEBHOOK_URL - this will fail but tests error handling
				if err := os.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/test/test/test"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
		{
			name:      "Email Processor With SMTP",
			processor: "email-test@example.com",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: false,
			setup: func() error {
				// Set SMTP_SERVER
				if err := os.Setenv("SMTP_SERVER", "test-server"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
		{
			name:      "Webhook Processor With URL",
			processor: "webhook-monitoring",
			response: &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-execution",
			},
			templateName:  "",
			expectedError: false,
			setup: func() error {
				// Set webhook URL
				if err := os.Setenv("WEBHOOK_URL_MONITORING", "https://example.com/webhook"); err != nil {
					return err
				}
				return nil
			},
			validate: func() error { return nil },
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test
			err := tc.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Run the processor
			err = ProcessResponse(tc.processor, tc.response, tc.templateName)

			// Check error result
			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Additional validation if no error occurred and not expected
			if !tc.expectedError && err == nil {
				err = tc.validate()
				if err != nil {
					t.Errorf("Validation failed: %v", err)
				}
			}
		})
	}
}

// Custom test error type
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// TestInitTemplates tests the template initialization function
func TestInitTemplates(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create a test template file
	testTemplate := "test_template.tmpl"
	testContent := "Template content: {{.Content}}"
	if err := os.WriteFile(filepath.Join(tempDir, testTemplate), []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	// Initialize templates
	err = InitTemplates(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize templates: %v", err)
	}

	// Verify template was loaded
	manager := template.GetManager()
	tmpl, err := manager.GetTemplate("test_template")
	if err != nil {
		t.Fatalf("Failed to get template: %v", err)
	}

	// Verify template rendering
	sampleData := map[string]interface{}{"Content": "Sample content"}
	var renderedOutput strings.Builder
	if err := tmpl.Execute(&renderedOutput, sampleData); err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
	expectedOutput := "Template content: Sample content"
	if renderedOutput.String() != expectedOutput {
		t.Errorf("Rendered output mismatch. Expected: %q, Got: %q", expectedOutput, renderedOutput.String())
	}
}

func TestTemplateValidation(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "cronai-test-templates")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to clean up temp directory: %v", err)
		}
	}()

	// Create test template files
	validTemplate := filepath.Join(tempDir, "valid.tmpl")
	err = os.WriteFile(validTemplate, []byte("Test template: {{.Content}}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid template: %v", err)
	}

	invalidTemplate := filepath.Join(tempDir, "invalid.tmpl")
	err = os.WriteFile(invalidTemplate, []byte("Invalid template: {{.Content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid template: %v", err)
	}

	// Get template manager
	manager := template.GetManager()

	// Test validating a valid template
	err = manager.ValidateTemplate(validTemplate)
	if err != nil {
		t.Errorf("Failed to validate valid template: %v", err)
	}

	// Test validating an invalid template
	err = manager.ValidateTemplate(invalidTemplate)
	if err == nil {
		t.Errorf("Expected error validating invalid template but got nil")
	}

	// Test validating templates in directory
	results, err := manager.ValidateTemplatesInDir(tempDir)
	if err != nil {
		t.Errorf("Failed to validate templates in directory: %v", err)
	}

	// Check that validation results are as expected
	if results["valid.tmpl"] != nil {
		t.Errorf("Expected valid template to pass validation but got error: %v", results["valid.tmpl"])
	}
	if results["invalid.tmpl"] == nil {
		t.Errorf("Expected invalid template to fail validation but got nil error")
	}
}

func TestTemplateLoading(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := os.MkdirTemp("", "cronai-test-templates")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to clean up temp directory: %v", err)
		}
	}()

	// Create a test template file
	testTemplate := filepath.Join(tempDir, "test_template.tmpl")
	templateContent := "Template content: {{.Content}}"
	err = os.WriteFile(testTemplate, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Initialize templates with test directory
	err = InitTemplates(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize templates: %v", err)
	}

	// Get template manager
	manager := template.GetManager()

	// Check that test template was loaded
	if !manager.TemplateExists("test_template") {
		t.Errorf("Expected test_template to be loaded but it wasn't")
	}

	// Execute the loaded template
	data := template.Data{
		Content: "Test content",
	}
	result, err := manager.Execute("test_template", data)
	if err != nil {
		t.Errorf("Failed to execute loaded template: %v", err)
	}

	expected := "Template content: Test content"
	if result != expected {
		t.Errorf("Expected template result %q but got %q", expected, result)
	}
}
