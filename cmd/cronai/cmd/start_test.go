package cmd

import (
	"strings"
	"testing"
)

func TestStartCommand(t *testing.T) {
	// Test that start command is properly configured
	if startCmd.Use != "start" {
		t.Errorf("Expected start command Use to be 'start', got %s", startCmd.Use)
	}

	if startCmd.Short != "Start the CronAI service" {
		t.Errorf("Unexpected short description: %s", startCmd.Short)
	}

	// Verify it's added to root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "start" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Start command not found in root command")
	}
}

func TestStartCommandExecution(t *testing.T) {
	// Skip this test since it depends on the actual implementation
	// and we don't want to execute the service during tests
	t.Skip("Skipping start command execution test to avoid service execution")
}

func TestStartCommandLongDescription(t *testing.T) {
	expectedStrings := []string{
		"Start the CronAI agent service",
		"Configuration Format:",
		"schedule model prompt processor",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(startCmd.Long, expected) {
			t.Errorf("Expected Long description to contain '%s'", expected)
		}
	}
}

func TestStartCommandExamples(t *testing.T) {
	expectedExamples := []string{
		"cronai start",
		"cronai start --config=/etc/cronai/production.config",
		"sudo systemctl start cronai",
	}

	for _, expected := range expectedExamples {
		if !strings.Contains(startCmd.Example, expected) {
			t.Errorf("Expected Example to contain '%s'", expected)
		}
	}
}
