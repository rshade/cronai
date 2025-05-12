package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LoadPrompt loads a prompt from the cron_prompts directory
func LoadPrompt(promptName string) (string, error) {
	// Add .md extension if not present
	if !strings.HasSuffix(promptName, ".md") {
		promptName = promptName + ".md"
	}

	// Try different paths for the prompt file
	// First check relative to the current directory
	promptPath := filepath.Join("cron_prompts", promptName)

	// If not found, try relative to the project root
	if _, err := os.Stat(promptPath); os.IsNotExist(err) {
		// Try with project root path
		promptPath = filepath.Join("..", "..", "cron_prompts", promptName)
	}

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

// LoadPromptWithVariables loads a prompt and replaces variables with their values
func LoadPromptWithVariables(promptName string, variables map[string]string) (string, error) {
	// Load the base prompt
	promptContent, err := LoadPrompt(promptName)
	if err != nil {
		return "", err
	}

	// Apply variables to the prompt content
	return ApplyVariables(promptContent, variables), nil
}

// ApplyVariables replaces variable placeholders in the format {{variable_name}} with their values
func ApplyVariables(content string, variables map[string]string) string {
	// If no variables provided, return the original content
	if variables == nil {
		return content
	}

	// Regular expression to match {{variable}} patterns
	variablePattern := regexp.MustCompile(`\{\{(\w+)\}\}`)

	// Replace all variables in the content
	result := variablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the variable name (remove the {{ and }})
		varName := match[2 : len(match)-2]

		// Look up the value in the variables map
		if value, exists := variables[varName]; exists {
			return value
		}

		// If the variable doesn't exist, leave it as is
		return match
	})

	return result
}
