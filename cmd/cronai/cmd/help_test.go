package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHelpCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
	}{
		{
			name: "general help",
			args: []string{},
			expectedOutput: []string{
				"CronAI - AI Agent for Scheduled Prompt Execution",
				"Available Commands:",
				"start",
				"run",
				"list",
				"prompt",
				"validate",
				"help",
			},
		},
		{
			name: "start command help",
			args: []string{"start"},
			expectedOutput: []string{
				"Start Command - Launch CronAI Service",
				"Usage:",
				"Flags:",
				"Examples:",
			},
		},
		{
			name: "run command help",
			args: []string{"run"},
			expectedOutput: []string{
				"Run Command - Execute Single AI Task",
				"Required Flags:",
				"Optional Flags:",
				"Special Variables:",
				"Examples:",
			},
		},
		{
			name: "list command help",
			args: []string{"list"},
			expectedOutput: []string{
				"List Command - Display Scheduled Tasks",
				"Usage:",
				"Output Format:",
				"Examples:",
			},
		},
		{
			name: "prompt command help",
			args: []string{"prompt"},
			expectedOutput: []string{
				"Prompt Command - Manage AI Prompts",
				"Subcommands:",
				"list",
				"search",
				"show",
				"preview",
			},
		},
		{
			name: "validate command help",
			args: []string{"validate"},
			expectedOutput: []string{
				"Validate Command - Check Template Syntax",
				"Usage:",
				"Flags:",
				"Examples:",
			},
		},
		{
			name: "unknown command",
			args: []string{"unknown"},
			expectedOutput: []string{
				"Unknown command: unknown",
				"Use 'cronai help' to see available commands",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output by redirecting stdout
			originalStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			os.Stdout = w

			// Run the help command directly
			showHelp(helpCmd, tt.args)

			// Restore stdout and read the output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = originalStdout
			output, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			outputStr := string(output)

			// Check for expected output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got %q", expected, outputStr)
				}
			}
		})
	}
}

// TestHelpCommandPrintsToStdout verifies output is printed correctly
func TestHelpCommandPrintsToStdout(t *testing.T) {
	// This is a simple test to ensure fmt.Print is working as expected
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	fmt.Print("test output")

	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = originalStdout
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	if string(output) != "test output" {
		t.Errorf("Expected 'test output', got %q", string(output))
	}
}
