package prompt

import (
	"reflect"
	"testing"
)

func TestExtractMetadata(t *testing.T) {
	// Test content with metadata
	contentWithMetadata := `---
name: Test Prompt
description: A test prompt for unit testing
author: Test Author
version: 1.0
category: testing
tags: test, unit, metadata
variables:
  - name: testVar1
    description: Test variable 1
  - name: testVar2
    description: Test variable 2
---

# Test Prompt

This is the actual prompt content.`

	expectedMetadata := &Metadata{
		Name:        "Test Prompt",
		Description: "A test prompt for unit testing",
		Author:      "Test Author",
		Version:     "1.0",
		Category:    "testing",
		Tags:        []string{"test", "unit", "metadata"},
		Variables: []Variable{
			{Name: "testVar1", Description: "Test variable 1"},
			{Name: "testVar2", Description: "Test variable 2"},
		},
		Path: "test_path",
	}

	expectedContent := `# Test Prompt

This is the actual prompt content.`

	metadata, content, err := ExtractMetadata(contentWithMetadata, "test_path")
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	// Debug output for first test
	t.Logf("FIRST TEST - Found variables: %+v", metadata.Variables)
	t.Logf("FIRST TEST - Number of variables: %d", len(metadata.Variables))
	for i, v := range metadata.Variables {
		t.Logf("FIRST TEST - Variable %d: %s - %s", i, v.Name, v.Description)
	}
	t.Logf("FIRST TEST - Expected variables: %+v", expectedMetadata.Variables)
	t.Logf("FIRST TEST - Number of expected variables: %d", len(expectedMetadata.Variables))

	// Compare metadata
	if metadata.Name != expectedMetadata.Name {
		t.Errorf("Expected name %q, got %q", expectedMetadata.Name, metadata.Name)
	}
	if metadata.Description != expectedMetadata.Description {
		t.Errorf("Expected description %q, got %q", expectedMetadata.Description, metadata.Description)
	}
	if metadata.Author != expectedMetadata.Author {
		t.Errorf("Expected author %q, got %q", expectedMetadata.Author, metadata.Author)
	}
	if metadata.Version != expectedMetadata.Version {
		t.Errorf("Expected version %q, got %q", expectedMetadata.Version, metadata.Version)
	}
	if metadata.Category != expectedMetadata.Category {
		t.Errorf("Expected category %q, got %q", expectedMetadata.Category, metadata.Category)
	}
	if metadata.Path != expectedMetadata.Path {
		t.Errorf("Expected path %q, got %q", expectedMetadata.Path, metadata.Path)
	}

	// Compare tags
	if !reflect.DeepEqual(metadata.Tags, expectedMetadata.Tags) {
		t.Errorf("Expected tags %v, got %v", expectedMetadata.Tags, metadata.Tags)
	}

	// Compare variables
	if len(metadata.Variables) != len(expectedMetadata.Variables) {
		t.Errorf("Expected %d variables, got %d", len(expectedMetadata.Variables), len(metadata.Variables))
	} else {
		for i, v := range metadata.Variables {
			if v.Name != expectedMetadata.Variables[i].Name {
				t.Errorf("Variable %d: expected name %q, got %q", i, expectedMetadata.Variables[i].Name, v.Name)
			}
			if v.Description != expectedMetadata.Variables[i].Description {
				t.Errorf("Variable %d: expected description %q, got %q", i, expectedMetadata.Variables[i].Description, v.Description)
			}
		}
	}

	// Compare content
	if content != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, content)
	}

	// Test content without metadata
	contentWithoutMetadata := `# Test Prompt

This is a prompt without metadata.`

	metadata, content, err = ExtractMetadata(contentWithoutMetadata, "test_path")
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	// Check that we get an empty metadata object
	if metadata.Name != "" || metadata.Description != "" || metadata.Author != "" ||
		metadata.Version != "" || metadata.Category != "" || len(metadata.Tags) != 0 ||
		len(metadata.Variables) != 0 {
		t.Errorf("Expected empty metadata for content without metadata section, got %+v", metadata)
	}

	// Check that path is still set
	if metadata.Path != "test_path" {
		t.Errorf("Expected path %q, got %q", "test_path", metadata.Path)
	}

	// Check that content is unchanged
	if content != contentWithoutMetadata {
		t.Errorf("Expected content %q, got %q", contentWithoutMetadata, content)
	}

	// Test with valid metadata
	validContent := `---
name: Test Prompt
description: A test prompt
author: Test Author
version: 1.0
category: test
tags: [test, example]
variables:
  - name: var1
    description: First variable
  - name: var2
    description: Second variable
extends: base
---
Content here`
	expectedMetadata = &Metadata{
		Name:        "Test Prompt",
		Description: "A test prompt",
		Author:      "Test Author",
		Version:     "1.0",
		Category:    "test",
		Tags:        []string{"test", "example"},
		Variables: []Variable{
			{
				Name:        "var1",
				Description: "First variable",
			},
			{
				Name:        "var2",
				Description: "Second variable",
			},
		},
		Extends: "base",
	}

	// Test the valid content
	metadata, content, err = ExtractMetadata(validContent, "test_path")
	if err != nil {
		t.Fatalf("ExtractMetadata failed: %v", err)
	}

	// Debug output for second test
	t.Logf("SECOND TEST - Found variables: %+v", metadata.Variables)
	t.Logf("SECOND TEST - Number of variables: %d", len(metadata.Variables))
	for i, v := range metadata.Variables {
		t.Logf("SECOND TEST - Variable %d: %s - %s", i, v.Name, v.Description)
	}
	t.Logf("SECOND TEST - Expected variables: %+v", expectedMetadata.Variables)
	t.Logf("SECOND TEST - Number of expected variables: %d", len(expectedMetadata.Variables))

	// Compare metadata
	if metadata.Name != expectedMetadata.Name {
		t.Errorf("Expected name %q, got %q", expectedMetadata.Name, metadata.Name)
	}
	if metadata.Description != expectedMetadata.Description {
		t.Errorf("Expected description %q, got %q", expectedMetadata.Description, metadata.Description)
	}
	if metadata.Author != expectedMetadata.Author {
		t.Errorf("Expected author %q, got %q", expectedMetadata.Author, metadata.Author)
	}
	if metadata.Version != expectedMetadata.Version {
		t.Errorf("Expected version %q, got %q", expectedMetadata.Version, metadata.Version)
	}
	if metadata.Category != expectedMetadata.Category {
		t.Errorf("Expected category %q, got %q", expectedMetadata.Category, metadata.Category)
	}
	if metadata.Extends != expectedMetadata.Extends {
		t.Errorf("Expected extends %q, got %q", expectedMetadata.Extends, metadata.Extends)
	}

	// Compare tags
	if !reflect.DeepEqual(metadata.Tags, expectedMetadata.Tags) {
		t.Errorf("Expected tags %v, got %v", expectedMetadata.Tags, metadata.Tags)
	}

	// Compare variables
	if len(metadata.Variables) != len(expectedMetadata.Variables) {
		t.Errorf("Expected %d variables, got %d", len(expectedMetadata.Variables), len(metadata.Variables))
	} else {
		for i, v := range metadata.Variables {
			if v.Name != expectedMetadata.Variables[i].Name {
				t.Errorf("Variable %d: expected name %q, got %q", i, expectedMetadata.Variables[i].Name, v.Name)
			}
			if v.Description != expectedMetadata.Variables[i].Description {
				t.Errorf("Variable %d: expected description %q, got %q", i, expectedMetadata.Variables[i].Description, v.Description)
			}
		}
	}

	// Compare content
	if content != "Content here" {
		t.Errorf("Expected content %q, got %q", "Content here", content)
	}
}
