package template

import (
	"testing"
	"time"
)

func TestTemplateRegistration(t *testing.T) {
	manager := GetManager()

	// Register a test template
	err := manager.RegisterTemplate("test_template", "Hello, {{.Model}}!")
	if err != nil {
		t.Errorf("Failed to register template: %v", err)
	}

	// Try to register an invalid template
	err = manager.RegisterTemplate("invalid_template", "Hello, {{.Model")
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}

func TestTemplateExecution(t *testing.T) {
	manager := GetManager()

	// Register a test template
	templateContent := "Model: {{.Model}}, Content: {{.Content}}"
	err := manager.RegisterTemplate("exec_test", templateContent)
	if err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	// Create template data
	data := TemplateData{
		Model:     "TestModel",
		Content:   "Test content",
		Timestamp: time.Now(),
	}

	// Execute the template
	result, err := manager.Execute("exec_test", data)
	if err != nil {
		t.Errorf("Failed to execute template: %v", err)
	}

	expected := "Model: TestModel, Content: Test content"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test non-existent template
	_, err = manager.Execute("non_existent", data)
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestSafeExecute(t *testing.T) {
	manager := GetManager()

	// Create template data
	data := TemplateData{
		Model:     "TestModel",
		Content:   "Fallback content",
		Timestamp: time.Now(),
	}

	// Test fallback to raw content
	result := manager.SafeExecute("non_existent", data)
	if result != "Fallback content" {
		t.Errorf("Expected fallback to content, got %q", result)
	}

	// Register a default template
	err := manager.RegisterTemplate("default_test", "Default: {{.Model}}")
	if err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	// Test fallback to default template
	result = manager.SafeExecute("test_specific", data)
	if result != "Default: TestModel" {
		t.Errorf("Expected fallback to default template, got %q", result)
	}
}
