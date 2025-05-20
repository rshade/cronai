package processor

import (
	"os"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGitHubProcessor(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Test with valid config
	config := Config{
		Type:   "github",
		Target: "issue:owner/repo",
	}
	processor, err := NewGitHubProcessor(config)
	if err != nil {
		t.Errorf("NewGitHubProcessor failed: %v", err)
	}
	if processor == nil {
		t.Error("NewGitHubProcessor returned nil processor")
	}

	// Test with invalid config
	config = Config{
		Type: "github",
	}
	processor, err = NewGitHubProcessor(config)
	if err == nil {
		t.Error("NewGitHubProcessor should fail with invalid config")
	}
	if processor != nil {
		t.Error("NewGitHubProcessor should return nil processor with invalid config")
	}
}

func TestGitHubProcessor_Validate(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Create test processor
	config := Config{
		Type:   "github",
		Target: "issue:owner/repo",
	}
	processor, err := NewGitHubProcessor(config)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Test validation
	err = processor.Validate()
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestGitHubProcessor_Process(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Set a dummy token for the test
	oldToken := os.Getenv("GITHUB_TOKEN")
	if err := os.Setenv("GITHUB_TOKEN", "dummy-token-for-testing"); err != nil {
		t.Fatalf("Failed to set GITHUB_TOKEN: %v", err)
	}
	defer func() {
		if err := os.Setenv("GITHUB_TOKEN", oldToken); err != nil {
			t.Errorf("Failed to restore GITHUB_TOKEN: %v", err)
		}
	}()

	// Create test processor
	config := Config{
		Type:   "github",
		Target: "issue:owner/repo",
	}
	processor, err := NewGitHubProcessor(config)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Test processing response
	response := &models.ModelResponse{
		Content:    "Test content",
		Model:      "test-model",
		Timestamp:  time.Now(),
		PromptName: "test-prompt",
	}
	err = processor.Process(response, "")
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}
}

func TestGitHubProcessor_GetType(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Create test processor
	config := Config{
		Type:   "github",
		Target: "issue:owner/repo",
	}
	processor, err := NewGitHubProcessor(config)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Test GetType
	processorType := processor.GetType()
	if processorType != "github" {
		t.Errorf("GetType returned wrong type: got %s, want github", processorType)
	}
}

func TestGitHubProcessor_GetConfig(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Create test processor
	config := Config{
		Type:   "github",
		Target: "owner/repo",
	}
	processor, err := NewGitHubProcessor(config)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Test GetConfig
	config = processor.GetConfig()
	if config.Type != "github" || config.Target != "owner/repo" {
		t.Errorf("GetConfig returned wrong config: got %+v, want {Type:github Target:owner/repo}", config)
	}
}

func TestGitHubProcessor_TargetParsing(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	defer func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	}()

	// Save and restore environment variable
	oldToken := os.Getenv(EnvGitHubToken)
	defer func() {
		if err := os.Setenv(EnvGitHubToken, oldToken); err != nil {
			t.Errorf("Failed to restore GITHUB_TOKEN: %v", err)
		}
	}()
	if err := os.Setenv(EnvGitHubToken, "test-token"); err != nil {
		t.Fatalf("Failed to set GITHUB_TOKEN: %v", err)
	}

	tests := []struct {
		name    string
		target  string
		action  string
		repo    string
		wantErr bool
	}{
		{
			name:    "Issue target",
			target:  "issue:owner/repo",
			action:  "issue",
			repo:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "Comment target",
			target:  "comment:owner/repo#123",
			action:  "comment",
			repo:    "owner/repo#123",
			wantErr: false,
		},
		{
			name:    "PR target",
			target:  "pr:owner/repo",
			action:  "pr",
			repo:    "owner/repo",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			target:  "invalid",
			wantErr: true,
		},
		{
			name:    "Invalid action",
			target:  "unknown:owner/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := NewGitHubProcessor(Config{
				Type:   "github",
				Target: tt.target,
			})
			require.NoError(t, err)

			response := &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "test-prompt",
				ExecutionID: "test-execution-id",
			}

			err = processor.Process(response, "")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// In MVP implementation, should succeed as we're just logging
				assert.NoError(t, err)
			}
		})
	}
}
