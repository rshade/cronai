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

// PromptCategories defines standard categories for prompt organization
var PromptCategories = []string{
	"general",
	"email",
	"slack",
	"webhook",
	"file",
	"github",
	"custom",
}

// LoadPrompt loads a prompt from the cron_prompts directory
func LoadPrompt(promptName string) (string, error) {
	// Add .md extension if not present
	if !strings.HasSuffix(promptName, ".md") {
		promptName = promptName + ".md"
	}

	// Try different paths for the prompt file
	var paths []string

	// First check if CRON_PROMPTS_DIR environment variable is set
	if dir := os.Getenv("CRON_PROMPTS_DIR"); dir != "" {
		// Try the environment variable path first
		paths = append(paths, filepath.Join(dir, promptName))
		// Try category subdirectories under env path
		if !strings.Contains(promptName, string(os.PathSeparator)) {
			baseFilename := filepath.Base(promptName)
			for _, category := range PromptCategories {
				paths = append(paths, filepath.Join(dir, category, baseFilename))
			}
			// Also try a "general" category directory if not in the standard categories
			paths = append(paths, filepath.Join(dir, "general", baseFilename))
		}
	}

	// Fall back to default paths
	paths = append(paths,
		// First try the exact path provided (which might include a category subdirectory)
		filepath.Join("cron_prompts", promptName),
		// Then try from project root
		filepath.Join("..", "..", "cron_prompts", promptName),
	)

	// If promptName doesn't contain a directory separator, also try category subdirectories
	if !strings.Contains(promptName, string(os.PathSeparator)) {
		baseFilename := filepath.Base(promptName)
		for _, category := range PromptCategories {
			// Try category subdirectories relative to current directory
			paths = append(paths, filepath.Join("cron_prompts", category, baseFilename))
			// Try category subdirectories relative to project root
			paths = append(paths, filepath.Join("..", "..", "cron_prompts", category, baseFilename))
		}
	}

	// Try each path until we find the file
	var promptPath string
	var fileExists bool

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			promptPath = path
			fileExists = true
			break
		}
	}

	if !fileExists {
		return "", fmt.Errorf("prompt file not found: %s (tried all category directories)", promptName)
	}

	// Read the prompt file
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %w", err)
	}

	// Remove trailing newlines to ensure tests pass
	return strings.TrimRight(string(promptContent), "\n"), nil
}

// GetPromptPath resolves the file path to a prompt file
func GetPromptPath(promptName string) (string, error) {
	// Add .md extension if not present
	if !strings.HasSuffix(promptName, ".md") {
		promptName = promptName + ".md"
	}

	// Try different paths for the prompt file
	paths := []string{
		// First try the exact path provided (which might include a category subdirectory)
		filepath.Join("cron_prompts", promptName),
		// Then try from project root
		filepath.Join("..", "..", "cron_prompts", promptName),
	}

	// If promptName doesn't contain a directory separator, also try category subdirectories
	if !strings.Contains(promptName, string(os.PathSeparator)) {
		baseFilename := filepath.Base(promptName)
		for _, category := range PromptCategories {
			// Try category subdirectories relative to current directory
			paths = append(paths, filepath.Join("cron_prompts", category, baseFilename))
			// Try category subdirectories relative to project root
			paths = append(paths, filepath.Join("..", "..", "cron_prompts", category, baseFilename))
		}
	}

	// Try each path until we find the file
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("prompt file not found: %s (tried all category directories)", promptName)
}

