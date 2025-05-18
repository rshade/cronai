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
	oldExit := osExit
	defer func() {
		os.Args = oldArgs
		osExit = oldExit
	}()

	// Test with no arguments (should show usage)
	os.Args = []string{"search_prompt"}

	// Reset flag.CommandLine for each test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Redirect stdout to avoid test output pollution
	oldStdout := os.Stdout
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = oldStdout
		if err := devNull.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Track if exit was called with correct code
	exitCalled := false
	osExit = func(code int) {
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
		exitCalled = true
		panic("os.Exit called")
	}

	// Expect panic from our mock os.Exit
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic from mocked os.Exit, but didn't get one")
		} else if !exitCalled {
			t.Error("os.Exit was not called")
		}
	}()

	main()
}

func TestMainWithQuery(t *testing.T) {
	// Create test prompt directory
	tmpDir := t.TempDir()
	oldPromptsDir := os.Getenv("CRON_PROMPTS_DIR")
	if err := os.Setenv("CRON_PROMPTS_DIR", tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Setenv("CRON_PROMPTS_DIR", oldPromptsDir); err != nil {
			t.Fatal(err)
		}
	}()

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
			args: []string{"search_prompt", "-content", "-query", "keywords"},
			expectOut: []string{
				"search_test",
			},
		},
		{
			name: "no results",
			args: []string{"search_prompt", "-query", "nonexistent"},
			expectOut: []string{
				"No prompts found matching the search criteria",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original values
			oldArgs := os.Args
			oldStdout := os.Stdout
			oldExit := osExit
			defer func() {
				os.Args = oldArgs
				os.Stdout = oldStdout
				osExit = oldExit
			}()

			// Reset flag.CommandLine
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Set test args
			os.Args = tt.args

			// Capture output
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = w

			// Mock os.Exit
			exitCode := -1
			osExit = func(code int) {
				exitCode = code
				panic("os.Exit called")
			}

			// Run main in a defer/recover to catch the panic from os.Exit
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Expected panic from os.Exit, no action needed
						t.Log("Caught expected panic from os.Exit")
					}
				}()
				main()
			}()

			// Restore stdout and read output
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
			os.Stdout = oldStdout
			buf := make([]byte, 1024*1024)
			n, err := r.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
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
			if !tt.expectErr && exitCode != 0 && exitCode != -1 {
				t.Errorf("Expected exit code 0 for success, got %d", exitCode)
			}
		})
	}
}

func TestMainWithPositionalArgs(t *testing.T) {
	// Store original values
	oldArgs := os.Args
	oldExit := osExit
	defer func() {
		os.Args = oldArgs
		osExit = oldExit
	}()

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Test with positional arguments (query as args instead of flag)
	os.Args = []string{"search_prompt", "test", "query"}

	// Mock os.Exit
	osExit = func(_ int) {
		panic("os.Exit called")
	}

	// Capture output to avoid pollution
	oldStdout := os.Stdout
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = oldStdout
		if err := devNull.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Run main and catch the expected panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic from mocked os.Exit, but didn't get one")
		} else if r != "os.Exit called" {
			t.Errorf("Unexpected panic value: %v", r)
		}
	}()

	main()
}
