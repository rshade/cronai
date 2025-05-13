package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PromptInfo represents basic information about a prompt file
type PromptInfo struct {
	Name        string
	Path        string
	Category    string
	Description string
	HasMetadata bool
}

// ListPrompts lists all prompt files in the cron_prompts directory
func ListPrompts() ([]PromptInfo, error) {
	// Find the base cron_prompts directory
	basePaths := []string{
		"cron_prompts",
		filepath.Join("..", "..", "cron_prompts"),
	}

	var promptsDir string
	for _, path := range basePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			promptsDir = path
			break
		}
	}

	if promptsDir == "" {
		return nil, fmt.Errorf("cron_prompts directory not found")
	}

	var prompts []PromptInfo

	// Walk the directory and gather prompt files
	err := filepath.Walk(promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .md files
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip README.md files
		if strings.EqualFold(info.Name(), "README.md") {
			return nil
		}

		// Determine the category based on directory structure
		relPath, err := filepath.Rel(promptsDir, path)
		if err != nil {
			return err
		}

		// Get category from path
		category := "root"
		if dir := filepath.Dir(relPath); dir != "." {
			category = dir
		}

		// Create a basic prompt info
		promptInfo := PromptInfo{
			Name:     strings.TrimSuffix(filepath.Base(path), ".md"),
			Path:     relPath,
			Category: category,
		}

		// Try to read metadata for more information
		if content, err := os.ReadFile(path); err == nil {
			if metadata, _, err := ExtractMetadata(string(content), relPath); err == nil && metadata != nil {
				promptInfo.HasMetadata = true
				if metadata.Name != "" {
					promptInfo.Name = metadata.Name
				}
				promptInfo.Description = metadata.Description
			}
		}

		prompts = append(prompts, promptInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	return prompts, nil
}

// SearchPrompts searches for prompts matching the given criteria in metadata (name, path, description)
func SearchPrompts(query string, category string) ([]PromptInfo, error) {
	// Get all prompts
	allPrompts, err := ListPrompts()
	if err != nil {
		return nil, err
	}

	// Filter by criteria
	var filteredPrompts []PromptInfo

	queryLower := strings.ToLower(query)
	for _, prompt := range allPrompts {
		// Filter by category if specified
		if category != "" && !strings.EqualFold(prompt.Category, category) {
			continue
		}

		// Filter by query if specified
		if query != "" {
			// Check if query matches name, path, or description
			nameLower := strings.ToLower(prompt.Name)
			pathLower := strings.ToLower(prompt.Path)
			descLower := strings.ToLower(prompt.Description)

			if !strings.Contains(nameLower, queryLower) &&
				!strings.Contains(pathLower, queryLower) &&
				!strings.Contains(descLower, queryLower) {
				continue
			}
		}

		filteredPrompts = append(filteredPrompts, prompt)
	}

	return filteredPrompts, nil
}

// SearchPromptContent searches for prompts with content matching the given query
func SearchPromptContent(query string, category string) ([]PromptInfo, error) {
	// If empty query, return all prompts in the category
	if query == "" {
		return SearchPrompts("", category)
	}

	// Find the base cron_prompts directory
	basePaths := []string{
		"cron_prompts",
		filepath.Join("..", "..", "cron_prompts"),
	}

	var promptsDir string
	for _, path := range basePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			promptsDir = path
			break
		}
	}

	if promptsDir == "" {
		return nil, fmt.Errorf("cron_prompts directory not found")
	}

	var matchingPrompts []PromptInfo
	queryLower := strings.ToLower(query)

	// Walk the directory and search prompt contents
	err := filepath.Walk(promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include .md files
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// Skip README.md files
		if strings.EqualFold(info.Name(), "README.md") {
			return nil
		}

		// Determine the category based on directory structure
		relPath, err := filepath.Rel(promptsDir, path)
		if err != nil {
			return err
		}

		// Get category from path
		fileCategory := "root"
		if dir := filepath.Dir(relPath); dir != "." {
			fileCategory = dir
		}

		// Filter by category if specified
		if category != "" && !strings.EqualFold(fileCategory, category) {
			return nil
		}

		// Read the file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		// Check if content contains the query
		if !strings.Contains(strings.ToLower(string(content)), queryLower) {
			return nil
		}

		// Create a prompt info object for the match
		promptInfo := PromptInfo{
			Name:     strings.TrimSuffix(filepath.Base(path), ".md"),
			Path:     relPath,
			Category: fileCategory,
		}

		// Try to read metadata for more information
		if metadata, _, err := ExtractMetadata(string(content), relPath); err == nil && metadata != nil {
			promptInfo.HasMetadata = true
			if metadata.Name != "" {
				promptInfo.Name = metadata.Name
			}
			promptInfo.Description = metadata.Description
		}

		matchingPrompts = append(matchingPrompts, promptInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search prompt content: %w", err)
	}

	return matchingPrompts, nil
}

