// Package router provides event routing functionality for processing GitHub webhook events.
package router

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
)

// mockProcessor implements processor.Processor for testing
type mockProcessor struct {
	processed    bool
	lastResponse *models.ModelResponse
	lastTemplate string
	err          error
}

func (m *mockProcessor) Process(response *models.ModelResponse, templateName string) error {
	m.processed = true
	m.lastResponse = response
	m.lastTemplate = templateName
	return m.err
}

func (m *mockProcessor) Validate() error {
	return nil
}

func (m *mockProcessor) GetType() string {
	return "mock"
}

func (m *mockProcessor) GetConfig() processor.Config {
	return processor.Config{Type: "mock"}
}

// mockModel implements models.ModelClient for testing
type mockModel struct {
	response string
	err      error
}

func (m *mockModel) Execute(_ string) (*models.ModelResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &models.ModelResponse{
		Content: m.response,
		Model:   "mock",
	}, nil
}

func TestIssuesHandler(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantAction string
		modelErr   error
		wantErr    bool
	}{
		{
			name: "issue opened",
			payload: `{
				"action": "opened",
				"issue": {
					"number": 123,
					"title": "Test Issue",
					"body": "Test body",
					"user": {"login": "testuser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "opened",
			wantErr:    false,
		},
		{
			name: "issue closed",
			payload: `{
				"action": "closed",
				"issue": {
					"number": 456,
					"title": "Another Issue",
					"state": "closed",
					"user": {"login": "anotheruser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "closed",
			wantErr:    false,
		},
		{
			name:       "invalid payload",
			payload:    `{invalid json`,
			wantAction: "",
			wantErr:    true,
		},
		{
			name: "model error",
			payload: `{
				"action": "opened",
				"issue": {"number": 789},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			modelErr: fmt.Errorf("model failed"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &mockModel{response: "AI response", err: tt.modelErr}
			processor := &mockProcessor{}
			handler := NewIssuesHandler(model, processor)

			event := Event{
				Type:    "issues",
				Action:  tt.wantAction,
				Payload: json.RawMessage(tt.payload),
			}

			err := handler.Handle(event)

			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr || tt.modelErr != nil {
				return
			}

			if !processor.processed {
				t.Error("Handle() did not call processor.Process")
			}

			if processor.lastTemplate != "issues_template" {
				t.Errorf("Handle() template = %v, want issues_template", processor.lastTemplate)
			}

			// Check that metadata was added to response
			if processor.lastResponse != nil && processor.lastResponse.Variables != nil {
				if processor.lastResponse.Variables["action"] != tt.wantAction {
					t.Errorf("Handle() action in variables = %v, want %v", processor.lastResponse.Variables["action"], tt.wantAction)
				}
			}
		})
	}
}

func TestPullRequestHandler(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantAction string
		wantErr    bool
	}{
		{
			name: "pull_request opened",
			payload: `{
				"action": "opened",
				"number": 42,
				"pull_request": {
					"title": "Test PR",
					"body": "Test PR body",
					"state": "open",
					"head": {"ref": "feature-branch"},
					"base": {"ref": "main"},
					"user": {"login": "testuser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "opened",
			wantErr:    false,
		},
		{
			name: "pull_request synchronize",
			payload: `{
				"action": "synchronize",
				"number": 43,
				"pull_request": {
					"title": "Updated PR",
					"state": "open",
					"head": {"ref": "feature-branch"},
					"base": {"ref": "main"},
					"user": {"login": "testuser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "synchronize",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &mockModel{response: "AI response"}
			processor := &mockProcessor{}
			handler := NewPullRequestHandler(model, processor)

			event := Event{
				Type:    "pull_request",
				Action:  tt.wantAction,
				Payload: json.RawMessage(tt.payload),
			}

			err := handler.Handle(event)

			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !processor.processed {
				t.Error("Handle() did not call processor.Process")
			}

			if processor.lastTemplate != "pull_request_template" {
				t.Errorf("Handle() template = %v, want pull_request_template", processor.lastTemplate)
			}
		})
	}
}

func TestPushHandler(t *testing.T) {
	model := &mockModel{response: "AI response"}
	processor := &mockProcessor{}
	handler := NewPushHandler(model, processor)

	payload := `{
		"ref": "refs/heads/main",
		"before": "abc123",
		"after": "def456",
		"commits": [{
			"id": "def456",
			"message": "Test commit",
			"author": {
				"name": "Test Author",
				"email": "test@example.com"
			}
		}],
		"repository": {
			"name": "testrepo",
			"owner": {"login": "testowner"}
		}
	}`

	event := Event{
		Type:    "push",
		Payload: json.RawMessage(payload),
	}

	err := handler.Handle(event)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if !processor.processed {
		t.Error("Handle() did not call processor.Process")
	}

	if processor.lastTemplate != "push_template" {
		t.Errorf("Handle() template = %v, want push_template", processor.lastTemplate)
	}

	// Check that metadata was added
	if processor.lastResponse != nil && processor.lastResponse.Variables != nil {
		if processor.lastResponse.Variables["ref"] != "refs/heads/main" {
			t.Errorf("Handle() ref in variables = %v, want refs/heads/main", processor.lastResponse.Variables["ref"])
		}
	}
}

func TestReleaseHandler(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantAction string
		wantErr    bool
	}{
		{
			name: "release created",
			payload: `{
				"action": "created",
				"release": {
					"tag_name": "v1.0.0",
					"name": "Version 1.0.0",
					"body": "First release",
					"draft": false,
					"prerelease": false,
					"author": {"login": "testuser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "created",
			wantErr:    false,
		},
		{
			name: "release published",
			payload: `{
				"action": "published",
				"release": {
					"tag_name": "v2.0.0",
					"name": "Version 2.0.0",
					"body": "Major release",
					"draft": false,
					"prerelease": false,
					"author": {"login": "testuser"}
				},
				"repository": {
					"name": "testrepo",
					"owner": {"login": "testowner"}
				}
			}`,
			wantAction: "published",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &mockModel{response: "AI response"}
			processor := &mockProcessor{}
			handler := NewReleaseHandler(model, processor)

			event := Event{
				Type:    "release",
				Action:  tt.wantAction,
				Payload: json.RawMessage(tt.payload),
			}

			err := handler.Handle(event)

			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !processor.processed {
				t.Error("Handle() did not call processor.Process")
			}

			if processor.lastTemplate != "release_template" {
				t.Errorf("Handle() template = %v, want release_template", processor.lastTemplate)
			}
		})
	}
}

func TestHandlerProcessorError(t *testing.T) {
	tests := []struct {
		name        string
		handlerType string
		payload     string
	}{
		{
			name:        "issues handler",
			handlerType: "issues",
			payload: `{
				"action": "opened",
				"issue": {"number": 1},
				"repository": {"name": "test", "owner": {"login": "owner"}}
			}`,
		},
		{
			name:        "pull_request handler",
			handlerType: "pull_request",
			payload: `{
				"action": "opened",
				"number": 1,
				"pull_request": {"title": "Test"},
				"repository": {"name": "test", "owner": {"login": "owner"}}
			}`,
		},
		{
			name:        "push handler",
			handlerType: "push",
			payload: `{
				"ref": "refs/heads/main",
				"commits": [],
				"repository": {"name": "test", "owner": {"login": "owner"}}
			}`,
		},
		{
			name:        "release handler",
			handlerType: "release",
			payload: `{
				"action": "created",
				"release": {"tag_name": "v1.0.0"},
				"repository": {"name": "test", "owner": {"login": "owner"}}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &mockModel{response: "AI response"}
			processor := &mockProcessor{err: fmt.Errorf("processor error")}

			var handler EventHandler
			switch tt.handlerType {
			case "issues":
				handler = NewIssuesHandler(model, processor)
			case "pull_request":
				handler = NewPullRequestHandler(model, processor)
			case "push":
				handler = NewPushHandler(model, processor)
			case "release":
				handler = NewReleaseHandler(model, processor)
			}

			event := Event{
				Type:    tt.handlerType,
				Action:  "test",
				Payload: json.RawMessage(tt.payload),
			}

			err := handler.Handle(event)
			if err == nil {
				t.Error("Handle() expected error, got nil")
			}

			if err.Error() != "processor failed: processor error" {
				t.Errorf("Handle() error = %v, want 'processor failed: processor error'", err)
			}
		})
	}
}
