package cmd

import (
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	// Test that list command is properly configured
	if listCmd.Use != "list" {
		t.Errorf("Expected list command Use to be 'list', got %s", listCmd.Use)
	}

	if listCmd.Short != "List all scheduled AI tasks" {
		t.Errorf("Unexpected short description: %s", listCmd.Short)
	}

	// Verify it's added to root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}
	if !found {
		t.Error("List command not found in root command")
	}
}

func TestListCommandExecution(t *testing.T) {
	// Skip this test for now as it requires mock implementation of the list command
	t.Skip("Skipping list command execution test to avoid file system dependencies")
}

func TestListCommandLongDescription(t *testing.T) {
	expectedStrings := []string{
		"List all scheduled AI tasks",
		"configuration file",
		"schedule, AI model, prompt name, and processor",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(listCmd.Long, expected) {
			t.Errorf("Expected Long description to contain '%s'", expected)
		}
	}
}

func TestListCommandExamples(t *testing.T) {
	expectedExamples := []string{
		"cronai list",
		"cronai list --config=/etc/cronai/production.config",
	}

	for _, expected := range expectedExamples {
		if !strings.Contains(listCmd.Example, expected) {
			t.Errorf("Expected Example to contain '%s'", expected)
		}
	}
}
