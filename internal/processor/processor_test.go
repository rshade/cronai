package processor

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
)

func TestProcessResponse(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := filepath.Join(os.TempDir(), "cronai-test-"+time.Now().Format("20060102150405"))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to clean up temp directory: %v", err)
		}
	}() // Clean up after test

	// Create a test response
	response := &models.ModelResponse{
		Content:     "This is a test response",
		Model:       "test-model",
		PromptName:  "test-prompt",
		ExecutionID: "test-execution-1234",
		Variables: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Set environment variables for testing
	if err := os.Setenv("SLACK_TOKEN", "test-slack-token"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("SMTP_SERVER", "test-smtp-server"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("WEBHOOK_URL", "https://example.com/webhook"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := os.Setenv("WEBHOOK_URL_MONITORING", "https://example.com/monitoring"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Initialize templates with test directory
	err = InitTemplates("")
	if err != nil {
		t.Fatalf("Failed to initialize templates: %v", err)
	}

	// Test cases for different processors
	testCases := []struct {
		name          string
		processorName string
		templateName  string
		shouldError   bool
	}{
		{
			name:          "Slack processor",
			processorName: "slack-test-channel",
			templateName:  "",
			shouldError:   false,
		},
		{
			name:          "Slack processor with custom template",
			processorName: "slack-test-channel",
			templateName:  "custom_slack",
			shouldError:   false,
		},
		{
			name:          "Email processor",
			processorName: "email-test@example.com",
			templateName:  "",
			shouldError:   false,
		},
		{
			name:          "Webhook processor",
			processorName: "webhook-test",
			templateName:  "",
			shouldError:   false,
		},
		{
			name:          "Monitoring webhook processor",
			processorName: "webhook-monitoring",
			templateName:  "",
			shouldError:   false,
		},
		{
			name:          "File processor",
			processorName: "file",
			templateName:  "",
			shouldError:   false,
		},
		{
			name:          "Unsupported processor",
			processorName: "unsupported-processor",
			templateName:  "",
			shouldError:   true,
		},
	}

	// Register a custom test template for testing
	manager := template.GetManager()
	err = manager.RegisterTemplate("custom_slack", `{"blocks":[{"type":"section","text":{"type":"mrkdwn","text":"Custom template: {{.Content}}"}}]}`)
	if err != nil {
		t.Fatalf("Failed to register custom template: %v", err)
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ProcessResponse(tc.processorName, response, tc.templateName)
			if tc.shouldError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTemplateValidation(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir := filepath.Join(os.TempDir(), "cronai-test-templates-"+time.Now().Format("20060102150405"))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to clean up temp directory: %v", err)
		}
	}() // Clean up after test

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
	tempDir := filepath.Join(os.TempDir(), "cronai-test-templates-"+time.Now().Format("20060102150405"))
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: Failed to clean up temp directory: %v", err)
		}
	}() // Clean up after test

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
	data := template.TemplateData{
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
