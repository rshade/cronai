package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	// Capture output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Temporarily replace os.Args
	oldArgs := os.Args
	os.Args = []string{"cronai", "--help"}
	defer func() { os.Args = oldArgs }()

	// Test execute with help flag
	Execute()

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe writer: %v", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Check for expected output from Long description
	expectedStrings := []string{
		"CronAI - Your Automated AI Assistant",
		"Schedule AI prompts using cron syntax",
		"Support for multiple AI models",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestRootCmd(t *testing.T) {
	// Test root command setup
	if rootCmd.Use != "cronai" {
		t.Errorf("Expected root command Use to be 'cronai', got %s", rootCmd.Use)
	}

	if rootCmd.Short != "AI agent for scheduled prompt execution" {
		t.Errorf("Unexpected short description: %s", rootCmd.Short)
	}

	// Check if config flag exists
	configFlag := rootCmd.PersistentFlags().Lookup("config")
	if configFlag == nil {
		t.Error("Expected config flag to be defined")
	}

	// Check if version flag exists
	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Error("Expected version flag to be defined")
	}
}

func TestInitConfig(t *testing.T) {
	// Create temporary .env file
	tmpDir := t.TempDir()
	envFile := tmpDir + "/.env"
	content := []byte("TEST_VAR=test_value")
	if err := os.WriteFile(envFile, content, 0644); err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Change to temporary directory
	oldPwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldPwd); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Test config init
	initConfig()

	// Since godotenv loads to environment, we should not see an error in the test output
	// The function prints a warning if no .env found, but doesn't fail
	// This test primarily ensures initConfig runs without panicking
}

func TestInitConfigNoEnvFile(t *testing.T) {
	// Change to temporary directory without .env file
	tmpDir := t.TempDir()
	oldPwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldPwd); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// Test that initConfig runs without error even without .env file
	initConfig()
}

func TestRootCommandLongDescription(t *testing.T) {
	// Test comprehensive long description contains key features
	expectedStrings := []string{
		"CronAI - Your Automated AI Assistant",
		"Schedule AI prompts using cron syntax",
		"Support for multiple AI models",
		"Template-based response formatting",
		"Model fallback and error handling",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(rootCmd.Long, expected) {
			t.Errorf("Expected Long description to contain '%s'", expected)
		}
	}
}

func TestRootCommandExamples(t *testing.T) {
	expectedExamples := []string{
		"cronai start",
		"cronai run --model=claude",
		"cronai list",
		"cronai prompt search",
	}

	for _, expected := range expectedExamples {
		if !strings.Contains(rootCmd.Example, expected) {
			t.Errorf("Expected Example to contain '%s'", expected)
		}
	}
}

func TestRootCommandVersionFlag(t *testing.T) {
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

	// Create a command with version flag set
	cmd := rootCmd
	flags := cmd.Flags()
	err = flags.Set("version", "true")
	if err != nil {
		t.Fatalf("Failed to set version flag: %v", err)
	}

	// Run the root command with version flag
	rootCmd.Run(cmd, []string{})

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

	// Reset the flag for other tests
	err = flags.Set("version", "false")
	if err != nil {
		t.Fatalf("Failed to reset version flag: %v", err)
	}
}

func TestRootCommandHelpOutput(t *testing.T) {
	// Capture output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run the root command without flags (should show help)
	rootCmd.Run(rootCmd, []string{})

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

	// Check that help output contains expected content
	if !strings.Contains(output, "CronAI - Your Automated AI Assistant") {
		t.Errorf("Expected help output to contain description, got: %s", output)
	}
}
