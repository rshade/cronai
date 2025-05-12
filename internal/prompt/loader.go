package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadPrompt loads a prompt from the cron_prompts directory
func LoadPrompt(promptName string) (string, error) {
	// Add .md extension if not present
	if !strings.HasSuffix(promptName, ".md") {
		promptName = promptName + ".md"
	}

	// Build the prompt file path
	promptPath := filepath.Join("cron_prompts", promptName)

	// Check if file exists
	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("prompt file not found: %s", promptPath)
	}

	// Read the prompt file
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %w", err)
	}

	return string(promptContent), nil
}
