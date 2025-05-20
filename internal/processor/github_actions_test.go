package processor

import (
	"os"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGitHubTest(t *testing.T) {
	// Set GO_TEST environment variable to skip token validation
	oldTestEnv := os.Getenv("GO_TEST")
	if err := os.Setenv("GO_TEST", "1"); err != nil {
		t.Fatalf("Failed to set GO_TEST: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Setenv("GO_TEST", oldTestEnv); err != nil {
			t.Errorf("Failed to restore GO_TEST: %v", err)
		}
	})

	// Set a dummy GitHub token for tests
	oldToken := os.Getenv("GITHUB_TOKEN")
	if err := os.Setenv("GITHUB_TOKEN", "test-token"); err != nil {
		t.Fatalf("Failed to set GITHUB_TOKEN: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Setenv("GITHUB_TOKEN", oldToken); err != nil {
			t.Errorf("Failed to restore GITHUB_TOKEN: %v", err)
		}
	})
}

func TestGitHubProcessor_ProcessGitHubWithTemplate(t *testing.T) {
	setupGitHubTest(t)

	tests := []struct {
		name         string
		target       string
		templateName string
		wantErr      bool
	}{
		{
			name:         "issue template",
			target:       "issue:owner/repo",
			templateName: "default_github_issue",
			wantErr:      false,
		},
		{
			name:         "comment template",
			target:       "comment:owner/repo#123",
			templateName: "default_github_comment",
			wantErr:      false,
		},
		{
			name:         "PR template",
			target:       "pr:owner/repo",
			templateName: "default_github_pr",
			wantErr:      false,
		},
		{
			name:         "invalid target",
			target:       "invalid",
			templateName: "default_github_issue",
			wantErr:      true,
		},
		{
			name:         "unsupported action",
			target:       "unsupported:owner/repo",
			templateName: "default_github_issue",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a GitHub processor
			processor, err := NewGitHubProcessor(Config{
				Type:   "github",
				Target: tt.target,
			})
			require.NoError(t, err)

			// Prepare template data
			testTime := time.Now()
			data := template.Data{
				Content:     "Test content",
				Model:       "test-model",
				Timestamp:   testTime,
				PromptName:  "test-prompt",
				ExecutionID: "test-execution-id",
				Variables: map[string]string{
					"head_branch": "feature-branch",
					"base_branch": "main",
				},
				Metadata: map[string]string{
					"date": testTime.Format("2006-01-02"),
				},
			}

			// Test process with template
			gitHubProcessor, ok := processor.(*GitHubProcessor)
			require.True(t, ok, "Failed to cast processor to GitHubProcessor")
			err = gitHubProcessor.processGitHubWithTemplate(tt.target, data, tt.templateName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitHubProcessor_Issue(t *testing.T) {
	setupGitHubTest(t)

	// Create a GitHub processor
	processor, err := NewGitHubProcessor(Config{
		Type:   "github",
		Target: "issue:owner/repo",
	})
	require.NoError(t, err)
	gitHubProcessor, ok := processor.(*GitHubProcessor)
	require.True(t, ok, "Failed to cast processor to GitHubProcessor")

	// Test issue creation payload
	payload := map[string]interface{}{
		"title":  "Test Issue",
		"body":   "Test Body",
		"labels": []interface{}{"label1", "label2"},
	}

	err = gitHubProcessor.processGitHubIssue("owner/repo", payload)
	assert.NoError(t, err)

	// Test with missing title
	invalidPayload := map[string]interface{}{
		"body": "Test Body",
	}
	err = gitHubProcessor.processGitHubIssue("owner/repo", invalidPayload)
	assert.Error(t, err)

	// Test with invalid repo format
	err = gitHubProcessor.processGitHubIssue("invalid", payload)
	assert.Error(t, err)
}

func TestGitHubProcessor_Comment(t *testing.T) {
	setupGitHubTest(t)

	// Create a GitHub processor
	processor, err := NewGitHubProcessor(Config{
		Type:   "github",
		Target: "comment:owner/repo#123",
	})
	require.NoError(t, err)
	gitHubProcessor, ok := processor.(*GitHubProcessor)
	require.True(t, ok, "Failed to cast processor to GitHubProcessor")

	// Test comment creation payload
	payload := map[string]interface{}{
		"body": "Test Comment",
	}

	err = gitHubProcessor.processGitHubComment("owner/repo#123", payload)
	assert.NoError(t, err)

	// Test with missing body
	invalidPayload := map[string]interface{}{}
	err = gitHubProcessor.processGitHubComment("owner/repo#123", invalidPayload)
	assert.Error(t, err)

	// Test with invalid repo format
	err = gitHubProcessor.processGitHubComment("invalid", payload)
	assert.Error(t, err)

	// Test with invalid issue number
	err = gitHubProcessor.processGitHubComment("owner/repo#abc", payload)
	assert.Error(t, err)
}

func TestGitHubProcessor_PR(t *testing.T) {
	setupGitHubTest(t)

	// Create a GitHub processor
	processor, err := NewGitHubProcessor(Config{
		Type:   "github",
		Target: "pr:owner/repo",
	})
	require.NoError(t, err)
	gitHubProcessor, ok := processor.(*GitHubProcessor)
	require.True(t, ok, "Failed to cast processor to GitHubProcessor")

	// Test PR creation payload
	payload := map[string]interface{}{
		"title": "Test PR",
		"body":  "Test Body",
		"head":  "feature-branch",
		"base":  "main",
	}

	err = gitHubProcessor.processGitHubPR("owner/repo", payload)
	assert.NoError(t, err)

	// Test with missing title
	invalidPayload := map[string]interface{}{
		"body": "Test Body",
		"head": "feature-branch",
		"base": "main",
	}
	err = gitHubProcessor.processGitHubPR("owner/repo", invalidPayload)
	assert.Error(t, err)

	// Test with missing head branch
	invalidPayload = map[string]interface{}{
		"title": "Test PR",
		"body":  "Test Body",
		"base":  "main",
	}
	err = gitHubProcessor.processGitHubPR("owner/repo", invalidPayload)
	assert.Error(t, err)

	// Test with invalid repo format
	err = gitHubProcessor.processGitHubPR("invalid", payload)
	assert.Error(t, err)
}

func TestGitHubProcessor_Process_Integration(t *testing.T) {
	setupGitHubTest(t)

	// Test with different targets
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:    "issue target",
			target:  "issue:owner/repo",
			wantErr: false,
		},
		{
			name:    "comment target",
			target:  "comment:owner/repo#123",
			wantErr: false,
		},
		{
			name:    "PR target",
			target:  "pr:owner/repo",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a GitHub processor
			processor, err := NewGitHubProcessor(Config{
				Type:   "github",
				Target: tt.target,
			})
			require.NoError(t, err)

			// Create a model response
			response := &models.ModelResponse{
				Content:     "Test content",
				Model:       "test-model",
				Timestamp:   time.Now(),
				PromptName:  "test-prompt",
				ExecutionID: "test-execution-id",
				Variables: map[string]string{
					"head_branch": "feature-branch",
					"base_branch": "main",
				},
			}

			// Process the response
			err = processor.Process(response, "")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
