// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// mockConsumer implements the Consumer interface for testing
type mockConsumer struct {
	name      string
	connected bool
	messages  chan *Message
	errors    chan error
}

func (m *mockConsumer) Connect(_ context.Context) error {
	m.connected = true
	return nil
}

func (m *mockConsumer) Disconnect(_ context.Context) error {
	m.connected = false
	return nil
}

func (m *mockConsumer) Consume(_ context.Context) (<-chan *Message, <-chan error) {
	return m.messages, m.errors
}

func (m *mockConsumer) Acknowledge(_ context.Context, _ *Message) error {
	return nil
}

func (m *mockConsumer) Reject(_ context.Context, _ *Message, _ bool) error {
	return nil
}

func (m *mockConsumer) Name() string {
	return m.name
}

func (m *mockConsumer) Validate() error {
	if m.name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

// mockFactory creates mock consumers
func mockFactory(config *ConsumerConfig) (Consumer, error) {
	if config.Type == "error" {
		return nil, fmt.Errorf("mock error")
	}
	return &mockConsumer{
		name:     config.Type,
		messages: make(chan *Message),
		errors:   make(chan error),
	}, nil
}

func TestRegistry_Register(t *testing.T) {
	registry := &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	tests := []struct {
		name      string
		queueType string
		factory   ConsumerFactory
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful registration",
			queueType: "test",
			factory:   mockFactory,
			wantErr:   false,
		},
		{
			name:      "empty queue type",
			queueType: "",
			factory:   mockFactory,
			wantErr:   true,
			errMsg:    "queue type cannot be empty",
		},
		{
			name:      "nil factory",
			queueType: "test-nil",
			factory:   nil,
			wantErr:   true,
			errMsg:    "factory cannot be nil",
		},
		{
			name:      "duplicate registration",
			queueType: "duplicate",
			factory:   mockFactory,
			wantErr:   true,
			errMsg:    "queue type duplicate already registered",
		},
	}

	// First register duplicate
	if err := registry.Register("duplicate", mockFactory); err != nil {
		t.Fatalf("failed to register factory: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Register(tt.queueType, tt.factory)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	// Register a factory
	if err := registry.Register("test", mockFactory); err != nil {
		t.Fatalf("failed to register factory: %v", err)
	}

	tests := []struct {
		name      string
		queueType string
		wantErr   bool
	}{
		{
			name:      "existing type",
			queueType: "test",
			wantErr:   false,
		},
		{
			name:      "non-existing type",
			queueType: "notfound",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, err := registry.Get(tt.queueType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if factory != nil {
					t.Errorf("expected nil factory but got %v", factory)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if factory == nil {
					t.Errorf("expected factory but got nil")
				}
			}
		})
	}
}

func TestRegistry_List(t *testing.T) {
	registry := &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	// Register multiple factories
	types := []string{"rabbitmq", "sqs", "servicebus"}
	for _, qType := range types {
		if err := registry.Register(qType, mockFactory); err != nil {
			t.Fatalf("failed to register factory %s: %v", qType, err)
		}
	}

	list := registry.List()
	if len(list) != len(types) {
		t.Errorf("expected %d types, got %d", len(types), len(list))
	}

	// Check all types are present
	typeMap := make(map[string]bool)
	for _, qType := range list {
		typeMap[qType] = true
	}

	for _, expectedType := range types {
		if !typeMap[expectedType] {
			t.Errorf("expected type %s not found in list", expectedType)
		}
	}
}

func TestCreateConsumer(t *testing.T) {
	// Save and restore global registry
	old := globalRegistry
	defer func() { globalRegistry = old }()

	// Reset global registry for testing
	globalRegistry = &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	// Register mock factory
	if err := Register("mock", mockFactory); err != nil {
		t.Fatalf("failed to register mock factory: %v", err)
	}
	if err := Register("error", mockFactory); err != nil {
		t.Fatalf("failed to register error factory: %v", err)
	}

	tests := []struct {
		name    string
		config  *ConsumerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful creation",
			config: &ConsumerConfig{
				Type:       "mock",
				Connection: "test://localhost",
				Queue:      "test-queue",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
			errMsg:  "config cannot be nil",
		},
		{
			name: "unregistered type",
			config: &ConsumerConfig{
				Type:       "unknown",
				Connection: "test://localhost",
				Queue:      "test-queue",
			},
			wantErr: true,
		},
		{
			name: "factory error",
			config: &ConsumerConfig{
				Type:       "error",
				Connection: "test://localhost",
				Queue:      "test-queue",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer, err := CreateConsumer(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				if consumer != nil {
					t.Errorf("expected nil consumer but got %v", consumer)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if consumer == nil {
					t.Errorf("expected consumer but got nil")
				}
			}
		})
	}
}

// Helper function
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestCreateConsumer_ValidationFailure tests the validation failure path
func TestCreateConsumer_ValidationFailure(t *testing.T) {
	// Save original global registry and restore after test
	originalRegistry := globalRegistry
	defer func() {
		globalRegistry = originalRegistry
	}()

	// Reset global registry for testing
	globalRegistry = &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	// Register a factory that creates consumers with empty names (which fail validation)
	validationFailFactory := func(_ *ConsumerConfig) (Consumer, error) {
		return &mockConsumer{
			name:     "", // Empty name will fail validation
			messages: make(chan *Message),
			errors:   make(chan error),
		}, nil
	}

	if err := Register("validation-fail", validationFailFactory); err != nil {
		t.Fatalf("failed to register validation-fail factory: %v", err)
	}

	// Create consumer config that will result in validation failure
	config := &ConsumerConfig{
		Type:       "validation-fail",
		Connection: "test://localhost",
		Queue:      "test-queue",
	}

	// Attempt to create consumer
	consumer, err := CreateConsumer(config)

	// Assert error is returned
	if err == nil {
		t.Error("expected error but got nil")
	}
	if consumer != nil {
		t.Errorf("expected nil consumer but got %v", consumer)
	}
}

// TestCreateConsumer_ConfigValidation tests config validation
func TestCreateConsumer_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *ConsumerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
			errMsg:  "config cannot be nil",
		},
		{
			name: "empty type",
			config: &ConsumerConfig{
				Type:       "",
				Connection: "test://localhost",
				Queue:      "test-queue",
			},
			wantErr: true,
			errMsg:  "consumer type cannot be empty",
		},
		{
			name: "empty connection",
			config: &ConsumerConfig{
				Type:       "test",
				Connection: "",
				Queue:      "test-queue",
			},
			wantErr: true,
			errMsg:  "consumer connection cannot be empty",
		},
		{
			name: "empty queue",
			config: &ConsumerConfig{
				Type:       "test",
				Connection: "test://localhost",
				Queue:      "",
			},
			wantErr: true,
			errMsg:  "consumer queue name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer, err := CreateConsumer(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				if consumer != nil {
					t.Errorf("expected nil consumer but got %v", consumer)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if consumer == nil {
					t.Error("expected consumer but got nil")
				}
			}
		})
	}

}

