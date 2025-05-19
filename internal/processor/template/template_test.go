package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestManager tests the template manager functionality
func TestManager(t *testing.T) {
	// Test singleton instance
	t.Run("singleton instance", func(t *testing.T) {
		m1 := GetManager()
		m2 := GetManager()
		assert.Same(t, m1, m2, "GetInstance should return the same instance")
	})

	// Test basic operations
	t.Run("basic operations", func(t *testing.T) {
		manager := &Manager{
			templates:   make(map[string]*template.Template),
			inheritance: make(map[string]*templateInheritance),
		}

		// Register a template
		templateContent := "Hello {{.Content}}!"
		err := manager.Register("test", templateContent)
		assert.NoError(t, err)
		assert.True(t, manager.Has("test"))

		// Execute the template
		data := TemplateData{
			Content: "World",
		}
		result, err := manager.Execute("test", data)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World!", result)

		// Try non-existent template
		_, err = manager.Execute("non_existent", data)
		assert.Error(t, err)
		assert.False(t, manager.Has("non_existent"))
	})
}

// TestValidate tests the Validate function
func TestValidate(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	tests := []struct {
		name         string
		templateName string
		content      string
		expectErr    bool
		errMsg       string
	}{
		{
			name:         "valid template",
			templateName: "valid",
			content:      "Hello {{.Content}}!",
			expectErr:    false,
		},
		{
			name:         "invalid template syntax",
			templateName: "invalid",
			content:      "Hello {{.Content}",
			expectErr:    true,
			errMsg:       "unclosed action",
		},
		{
			name:         "empty template name",
			templateName: "",
			content:      "Hello",
			expectErr:    true,
			errMsg:       "template name cannot be empty",
		},
		{
			name:         "undefined variable",
			templateName: "undefined",
			content:      "Hello {{.UndefinedField}}!",
			expectErr:    false, // Template is valid, error occurs during execution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Validate(tt.templateName, tt.content)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLoadTemplatesFromDir tests the LoadTemplatesFromDir function
func TestLoadTemplatesFromDir(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	templatesDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templatesDir, 0755))

	// Create test template files
	templates := []struct {
		filename string
		content  string
	}{
		{
			filename: "email.tmpl",
			content:  "Subject: {{.Subject}}\n\n{{.Content}}",
		},
		{
			filename: "slack.tmpl",
			content:  "*{{.Title}}*\n\n{{.Content}}",
		},
		{
			filename: "invalid.tmpl",
			content:  "{{.Title}", // Invalid template
		},
		{
			filename: "not_template.txt",
			content:  "This is not a template file",
		},
	}

	for _, tmpl := range templates {
		path := filepath.Join(templatesDir, tmpl.filename)
		require.NoError(t, os.WriteFile(path, []byte(tmpl.content), 0644))
	}

	// Create subdirectory
	subDir := filepath.Join(templatesDir, "sub")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	subTemplate := filepath.Join(subDir, "sub.tmpl")
	require.NoError(t, os.WriteFile(subTemplate, []byte("Sub: {{.Content}}"), 0644))

	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	// Load templates from directory
	err := manager.LoadTemplatesFromDir(templatesDir)
	assert.NoError(t, err)

	// Check loaded templates
	assert.True(t, manager.Has("email"))
	assert.True(t, manager.Has("slack"))
	assert.False(t, manager.Has("invalid"))      // Should skip invalid templates
	assert.False(t, manager.Has("not_template")) // Should skip non-.tmpl files
	assert.True(t, manager.Has("sub"))           // Should load from subdirectory

	// Test non-existent directory
	err = manager.LoadTemplatesFromDir(filepath.Join(tempDir, "non_existent"))
	assert.Error(t, err)
}

// TestLoadLibraryTemplates tests the LoadLibraryTemplates function
func TestLoadLibraryTemplates(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	templatesDir := filepath.Join(tempDir, "templates")
	libDir := filepath.Join(templatesDir, "library")
	require.NoError(t, os.MkdirAll(libDir, 0755))

	// Create library template files
	templates := []struct {
		filename string
		content  string
	}{
		{
			filename: "header.tmpl",
			content:  "<header>{{.Title}}</header>",
		},
		{
			filename: "footer.tmpl",
			content:  "<footer>{{.Copyright}}</footer>",
		},
		{
			filename: "base.tmpl",
			content:  "{{.Header}}\n{{.Content}}\n{{.Footer}}",
		},
	}

	for _, tmpl := range templates {
		path := filepath.Join(libDir, tmpl.filename)
		require.NoError(t, os.WriteFile(path, []byte(tmpl.content), 0644))
	}

	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	// Load library templates
	err := manager.LoadLibraryTemplates(templatesDir)
	assert.NoError(t, err)

	// Check loaded templates (library templates have "library." prefix)
	assert.True(t, manager.Has("library.header"))
	assert.True(t, manager.Has("library.footer"))
	assert.True(t, manager.Has("library.base"))

	// Test with no library directory
	noLibDir := filepath.Join(tempDir, "no_library")
	require.NoError(t, os.MkdirAll(noLibDir, 0755))
	err = manager.LoadLibraryTemplates(noLibDir)
	assert.NoError(t, err) // Should succeed even if library dir doesn't exist
}

// TestValidateTemplate tests the ValidateTemplate function
func TestValidateTemplate(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	tests := []struct {
		name      string
		content   string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid template",
			content:   "Hello {{.Name}}!",
			expectErr: false,
		},
		{
			name:      "invalid syntax",
			content:   "Hello {{.Name}",
			expectErr: true,
			errMsg:    "unclosed action",
		},
		{
			name:      "complex template",
			content:   "{{if .Show}}{{.Content}}{{else}}Hidden{{end}}",
			expectErr: false,
		},
		{
			name:      "with range",
			content:   "{{range .Items}}{{.}}{{end}}",
			expectErr: false,
		},
		{
			name:      "undefined action",
			content:   "{{invalid_action}}",
			expectErr: true,
			errMsg:    "function \"invalid_action\" not defined",
		},
		{
			name:      "empty template",
			content:   "",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateTemplate(tt.content)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestProcessInheritance tests the ProcessInheritance function
func TestProcessInheritance(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedBase   string
		expectedBlocks map[string]string
		expectErr      bool
	}{
		{
			name: "simple inheritance",
			content: `{{extends "base"}}
{{block "header"}}
Custom Header
{{endblock}}`,
			expectedBase: "base",
			expectedBlocks: map[string]string{
				"header": "\nCustom Header\n",
			},
		},
		{
			name: "multiple blocks",
			content: `{{extends "layout"}}
{{block "header"}}
Header Content
{{endblock}}
{{block "footer"}}
Footer Content
{{endblock}}`,
			expectedBase: "layout",
			expectedBlocks: map[string]string{
				"header": "\nHeader Content\n",
				"footer": "\nFooter Content\n",
			},
		},
		{
			name: "no inheritance",
			content: `Just regular content
without any extends directive`,
			expectedBase:   "",
			expectedBlocks: map[string]string{},
		},
		{
			name: "empty extends",
			content: `{{extends ""}}
{{block "content"}}
Test
{{endblock}}`,
			expectErr: true,
		},
		{
			name: "block without extends",
			content: `{{block "content"}}
Test
{{endblock}}`,
			expectErr: true,
		},
		{
			name: "unclosed block",
			content: `{{extends "base"}}
{{block "content"}}
No endblock`,
			expectErr: true,
		},
		{
			name: "duplicate blocks",
			content: `{{extends "base"}}
{{block "content"}}
First
{{endblock}}
{{block "content"}}
Second
{{endblock}}`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessInheritance(tt.content)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBase, result.Parent)
				assert.Equal(t, tt.expectedBlocks, result.Blocks)
			}
		})
	}
}

