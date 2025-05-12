package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrompt(t *testing.T) {
	// Create a temporary test prompt file
	testDir := "../../cron_prompts"
	testPrompt := "test_prompt"
	testPromptContent := "# Test Prompt\n\nThis is a test prompt."
	
	// Ensure the directory exists
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create the test prompt file
	testPromptPath := filepath.Join(testDir, testPrompt+".md")
	err = os.WriteFile(testPromptPath, []byte(testPromptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test prompt file: %v", err)
	}

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