// LoadPromptWithVariables loads a prompt and processes it as a template with variables
func LoadPromptWithVariables(promptName string, variables map[string]string) (string, error) {
	// Load the base prompt
	promptContent, err := LoadPrompt(promptName)
	if err != nil {
		return "", err
	}

	// Extract metadata and content
	metadata, content, err := ExtractMetadata(promptContent, promptName)
	if err != nil {
		return "", fmt.Errorf("failed to extract content: %w", err)
	}

	// Check if this prompt extends another template (inheritance)
	if metadata.Extends != "" {
		// Process the prompt with inheritance support
		_, finalContent, err := ProcessPromptWithInheritance(promptName, promptContent, variables)
		if err != nil {
			return "", fmt.Errorf("failed to process prompt with inheritance: %w", err)
		}
		return finalContent, nil
	}

	// Process includes for non-inheritance templates
	processedContent, err := ProcessIncludes(content)
	if err != nil {
		return "", err
	}

	// Check if the processed content contains template directives
	if containsTemplateDirectives(processedContent) {
		// Validate template syntax
		if err := ValidatePromptTemplate(processedContent, promptName); err != nil {
			return "", fmt.Errorf("prompt '%s' contains invalid template syntax: %w", promptName, err)
		}

		// Process as a template with the template engine
		return processPromptAsTemplate(processedContent, promptName, variables)
	}

	// Fallback to simple variable substitution for backward compatibility
	return ApplyVariables(processedContent, variables), nil
}

// containsTemplateDirectives checks if the content contains template directives like {{if}}, {{else}}, etc.
func containsTemplateDirectives(content string) bool {
	templateDirectivePattern := regexp.MustCompile(`\{\{\s*(if|else|end|range|with|define|template|block)\b`)
	return templateDirectivePattern.MatchString(content)
}

// ValidatePromptTemplate validates that a prompt contains valid template syntax
func ValidatePromptTemplate(content string, promptName string) error {
	// Get template manager
	manager := template.GetManager()

	// Use the template manager's Validate method
	// Pass the prompt name as the template name for validation
	if err := manager.Validate("prompt_"+promptName, content); err != nil {
		return fmt.Errorf("invalid template syntax in prompt: %w", err)
	}

	return nil
}

