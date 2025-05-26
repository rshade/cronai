// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/prompt"
)

// testPromptManager wraps the package-level prompt functions
type testPromptManager struct{}

func (m *testPromptManager) LoadPrompt(promptName string) (string, error) {
	return prompt.LoadPrompt(promptName)
}

func (m *testPromptManager) LoadPromptWithVariables(promptName string, variables map[string]string) (string, error) {
	return prompt.LoadPromptWithVariables(promptName, variables)
}

func (m *testPromptManager) ListPrompts() ([]prompt.Info, error) {
	return prompt.ListPrompts()
}

func (m *testPromptManager) GetPrompt(_ string) (prompt.Info, error) {
	return prompt.Info{}, fmt.Errorf("not implemented")
}

func (m *testPromptManager) GetPromptMetadata(_ string) (prompt.Metadata, error) {
	return prompt.Metadata{}, fmt.Errorf("not implemented")
}

func (m *testPromptManager) GetPromptContent(name string) (string, error) {
	return prompt.LoadPrompt(name)
}

func (m *testPromptManager) GetPromptVariables(_ string) ([]prompt.Variable, error) {
	return nil, fmt.Errorf("not implemented")
}

// Setup test environment
func setupTestEnvironment(t *testing.T) (string, func()) {
	// Save original CRON_PROMPTS_DIR
	originalPromptDir := os.Getenv("CRON_PROMPTS_DIR")
	// Create a temporary directory for test prompts
	tmpDir, err := os.MkdirTemp("", "cronai-queue-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create cron_prompts directory
	promptsDir := tmpDir + "/cron_prompts"
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
		t.Fatalf("failed to create prompts dir: %v", err)
	}

	// Create a test prompt file
	testPrompt := `# Test Prompt

Generate a report for {{project}} on {{date}}.

Include the following:
- Summary
- Key metrics
- Recommendations`

	if err := os.WriteFile(filepath.Join(promptsDir, "test_prompt.md"), []byte(testPrompt), 0644); err != nil {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
		t.Fatalf("failed to create test prompt: %v", err)
	}

	// Log the prompt directory for debugging
	t.Logf("Test prompts directory: %s", promptsDir)

	// List files in the prompts directory
	files, err := os.ReadDir(promptsDir)
	if err != nil {
		t.Logf("Failed to read prompts directory: %v", err)
	} else {
		t.Logf("Files in prompts directory:")
		for _, f := range files {
			t.Logf("  - %s", f.Name())
		}
	}

	// Don't change directories - just set the environment variable

	// Set environment variable for prompts directory
	if err := os.Setenv("CRON_PROMPTS_DIR", promptsDir); err != nil {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to clean up temp dir: %v", err)
		}
		t.Fatalf("failed to set environment variable: %v", err)
	}

	// Setup mock model execution
	originalExecute := executeModel
	executeModel = func(_, _ string, variables map[string]string, _ string) (*models.ModelResponse, error) {
		return &models.ModelResponse{
			Model:      variables["model"],
			PromptName: variables["promptName"],
			Content:    "Mock response for test",
		}, nil
	}

	// Console processor is already registered by default in GetRegistry()

	// Log environment variable
	t.Logf("CRON_PROMPTS_DIR environment variable: %s", os.Getenv("CRON_PROMPTS_DIR"))

	cleanup := func() {
		executeModel = originalExecute
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to clean up temp dir: %v", err)
		}
		// Restore original CRON_PROMPTS_DIR
		if originalPromptDir != "" {
			if err := os.Setenv("CRON_PROMPTS_DIR", originalPromptDir); err != nil {
				t.Logf("failed to restore CRON_PROMPTS_DIR: %v", err)
			}
		} else {
			if err := os.Unsetenv("CRON_PROMPTS_DIR"); err != nil {
				t.Logf("failed to unset CRON_PROMPTS_DIR: %v", err)
			}
		}
	}

	return tmpDir, cleanup
}

func TestNewTaskProcessor(t *testing.T) {
	processor := NewTaskProcessor()
	if processor == nil {
		t.Error("expected non-nil processor")
	}

	// Check if it's the right type
	if _, ok := processor.(*DefaultTaskProcessor); !ok {
		t.Error("expected DefaultTaskProcessor type")
	}
}