// TestTemplateInheritance tests template inheritance functionality
func TestTemplateInheritance(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	// Register base template
	baseTemplate := `<html>
<head><title>{{.Title}}</title></head>
<body>
{{block "header"}}Default Header{{endblock}}
{{block "content"}}Default Content{{endblock}}
{{block "footer"}}Default Footer{{endblock}}
</body>
</html>`
	err := manager.Register("base", baseTemplate)
	require.NoError(t, err)

	// Register child template
	childTemplate := `{{extends "base"}}
{{block "header"}}
<h1>Custom Header</h1>
{{endblock}}
{{block "content"}}
<p>Custom Content: {{.Message}}</p>
{{endblock}}`
	err = manager.Register("page", childTemplate)
	require.NoError(t, err)

	// Execute child template
	data := TemplateData{
		Variables: map[string]string{
			"Title":   "Test Page",
			"Message": "Hello World",
		},
	}
	result, err := manager.Execute("page", data)
	require.NoError(t, err)

	// Check result contains expected content
	assert.Contains(t, result, "<title>Test Page</title>")
	assert.Contains(t, result, "<h1>Custom Header</h1>")
	assert.Contains(t, result, "<p>Custom Content: Hello World</p>")
	assert.Contains(t, result, "Default Footer") // Should use default footer
}

