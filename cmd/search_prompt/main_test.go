package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainUsage(t *testing.T) {
	// Store original values
	oldArgs := os.Args

	// Test with no arguments (should show usage)
	os.Args = []string{"search_prompt"}

	// Reset flag.CommandLine for each test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Redirect stdout to avoid test output pollution
	oldStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stdout = devNull

	// We expect this to exit, so we'll use a goroutine
	done := make(chan bool)
	go func() {
		defer func() {
			// Recover from os.Exit
			if r := recover(); r != nil {
				done <- true
			}
		}()
		// Mock os.Exit to prevent actual exit during tests
		osExit = func(code int) {
			if code != 1 {
				t.Errorf("Expected exit code 1, got %d", code)
			}
			panic("os.Exit called")
		}
		main()
		done <- false
	}()

	// Wait for the goroutine
	<-done

	// Restore everything
	os.Stdout = oldStdout
	os.Args = oldArgs
	osExit = os.Exit
}

func TestMainWithQuery(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldPromptsDir := os.Getenv("CRON_PROMPTS_DIR")
	os.Setenv("CRON_PROMPTS_DIR", tmpDir)
	defer os.Setenv("CRON_PROMPTS_DIR", oldPromptsDir)

	// Create a test prompt
	testPrompt := `---
title: Test Prompt
description: A test prompt for search
category: test
---
This is test content.`

	promptPath := filepath.Join(tmpDir, "test_prompt.md")
	if err := os.WriteFile(promptPath, []byte(testPrompt), 0644); err != nil {
		t.Fatalf("Failed to create test prompt: %v", err)
	}

	tests := []struct {
		name      string
		args      []string
		expectOut []string
		expectErr bool
	}{
		{
			name: "search by query",
			args: []string{"search_prompt", "-query", "test"},
			expectOut: []string{
				"CATEGORY",
				"NAME",
				"DESCRIPTION",
				"PATH",
			},
		},
		{
			name: "search by category",
			args: []string{"search_prompt", "-category", "test"},
			expectOut: []string{
				"test_prompt",
			},
		},
		{
			name: "search in content",
			args: []string{"search_prompt", "-content", "-query", "content"},
			expectOut: []string{
				"test_prompt",
			},
		},
		{
			name: "no results",
			args: []string{"search_prompt", "-query", "nonexistent"},
			expectOut: []string{
				"No prompts found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original args
			oldArgs := os.Args
			oldStdout := os.Stdout

			// Reset flag.CommandLine
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Set test args
			os.Args = tt.args

			// Capture output
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Mock os.Exit
			exitCode := -1
			osExit = func(code int) {
				exitCode = code
				panic("os.Exit called")
			}

			// Run main in a goroutine to catch the panic from os.Exit
			done := make(chan bool)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						done <- true
					}
				}()
				main()
				done <- false
			}()

			// Wait for completion
			<-done

			// Restore stdout and read output
			w.Close()
			os.Stdout = oldStdout
			buf := make([]byte, 1024*1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Check output
			for _, expected := range tt.expectOut {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check exit code
			if tt.expectErr && exitCode != 1 {
				t.Errorf("Expected exit code 1 for error, got %d", exitCode)
			}
			if !tt.expectErr && exitCode != 0 {
				t.Errorf("Expected exit code 0 for success, got %d", exitCode)
			}

			// Restore
			os.Args = oldArgs
			osExit = os.Exit
		})
	}
}

// Mock os.Exit for testing
var osExit = os.Exit

func TestMainWithPositionalArgs(t *testing.T) {
	// Store original values
	oldArgs := os.Args

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Test with positional arguments (query as args instead of flag)
	os.Args = []string{"search_prompt", "test", "query"}

	// We expect this to search for "test query"
	// This is a simplified test just to ensure the code path works
	// In a real test, we'd check the actual search behavior

	// Mock os.Exit
	osExit = func(code int) {
		// We expect it to exit with 1 since there are no prompts
		if code != 1 && code != 0 {
			t.Errorf("Unexpected exit code: %d", code)
		}
	}

	// Capture output to avoid pollution
	oldStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stdout = devNull

	// Run in a defer/recover to catch panic from os.Exit mock
	defer func() {
		if r := recover(); r == nil {
			// If we didn't panic, that's also OK (means main() completed)
		}
		// Restore
		os.Stdout = oldStdout
		os.Args = oldArgs
		osExit = os.Exit
	}()

	// This should combine the args into a query
	main()
}
