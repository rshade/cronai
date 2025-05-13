package template

import (
	"testing"
)

func TestConditionalTemplates(t *testing.T) {
	manager := GetManager()

	// Test cases for different types of conditional logic
	tests := []struct {
		name     string
		template string
		data     TemplateData
		expected string
	}{
		{
			name:     "basic if condition - true case",
			template: "{{if eq .Model \"TestModel\"}}Correct model{{else}}Wrong model{{end}}",
			data: TemplateData{
				Model: "TestModel",
			},
			expected: "Correct model",
		},
		{
			name:     "basic if condition - false case",
			template: "{{if eq .Model \"OtherModel\"}}Wrong model{{else}}Not the expected model{{end}}",
			data: TemplateData{
				Model: "TestModel",
			},
			expected: "Not the expected model",
		},
		{
			name:     "variable existence check",
			template: "{{if hasVar .Variables \"feature\"}}Feature exists{{else}}Feature not found{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"feature": "enabled",
				},
			},
			expected: "Feature exists",
		},
		{
			name:     "variable existence check - missing variable",
			template: "{{if hasVar .Variables \"missing\"}}Variable exists{{else}}Variable not found{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"feature": "enabled",
				},
			},
			expected: "Variable not found",
		},
		{
			name:     "get variable with default",
			template: "Status: {{getVar .Variables \"status\" \"unknown\"}}",
			data: TemplateData{
				Variables: map[string]string{
					"feature": "enabled",
				},
			},
			expected: "Status: unknown",
		},
		{
			name:     "get variable with default - variable exists",
			template: "Status: {{getVar .Variables \"status\" \"unknown\"}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "active",
				},
			},
			expected: "Status: active",
		},
		{
			name:     "string comparison - equals",
			template: "{{if eq (getVar .Variables \"status\" \"\") \"active\"}}System active{{else}}System inactive{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "active",
				},
			},
			expected: "System active",
		},
		{
			name:     "string comparison - not equals",
			template: "{{if ne (getVar .Variables \"status\" \"\") \"active\"}}System inactive{{else}}System active{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "inactive",
				},
			},
			expected: "System inactive",
		},
		{
			name:     "string contains",
			template: "{{if contains (getVar .Variables \"message\" \"\") \"error\"}}Error detected{{else}}No errors{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"message": "system error occurred",
				},
			},
			expected: "Error detected",
		},
		{
			name:     "numeric comparison - less than",
			template: "{{if lt (getVar .Variables \"count\" \"0\") \"10\"}}Count is small{{else}}Count is large{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"count": "5",
				},
			},
			expected: "Count is small",
		},
		{
			name:     "numeric comparison - greater than",
			template: "{{if gt (getVar .Variables \"count\" \"0\") \"10\"}}Count is large{{else}}Count is small{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"count": "15",
				},
			},
			expected: "Count is large",
		},
		{
			name:     "logical AND",
			template: "{{if and (hasVar .Variables \"feature\") (eq (getVar .Variables \"feature\" \"\") \"enabled\")}}Feature is enabled{{else}}Feature is not enabled{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"feature": "enabled",
				},
			},
			expected: "Feature is enabled",
		},
		{
			name:     "logical OR",
			template: "{{if or (eq (getVar .Variables \"status\" \"\") \"active\") (eq (getVar .Variables \"mode\" \"\") \"testing\")}}System is running{{else}}System is stopped{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "inactive",
					"mode":   "testing",
				},
			},
			expected: "System is running",
		},
		{
			name:     "logical NOT",
			template: "{{if not (hasVar .Variables \"error\")}}No errors{{else}}Errors detected{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "active",
				},
			},
			expected: "No errors",
		},
		{
			name:     "nested conditionals",
			template: "{{if hasVar .Variables \"status\"}}{{if eq .Variables.status \"active\"}}Active system{{else}}Inactive system{{end}}{{else}}Unknown status{{end}}",
			data: TemplateData{
				Variables: map[string]string{
					"status": "active",
				},
			},
			expected: "Active system",
		},
		{
			name: "complex conditional with multiple levels",
			template: `{{if hasVar .Variables "environment"}}
  {{if eq .Variables.environment "production"}}
    {{if hasVar .Variables "status"}}
      {{if eq .Variables.status "healthy"}}Production is healthy{{else}}Production has issues{{end}}
    {{else}}Production status unknown{{end}}
  {{else if eq .Variables.environment "staging"}}
    Staging environment
  {{else}}
    Development environment
  {{end}}
{{else}}
  Environment not specified
{{end}}`,
			data: TemplateData{
				Variables: map[string]string{
					"environment": "production",
					"status":      "healthy",
				},
			},
			expected: "\n  \n    \n      Production is healthy\n    \n  \n",
		},
	}

	// Register and test each template
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Register the template
			templateName := "test_" + tc.name
			err := manager.RegisterTemplate(templateName, tc.template)
			if err != nil {
				t.Fatalf("Failed to register template: %v", err)
			}

			// Execute the template
			result, err := manager.Execute(templateName, tc.data)
			if err != nil {
				t.Fatalf("Failed to execute template: %v", err)
			}

			// Check the result
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestConditionalErrorHandling(t *testing.T) {
	manager := GetManager()

	// Invalid templates that should fail validation
	invalidTemplates := []struct {
		name     string
		template string
	}{
		{
			name:     "unclosed if block",
			template: "{{if eq .Model \"TestModel\"}}Unclosed block",
		},
		{
			name:     "mismatched blocks",
			template: "{{if eq .Model \"TestModel\"}}{{else}}{{if eq .Model \"OtherModel\"}}Nested{{end}}",
		},
		{
			name:     "syntax error",
			template: "{{if .Model {{end}}",
		},
	}

	for _, tc := range invalidTemplates {
		t.Run(tc.name, func(t *testing.T) {
			// Try to register the template
			err := manager.RegisterTemplate("invalid_"+tc.name, tc.template)
			if err == nil {
				t.Errorf("Expected error for invalid template %s, got nil", tc.name)
			}
		})
	}
}

