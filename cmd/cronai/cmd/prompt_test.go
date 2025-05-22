// Package cmd implements the command-line interface for the cronai application.
// It provides commands for managing AI prompts, including listing, searching,
// showing details, and previewing prompts with variables.
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/rshade/cronai/internal/prompt"
	"github.com/spf13/cobra"
)

func TestPromptCommand(t *testing.T) {
	// Test that prompt command is properly configured
	if promptCmd.Use != "prompt" {
		t.Errorf("Expected prompt command Use to be 'prompt', got %s", promptCmd.Use)
	}

	if promptCmd.Short != "Manage AI prompt templates" {
		t.Errorf("Unexpected short description: %s", promptCmd.Short)
	}

	// Verify subcommands exist
	subcommands := []string{"list", "search", "show", "preview"}
	for _, subCmd := range subcommands {
		found := false
		for _, cmd := range promptCmd.Commands() {
			if cmd.Name() == subCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand '%s' not found in prompt command", subCmd)
		}
	}
}

func TestPromptListCommandReal(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatalf("Failed to set CRON_PROMPTS_DIR: %v", err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldDir); err != nil {
			t.Errorf("Failed to restore CRON_PROMPTS_DIR: %v", err)
		}
	}()

	// Create test category directories
	testDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	reportDir := filepath.Join(tmpDir, "reports")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		t.Fatalf("Failed to create reports directory: %v", err)
	}

	// Create test prompts
	testPrompt1 := `---
title: Test Prompt
description: A test prompt
category: test
---
This is a test prompt`

	testPrompt2 := `---
title: Report Template
description: A report template
category: reports
---
This is a report template`

	if err := os.WriteFile(filepath.Join(testDir, "test_prompt.md"), []byte(testPrompt1), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	if err := os.WriteFile(filepath.Join(reportDir, "report_template.md"), []byte(testPrompt2), 0644); err != nil {
		t.Fatalf("Failed to create report prompt: %v", err)
	}

	// Capture output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run the actual prompt list command
	promptListCmd.Run(nil, []string{})

	// Restore stdout and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = old
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := string(out)

	// Check output contains both categories and prompts
	expectedOutput := []string{
		"reports:",
		"report_template",
		"test:",
		"test_prompt",
	}

	for _, expected := range expectedOutput {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestPromptListCommand(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatalf("Failed to set CRON_PROMPTS_DIR: %v", err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldDir); err != nil {
			t.Errorf("Failed to restore CRON_PROMPTS_DIR: %v", err)
		}
	}()

	// Create test category directory
	testDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test prompt
	testPrompt := `---
title: Test Prompt
description: A test prompt
category: test
---
This is a test prompt`

	if err := os.WriteFile(filepath.Join(testDir, "test_prompt.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "list all prompts",
			args: []string{},
			expectedOutput: []string{
				"CATEGORY",
				"NAME",
				"DESCRIPTION",
				"PATH",
			},
		},
		{
			name: "list by category",
			args: []string{"--category", "test"},
			expectedOutput: []string{
				"test",
				"test_prompt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Create a fresh list command for each test
			cmd := &cobra.Command{
				Use:   "list",
				Short: "List available AI prompts",
				Run: func(cmd *cobra.Command, _ []string) {
					prompts, err := prompt.ListPrompts()
					if err != nil {
						fmt.Printf("Error listing prompts: %v\n", err)
						return
					}

					// Get category flag
					category, err := cmd.Flags().GetString("category")
					if err != nil {
						fmt.Printf("Error getting category flag: %v\n", err)
						return
					}

					// Filter by category if specified
					if category != "" {
						var filtered []prompt.Info
						for _, p := range prompts {
							if strings.EqualFold(p.Category, category) {
								filtered = append(filtered, p)
							}
						}
						prompts = filtered
					}

					if len(prompts) == 0 {
						fmt.Println("No prompts found")
						return
					}

					// Format output using tabwriter
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
					if _, err := fmt.Fprintln(w, "CATEGORY\tNAME\tDESCRIPTION\tPATH"); err != nil {
						fmt.Printf("Error writing to tabwriter: %v\n", err)
						return
					}
					for _, p := range prompts {
						if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Category, p.Name, p.Description, p.Path); err != nil {
							fmt.Printf("Error writing to tabwriter: %v\n", err)
							return
						}
					}
					if err := w.Flush(); err != nil {
						fmt.Printf("Error flushing tabwriter: %v\n", err)
					}
				},
			}

			// Add category flag
			cmd.Flags().StringP("category", "c", "", "Filter prompts by category")
			cmd.SetArgs(tt.args)

			// Run the command
			if err := cmd.Execute(); err != nil {
				if !tt.expectError {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.expectError {
				t.Error("Expected error but got none")
			}

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			out, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := string(out)

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPromptSearchCommand(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatalf("Failed to set CRON_PROMPTS_DIR: %v", err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldDir); err != nil {
			t.Errorf("Failed to restore CRON_PROMPTS_DIR: %v", err)
		}
	}()

	// Create test category directory
	testDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test prompt
	testPrompt := `---
title: Search Test
description: A searchable test prompt
category: test
---
This prompt contains searchable content with keywords.`

	if err := os.WriteFile(filepath.Join(testDir, "search_test.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	// Also add a test prompt at root level
	rootTestPrompt := `---
title: Root Test
description: A root level test prompt
category: root
---
This prompt at root level.`

	if err := os.WriteFile(filepath.Join(tmpDir, "root_test.md"), []byte(rootTestPrompt), 0644); err != nil {
		t.Fatalf("Failed to create root test prompt: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "search with query",
			args: []string{"searchable"},
			expectedOutput: []string{
				"CATEGORY",
				"NAME",
				"search_test",
			},
		},
		{
			name: "search in content",
			args: []string{"keywords", "--content"},
			expectedOutput: []string{
				"search_test",
			},
		},
		{
			name: "search with category filter",
			args: []string{"test", "--category", "test"},
			expectedOutput: []string{
				"test_prompt",
				"search_test",
			},
		},
		{
			name: "no results",
			args: []string{"nonexistent"},
			expectedOutput: []string{
				"No prompts found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Create a fresh search command for each test
			cmd := &cobra.Command{
				Use:   "search [query]",
				Short: "Search for prompts",
				Long:  `Search for prompts by name, description, or content.`,
				Run: func(cmd *cobra.Command, args []string) {
					// Get query from args if not provided as flag
					searchQuery, err := cmd.Flags().GetString("query")
					if err != nil {
						fmt.Printf("Error getting query flag: %v\n", err)
						return
					}
					if len(args) > 0 && searchQuery == "" {
						searchQuery = args[0]
					}

					var prompts []prompt.Info

					searchContent, err := cmd.Flags().GetBool("content")
					if err != nil {
						fmt.Printf("Error getting content flag: %v\n", err)
						return
					}
					category, err := cmd.Flags().GetString("category")
					if err != nil {
						fmt.Printf("Error getting category flag: %v\n", err)
						return
					}

					if searchContent {
						// Search in prompt content
						prompts, err = prompt.SearchPromptContent(searchQuery, category)
					} else {
						// Search in metadata only
						prompts, err = prompt.SearchPrompts(searchQuery, category)
					}

					if err != nil {
						fmt.Printf("Error searching prompts: %v\n", err)
						return
					}

					if len(prompts) == 0 {
						fmt.Println("No prompts found matching the search criteria")
						return
					}

					// Format output using tabwriter
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
					if _, err := fmt.Fprintln(w, "CATEGORY\tNAME\tDESCRIPTION\tPATH"); err != nil {
						fmt.Printf("Error writing to tabwriter: %v\n", err)
						return
					}
					for _, p := range prompts {
						if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Category, p.Name, p.Description, p.Path); err != nil {
							fmt.Printf("Error writing to tabwriter: %v\n", err)
							return
						}
					}
					if err := w.Flush(); err != nil {
						fmt.Printf("Error flushing tabwriter: %v\n", err)
					}
				},
			}

			// Add flags
			cmd.Flags().StringP("category", "c", "", "Filter prompts by category")
			cmd.Flags().StringP("query", "q", "", "Search query")
			cmd.Flags().BoolP("content", "t", false, "Search in prompt content")
			cmd.SetArgs(tt.args)

			// Run the command
			if err := cmd.Execute(); err != nil {
				if !tt.expectError {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.expectError {
				t.Error("Expected error but got none")
			}

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			out, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := string(out)

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPromptShowCommand(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatalf("Failed to set CRON_PROMPTS_DIR: %v", err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldDir); err != nil {
			t.Errorf("Failed to restore CRON_PROMPTS_DIR: %v", err)
		}
	}()

	// Create test category directory
	testDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test prompt
	testPrompt := `---
title: Show Test
description: A test prompt for show command
author: Test Author
version: 1.0.0
category: test
tags: [test, example]
variables:
  - name: var1
    description: First variable
  - name: var2
    description: Second variable
---
This is the content of the test prompt.`

	// Write both to the test directory and directly in tmpDir to cover different search paths
	if err := os.WriteFile(filepath.Join(testDir, "show_test.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt in test dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "show_test.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt in root dir: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name:           "show prompt",
			args:           []string{"show_test"},
			expectError:    false, // Should work now that we've added the file
			expectedOutput: []string{"Name:", "Description:", "Category:"},
		},
		{
			name:           "show prompt with variables",
			args:           []string{"show_test", "--vars"},
			expectError:    false, // Should work now that we've added the file
			expectedOutput: []string{"Variables:", "NAME", "DESCRIPTION", "var1", "var2"},
		},
		{
			name:        "non-existent prompt",
			args:        []string{"nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh show command for each test
			cmd := &cobra.Command{
				Use:   "show [promptName]",
				Short: "Show prompt details",
				Long:  `Show detailed information about a specific prompt.`,
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					promptName := args[0]

					// Get prompt metadata
					metadata, err := prompt.GetPromptMetadata(promptName)
					if err != nil {
						return fmt.Errorf("error getting prompt metadata: %w", err)
					}

					// Display prompt details
					fmt.Printf("Name: %s\n", metadata.Name)
					fmt.Printf("Description: %s\n", metadata.Description)
					fmt.Printf("Author: %s\n", metadata.Author)
					fmt.Printf("Version: %s\n", metadata.Version)
					fmt.Printf("Category: %s\n", metadata.Category)

					if len(metadata.Tags) > 0 {
						fmt.Printf("Tags: %s\n", strings.Join(metadata.Tags, ", "))
					}

					// Show variables if requested
					showVars, err := cmd.Flags().GetBool("vars")
					if err != nil {
						return fmt.Errorf("error getting vars flag: %w", err)
					}
					if showVars && len(metadata.Variables) > 0 {
						fmt.Println("\nVariables:")
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
						if _, err := fmt.Fprintln(w, "NAME\tDESCRIPTION"); err != nil {
							return fmt.Errorf("error writing to tabwriter: %w", err)
						}
						for _, v := range metadata.Variables {
							if _, err := fmt.Fprintf(w, "%s\t%s\n", v.Name, v.Description); err != nil {
								return fmt.Errorf("error writing to tabwriter: %w", err)
							}
						}
						if err := w.Flush(); err != nil {
							return fmt.Errorf("error flushing tabwriter: %w", err)
						}
					}

					// Get the prompt content
					content, err := prompt.LoadPrompt(promptName)
					if err != nil {
						return fmt.Errorf("error loading prompt content: %w", err)
					}

					fmt.Println("\nContent:")
					fmt.Println("----------------------------------------")
					fmt.Println(content)
					fmt.Println("----------------------------------------")
					return nil
				},
			}

			// Add the vars flag
			cmd.Flags().BoolP("vars", "v", false, "Show prompt variables")
			cmd.SetArgs(tt.args)

			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Run the command
			if err := cmd.Execute(); err != nil {
				if !tt.expectError {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.expectError {
				t.Error("Expected error but got none")
			}

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			out, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := string(out)

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}

func TestPromptPreviewCommand(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatalf("Failed to set CRON_PROMPTS_DIR: %v", err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldDir); err != nil {
			t.Errorf("Failed to restore CRON_PROMPTS_DIR: %v", err)
		}
	}()

	// Create general category directory
	generalDir := filepath.Join(tmpDir, "general")
	if err := os.MkdirAll(generalDir, 0755); err != nil {
		t.Fatalf("Failed to create general directory: %v", err)
	}

	// Create test prompt
	testPrompt := `---
title: Preview Test
description: A test prompt for preview command
variables:
  - name: name
    description: User name
  - name: role
    description: User role
---
Hello {{.Variables.name}}, your role is {{.Variables.role}}.
Current date: {{.Variables.CURRENT_DATE}}`

	// Write both to general directory and directly in tmpDir to cover different search paths
	if err := os.WriteFile(filepath.Join(generalDir, "preview_test.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt in general dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "preview_test.md"), []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt in root dir: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "preview without variables",
			args: []string{"preview_test"},
			expectedOutput: []string{
				"Preview:",
				"Hello {{.Variables.name}}",
				"Current date: 2025-05-12",
			},
		},
		{
			name: "preview with variables",
			args: []string{"preview_test", "--vars", "name=John,role=Developer"},
			expectedOutput: []string{
				"Preview:",
				"Hello John, your role is Developer",
				"Current date: 2025-05-12",
			},
		},
		{
			name:        "non-existent prompt",
			args:        []string{"nonexistent"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh prompt preview command for each test
			promptPreviewCmd := &cobra.Command{
				Use:   "preview [promptName]",
				Short: "Preview a prompt with variables",
				Long:  `Preview a prompt with variables and includes processed.`,
				Args:  cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					promptName := args[0]

					// Get variables from command line
					variablesFlag, err := cmd.Flags().GetString("vars")
					if err != nil {
						return fmt.Errorf("error getting vars flag: %w", err)
					}
					variables := make(map[string]string)

					// Parse variables from comma-separated list of key=value pairs
					if variablesFlag != "" {
						pairs := strings.Split(variablesFlag, ",")
						for _, pair := range pairs {
							kv := strings.SplitN(pair, "=", 2)
							if len(kv) == 2 {
								variables[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
							}
						}
					}

					// Add special variables
					variables["CURRENT_DATE"] = "2025-05-12"
					variables["CURRENT_TIME"] = "12:00:00"
					variables["CURRENT_DATETIME"] = "2025-05-12 12:00:00"

					// Load prompt with variables
					content, err := prompt.LoadPromptWithVariables(promptName, variables)
					if err != nil {
						return fmt.Errorf("error loading prompt with variables: %w", err)
					}

					fmt.Println("Preview:")
					fmt.Println("----------------------------------------")
					fmt.Println(content)
					fmt.Println("----------------------------------------")
					return nil
				},
			}

			// Add the vars flag
			promptPreviewCmd.Flags().String("vars", "", "Variables in format 'key1=value1,key2=value2'")

			// Capture output
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Set args
			promptPreviewCmd.SetArgs(tt.args)

			// Run the command
			if err := promptPreviewCmd.Execute(); err != nil {
				if !tt.expectError {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.expectError {
				t.Error("Expected error but got none")
			}

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = old
			out, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := string(out)

			// Check output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}
		})
	}
}
