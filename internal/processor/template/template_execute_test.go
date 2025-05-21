package template

import (
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteTemplate(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*Inheritance),
	}

	// Register a template with multiple template definitions
	templateContent := `
	{{define "header"}}Header: {{.Title}}{{end}}
	{{define "footer"}}Footer: {{.Copyright}}{{end}}
	Main template with {{.Content}}
	`
	err := manager.RegisterTemplate("multi", templateContent)
	require.NoError(t, err)

	// Create test data
	data := Data{
		Content: "Test Content",
		Variables: map[string]string{
			"Title":     "Page Title",
			"Copyright": "© 2024",
		},
	}

	// Test executing the main template
	result, err := manager.Execute("multi", data)
	require.NoError(t, err)
	assert.Contains(t, result, "Main template with Test Content")

	// Test executing a named template within the template
	headerResult, err := manager.ExecuteTemplate("multi", "header", data)
	require.NoError(t, err)
	assert.Equal(t, "Header: Page Title", headerResult)

	footerResult, err := manager.ExecuteTemplate("multi", "footer", data)
	require.NoError(t, err)
	assert.Equal(t, "Footer: © 2024", footerResult)

	// Test with non-existent template
	_, err = manager.ExecuteTemplate("non-existent", "header", data)
	assert.Error(t, err)

	// Test with non-existent named template
	_, err = manager.ExecuteTemplate("multi", "non-existent", data)
	assert.Error(t, err)
}

func TestExecuteTemplateWithInheritance(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*Inheritance),
	}

	// Register parent template
	parentTemplate := `
	{{define "header"}}Parent Header: {{.Title}}{{end}}
	{{define "content"}}Parent Content: {{.Content}}{{end}}
	{{define "footer"}}Parent Footer: {{.Copyright}}{{end}}
	`
	err := manager.RegisterTemplate("parent", parentTemplate)
	require.NoError(t, err)

	// Register child template that extends parent
	childTemplate := `
	{{extends "parent"}}
	{{define "header"}}Child Header: {{.Title}}{{end}}
	`
	err = manager.RegisterTemplate("child", childTemplate)
	require.NoError(t, err)

	// Create test data
	data := Data{
		Content: "Test Content",
		Variables: map[string]string{
			"Title":     "Page Title",
			"Copyright": "© 2024",
		},
	}

	// Execute the header template from the child template
	headerResult, err := manager.ExecuteTemplate("child", "header", data)
	require.NoError(t, err)
	assert.Equal(t, "Child Header: Page Title", headerResult)

	// Execute the content template (should use parent's)
	contentResult, err := manager.ExecuteTemplate("child", "content", data)
	require.NoError(t, err)
	assert.Equal(t, "Parent Content: Test Content", contentResult)

	// Execute the footer template (should use parent's)
	footerResult, err := manager.ExecuteTemplate("child", "footer", data)
	require.NoError(t, err)
	assert.Equal(t, "Parent Footer: © 2024", footerResult)
}

func TestExecuteTemplateWithComplexData(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*Inheritance),
	}

	// Register template with complex data access
	templateContent := `
	{{define "metadata"}}
	Model: {{.Model}}
	Timestamp: {{formatDate .Timestamp "2006-01-02 15:04:05"}}
	Execution ID: {{.ExecutionID}}
	{{end}}
	`
	err := manager.RegisterTemplate("complex", templateContent)
	require.NoError(t, err)

	// Create complex test data
	testTime := time.Date(2024, 3, 15, 12, 30, 45, 0, time.UTC)
	data := Data{
		Content:     "Complex test content",
		Model:       "test-model",
		Timestamp:   testTime,
		PromptName:  "test-prompt",
		ExecutionID: "exec-123456",
		Variables: map[string]string{
			"key1": "value1",
		},
		Metadata: map[string]string{
			"meta1": "metadata1",
		},
	}

	// Execute the metadata template
	result, err := manager.ExecuteTemplate("complex", "metadata", data)
	require.NoError(t, err)
	assert.Contains(t, result, "Model: test-model")
	assert.Contains(t, result, "Timestamp: 2024-03-15 12:30:45")
	assert.Contains(t, result, "Execution ID: exec-123456")
}

func TestTemplateValidation(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*Inheritance),
	}

	// Test ValidateTemplatesInDir with a temp directory
	tempDir := t.TempDir()
	results, err := manager.ValidateTemplatesInDir(tempDir)
	assert.NoError(t, err)
	assert.Empty(t, results)

	// Test non-existent directory
	_, err = manager.ValidateTemplatesInDir("/non/existent/directory")
	assert.Error(t, err)
}