// TestTemplateFunctions tests custom template functions
func TestTemplateFunctions(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	tests := []struct {
		name     string
		template string
		data     TemplateData
		expected string
	}{
		{
			name:     "upper function",
			template: "{{upper .Content}}",
			data:     TemplateData{Content: "hello"},
			expected: "HELLO",
		},
		{
			name:     "lower function",
			template: "{{lower .Content}}",
			data:     TemplateData{Content: "HELLO"},
			expected: "hello",
		},
		{
			name:     "title function",
			template: "{{title .Content}}",
			data:     TemplateData{Content: "hello world"},
			expected: "Hello World",
		},
		{
			name:     "trim function",
			template: "{{trim .Content}}",
			data:     TemplateData{Content: "  hello  "},
			expected: "hello",
		},
		{
			name:     "join function",
			template: `{{join .Variables.items ","}}`,
			data: TemplateData{
				Variables: map[string]string{
					"items": "a b c",
				},
			},
			expected: "a,b,c",
		},
		{
			name:     "formatDate function",
			template: `{{formatDate .Timestamp "2006-01-02"}}`,
			data: TemplateData{
				Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			expected: "2024-01-15",
		},
		{
			name:     "json function",
			template: `{{json .Variables}}`,
			data: TemplateData{
				Variables: map[string]string{
					"key": "value",
				},
			},
			expected: `{"key":"value"}`,
		},
		{
			name:     "default function",
			template: `{{default .Content "No content"}}`,
			data:     TemplateData{Content: ""},
			expected: "No content",
		},
		{
			name:     "contains function",
			template: `{{if contains .Content "world"}}Found{{else}}Not found{{end}}`,
			data:     TemplateData{Content: "hello world"},
			expected: "Found",
		},
		{
			name:     "replace function",
			template: `{{replace .Content "world" "universe"}}`,
			data:     TemplateData{Content: "hello world"},
			expected: "hello universe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Register(tt.name, tt.template)
			require.NoError(t, err)

			result, err := manager.Execute(tt.name, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSafeExecute tests the SafeExecute function
func TestSafeExecute(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	// Register a template
	err := manager.Register("test", "Hello {{.Name}}!")
	require.NoError(t, err)

	// Test successful execution
	data := TemplateData{
		Variables: map[string]string{
			"Name": "World",
		},
	}
	result := manager.SafeExecute("test", data)
	assert.Equal(t, "Hello World!", result)

	// Test non-existent template (should use fallback)
	result = manager.SafeExecute("non_existent", data)
	assert.Contains(t, result, "Response from") // Default fallback template

	// Test execution error with invalid data
	invalidData := TemplateData{} // Missing Name variable
	result = manager.SafeExecute("test", invalidData)
	assert.Equal(t, "Hello <no value>!", result) // Template handles missing values
}

// TestRegisterOrPanic tests the registerOrPanic function
func TestRegisterOrPanic(t *testing.T) {
	// This test needs to be careful about panics
	t.Run("successful registration", func(t *testing.T) {
		manager := &Manager{
			templates:   make(map[string]*template.Template),
			inheritance: make(map[string]*templateInheritance),
		}

		// Should not panic
		assert.NotPanics(t, func() {
			manager.registerOrPanic("test", "Hello {{.Content}}!")
		})
		assert.True(t, manager.Has("test"))
	})

	t.Run("registration with invalid template", func(t *testing.T) {
		manager := &Manager{
			templates:   make(map[string]*template.Template),
			inheritance: make(map[string]*templateInheritance),
		}

		// Should panic with invalid template
		assert.Panics(t, func() {
			manager.registerOrPanic("invalid", "{{.Content")
		})
	})
}

// TestConcurrency tests concurrent access to the template manager
func TestConcurrency(t *testing.T) {
	manager := &Manager{
		templates:   make(map[string]*template.Template),
		inheritance: make(map[string]*templateInheritance),
	}

	// Register initial templates
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("template%d", i)
		content := fmt.Sprintf("Template %d: {{.Content}}", i)
		err := manager.Register(name, content)
		require.NoError(t, err)
	}

	// Run concurrent operations
	done := make(chan bool)
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Perform various operations
			switch id % 4 {
			case 0:
				// Register new template
				name := fmt.Sprintf("concurrent%d", id)
				manager.Register(name, "Concurrent {{.ID}}")
			case 1:
				// Execute template
				data := TemplateData{Content: fmt.Sprintf("Content %d", id)}
				manager.Execute(fmt.Sprintf("template%d", id%10), data)
			case 2:
				// Check existence
				manager.Has(fmt.Sprintf("template%d", id%10))
			case 3:
				// Safe execute
				data := TemplateData{Content: fmt.Sprintf("Safe %d", id)}
				manager.SafeExecute(fmt.Sprintf("template%d", id%10), data)
			}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify state is consistent
	for i := 0; i < 10; i++ {
		assert.True(t, manager.Has(fmt.Sprintf("template%d", i)))
	}
}
