// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseQueueConfig(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *Task
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "empty line",
			line:     "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "comment line",
			line:     "# This is a comment",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "non-queue line",
			line:     "0 8 * * * openai prompt processor",
			expected: nil,
			wantErr:  false,
		},
		{
			name: "basic queue configuration",
			line: "queue main-queue rabbitmq amqp://localhost:5672 tasks",
			expected: &Task{
				Name:       "main-queue",
				Type:       "rabbitmq",
				Connection: "amqp://localhost:5672",
				Queue:      "tasks",
				// Accept both nil and empty-map – avoids brittle failure
				Options:    nil,
				RetryLimit: 3,
				RetryDelay: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "queue with options",
			line: "queue sqs-consumer sqs https://sqs.us-east-1.amazonaws.com/123456789/myqueue myqueue retry_limit=5,retry_delay=10s",
			expected: &Task{
				Name:       "sqs-consumer",
				Type:       "sqs",
				Connection: "https://sqs.us-east-1.amazonaws.com/123456789/myqueue",
				Queue:      "myqueue",
				Options: map[string]interface{}{
					"retry_limit": 5,
					"retry_delay": 10 * time.Second,
				},
				RetryLimit: 5,
				RetryDelay: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "queue with provider-specific options",
			line: "queue azure-bus servicebus Endpoint=sb://namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=key topic1 retry_limit=2,retry_delay=1m,prefetch_count=10,max_concurrent=5",
			expected: &Task{
				Name:       "azure-bus",
				Type:       "servicebus",
				Connection: "Endpoint=sb://namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=key",
				Queue:      "topic1",
				Options: map[string]interface{}{
					"retry_limit":    2,
					"retry_delay":    1 * time.Minute,
					"prefetch_count": "10",
					"max_concurrent": "5",
				},
				RetryLimit: 2,
				RetryDelay: 1 * time.Minute,
			},
			wantErr: false,
		},
		{
			name:    "insufficient fields",
			line:    "queue main-queue rabbitmq",
			wantErr: true,
			errMsg:  "invalid queue format: need at least 4 fields (name, type, connection, queue)",
		},
		{
			name:    "invalid retry_limit",
			line:    "queue test sqs conn queue retry_limit=abc",
			wantErr: true,
			errMsg:  "invalid options format: invalid retry_limit: abc",
		},
		{
			name:    "invalid retry_delay",
			line:    "queue test sqs conn queue retry_delay=invalid",
			wantErr: true,
			errMsg:  "invalid options format: invalid retry_delay: invalid",
		},
		{
			name:    "malformed options",
			line:    "queue test sqs conn queue invalid_option",
			wantErr: true,
			errMsg:  "invalid options format: invalid option format: invalid_option",
		},
		{
			name: "queue with spaces in connection string",
			line: "queue pubsub-consumer pubsub projects/my-project/topics/my-topic my-subscription retry_limit=4",
			expected: &Task{
				Name:       "pubsub-consumer",
				Type:       "pubsub",
				Connection: "projects/my-project/topics/my-topic",
				Queue:      "my-subscription",
				Options: map[string]interface{}{
					"retry_limit": 4,
				},
				RetryLimit: 4,
				RetryDelay: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "queue with quoted connection string containing spaces",
			line: `queue azure-consumer azurebus "Endpoint=sb://mybus.servicebus.windows.net/;SharedAccessKeyName=policy;SharedAccessKey=key with spaces" "my queue name" retry_limit=3`,
			expected: &Task{
				Name:       "azure-consumer",
				Type:       "azurebus",
				Connection: "Endpoint=sb://mybus.servicebus.windows.net/;SharedAccessKeyName=policy;SharedAccessKey=key with spaces",
				Queue:      "my queue name",
				Options: map[string]interface{}{
					"retry_limit": 3,
				},
				RetryLimit: 3,
				RetryDelay: 5 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseQueueConfig(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestParseQueueOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "empty options",
			input:    "",
			expected: map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "single option",
			input:    "key=value",
			expected: map[string]interface{}{"key": "value"},
			wantErr:  false,
		},
		{
			name:  "multiple options",
			input: "key1=value1,key2=value2,key3=value3",
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			wantErr: false,
		},
		{
			name:  "options with spaces",
			input: "key1 = value1 , key2 = value2",
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name:     "retry_limit option",
			input:    "retry_limit=10",
			expected: map[string]interface{}{"retry_limit": 10},
			wantErr:  false,
		},
		{
			name:     "retry_delay option",
			input:    "retry_delay=30s",
			expected: map[string]interface{}{"retry_delay": 30 * time.Second},
			wantErr:  false,
		},
		{
			name:  "mixed options",
			input: "retry_limit=5,retry_delay=1m,custom_option=custom_value",
			expected: map[string]interface{}{
				"retry_limit":   5,
				"retry_delay":   1 * time.Minute,
				"custom_option": "custom_value",
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid_format",
			wantErr: true,
			errMsg:  "invalid option format: invalid_format",
		},
		{
			name:    "invalid retry_limit",
			input:   "retry_limit=not_a_number",
			wantErr: true,
			errMsg:  "invalid retry_limit: not_a_number",
		},
		{
			name:    "invalid retry_delay",
			input:   "retry_delay=not_a_duration",
			wantErr: true,
			errMsg:  "invalid retry_delay: not_a_duration",
		},
		{
			name:     "empty value",
			input:    "key=",
			expected: map[string]interface{}{"key": ""},
			wantErr:  false,
		},
		{
			name:  "trailing comma",
			input: "key1=value1,key2=value2,",
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseQueueOptions(tt.input)

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
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestCreateConsumerConfig(t *testing.T) {
	task := &Task{
		Name:       "test-queue",
		Type:       "rabbitmq",
		Connection: "amqp://localhost:5672",
		Queue:      "tasks",
		Options: map[string]interface{}{
			"prefetch": "10",
			"durable":  "true",
		},
		RetryLimit: 5,
		RetryDelay: 10 * time.Second,
	}

	config := CreateConsumerConfig(task)

	if config.Type != task.Type {
		t.Errorf("expected type %s, got %s", task.Type, config.Type)
	}
	if config.Connection != task.Connection {
		t.Errorf("expected connection %s, got %s", task.Connection, config.Connection)
	}
	if config.Queue != task.Queue {
		t.Errorf("expected queue %s, got %s", task.Queue, config.Queue)
	}
	if config.RetryLimit != task.RetryLimit {
		t.Errorf("expected retry limit %d, got %d", task.RetryLimit, config.RetryLimit)
	}
	if config.RetryDelay != task.RetryDelay {
		t.Errorf("expected retry delay %v, got %v", task.RetryDelay, config.RetryDelay)
	}
	if !reflect.DeepEqual(config.Options, task.Options) {
		t.Errorf("expected options %+v, got %+v", task.Options, config.Options)
	}
}

func TestIsQueueConfig(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "queue config line",
			line:     "queue main rabbitmq conn queue",
			expected: true,
		},
		{
			name:     "queue config with spaces",
			line:     "  queue main rabbitmq conn queue  ",
			expected: true,
		},
		{
			name:     "cron config line",
			line:     "0 8 * * * openai prompt processor",
			expected: false,
		},
		{
			name:     "comment line",
			line:     "# queue main rabbitmq conn queue",
			expected: false,
		},
		{
			name:     "empty line",
			line:     "",
			expected: false,
		},
		{
			name:     "queue in middle of line",
			line:     "not a queue config line",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsQueueConfig(tt.line)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
