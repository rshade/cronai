package prompt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/processor/template"
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

// LoadPromptWithVariables loads a prompt and processes it as a template with variables
func LoadPromptWithVariables(promptName string, variables map[string]string) (string, error) {
	// Load the base prompt
	promptContent, err := LoadPrompt(promptName)
	if err != nil {
		return "", err
	}

	// First check if the prompt contains template directives
	if containsTemplateDirectives(promptContent) {
		// Validate template syntax
		if err := ValidatePromptTemplate(promptContent); err != nil {
			return "", fmt.Errorf("prompt '%s' contains invalid template syntax: %w", promptName, err)
		}

		// Process as a template with the template engine
		return processPromptAsTemplate(promptContent, promptName, variables)
	}

	// Fallback to simple variable substitution for backward compatibility
	return ApplyVariables(promptContent, variables), nil
}

// containsTemplateDirectives checks if the content contains template directives like {{if}}, {{else}}, etc.
func containsTemplateDirectives(content string) bool {
	templateDirectivePattern := regexp.MustCompile(`\{\{\s*(if|else|end|range|with|define|template|block)\b`)
	return templateDirectivePattern.MatchString(content)
}

// ValidatePromptTemplate validates that a prompt contains valid template syntax
func ValidatePromptTemplate(content string) error {
	// Get template manager
	manager := template.GetManager()

	// Use the template manager's Validate method
	if err := manager.Validate(content); err != nil {
		return fmt.Errorf("invalid template syntax in prompt: %w", err)
	}

	return nil
}

// processPromptAsTemplate processes the prompt content as a Go template with conditional logic
func processPromptAsTemplate(content string, promptName string, variables map[string]string) (string, error) {
	// Validate the template syntax first
	if err := ValidatePromptTemplate(content); err != nil {
		return "", err
	}

	// Get or create a template manager for prompts
	manager := template.GetManager()

	// Register the prompt as a template
	templateName := "prompt_" + promptName
	err := manager.RegisterTemplate(templateName, content)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt as template: %w", err)
	}

	// Prepare the template data with variables
	data := template.TemplateData{
		Variables: variables,
		// Add other fields that might be useful in prompt templates
		Timestamp: time.Now(),
	}

	// Execute the template
	var buf bytes.Buffer
	tmpl, err := manager.GetTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to get template: %w", err)
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ApplyVariables replaces variable placeholders in the format {{variable_name}} with their values
// while preserving conditional logic syntax
func ApplyVariables(content string, variables map[string]string) string {
	// If no variables provided, return the original content
	if variables == nil {
		return content
	}

	// Regular expression to match standalone {{variable}} patterns
	// Matches {{variable}} but we'll filter out template directives later
	variablePattern := regexp.MustCompile(`\{\{(\w+)\}\}`)

	// Template directive keywords that we should ignore
	templateDirectives := map[string]bool{
		"if":       true,
		"else":     true,
		"end":      true,
		"range":    true,
		"with":     true,
		"define":   true,
		"template": true,
		"block":    true,
	}

	// Replace all variables in the content
	result := variablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the variable name (remove the {{ and }})
		varName := match[2 : len(match)-2]

		// Skip template directive keywords
		if _, isDirective := templateDirectives[varName]; isDirective {
			return match
		}

		// Also skip if it's part of a template directive block
		// This is a simplistic check - it won't catch all cases but will handle most common ones
		parts := strings.Fields(varName)
		if len(parts) > 0 && templateDirectives[parts[0]] {
			return match
		}

		// Look up the value in the variables map
		if value, exists := variables[varName]; exists {
			return value
		}

		// If the variable doesn't exist, leave it as is
		return match
	})

	return result
}
