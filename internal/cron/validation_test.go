package cron

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestValidateTask(t *testing.T) {
	// Create a helper function to check for specific error messages
	containsErrorMessage := func(err error, message string) bool {
		if err == nil {
			return false
		}
		return strings.Contains(err.Error(), message)
	}

	tests := []struct {
		name          string
		task          Task
		expectError   bool
		errorMessages []string
	}{
		{
			name: "valid task",
			task: Task{
				Schedule:  "0 8 * * *",
				Model:     "claude",
				Prompt:    "test_prompt",
				Processor: "slack-test",
			},
			expectError: false,
		},
		{
			name: "invalid cron schedule",
			task: Task{
				Schedule:  "invalid-schedule",
				Model:     "claude",
				Prompt:    "test_prompt",
				Processor: "slack-test",
			},
			expectError:   true,
			errorMessages: []string{"invalid cron schedule"},
		},
		{
			name: "invalid model",
			task: Task{
				Schedule:  "0 8 * * *",
				Model:     "unsupported-model",
				Prompt:    "test_prompt",
				Processor: "slack-test",
			},
			expectError:   true,
			errorMessages: []string{"unsupported model"},
		},
		{
			name: "invalid prompt file",
			task: Task{
				Schedule:  "0 8 * * *",
				Model:     "claude",
				Prompt:    "non_existent_prompt",
				Processor: "slack-test",
			},
			expectError:   true,
			errorMessages: []string{"prompt file"},
		},
		{
			name: "invalid processor",
			task: Task{
				Schedule:  "0 8 * * *",
				Model:     "claude",
				Prompt:    "test_prompt",
				Processor: "invalid-processor",
			},
			expectError:   true,
			errorMessages: []string{"invalid processor"},
		},
		{
			name: "invalid template",
			task: Task{
				Schedule:  "0 8 * * *",
				Model:     "claude",
				Prompt:    "test_prompt",
				Processor: "slack-test",
				Template:  "non_existent_template",
			},
			expectError:   true,
			errorMessages: []string{"template"},
		},
		{
			name: "invalid model parameters",
			task: Task{
				Schedule:    "0 8 * * *",
				Model:       "claude",
				Prompt:      "test_prompt",
				Processor:   "slack-test",
				ModelParams: "temperature=invalid",
			},
			expectError:   true,
			errorMessages: []string{"model parameters"},
		},
		{
			name: "invalid model parameters values",
			task: Task{
				Schedule:    "0 8 * * *",
				Model:       "claude",
				Prompt:      "test_prompt",
				Processor:   "slack-test",
				ModelParams: "temperature=2.0",
			},
			expectError:   true,
			errorMessages: []string{"temperature must be between 0 and 1"},
		},
		{
			name: "multiple validation errors",
			task: Task{
				Schedule:    "invalid-schedule",
				Model:       "unsupported-model",
				Prompt:      "non_existent_prompt",
				Processor:   "invalid-processor",
				Template:    "non_existent_template",
				ModelParams: "temperature=2.0",
			},
			expectError: true,
			errorMessages: []string{
				"invalid cron schedule",
				"unsupported model",
				"prompt file",
				"invalid processor",
				"template",
				"temperature must be between 0 and 1",
			},
		},
	}

	// Setup test prompt file
	// Create a directory and file for testing
	if err := setupTestPromptFile(t); err != nil {
		t.Fatalf("Failed to setup test prompt file: %v", err)
	}
	defer cleanupTestPromptFile(t)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTask(tc.task, 1)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for %s but got nil", tc.name)
				return
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for %s but got: %v", tc.name, err)
				return
			}

			// If we expect an error, check that it contains the expected messages
			if tc.expectError {
				if !strings.Contains(err.Error(), tc.errorMessages[0]) {
					// Try to convert to multierror to check individual messages
					if multiErr, ok := err.(*multierror.Error); ok {
						for _, expectedMsg := range tc.errorMessages {
							found := false
							for _, subErr := range multiErr.Errors {
								if containsErrorMessage(subErr, expectedMsg) {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("Expected error containing '%s' but not found in: %v", expectedMsg, err)
							}
						}
					} else {
						// Not a multierror, check if it contains the main expected message
						if !containsErrorMessage(err, tc.errorMessages[0]) {
							t.Errorf("Expected error containing '%s' but got: %v", tc.errorMessages[0], err)
						}
					}
				}
			}
		})
	}
}

// Helper functions for creating test files
func setupTestPromptFile(_ *testing.T) error {
	// Ensure cron_prompts directory exists
	if err := createTestDirectory("cron_prompts"); err != nil {
		return err
	}

	// Create test prompt file
	return createTestFile("cron_prompts/test_prompt.md", "This is a test prompt")
}

func cleanupTestPromptFile(t *testing.T) {
	// Cleanup test files
	if err := removeTestFile("cron_prompts/test_prompt.md"); err != nil {
		t.Logf("Warning: Failed to remove test prompt file: %v", err)
	}
}

// Helper functions for file operations
func createTestDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func createTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func removeTestFile(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return os.Remove(path)
	}
	return nil
}

// TestValidateConfig tests the config validation in parseConfigFile function
func TestValidateConfig(t *testing.T) {
	// Create a temporary test config file with valid and invalid entries
	testConfigPath := "test_config_validation.tmp"
	testConfigContent := `# Test config file with validation issues
# Valid entries
0 8 * * * claude test_prompt slack-pm-channel
0 9 * * 1 openai test_prompt email-team@company.com

# Invalid entries
invalid-cron * * * * claude test_prompt slack-channel
0 8 * * * invalid-model test_prompt slack-channel
0 8 * * * claude non-existent-prompt slack-channel
0 8 * * * claude test_prompt invalid-processor
0 8 * * * claude test_prompt slack-channel non-existent-template
0 8 * * * claude test_prompt slack-channel model_params=temperature=2.0
`

	// Setup test environment
	if err := setupTestPromptFile(t); err != nil {
		t.Fatalf("Failed to setup test prompt file: %v", err)
	}
	defer cleanupTestPromptFile(t)

	// Create the test config file
	if err := createTestFile(testConfigPath, testConfigContent); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer func() {
		if err := removeTestFile(testConfigPath); err != nil {
			t.Logf("Warning: Failed to remove test config file: %v", err)
		}
	}()

	// Parse the config file - should get some valid tasks and validation errors
	tasks, err := parseConfigFile(testConfigPath)

	// We should have some valid tasks
	if len(tasks) < 1 {
		t.Errorf("Expected at least one valid task, got %d", len(tasks))
	}

	// We should also have validation errors
	if err == nil {
		t.Error("Expected validation errors, got nil")
	}

	// Check that the error contains expected validation messages
	expectedErrorMessages := []string{
		"invalid cron schedule",
		"invalid model",
		"prompt file",
		"invalid processor",
		"invalid variable format",
	}

	for _, expected := range expectedErrorMessages {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected error to contain '%s', but got: %v", expected, err)
		}
	}
}