func TestDefaultTaskProcessor_Process(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	processor := &DefaultTaskProcessor{
		promptManager: &testPromptManager{},
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		task    *TaskMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful file-based prompt",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "console",
				Variables: map[string]string{
					"project": "CronAI",
					"date":    "2024-01-01",
				},
				IsInline: false,
			},
			wantErr: false,
		},
		{
			name: "successful inline prompt",
			task: &TaskMessage{
				Model:     "claude",
				Prompt:    "Generate a summary report for the project.",
				Processor: "console",
				Variables: map[string]string{
					"format": "markdown",
				},
				IsInline: true,
			},
			wantErr: false,
		},
		{
			name: "file prompt without variables",
			task: &TaskMessage{
				Model:     "gemini",
				Prompt:    "test_prompt",
				Processor: "console",
				IsInline:  false,
			},
			wantErr: false,
		},
		{
			name: "non-existent prompt file",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "non_existent",
				Processor: "console",
				IsInline:  false,
			},
			wantErr: true,
			errMsg:  "failed to load prompt",
		},
		{
			name: "invalid processor defaults to console",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "invalid-processor",
				IsInline:  false,
			},
			wantErr: false, // defaults to console processor
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processor.Process(ctx, tt.task)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !containsError(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParseProcessorName(t *testing.T) {
	tests := []struct {
		name           string
		processorName  string
		expectedType   string
		expectedTarget string
	}{
		{
			name:           "slack processor",
			processorName:  "slack-channel-name",
			expectedType:   "slack",
			expectedTarget: "channel-name",
		},
		{
			name:           "email processor",
			processorName:  "email-admin@example.com",
			expectedType:   "email",
			expectedTarget: "admin@example.com",
		},
		{
			name:           "webhook processor",
			processorName:  "webhook-https://example.com/hook",
			expectedType:   "webhook",
			expectedTarget: "https://example.com/hook",
		},
		{
			name:           "github processor",
			processorName:  "github-owner/repo#123",
			expectedType:   "github",
			expectedTarget: "owner/repo#123",
		},
		{
			name:           "file processor",
			processorName:  "file-output.txt",
			expectedType:   "file",
			expectedTarget: "output.txt",
		},
		{
			name:           "log-to-file special case",
			processorName:  "log-to-file",
			expectedType:   "file",
			expectedTarget: "to-file",
		},
		{
			name:           "console processor",
			processorName:  "console",
			expectedType:   "console",
			expectedTarget: "",
		},
		{
			name:           "unknown processor",
			processorName:  "unknown",
			expectedType:   "console",
			expectedTarget: "",
		},
		{
			name:           "empty processor",
			processorName:  "",
			expectedType:   "console",
			expectedTarget: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			procType, target := parseProcessorName(tt.processorName)

			if procType != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, procType)
			}
			if target != tt.expectedTarget {
				t.Errorf("expected target %q, got %q", tt.expectedTarget, target)
			}
		})
	}
}

func TestDefaultTaskProcessor_ProcessWithModelError(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Override model execution to return error
	executeModel = func(_, _ string, _ map[string]string, _ string) (*models.ModelResponse, error) {
		return nil, fmt.Errorf("model execution failed")
	}

	processor := &DefaultTaskProcessor{
		promptManager: &testPromptManager{},
	}

	ctx := context.Background()
	task := &TaskMessage{
		Model:     "openai",
		Prompt:    "test_prompt",
		Processor: "console",
		IsInline:  false,
	}

	err := processor.Process(ctx, task)
	if err == nil {
		t.Error("expected error but got nil")
	}
	if !containsError(err.Error(), "failed to execute model") {
		t.Errorf("expected model execution error, got: %v", err)
	}
}

// Helper function
func containsError(err, substr string) bool {
	for i := 0; i <= len(err)-len(substr); i++ {
		if err[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestDefaultTaskProcessor_PromptNameVariable tests that promptName variable is available during template substitution
func TestDefaultTaskProcessor_PromptNameVariable(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test prompt that uses the promptName variable
	promptContent := `This is a test prompt.
The prompt name is: {{promptName}}
End of prompt.`

	// Get the prompts directory from environment
	promptsDir := os.Getenv("CRON_PROMPTS_DIR")
	promptPath := filepath.Join(promptsDir, "prompt_with_name.md")
	if err := os.WriteFile(promptPath, []byte(promptContent), 0644); err != nil {
		t.Fatalf("failed to create test prompt: %v", err)
	}

	// Variable to capture the processed prompt content
	var capturedPrompt string

	// Override model execution to capture the prompt
	executeModel = func(_, prompt string, _ map[string]string, _ string) (*models.ModelResponse, error) {
		capturedPrompt = prompt
		return &models.ModelResponse{
			Content: "test response",
		}, nil
	}

	processor := &DefaultTaskProcessor{
		promptManager: &testPromptManager{},
	}

	ctx := context.Background()
	task := &TaskMessage{
		Model:     "openai",
		Prompt:    "prompt_with_name",
		Processor: "console",
		IsInline:  false,
	}

	err := processor.Process(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that {{promptName}} was replaced with the actual prompt name
	expectedContent := `This is a test prompt.
The prompt name is: prompt_with_name
End of prompt.`

	if capturedPrompt != expectedContent {
		t.Errorf("promptName variable not properly substituted.\nExpected:\n%s\nGot:\n%s", expectedContent, capturedPrompt)
	}
}

// TestDefaultTaskProcessor_InlinePromptNameVariable tests that inline tasks also get promptName in variables
func TestDefaultTaskProcessor_InlinePromptNameVariable(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Variable to capture the task variables
	var capturedVariables map[string]string

	// Override model execution to capture the variables
	executeModel = func(_, _ string, variables map[string]string, _ string) (*models.ModelResponse, error) {
		capturedVariables = variables
		return &models.ModelResponse{
			Content: "test response",
		}, nil
	}

	processor := &DefaultTaskProcessor{
		promptManager: &testPromptManager{},
	}

	ctx := context.Background()
	task := &TaskMessage{
		Model:     "openai",
		Prompt:    "This is an inline prompt for testing.",
		Processor: "console",
		IsInline:  true,
	}

	err := processor.Process(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that promptName is in the variables
	if capturedVariables == nil {
		t.Fatal("expected variables to be captured but got nil")
	}

	promptName, exists := capturedVariables["promptName"]
	if !exists {
		t.Error("expected promptName to exist in variables for inline task")
	}

	if promptName != task.Prompt {
		t.Errorf("expected promptName to be %q, got %q", task.Prompt, promptName)
	}
}
