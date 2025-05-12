package prompt

import (
	"testing"
)

func TestLoadPrompt(t *testing.T) {
	// Test using the existing test prompt file
	testPrompt := "test_prompt"
	testPromptContent := "# Test Prompt\n\nThis is a test prompt."

	// Test loading the prompt without .md extension
	content, err := LoadPrompt(testPrompt)
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	if content != testPromptContent {
		t.Errorf("Expected prompt content %q, got %q", testPromptContent, content)
	}

	// Test loading the prompt with .md extension
	content, err = LoadPrompt(testPrompt + ".md")
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	if content != testPromptContent {
		t.Errorf("Expected prompt content %q, got %q", testPromptContent, content)
	}

	// Test loading a non-existent prompt
	_, err = LoadPrompt("non_existent_prompt")
	if err == nil {
		t.Error("Expected error when loading non-existent prompt, got nil")
	}
}

func TestLoadPromptWithVariables(t *testing.T) {
	// Use the existing test prompt file with variables
	testPrompt := "test_prompt_vars"
	testPromptContent := "# Test Prompt\n\nHello {{name}},\n\nThis is a test prompt for {{project}} with {{variable}} that doesn't exist."

	// Define test variables
	variables := map[string]string{
		"name":    "User",
		"project": "CronAI",
	}

	// Expected output after variable substitution
	expectedOutput := "# Test Prompt\n\nHello User,\n\nThis is a test prompt for CronAI with {{variable}} that doesn't exist."

	// Test loading the prompt with variables
	content, err := LoadPromptWithVariables(testPrompt, variables)
	if err != nil {
		t.Fatalf("Failed to load prompt with variables: %v", err)
	}

	if content != expectedOutput {
		t.Errorf("Expected prompt content with variables\n%q, got\n%q", expectedOutput, content)
	}

	// Test loading the prompt with empty variables map
	content, err = LoadPromptWithVariables(testPrompt, map[string]string{})
	if err != nil {
		t.Fatalf("Failed to load prompt with empty variables: %v", err)
	}

	if content != testPromptContent {
		t.Errorf("Expected original prompt content with empty variables map, got %q", content)
	}

	// Test loading the prompt with nil variables map
	content, err = LoadPromptWithVariables(testPrompt, nil)
	if err != nil {
		t.Fatalf("Failed to load prompt with nil variables: %v", err)
	}

	if content != testPromptContent {
		t.Errorf("Expected original prompt content with nil variables map, got %q", content)
	}
}

func TestApplyVariables(t *testing.T) {
	testContent := "Hello {{name}}, welcome to {{project}}. Your ID is {{id}}."
	variables := map[string]string{
		"name":    "User",
		"project": "CronAI",
		"id":      "12345",
	}

	expected := "Hello User, welcome to CronAI. Your ID is 12345."
	result := ApplyVariables(testContent, variables)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with missing variables
	partialVars := map[string]string{
		"name": "User",
	}
	expectedPartial := "Hello User, welcome to {{project}}. Your ID is {{id}}."
	resultPartial := ApplyVariables(testContent, partialVars)

	if resultPartial != expectedPartial {
		t.Errorf("Expected %q, got %q", expectedPartial, resultPartial)
	}

	// Test with nil variables
	nilResult := ApplyVariables(testContent, nil)
	if nilResult != testContent {
		t.Errorf("Expected original content with nil variables, got %q", nilResult)
	}
}