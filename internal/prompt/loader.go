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

// Standard categories for prompt organization
var PromptCategories = []string{
	"system",
	"monitoring",
	"reports",
	"templates",
}

// LoadPrompt loads a prompt from the cron_prompts directory
func LoadPrompt(promptName string) (string, error) {
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

	return string(promptContent), nil
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
		if err := ValidatePromptTemplate(processedContent); err != nil {
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

// ProcessIncludes processes {{include "template_name"}} directives in the prompt content
func ProcessIncludes(content string) (string, error) {
	// Regular expression to match {{include "template_name"}} patterns
	includePattern := regexp.MustCompile(`\{\{include\s+"([^"]+)"(\s+.*)?\}\}`)

	// Find all includes in the content
	includes := includePattern.FindAllStringSubmatch(content, -1)
	if len(includes) == 0 {
		return content, nil
	}

	// Process each include
	result := content
	for _, includeMatch := range includes {
		if len(includeMatch) < 2 {
			continue
		}

		includePath := includeMatch[1]
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
						return "", fmt.Errorf("failed to load include %q: %w", includePath, err)
					}
				}
			} else {
				// Regular prompt loading for non-absolute paths
				includeContent, err = LoadPrompt(includePath)
				if err != nil {
					return "", fmt.Errorf("failed to load include %q: %w", includePath, err)
				}
			}
		}

		// Extract included content without its metadata
		_, parsedContent, err := ExtractMetadata(includeContent, includePath)
		if err != nil {
			return "", fmt.Errorf("failed to extract content from include %q: %w", includePath, err)
		}

		// Replace the include directive with the content
		// Trim any extra whitespace and ensure proper newlines
		parsedContent = strings.TrimSpace(parsedContent)
		originalInclude := includeMatch[0]

		// Replace the include directive while maintaining newlines correctly
		if strings.HasPrefix(originalInclude, "\n") {
			parsedContent = "\n" + parsedContent
		}

		// Ensure the replacement ends with a newline if the original did
		if strings.HasSuffix(originalInclude, "\n") {
			parsedContent = parsedContent + "\n"
		} else {
			parsedContent = parsedContent + "\n"
		}

		result = strings.Replace(result, originalInclude, parsedContent, 1)
	}

	// Check if there are nested includes that need to be processed
	if includePattern.MatchString(result) {
		return ProcessIncludes(result)
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
		data := template.TemplateData{
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