// TestGlobalRegistryFunctions tests the global registry functions
func TestGlobalRegistryFunctions(t *testing.T) {
	// Save original global registry and restore after test
	originalRegistry := globalRegistry
	defer func() {
		globalRegistry = originalRegistry
	}()

	// Reset global registry for testing
	globalRegistry = &Registry{
		factories: make(map[string]ConsumerFactory),
	}

	// Test Register function
	err := Register("test-queue", mockFactory)
	if err != nil {
		t.Errorf("Register() unexpected error: %v", err)
	}

	// Test Get function
	factory, err := Get("test-queue")
	if err != nil {
		t.Errorf("Get() unexpected error: %v", err)
	}
	if factory == nil {
		t.Error("Get() expected factory but got nil")
	}

	// Test List function
	err = Register("another-queue", mockFactory)
	if err != nil {
		t.Errorf("Register() unexpected error: %v", err)
	}

	list := List()
	if len(list) != 2 {
		t.Errorf("List() expected 2 items, got %d", len(list))
	}

	// Verify both queue types are in the list
	found := make(map[string]bool)
	for _, qType := range list {
		found[qType] = true
	}

	if !found["test-queue"] {
		t.Error("List() missing 'test-queue'")
	}
	if !found["another-queue"] {
		t.Error("List() missing 'another-queue'")
	}
}
