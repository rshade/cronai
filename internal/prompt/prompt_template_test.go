package prompt

import (
	"os"
	"strings"
	"testing"
)

func TestRoadmapValidationTemplate(t *testing.T) {
	// Load the template directly since we're in the same package
	promptContent, err := os.ReadFile("testdata/roadmap_validation.md")
	if err != nil {
		t.Fatalf("Failed to load roadmap validation template: %v", err)
	}

	prompt := string(promptContent)

	// Verify the prompt contains the expected sections
	expectedSections := []string{
		"System Instructions",
		"Context",
		"Current Project Status",
		"Current Roadmap",
		"Validation Tasks",
		"Output Format",
	}

	for _, section := range expectedSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("Template is missing expected section: %s", section)
		}
	}

	// Verify the prompt contains the variable placeholders
	expectedVars := []string{
		"{{project_status}}",
		"{{roadmap}}",
	}

	for _, variable := range expectedVars {
		if !strings.Contains(prompt, variable) {
			t.Errorf("Template is missing expected variable placeholder: %s", variable)
		}
	}

	// Test variable replacement
	vars := map[string]string{
		"project_status": "Project is 75% complete with core functionality implemented.",
		"roadmap":        "Q2 2025: MVP Release\nQ3 2025: Enhanced Usability\nQ4 2025: Integration & Scale",
	}

	// Apply variables directly
	promptWithVars := ApplyVariables(prompt, vars)

	// Verify that variable replacement worked
	if strings.Contains(promptWithVars, "{{project_status}}") {
		t.Error("Variable replacement failed: {{project_status}} still present in output")
	}

	if !strings.Contains(promptWithVars, vars["project_status"]) {
		t.Error("Variable replacement failed: project_status value not found in output")
	}

	if strings.Contains(promptWithVars, "{{roadmap}}") {
		t.Error("Variable replacement failed: {{roadmap}} still present in output")
	}

	if !strings.Contains(promptWithVars, vars["roadmap"]) {
		t.Error("Variable replacement failed: roadmap value not found in output")
	}
}
