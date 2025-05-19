package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreatePromptWithMetadata tests the CreatePromptWithMetadata function
func TestCreatePromptWithMetadata(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	content := "This is a test prompt content."
	metadata := PromptMetadata{
		Name:        "test_prompt",
		Description: "A test prompt",
		Category:    "testing",
		Tags:        []string{"test", "demo"},
		Author:      "test_user",
		Version:     "1.0.0",
	}

	// Test creating a new prompt
	err := CreatePromptWithMetadata("testing", "test_prompt", &metadata, content)
	assert.NoError(t, err)

	// Verify the file was created
	promptPath := filepath.Join(promptsDir, "test_prompt.md")
	assert.FileExists(t, promptPath)

	// Read the file and verify content
	fileContent, err := os.ReadFile(promptPath)
	require.NoError(t, err)

	// Should contain metadata and content
	assert.Contains(t, string(fileContent), "---")
	assert.Contains(t, string(fileContent), "description: A test prompt")
	assert.Contains(t, string(fileContent), "category: testing")
	assert.Contains(t, string(fileContent), content)

	// Test creating in a category subdirectory
	systemDir := filepath.Join(promptsDir, "system")
	require.NoError(t, os.MkdirAll(systemDir, 0755))

	err = CreatePromptWithMetadata("system", "system_prompt", &metadata, "System prompt content")
	assert.NoError(t, err)

	systemPath := filepath.Join(promptsDir, "system", "system_prompt.md")
	assert.FileExists(t, systemPath)

	// Test overwriting existing file (should fail)
	err = CreatePromptWithMetadata("testing", "test_prompt", &metadata, "New content")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// TestGetPromptInfo tests the GetPromptInfo function
func TestGetPromptInfo(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Create test prompt with metadata
	testPrompt := filepath.Join(promptsDir, "test.md")
	content := `---
description: Test prompt description
category: testing
tags: [test, demo]
author: test_user
version: 1.0.0
---

# Test Prompt
This is a test prompt.`
	require.NoError(t, os.WriteFile(testPrompt, []byte(content), 0644))

	// Create prompt without metadata
	noMetaPrompt := filepath.Join(promptsDir, "no_meta.md")
	noMetaContent := "# No Metadata\nThis prompt has no metadata."
	require.NoError(t, os.WriteFile(noMetaPrompt, []byte(noMetaContent), 0644))

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name            string
		prompt          string
		expectedDesc    string
		expectedCat     string
		expectedHasMeta bool
		expectErr       bool
	}{
		{
			name:            "prompt with metadata",
			prompt:          "test",
			expectedDesc:    "Test prompt description",
			expectedCat:     "testing",
			expectedHasMeta: true,
		},
		{
			name:            "prompt without metadata",
			prompt:          "no_meta",
			expectedDesc:    "This prompt has no metadata.", // Extracted from content
			expectedCat:     "",
			expectedHasMeta: false,
		},
		{
			name:      "non-existent prompt",
			prompt:    "non_existent",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetPromptInfo(tt.prompt)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDesc, info.Description)
				assert.Equal(t, tt.expectedCat, info.Category)
				// HasMetadata field no longer exists in PromptMetadata
			}
		})
	}
}

// TestListPrompts tests the ListPrompts function
func TestListPrompts(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Create category directories
	systemDir := filepath.Join(promptsDir, "system")
	reportsDir := filepath.Join(promptsDir, "reports")
	require.NoError(t, os.MkdirAll(systemDir, 0755))
	require.NoError(t, os.MkdirAll(reportsDir, 0755))

	// Create test prompts
	prompts := []struct {
		path     string
		content  string
		category string
	}{
		{
			path:     filepath.Join(promptsDir, "root_prompt.md"),
			content:  "# Root Prompt",
			category: "",
		},
		{
			path:     filepath.Join(systemDir, "system_check.md"),
			content:  "# System Check",
			category: "system",
		},
		{
			path:     filepath.Join(reportsDir, "monthly_report.md"),
			content:  "# Monthly Report",
			category: "reports",
		},
		{
			path:     filepath.Join(promptsDir, "test.txt"),
			content:  "Not a markdown file",
			category: "",
		},
	}

	for _, p := range prompts {
		require.NoError(t, os.WriteFile(p.path, []byte(p.content), 0644))
	}

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// List prompts
	list, err := ListPrompts()
	require.NoError(t, err)

	// Should find 3 markdown files
	assert.Len(t, list, 3)

	// Verify the prompts
	foundPrompts := make(map[string]bool)
	for _, info := range list {
		foundPrompts[info.Name] = true

		switch info.Name {
		case "root_prompt":
			assert.Equal(t, "", info.Category)
		case "system_check":
			assert.Equal(t, "system", info.Category)
		case "monthly_report":
			assert.Equal(t, "reports", info.Category)
		}
	}

	assert.True(t, foundPrompts["root_prompt"])
	assert.True(t, foundPrompts["system_check"])
	assert.True(t, foundPrompts["monthly_report"])
}

