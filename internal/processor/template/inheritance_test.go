package template

import (
	"testing"
	"time"
)

func TestTemplateInheritanceParsing(t *testing.T) {
	manager := GetManager()

	// Test a template with inheritance
	childTemplate := `{{extends "parent_template"}}
{{block "header" .}}Child Header{{end}}
{{block "footer" .}}Child Footer{{end}}
`

	inheritance, _, err := manager.ParseInheritance("child_template", childTemplate)
	if err != nil {
		t.Errorf("ParseInheritance failed: %v", err)
	}

	if inheritance == nil {
		t.Fatal("Expected inheritance struct, got nil")
	}

	if inheritance.Parent != "parent_template" {
		t.Errorf("Expected parent_template, got %s", inheritance.Parent)
	}

	if len(inheritance.Blocks) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(inheritance.Blocks))
	}

	if _, ok := inheritance.Blocks["header"]; !ok {
		t.Error("Expected header block, not found")
	}

	if _, ok := inheritance.Blocks["footer"]; !ok {
		t.Error("Expected footer block, not found")
	}

	// Test a template without inheritance
	regularTemplate := `Regular template with no inheritance`

	inheritance2, processedContent2, err := manager.ParseInheritance("regular_template", regularTemplate)
	if err != nil {
		t.Errorf("ParseInheritance failed: %v", err)
	}

	if inheritance2 != nil {
		t.Error("Expected nil inheritance for regular template")
	}

	if processedContent2 != regularTemplate {
		t.Errorf("Expected unchanged content, got %s", processedContent2)
	}
}

func TestTemplateInheritanceExecution(t *testing.T) {
	manager := GetManager()

	// Register a parent template
	parentTemplate := `Parent start
{{template "header" .}}
Middle content
{{template "footer" .}}
Parent end

{{define "header"}}Parent Header{{end}}
{{define "footer"}}Parent Footer{{end}}`

	err := manager.RegisterTemplate("parent_test", parentTemplate)
	if err != nil {
		t.Fatalf("Failed to register parent template: %v", err)
	}

	// Register a child template that extends the parent
	childTemplate := `{{extends "parent_test"}}
{{define "header"}}Child Header Override{{end}}`

	err = manager.RegisterTemplate("child_test", childTemplate)
	if err != nil {
		t.Fatalf("Failed to register child template: %v", err)
	}

	// Create template data
	data := TemplateData{
		Model:     "TestModel",
		Content:   "Test content",
		Timestamp: time.Now(),
	}

	// Execute the child template
	result, err := manager.Execute("child_test", data)
	if err != nil {
		t.Errorf("Failed to execute template with inheritance: %v", err)
	}

	// Check that the child's header override was used
	if !contains(result, "Child Header Override") {
		t.Error("Child header override not found in result")
	}

	// Check that the parent's footer was used (not overridden)
	if !contains(result, "Parent Footer") {
		t.Error("Parent footer not found in result")
	}

	// Register another child with multiple block overrides
	childTemplate2 := `{{extends "parent_test"}}
{{define "header"}}Child 2 Header{{end}}
{{define "footer"}}Child 2 Footer{{end}}`

	err = manager.RegisterTemplate("child_test2", childTemplate2)
	if err != nil {
		t.Fatalf("Failed to register second child template: %v", err)
	}

	// Execute the second child template
	result2, err := manager.Execute("child_test2", data)
	if err != nil {
		t.Errorf("Failed to execute second template with inheritance: %v", err)
	}

	// Check that both child blocks were used
	if !contains(result2, "Child 2 Header") {
		t.Error("Child 2 header not found in result")
	}

	if !contains(result2, "Child 2 Footer") {
		t.Error("Child 2 footer not found in result")
	}
}

