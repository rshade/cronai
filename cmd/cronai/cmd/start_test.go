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
		"cronai start --mode cron",
		"cronai start --mode bot",
		"cronai start --mode queue",
	}

	for _, expected := range expectedExamples {
		if !strings.Contains(startCmd.Example, expected) {
			t.Errorf("Expected Example to contain '%s'", expected)
		}
	}
}

func TestValidateMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid cron mode",
			mode:    "cron",
			wantErr: false,
		},
		{
			name:    "bot mode not yet implemented",
			mode:    "bot",
			wantErr: true,
			errMsg:  "mode 'bot' is not yet implemented (coming in future releases)",
		},
		{
			name:    "queue mode not yet implemented",
			mode:    "queue",
			wantErr: true,
			errMsg:  "mode 'queue' is not yet implemented (coming in future releases)",
		},
		{
			name:    "invalid mode",
			mode:    "invalid",
			wantErr: true,
			errMsg:  "invalid mode 'invalid': must be one of: cron, bot, queue",
		},
		{
			name:    "empty mode should be invalid",
			mode:    "",
			wantErr: true,
			errMsg:  "invalid mode '': must be one of: cron, bot, queue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("validateMode() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestModeFlagDefault(t *testing.T) {
	// Reset the flag value to ensure test isolation
	operationMode = ""

	// Check that operationMode has the correct default value after flag parsing
	flag := startCmd.Flags().Lookup("mode")
	if flag == nil {
		t.Fatal("Expected 'mode' flag to exist")
	}

	if flag.DefValue != "cron" {
		t.Errorf("Expected default mode to be 'cron', got %s", flag.DefValue)
	}
}

func TestModeFlagExists(t *testing.T) {
	// Check that the mode flag exists
	flag := startCmd.Flags().Lookup("mode")
	if flag == nil {
		t.Fatal("Expected 'mode' flag to exist on start command")
	}

	// Check flag properties
	if flag.Name != "mode" {
		t.Errorf("Expected flag name to be 'mode', got %s", flag.Name)
	}

	if flag.DefValue != "cron" {
		t.Errorf("Expected default value to be 'cron', got %s", flag.DefValue)
	}

	// Check that the usage string mentions all modes
	expectedUsageStrings := []string{
		"Operation mode",
		"cron",
		"bot",
		"queue",
	}

	for _, expected := range expectedUsageStrings {
		if !strings.Contains(flag.Usage, expected) {
			t.Errorf("Expected flag usage to contain '%s', usage: %s", expected, flag.Usage)
		}
	}
}