// GetPromptInfo gets detailed information about a specific prompt
func GetPromptInfo(promptName string) (*PromptMetadata, error) {
	// Load the prompt and extract metadata
	metadata, err := GetPromptMetadata(promptName)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// CreatePromptWithMetadata creates a new prompt file with metadata
func CreatePromptWithMetadata(category, promptName string, metadata *PromptMetadata, content string) error {
	// Ensure prompt has .md extension
	if !strings.HasSuffix(promptName, ".md") {
		promptName = promptName + ".md"
	}

	// Find the base cron_prompts directory
	basePaths := []string{
		"cron_prompts",
		filepath.Join("..", "..", "cron_prompts"),
	}

	var promptsDir string
	for _, path := range basePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			promptsDir = path
			break
		}
	}

	if promptsDir == "" {
		return fmt.Errorf("cron_prompts directory not found")
	}

	// Determine the full path
	var fullPath string
	if category != "" {
		fullPath = filepath.Join(promptsDir, category, promptName)

		// Ensure category directory exists
		categoryDir := filepath.Join(promptsDir, category)
		if err := os.MkdirAll(categoryDir, 0755); err != nil {
			return fmt.Errorf("failed to create category directory: %w", err)
		}
	} else {
		fullPath = filepath.Join(promptsDir, promptName)
	}

	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("prompt file already exists: %s", fullPath)
	}

	// Build the content with metadata
	var fileContent strings.Builder

	// Add metadata section if provided
	if metadata != nil {
		fileContent.WriteString("---\n")
		if metadata.Name != "" {
			fileContent.WriteString(fmt.Sprintf("name: %s\n", metadata.Name))
		}
		if metadata.Description != "" {
			fileContent.WriteString(fmt.Sprintf("description: %s\n", metadata.Description))
		}
		if metadata.Author != "" {
			fileContent.WriteString(fmt.Sprintf("author: %s\n", metadata.Author))
		}
		if metadata.Version != "" {
			fileContent.WriteString(fmt.Sprintf("version: %s\n", metadata.Version))
		}
		if metadata.Category != "" {
			fileContent.WriteString(fmt.Sprintf("category: %s\n", metadata.Category))
		}
		if len(metadata.Tags) > 0 {
			fileContent.WriteString(fmt.Sprintf("tags: %s\n", strings.Join(metadata.Tags, ", ")))
		}
		if len(metadata.Variables) > 0 {
			fileContent.WriteString("variables:\n")
			for _, v := range metadata.Variables {
				fileContent.WriteString(fmt.Sprintf("  - name: %s\n", v.Name))
				fileContent.WriteString(fmt.Sprintf("    description: %s\n", v.Description))
			}
		}
		fileContent.WriteString("---\n\n")
	}

	// Add the content
	fileContent.WriteString(content)

	// Write the file
	err := os.WriteFile(fullPath, []byte(fileContent.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write prompt file: %w", err)
	}

	return nil
}
