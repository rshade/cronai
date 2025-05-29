// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDefaultMessageParser_Parse(t *testing.T) {
	parser := NewMessageParser()

	tests := []struct {
		name     string
		message  *Message
		expected *TaskMessage
		wantErr  bool
		errMsg   string
	}{
		{
			name:    "nil message",
			message: nil,
			wantErr: true,
			errMsg:  "message cannot be nil",
		},
		{
			name: "empty message body",
			message: &Message{
				ID:   "test-1",
				Body: []byte{},
			},
			wantErr: true,
			errMsg:  "message body is empty",
		},
		{
			name: "comprehensive message format",
			message: &Message{
				ID: "test-2",
				Body: []byte(`{
					"model": "openai",
					"prompt": "test_prompt",
					"processor": "email-admin@example.com",
					"variables": {
						"key1": "value1",
						"key2": "value2"
					},
					"is_inline": false
				}`),
			},
			expected: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "email-admin@example.com",
				Variables: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				IsInline: false,
			},
			wantErr: false,
		},
		{
			name: "comprehensive message with inline prompt",
			message: &Message{
				ID: "test-3",
				Body: []byte(`{
					"model": "claude",
					"prompt": "Generate a report for {{date}}",
					"processor": "slack-reports",
					"is_inline": true
				}`),
			},
			expected: &TaskMessage{
				Model:     "claude",
				Prompt:    "Generate a report for {{date}}",
				Processor: "slack-reports",
				Variables: nil,
				IsInline:  true,
			},
			wantErr: false,
		},
		{
			name: "minimal message format",
			message: &Message{
				ID: "test-4",
				Body: []byte(`{
					"variables": {
						"date": "2024-01-01",
						"project": "CronAI"
					}
				}`),
				Attributes: map[string]string{
					"model":     "gemini",
					"prompt":    "weekly_report",
					"processor": "webhook-https://example.com/hook",
				},
			},
			expected: &TaskMessage{
				Model:     "gemini",
				Prompt:    "weekly_report",
				Processor: "webhook-https://example.com/hook",
				Variables: map[string]string{
					"date":    "2024-01-01",
					"project": "CronAI",
				},
				IsInline: false,
			},
			wantErr: false,
		},
		{
			name: "minimal format missing attributes",
			message: &Message{
				ID: "test-5",
				Body: []byte(`{
					"variables": {
						"key": "value"
					}
				}`),
				Attributes: map[string]string{
					"model": "openai",
					// Missing prompt and processor
				},
			},
			wantErr: true,
			errMsg:  "minimal message format requires model, prompt, and processor in message attributes",
		},
		{
			name: "invalid JSON",
			message: &Message{
				ID:   "test-6",
				Body: []byte(`{invalid json`),
			},
			wantErr: true,
			errMsg:  "unable to parse message: invalid format",
		},
		{
			name: "comprehensive format with typo in field name",
			message: &Message{
				ID: "test-8",
				Body: []byte(`{
					"model": "openai",
					"prompts": "test_prompt",
					"processor": "console"
				}`),
			},
			wantErr: true,
			errMsg:  "comprehensive message format missing required fields: model, prompt, processor",
		},
		{
			name: "comprehensive format with all fields missing",
			message: &Message{
				ID: "test-9",
				Body: []byte(`{
					"is_inline": true
				}`),
			},
			wantErr: true,
			errMsg:  "comprehensive message format missing required fields: model, prompt, processor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.message)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result but got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				} else if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestDefaultMessageParser_Validate(t *testing.T) {
	parser := NewMessageParser()

	tests := []struct {
		name    string
		task    *TaskMessage
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil task",
			task:    nil,
			wantErr: true,
			errMsg:  "task message cannot be nil",
		},
		{
			name: "valid task",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "console",
				Variables: map[string]string{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "empty model",
			task: &TaskMessage{
				Model:     "",
				Prompt:    "test_prompt",
				Processor: "console",
			},
			wantErr: true,
			errMsg:  "model cannot be empty",
		},
		{
			name: "whitespace model",
			task: &TaskMessage{
				Model:     "   ",
				Prompt:    "test_prompt",
				Processor: "console",
			},
			wantErr: true,
			errMsg:  "model cannot be empty",
		},
		{
			name: "empty prompt",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "",
				Processor: "console",
			},
			wantErr: true,
			errMsg:  "prompt cannot be empty",
		},
		{
			name: "empty processor",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "",
			},
			wantErr: true,
			errMsg:  "processor cannot be empty",
		},
		{
			name: "unsupported model",
			task: &TaskMessage{
				Model:     "unsupported",
				Prompt:    "test_prompt",
				Processor: "console",
			},
			wantErr: true,
			errMsg:  "unsupported model: unsupported",
		},
		{
			name: "file reference with newline",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt\nwith_newline",
				Processor: "console",
				IsInline:  false,
			},
			wantErr: true,
			errMsg:  "prompt file reference cannot contain newlines",
		},
		{
			name: "inline prompt with newline",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "Generate a report.\nInclude all details.",
				Processor: "console",
				IsInline:  true,
			},
			wantErr: false,
		},
		{
			name: "variable with whitespace key",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "console",
				Variables: map[string]string{
					"key with space": "value",
				},
			},
			wantErr: true,
			errMsg:  "variable key 'key with space' contains whitespace",
		},
		{
			name: "empty variable key",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "console",
				Variables: map[string]string{
					"": "value",
				},
			},
			wantErr: true,
			errMsg:  "variable key cannot be empty",
		},
		{
			name: "empty variable value",
			task: &TaskMessage{
				Model:     "openai",
				Prompt:    "test_prompt",
				Processor: "console",
				Variables: map[string]string{
					"key": "",
				},
			},
			wantErr: true,
			errMsg:  "variable value for key 'key' cannot be empty",
		},
		{
			name: "case insensitive model",
			task: &TaskMessage{
				Model:     "OPENAI",
				Prompt:    "test_prompt",
				Processor: "console",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.Validate(tt.task)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
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

func TestMessageParsingIntegration(t *testing.T) {
	parser := NewMessageParser()

	// Test roundtrip parsing
	original := &ComprehensiveMessage{
		Model:     "claude",
		Prompt:    "Generate weekly report",
		Processor: "email-team@example.com",
		Variables: map[string]string{
			"week":    "2024-W01",
			"project": "CronAI",
		},
		IsInline: true,
	}

	// Marshal to JSON
	body, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	// Create message
	msg := &Message{
		ID:         "integration-test",
		Body:       body,
		ReceivedAt: time.Now(),
	}

	// Parse
	task, err := parser.Parse(msg)
	if err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	// Validate
	if err := parser.Validate(task); err != nil {
		t.Fatalf("failed to validate task: %v", err)
	}

	// Check results
	if task.Model != original.Model {
		t.Errorf("expected model %s, got %s", original.Model, task.Model)
	}
	if task.Prompt != original.Prompt {
		t.Errorf("expected prompt %s, got %s", original.Prompt, task.Prompt)
	}
	if task.Processor != original.Processor {
		t.Errorf("expected processor %s, got %s", original.Processor, task.Processor)
	}
	if task.IsInline != original.IsInline {
		t.Errorf("expected IsInline %v, got %v", original.IsInline, task.IsInline)
	}
	if !reflect.DeepEqual(task.Variables, original.Variables) {
		t.Errorf("expected variables %v, got %v", original.Variables, task.Variables)
	}
}