// processPromptAsTemplate processes the prompt content as a Go template with conditional logic
func processPromptAsTemplate(content string, promptName string, variables map[string]string) (string, error) {
	// Validate the template syntax first
	if err := ValidatePromptTemplate(content, promptName); err != nil {
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
	data := template.Data{
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

	// Regular expression to match both {{variable}} and {{.Variables.variable}} patterns
	variablePattern := regexp.MustCompile(`\{\{(?:\.Variables\.)?(\w+)\}\}`)

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
		// Extract the variable name
		// Handle both {{name}} and {{.Variables.name}} patterns
		var varName string
		if strings.Contains(match, ".Variables.") {
			// Extract from {{.Variables.name}}
			parts := strings.Split(match, ".")
			if len(parts) >= 3 {
				varName = strings.TrimSuffix(parts[2], "}}")
			}
		} else {
			// Extract from {{name}}
			varName = match[2 : len(match)-2]
		}

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

// ProcessIncludes processes {{include "template_name"}} directives in the prompt content
func ProcessIncludes(content string) (string, error) {
	return processIncludesWithDepth(content, 0)
}

// processIncludesWithDepth processes includes with recursion depth tracking
func processIncludesWithDepth(content string, depth int) (string, error) {
	const maxDepth = 10
	if depth > maxDepth {
		return "", fmt.Errorf("maximum recursion depth exceeded (%d)", maxDepth)
	}

	// Regular expression to match {{include "template_name"}} or {{include template_name}} patterns
	includePattern := regexp.MustCompile(`\{\{include\s+(?:"([^"]+)"|([^}\s]+))(?:\s+.*)?\}\}`)

	// Find all includes in the content
	includes := includePattern.FindAllStringSubmatch(content, -1)
	if len(includes) == 0 {
		return content, nil
	}

	// Process each include
	result := content
	for _, includeMatch := range includes {
		if len(includeMatch) < 3 {
			continue
		}

		// Check which capture group has the match (quoted or unquoted)
		includePath := includeMatch[1]
		if includePath == "" {
			includePath = includeMatch[2]
		}
		// We'll preserve parameter parsing for future enhancements
		// but won't use it for now
		// if len(includeMatch) > 2 && includeMatch[2] != "" {
		//	includeParams = includeMatch[2]
		// }

		var includeContent string
		var err error

		// First check library path for component templates
		libraryPaths := []string{
			filepath.Join("templates", "library", includePath+".tmpl"),
			filepath.Join("templates", "library", includePath+".md"),
			filepath.Join("..", "templates", "library", includePath+".tmpl"),
			filepath.Join("..", "templates", "library", includePath+".md"),
			filepath.Join("..", "..", "templates", "library", includePath+".tmpl"),
			filepath.Join("..", "..", "templates", "library", includePath+".md"),
		}

		foundInLibrary := false
		for _, libPath := range libraryPaths {
			if _, statErr := os.Stat(libPath); statErr == nil {
				data, readErr := os.ReadFile(libPath)
				if readErr != nil {
					return "", fmt.Errorf("failed to read library template %q: %w", libPath, readErr)
				}
				includeContent = string(data)
				foundInLibrary = true
				break
			}
		}

		if !foundInLibrary {
			// Check if the includePath is an absolute path that exists (for tests)
			if filepath.IsAbs(includePath) {
				if _, statErr := os.Stat(includePath); statErr == nil {
					data, readErr := os.ReadFile(includePath)
					if readErr != nil {
						return "", fmt.Errorf("failed to read include file %q: %w", includePath, readErr)
					}
					includeContent = string(data)
				} else {
					// If the absolute path doesn't exist, try to load it as a regular prompt
					includeContent, err = LoadPrompt(includePath)
					if err != nil {
						return "", fmt.Errorf("failed to include %q: %w", includePath, err)
					}
				}
			} else {
				// Regular prompt loading for non-absolute paths
				includeContent, err = LoadPrompt(includePath)
				if err != nil {
					return "", fmt.Errorf("failed to include %q: %w", includePath, err)
				}
			}
		}

		// Extract included content without its metadata
		_, parsedContent, err := ExtractMetadata(includeContent, includePath)
		if err != nil {
			return "", fmt.Errorf("failed to extract content from include %q: %w", includePath, err)
		}

		// Replace the include directive with the content
		// Trim any extra whitespace
		parsedContent = strings.TrimSpace(parsedContent)
		originalInclude := includeMatch[0]

		// Simply replace the include directive with the parsed content
		result = strings.Replace(result, originalInclude, parsedContent, 1)
	}

	// Check if there are nested includes that need to be processed
	if includePattern.MatchString(result) {
		return processIncludesWithDepth(result, depth+1)
	}

	return result, nil
}

// LoadPromptWithIncludes loads a prompt and processes any {{include}} directives
func LoadPromptWithIncludes(promptName string) (string, error) {
	// Load the base prompt
	promptContent, err := LoadPrompt(promptName)
	if err != nil {
		return "", err
	}

	// Extract metadata and content
	_, content, err := ExtractMetadata(promptContent, promptName)
	if err != nil {
		return "", fmt.Errorf("failed to extract content: %w", err)
	}

	// Process includes
	processedContent, err := ProcessIncludes(content)
	if err != nil {
		return "", err
	}

	return processedContent, nil
}

// ProcessPromptWithInheritance processes a prompt with template inheritance
func ProcessPromptWithInheritance(path, content string, variables map[string]string) (map[string]string, string, error) {
	// Extract metadata to check for 'extends' property
	metadata, extractedContent, err := ExtractMetadata(content, path)
	if err != nil {
		return variables, "", fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Check if this prompt extends another
	if metadata.Extends != "" {
		// Try to load the parent prompt
		parentPrompt, err := LoadPrompt(metadata.Extends)
		if err != nil {
			return variables, "", fmt.Errorf("failed to load parent prompt %q: %w", metadata.Extends, err)
		}

		// Extract parent content
		_, parentContent, err := ExtractMetadata(parentPrompt, metadata.Extends)
		if err != nil {
			return variables, "", fmt.Errorf("failed to extract parent content: %w", err)
		}

		// Process includes in parent content
		processedParentContent, err := ProcessIncludes(parentContent)
		if err != nil {
			return variables, "", err
		}

		// Process includes in child content
		processedChildContent, err := ProcessIncludes(extractedContent)
		if err != nil {
			return variables, "", err
		}

		// Use the template engine to handle block overrides
		// This would merge the child's blocks into the parent template
		tmplManager := template.GetManager()

		// Register parent template first
		parentName := "parent_" + filepath.Base(path)
		err = tmplManager.RegisterTemplate(parentName, processedParentContent)
		if err != nil {
			return variables, "", fmt.Errorf("failed to register parent template: %w", err)
		}

		// Register child template with extends directive
		childName := "child_" + filepath.Base(path)
		childContentWithExtends := "{{extends \"" + parentName + "\"}}\n" + processedChildContent
		err = tmplManager.RegisterTemplate(childName, childContentWithExtends)
		if err != nil {
			return variables, "", fmt.Errorf("failed to register child template: %w", err)
		}

		// Create template data with variables
		data := template.Data{
			Variables: variables,
			Timestamp: time.Now(),
		}

		// Execute the child template (which will inherit from parent)
		result, err := tmplManager.Execute(childName, data)
		if err != nil {
			return variables, "", fmt.Errorf("failed to execute template with inheritance: %w", err)
		}

		return variables, result, nil
	}

	// If no inheritance, just process normally
	processedContent, err := ProcessIncludes(extractedContent)
	if err != nil {
		return variables, "", err
	}

	return variables, ApplyVariables(processedContent, variables), nil
}

// CreatePromptWithMetadata creates a new prompt file with the given metadata
func CreatePromptWithMetadata(category, promptName string, metadata *Metadata, content string) error {
	// Ensure the prompts directory exists
	promptsDir := "cron_prompts"
	if dir := os.Getenv("CRON_PROMPTS_DIR"); dir != "" {
		promptsDir = dir
	}

	// If category is specified, create the category subdirectory
	if category != "" {
		promptsDir = filepath.Join(promptsDir, category)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create prompts directory: %w", err)
	}

	// Create the prompt file path
	filePath := filepath.Join(promptsDir, promptName+".md")

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("prompt file %s already exists", filePath)
	}

	// Format the metadata section
	metadataStr := "---\n"
	if metadata.Name != "" {
		metadataStr += fmt.Sprintf("name: %s\n", metadata.Name)
	}
	if metadata.Description != "" {
		metadataStr += fmt.Sprintf("description: %s\n", metadata.Description)
	}
	if metadata.Author != "" {
		metadataStr += fmt.Sprintf("author: %s\n", metadata.Author)
	}
	if metadata.Version != "" {
		metadataStr += fmt.Sprintf("version: %s\n", metadata.Version)
	}
	if category != "" {
		metadataStr += fmt.Sprintf("category: %s\n", category)
	} else if metadata.Category != "" {
		metadataStr += fmt.Sprintf("category: %s\n", metadata.Category)
	}
	if len(metadata.Tags) > 0 {
		metadataStr += fmt.Sprintf("tags: %s\n", strings.Join(metadata.Tags, ", "))
	}
	if len(metadata.Variables) > 0 {
		metadataStr += "variables:\n"
		for _, v := range metadata.Variables {
			metadataStr += fmt.Sprintf("  - name: %s\n", v.Name)
			metadataStr += fmt.Sprintf("    description: %s\n", v.Description)
		}
	}
	if metadata.Extends != "" {
		metadataStr += fmt.Sprintf("extends: %s\n", metadata.Extends)
	}
	metadataStr += "---\n\n"

	// Combine metadata and content
	fullContent := metadataStr + content

	// Write the file
	if err := os.WriteFile(filePath, []byte(fullContent), 0644); err != nil {
		return fmt.Errorf("failed to write prompt file: %w", err)
	}

	return nil
}