// TestSearchPrompts tests the SearchPrompts function
func TestSearchPrompts(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	systemDir := filepath.Join(promptsDir, "system")
	require.NoError(t, os.MkdirAll(systemDir, 0755))

	// Create test prompts with metadata
	prompts := []struct {
		path    string
		content string
	}{
		{
			path: filepath.Join(promptsDir, "test1.md"),
			content: `---
description: Test prompt for unit testing
category: testing
tags: [test, unit]
author: test_user
---
# Test 1`,
		},
		{
			path: filepath.Join(systemDir, "system_check.md"),
			content: `---
description: System health check
category: system
tags: [monitoring, health]
author: admin
---
# System Check`,
		},
		{
			path: filepath.Join(promptsDir, "report.md"),
			content: `---
description: Monthly report generator
category: reports
tags: [reporting, monthly]
author: test_user
---
# Report`,
		},
		{
			path:    filepath.Join(promptsDir, "no_meta.md"),
			content: "# No Metadata",
		},
	}

	for _, p := range prompts {
		require.NoError(t, os.WriteFile(p.path, []byte(p.content), 0644))
	}

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name          string
		category      string
		tags          []string
		author        string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "search by category",
			category:      "system",
			expectedCount: 1,
			expectedNames: []string{"system_check"},
		},
		{
			name:          "search by author",
			author:        "test_user",
			expectedCount: 2,
			expectedNames: []string{"test1", "report"},
		},
		{
			name:          "search by tag",
			tags:          []string{"test"},
			expectedCount: 1,
			expectedNames: []string{"test1"},
		},
		{
			name:          "search by multiple criteria",
			category:      "system",
			tags:          []string{"monitoring"},
			expectedCount: 1,
			expectedNames: []string{"system_check"},
		},
		{
			name:          "search with no matches",
			category:      "non_existent",
			expectedCount: 0,
		},
		{
			name:          "search all (empty criteria)",
			expectedCount: 3, // Only prompts with metadata
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SearchPrompts now takes query and category only
			results, err := SearchPrompts("", tt.category)
			require.NoError(t, err)

			assert.Len(t, results, tt.expectedCount)

			if tt.expectedNames != nil {
				foundNames := make(map[string]bool)
				for _, info := range results {
					foundNames[info.Name] = true
				}

				for _, expectedName := range tt.expectedNames {
					assert.True(t, foundNames[expectedName], "Expected to find prompt: %s", expectedName)
				}
			}
		})
	}
}

// TestSearchPromptContent tests the SearchPromptContent function
func TestSearchPromptContent(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	// Create test prompts
	prompts := []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join(promptsDir, "test1.md"),
			content: "# Test Prompt\nThis is a test for searching content.",
		},
		{
			path:    filepath.Join(promptsDir, "test2.md"),
			content: "# Another Test\nThis contains the word searching.",
		},
		{
			path:    filepath.Join(promptsDir, "no_match.md"),
			content: "# No Match\nThis doesn't contain the target word.",
		},
	}

	for _, p := range prompts {
		require.NoError(t, os.WriteFile(p.path, []byte(p.content), 0644))
	}

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "search for 'searching'",
			query:         "searching",
			expectedCount: 2,
			expectedNames: []string{"test1", "test2"},
		},
		{
			name:          "search for 'Test Prompt'",
			query:         "Test Prompt",
			expectedCount: 1,
			expectedNames: []string{"test1"},
		},
		{
			name:          "search with no matches",
			query:         "non_existent_content",
			expectedCount: 0,
		},
		{
			name:          "case sensitive search",
			query:         "test",
			expectedCount: 2,
			expectedNames: []string{"test1", "test2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := SearchPromptContent(tt.query, "")
			require.NoError(t, err)

			assert.Len(t, results, tt.expectedCount)

			if tt.expectedNames != nil {
				foundNames := make(map[string]bool)
				for _, info := range results {
					foundNames[info.Name] = true
				}

				for _, expectedName := range tt.expectedNames {
					assert.True(t, foundNames[expectedName], "Expected to find prompt: %s", expectedName)
				}
			}
		})
	}
}