func TestIntegrationWithVariableSubstitution(t *testing.T) {
	// This tests how the template system works with variable substitution
	// when variables appear both in conditions and in text
	templateContent := `# Report for {{.Variables.project}}

{{if eq .Variables.environment "production"}}
## Production Environment Status
Current status: {{.Variables.status}}
{{if eq .Variables.status "healthy"}}
All systems operational.
{{else}}
Warning: System requires attention!
{{end}}
{{else}}
## Test Environment Status
This is a test environment.
{{end}}

Report generated on {{.Variables.date}}.`

	// Register the template
	manager := GetManager()
	templateName := "integration_test"
	err := manager.RegisterTemplate(templateName, templateContent)
	if err != nil {
		t.Fatalf("Failed to register template: %v", err)
	}

	// Create test data
	data := TemplateData{
		Variables: map[string]string{
			"project":     "CronAI",
			"environment": "production",
			"status":      "healthy",
			"date":        "2025-05-12",
		},
	}

	// Expected output for this data
	expected := `# Report for CronAI


## Production Environment Status
Current status: healthy

All systems operational.



Report generated on 2025-05-12.`

	// Execute the template
	result, err := manager.Execute(templateName, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Check the result
	if result != expected {
		t.Errorf("Expected:\n%q\n\nGot:\n%q", expected, result)
	}

	// Test with different values
	data.Variables["status"] = "degraded"
	expectedDegraded := `# Report for CronAI


## Production Environment Status
Current status: degraded

Warning: System requires attention!



Report generated on 2025-05-12.`

	resultDegraded, err := manager.Execute(templateName, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	if resultDegraded != expectedDegraded {
		t.Errorf("Expected degraded:\n%q\n\nGot:\n%q", expectedDegraded, resultDegraded)
	}

	// Test with different environment
	data.Variables["environment"] = "testing"
	expectedTesting := `# Report for CronAI


## Test Environment Status
This is a test environment.


Report generated on 2025-05-12.`

	resultTesting, err := manager.Execute(templateName, data)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	if resultTesting != expectedTesting {
		t.Errorf("Expected testing:\n%q\n\nGot:\n%q", expectedTesting, resultTesting)
	}
}
