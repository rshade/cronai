package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateInheritanceInPrompts(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()

	// Create cron_prompts directory structure
	cronPromptsDir := filepath.Join(tempDir, "cron_prompts")
	err := os.MkdirAll(cronPromptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cron_prompts directory: %v", err)
	}

	// Change to temp directory for the test to make relative paths work
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore original directory: %v", err)
		}
	}()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create a base template file
	baseTemplatePath := filepath.Join(cronPromptsDir, "base_template.md")
	baseContent := `---
name: Base Template
description: Base template with placeholders
version: 1.0
---

# Base Template

## Introduction
{{template "introduction" .}}

## Main Content
{{template "content" .}}

## Conclusion
{{template "conclusion" .}}

{{define "introduction"}}
Default introduction content.
{{end}}

{{define "content"}}
Default main content.
{{end}}

{{define "conclusion"}}
Default conclusion.
{{end}}`

	err = os.WriteFile(baseTemplatePath, []byte(baseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create base template file: %v", err)
	}

	// Create a derived template that extends the base
	derivedTemplatePath := filepath.Join(cronPromptsDir, "derived_template.md")
	derivedContent := `---
name: Derived Template
description: Template that extends the base
version: 1.0
extends: base_template
---

{{define "introduction"}}
This is a custom introduction that overrides the base template.
Project: {{.Variables.projectName}}
{{end}}

{{define "content"}}
This is the main content of the derived template.
Details: {{.Variables.details}}
{{end}}`

	err = os.WriteFile(derivedTemplatePath, []byte(derivedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create derived template file: %v", err)
	}

	// Process the derived template
	content, err := os.ReadFile(derivedTemplatePath)
	if err != nil {
		t.Fatalf("Failed to read derived template: %v", err)
	}

	// Extract metadata and content
	metadata, _, err := ExtractMetadata(string(content), derivedTemplatePath)
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	// Verify the metadata
	if metadata.Name != "Derived Template" {
		t.Errorf("Expected template name 'Derived Template', got '%s'", metadata.Name)
	}

	// Verify the extends attribute in metadata
	if metadata.Extends != "base_template" {
		t.Errorf("Expected extends attribute 'base_template', got '%s'", metadata.Extends)
	}

	// Apply variables to the template
	variables := map[string]string{
		"projectName": "Test Project",
		"details":     "Detailed information for the test",
	}

	// Process template with inheritance directives
	_, processed, err := testInheritanceHelper(derivedTemplatePath, string(content), variables)
	if err != nil {
		t.Fatalf("ProcessPromptWithInheritance failed: %v", err)
	}

	// Check that the processed content contains overrides from the derived template
	if !contains(processed, "This is a custom introduction") {
		t.Error("Custom introduction override not found in processed content")
	}

	if !contains(processed, "Project: Test Project") {
		t.Error("Variable replacement not found in processed content")
	}

	if !contains(processed, "This is the main content of the derived template") {
		t.Error("Custom content override not found in processed content")
	}

	if !contains(processed, "Details: Detailed information for the test") {
		t.Error("Variable replacement in content block not found")
	}

	// Check that the base template's default conclusion was used (not overridden)
	if !contains(processed, "Default conclusion") {
		t.Error("Default conclusion from base template not found")
	}
}

func TestPromptComposition(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()

	// Create templates library directory structure
	templatesLibDir := filepath.Join(tempDir, "templates", "library")
	err := os.MkdirAll(templatesLibDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates library directory: %v", err)
	}

	// Change to temp directory for the test to make relative paths work
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore original directory: %v", err)
		}
	}()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create component files
	headerPath := filepath.Join(templatesLibDir, "header.md")
	headerContent := `---
name: Header Component
description: Reusable header component
---
# {{title}}
Date: {{date}}
`

	err = os.WriteFile(headerPath, []byte(headerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create header component: %v", err)
	}

	footerPath := filepath.Join(templatesLibDir, "footer.md")
	footerContent := `---
name: Footer Component
description: Reusable footer component
---
---
Generated by: {{generator}}
Version: {{version}}
`

	err = os.WriteFile(footerPath, []byte(footerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create footer component: %v", err)
	}

	// Create cron_prompts directory for composed template
	cronPromptsDir := filepath.Join(tempDir, "cron_prompts")
	err = os.MkdirAll(cronPromptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create cron_prompts directory: %v", err)
	}

	// Create a template that uses components through inclusion
	composedTemplatePath := filepath.Join(cronPromptsDir, "composed_template.md")
	composedContent := `---
name: Composed Template
description: Template using component inclusion
---
{{include "header"}}

## Main Content
This is the main content of the composed template.
{{body}}

{{include "footer"}}
`

	err = os.WriteFile(composedTemplatePath, []byte(composedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create composed template: %v", err)
	}

	// Load and process the composed template
	templateContent, err := os.ReadFile(composedTemplatePath)
	if err != nil {
		t.Fatalf("Failed to read composed template: %v", err)
	}

	// Extract metadata
	_, extractedContent, err := ExtractMetadata(string(templateContent), composedTemplatePath)
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	// Process includes
	processedWithIncludes, err := ProcessIncludes(extractedContent)
	if err != nil {
		t.Fatalf("ProcessIncludes failed: %v", err)
	}

	// Apply variables
	variables := map[string]string{
		"title":     "Composed Document",
		"date":      "2025-05-13",
		"body":      "This is the variable content inserted in the body.",
		"generator": "CronAI Test",
		"version":   "2.0",
	}

	// Process the template with variables
	finalContent := ApplyVariables(processedWithIncludes, variables)

	// Debug: print the final content
	if testing.Verbose() {
		t.Logf("Processed with includes: %s", processedWithIncludes)
		t.Logf("Final content: %s", finalContent)
	}

	// Verify that included components and variables were properly processed
	if !contains(finalContent, "# Composed Document") {
		t.Error("Header component title not properly included or variables not applied")
	}

	if !contains(finalContent, "Date: 2025-05-13") {
		t.Error("Header component date not properly included or variables not applied")
	}

	if !contains(finalContent, "This is the variable content inserted in the body") {
		t.Error("Body variable not properly applied")
	}

	if !contains(finalContent, "Generated by: CronAI Test") {
		t.Error("Footer component generator not properly included or variables not applied")
	}

	if !contains(finalContent, "Version: 2.0") {
		t.Error("Footer component version not properly included or variables not applied")
	}
}

// Helper to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// testInheritanceHelper is a wrapper for ProcessPromptWithInheritance for testing purposes
func testInheritanceHelper(path, content string, variables map[string]string) (map[string]string, string, error) {
	// Use the actual ProcessPromptWithInheritance function
	return ProcessPromptWithInheritance(path, content, variables)
}