func TestNestedTemplateInheritance(t *testing.T) {
	manager := GetManager()

	// Register a base template - use define to make this more reliable
	baseTemplate := `{{define "base_template"}}Base start
{{template "section" .}}
Base end
{{end}}

{{define "section"}}Base Section{{end}}`

	err := manager.RegisterTemplate("base_test", baseTemplate)
	if err != nil {
		t.Fatalf("Failed to register base template: %v", err)
	}

	// Register a middle template extending the base
	middleTemplate := `{{extends "base_test"}}
{{define "middle_template"}}{{template "base_template" .}}{{end}}

{{define "section"}}
Middle Section Start
{{template "subsection" .}}
Middle Section End
{{end}}

{{define "subsection"}}Middle Subsection{{end}}`

	err = manager.RegisterTemplate("middle_test", middleTemplate)
	if err != nil {
		t.Fatalf("Failed to register middle template: %v", err)
	}

	// Register a leaf template extending the middle one
	leafTemplate := `{{extends "middle_test"}}
{{define "leaf_template"}}{{template "middle_template" .}}{{end}}

{{define "subsection"}}Leaf Subsection Override{{end}}`

	err = manager.RegisterTemplate("leaf_test", leafTemplate)
	if err != nil {
		t.Fatalf("Failed to register leaf template: %v", err)
	}

	// Create template data
	data := TemplateData{
		Model:     "TestModel",
		Content:   "Test content",
		Timestamp: time.Now(),
	}

	// Execute the leaf template
	// First make sure the nested templates are processed correctly
	if err := manager.RegisterTemplateWithIncludes("middle_test", ""); err != nil {
		t.Fatalf("Failed to process middle template includes: %v", err)
	}

	if err := manager.RegisterTemplateWithIncludes("base_test", ""); err != nil {
		t.Fatalf("Failed to process base template includes: %v", err)
	}

	// Now execute the leaf template with the leaf_template definition
	result, err := manager.ExecuteTemplate("leaf_test", "leaf_template", data)
	if err != nil {
		t.Errorf("Failed to execute template with nested inheritance: %v", err)
	}

	// Check the result contains content from all levels
	if !contains(result, "Base start") {
		t.Error("Base template start content not found")
	}

	if !contains(result, "Middle Section Start") {
		t.Error("Middle section start not found")
	}

	if !contains(result, "Leaf Subsection Override") {
		t.Error("Leaf subsection override not found")
	}

	if !contains(result, "Middle Section End") {
		t.Error("Middle section end not found")
	}

	if !contains(result, "Base end") {
		t.Error("Base template end content not found")
	}
}

func TestTemplateWithIncludes(t *testing.T) {
	manager := GetManager()

	// Register component templates with define blocks
	err := manager.RegisterTemplate("header_component", "{{define \"header_component\"}}HEADER: {{.Model}}{{end}}")
	if err != nil {
		t.Fatalf("Failed to register header component: %v", err)
	}

	err = manager.RegisterTemplate("footer_component", "{{define \"footer_component\"}}FOOTER: {{.Timestamp.Format \"2006-01-02\"}}{{end}}")
	if err != nil {
		t.Fatalf("Failed to register footer component: %v", err)
	}

	// Register template with includes
	templateWithIncludes := `{{define "template_with_includes"}}
{{template "header_component" .}}
Content: {{.Content}}
{{template "footer_component" .}}
{{end}}`

	// Use RegisterTemplateWithIncludes instead of RegisterTemplate
	err = manager.RegisterTemplateWithIncludes("template_with_includes", templateWithIncludes)
	if err != nil {
		t.Fatalf("Failed to register template with includes: %v", err)
	}

	// Create template data
	now := time.Now()
	data := TemplateData{
		Model:     "TestModel",
		Content:   "Test content",
		Timestamp: now,
	}

	// Execute the template using the named template
	result, err := manager.ExecuteTemplate("template_with_includes", "template_with_includes", data)
	if err != nil {
		t.Errorf("Failed to execute template with includes: %v", err)
	}

	// Check that included templates were rendered
	if !contains(result, "HEADER: TestModel") {
		t.Error("Header component not properly included")
	}

	if !contains(result, "Content: Test content") {
		t.Error("Main content not properly included")
	}

	expectedFooter := "FOOTER: " + now.Format("2006-01-02")
	if !contains(result, expectedFooter) {
		t.Error("Footer component not properly included")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return s != "" && s != substr && len(s) >= len(substr) && s != substr && substring(s, substr) >= 0
}

func substring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