// TestLoadPromptWithIncludes tests the LoadPromptWithIncludes function
func TestLoadPromptWithIncludes(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	promptsDir := filepath.Join(tempDir, "cron_prompts")
	templatesDir := filepath.Join(promptsDir, "templates")
	require.NoError(t, os.MkdirAll(templatesDir, 0755))

	// Create included file
	headerFile := filepath.Join(templatesDir, "header.md")
	headerContent := "## Common Header\nThis is a shared header."
	require.NoError(t, os.WriteFile(headerFile, []byte(headerContent), 0644))

	// Create footer file
	footerFile := filepath.Join(templatesDir, "footer.md")
	footerContent := "---\nCommon footer"
	require.NoError(t, os.WriteFile(footerFile, []byte(footerContent), 0644))

	// Create main prompt with includes
	mainPrompt := filepath.Join(promptsDir, "main.md")
	mainContent := `# Main Prompt

{{include templates/header.md}}

Main content goes here.

{{include templates/footer.md}}`
	require.NoError(t, os.WriteFile(mainPrompt, []byte(mainContent), 0644))

	// Create prompt with nested includes (3 levels)
	level3 := filepath.Join(templatesDir, "level3.md")
	require.NoError(t, os.WriteFile(level3, []byte("Level 3 content"), 0644))

	level2 := filepath.Join(templatesDir, "level2.md")
	require.NoError(t, os.WriteFile(level2, []byte("Level 2: {{include templates/level3.md}}"), 0644))

	level1 := filepath.Join(templatesDir, "level1.md")
	require.NoError(t, os.WriteFile(level1, []byte("Level 1: {{include templates/level2.md}}"), 0644))

	nestedPrompt := filepath.Join(promptsDir, "nested.md")
	nestedContent := "# Nested\n{{include templates/level1.md}}"
	require.NoError(t, os.WriteFile(nestedPrompt, []byte(nestedContent), 0644))

	// Create prompt with circular reference
	circularA := filepath.Join(templatesDir, "circular_a.md")
	circularB := filepath.Join(templatesDir, "circular_b.md")
	require.NoError(t, os.WriteFile(circularA, []byte("A: {{include templates/circular_b.md}}"), 0644))
	require.NoError(t, os.WriteFile(circularB, []byte("B: {{include templates/circular_a.md}}"), 0644))

	circularPrompt := filepath.Join(promptsDir, "circular.md")
	require.NoError(t, os.WriteFile(circularPrompt, []byte("# Circular\n{{include templates/circular_a.md}}"), 0644))

	// Change working directory for the test
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	tests := []struct {
		name      string
		prompt    string
		expected  string
		expectErr bool
		errMsg    string
	}{
		{
			name:   "load with includes",
			prompt: "main",
			expected: `# Main Prompt

## Common Header
This is a shared header.

Main content goes here.

---
Common footer`,
		},
		{
			name:   "load with nested includes",
			prompt: "nested",
			expected: `# Nested
Level 1: Level 2: Level 3 content`,
		},
		{
			name:      "max recursion depth exceeded",
			prompt:    "circular",
			expectErr: true,
			errMsg:    "maximum recursion depth",
		},
		{
			name:      "non-existent include",
			prompt:    "main",
			expectErr: true,
			errMsg:    "failed to include",
		},
	}

	// Test successful include
	t.Run(tests[0].name, func(t *testing.T) {
		content, err := LoadPromptWithIncludes(tests[0].prompt)
		require.NoError(t, err)
		assert.Equal(t, tests[0].expected, content)
	})

	// Test nested includes
	t.Run(tests[1].name, func(t *testing.T) {
		content, err := LoadPromptWithIncludes(tests[1].prompt)
		require.NoError(t, err)
		assert.Equal(t, tests[1].expected, content)
	})

	// Test circular reference
	t.Run(tests[2].name, func(t *testing.T) {
		_, err := LoadPromptWithIncludes(tests[2].prompt)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), tests[2].errMsg)
	})

	// Test missing include
	t.Run("missing include", func(t *testing.T) {
		// Create prompt with non-existent include
		badPrompt := filepath.Join(promptsDir, "bad.md")
		badContent := "# Bad\n{{include non_existent.md}}"
		require.NoError(t, os.WriteFile(badPrompt, []byte(badContent), 0644))

		_, err := LoadPromptWithIncludes("bad")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to include")
	})
}

// TestValidatePromptTemplate tests the ValidatePromptTemplate function
func TestValidatePromptTemplate(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
		errMsg    string
	}{
		{
			name:    "valid template",
			content: "Hello {{name}}, welcome to {{place}}!",
		},
		{
			name:    "nested braces",
			content: "Result: {{if .condition}}{{.value}}{{end}}",
		},
		{
			name:      "unclosed brace",
			content:   "Hello {{name, welcome!",
			expectErr: true,
			errMsg:    "unclosed '{{' at position",
		},
		{
			name:      "unopened brace",
			content:   "Hello name}}, welcome!",
			expectErr: true,
			errMsg:    "unexpected '}}' at position",
		},
		{
			name:      "mismatched braces",
			content:   "{{start}} content }}{{end",
			expectErr: true,
			errMsg:    "unclosed '{{' at position",
		},
		{
			name:    "empty template",
			content: "",
		},
		{
			name:    "no templates",
			content: "Just plain text without any templates.",
		},
		{
			name:    "multiple valid templates",
			content: "{{header}}\n\nContent: {{content}}\n\n{{footer}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePromptTemplate(tt.content)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
