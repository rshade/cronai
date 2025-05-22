// Package cmd implements the command line interface for CronAI.
package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Test that version command is properly configured
	if versionCmd.Use != "version" {
		t.Errorf("Expected version command Use to be 'version', got %s", versionCmd.Use)
	}

	if versionCmd.Short != "Print the version number of CronAI" {
		t.Errorf("Unexpected short description: %s", versionCmd.Short)
	}

	if !strings.Contains(versionCmd.Long, "Print the version number of CronAI") {
		t.Errorf("Expected Long description to contain version info: %s", versionCmd.Long)
	}
}

func TestVersionCommandExecution(t *testing.T) {
	// Set a test version
	oldVersion := Version
	Version = "test-version-1.2.3"
	defer func() { Version = oldVersion }()

	// Capture output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run the version command
	versionCmd.Run(nil, []string{})

	// Restore stdout and read output
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Check output
	expectedOutput := "CronAI version test-version-1.2.3"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain '%s', got: %s", expectedOutput, output)
	}
}

func TestVersionCommandInit(t *testing.T) {
	// Test that version command is added to root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Version command not found in root command")
	}
}
