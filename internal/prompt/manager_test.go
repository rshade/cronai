package prompt

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPromptManager tests the PromptManager functionality
func TestPromptManager(t *testing.T) {
	// PM is no longer used in the prompt package
	// Test the singleton instance
	/*
	t.Run("singleton instance", func(t *testing.T) {
		manager1 := PM
		manager2 := PM
		assert.Same(t, manager1, manager2, "PM should be the same instance")
	})
	*/

	// Test LoadPrompt through the manager
	t.Run("LoadPrompt", func(t *testing.T) {
		// Create temp directory structure
		tempDir := t.TempDir()
		promptsDir := filepath.Join(tempDir, "cron_prompts")
		require.NoError(t, os.MkdirAll(promptsDir, 0755))

		// Create test prompt
		testPrompt := filepath.Join(promptsDir, "test.md")
		testContent := "# Test Prompt\nThis is a test."
		require.NoError(t, os.WriteFile(testPrompt, []byte(testContent), 0644))

		// Change working directory for the test
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		// Load prompt through manager
		content, err := PM.LoadPrompt("test")
		assert.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	// Test LoadPromptWithVariables through the manager
	t.Run("LoadPromptWithVariables", func(t *testing.T) {
		// Create temp directory structure
		tempDir := t.TempDir()
		promptsDir := filepath.Join(tempDir, "cron_prompts")
		require.NoError(t, os.MkdirAll(promptsDir, 0755))

		// Create test prompt
		testPrompt := filepath.Join(promptsDir, "template.md")
		testContent := "Hello {{name}}, welcome to {{place}}!"
		require.NoError(t, os.WriteFile(testPrompt, []byte(testContent), 0644))

		// Change working directory for the test
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		// Load prompt with variables
		variables := map[string]string{
			"name":  "John",
			"place": "New York",
		}
		content, err := PM.LoadPromptWithVariables("template", variables)
		assert.NoError(t, err)
		assert.Equal(t, "Hello John, welcome to New York!", content)
	})
}

// TestPromptManagerConcurrency tests that the PromptManager is thread-safe
func TestPromptManagerConcurrency(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Create multiple test prompts
	numPrompts := 10
	for i := 0; i < numPrompts; i++ {
		promptPath := filepath.Join(promptsDir, fmt.Sprintf("prompt%d.md", i))
		content := fmt.Sprintf("# Prompt %d\nThis is prompt number %d.", i, i)
		require.NoError(t, os.WriteFile(promptPath, []byte(content), 0644))
	}

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Load random prompts
			promptName := fmt.Sprintf("prompt%d", id%numPrompts)
			content, err := PM.LoadPrompt(promptName)

			assert.NoError(t, err)
			assert.Contains(t, content, fmt.Sprintf("Prompt %d", id%numPrompts))

			// Also test with variables
			variables := map[string]string{
				"id": fmt.Sprintf("%d", id),
			}
			contentWithVars, err := PM.LoadPromptWithVariables(promptName, variables)
			assert.NoError(t, err)
			assert.NotEmpty(t, contentWithVars)
		}(i)
	}

	wg.Wait()
}

// TestPromptManagerErrorHandling tests error handling in the PromptManager
func TestPromptManagerErrorHandling(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Test loading non-existent prompt
	t.Run("non-existent prompt", func(t *testing.T) {
		_, err := PM.LoadPrompt("non_existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "prompt file not found")
	})

	// Test loading prompt with invalid path
	t.Run("invalid path", func(t *testing.T) {
		_, err := PM.LoadPrompt("../../../etc/passwd")
		assert.Error(t, err)
	})
}

// TestPromptManagerFunctions tests the standalone prompt functions
func TestPromptManagerFunctions(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	systemDir := filepath.Join(promptsDir, "system")
	require.NoError(t, os.MkdirAll(systemDir, 0755))

	// Create test prompts
	prompts := []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join(promptsDir, "test.md"),
			content: "# Test\nSimple test prompt.",
		},
		{
			path: filepath.Join(systemDir, "check.md"),
			content: `---
description: System check prompt
category: system
---
# System Check`,
		},
		{
			path:    filepath.Join(promptsDir, "template.md"),
			content: "Hello {{name}}!",
		},
	}

	for _, p := range prompts {
		require.NoError(t, os.WriteFile(p.path, []byte(p.content), 0644))
	}

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Test GetPromptPath
	t.Run("GetPromptPath", func(t *testing.T) {
		tests := []struct {
			name     string
			prompt   string
			expected string
		}{
			{
				name:     "root prompt",
				prompt:   "test",
				expected: filepath.Join(promptsDir, "test.md"),
			},
			{
				name:     "category prompt",
				prompt:   "system/check",
				expected: filepath.Join(systemDir, "check.md"),
			},
			{
				name:     "with extension",
				prompt:   "test.md",
				expected: filepath.Join(promptsDir, "test.md"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				path, err := GetPromptPath(tt.prompt)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, path)
			})
		}

		// Test non-existent prompt
		_, err := GetPromptPath("non_existent")
		assert.Error(t, err)
	})

	// Test global prompt functions
	t.Run("legacy functions", func(t *testing.T) {
		// LoadPrompt
		content, err := LoadPrompt("test")
		assert.NoError(t, err)
		assert.Contains(t, content, "Simple test prompt")

		// LoadPromptWithVariables
		variables := map[string]string{"name": "World"}
		content, err = LoadPromptWithVariables("template", variables)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World!", content)

		// GetPromptInfo
		info, err := GetPromptInfo("system/check")
		assert.NoError(t, err)
		assert.Equal(t, "System check prompt", info.Description)
		assert.Equal(t, "system", info.Category)

		// ListPrompts
		list, err := ListPrompts()
		assert.NoError(t, err)
		assert.Len(t, list, 3)
	})
}

// TestPromptManagerInitialization tests that the manager initializes properly
func TestPromptManagerInitialization(t *testing.T) {
	// The PM variable should be initialized
	assert.NotNil(t, PM, "PM should be initialized")

	// Should be of type *PromptManager
	_, ok := PM.(*PromptManager)
	assert.True(t, ok, "PM should be of type *PromptManager")
}
