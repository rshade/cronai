// Package processor provides functionality for processing model responses through various output channels
package processor

import (
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
)

// MockProcessor is a mock implementation of the Processor interface for testing
type MockProcessor struct {
	config         Config
	processError   error
	validateError  error
	processCalled  bool
	validateCalled bool
}

func (m *MockProcessor) Process(_ *models.ModelResponse, _ string) error {
	m.processCalled = true
	return m.processError
}

func (m *MockProcessor) Validate() error {
	m.validateCalled = true
	return m.validateError
}

func (m *MockProcessor) GetType() string {
	return m.config.Type
}

func (m *MockProcessor) GetConfig() Config {
	return m.config
}

func TestProcessorInterface(t *testing.T) {
	// Create test processor
	config := Config{
		Type:   "test",
		Target: "test-target",
	}
	processor := &MockProcessor{
		config: config,
	}

	// Test Process method
	response := &models.ModelResponse{
		Content:    "Test content",
		Model:      "test-model",
		Timestamp:  time.Now(),
		PromptName: "test-prompt",
	}
	err := processor.Process(response, "")
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}

	// Test Validate method
	err = processor.Validate()
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	// Test GetType method
	processorType := processor.GetType()
	if processorType != "test" {
		t.Errorf("GetType returned wrong type: got %s, want test", processorType)
	}

	// Test GetConfig method
	config = processor.GetConfig()
	if config.Type != "test" || config.Target != "test-target" {
		t.Errorf("GetConfig returned wrong config: got %+v, want {Type:test Target:test-target}", config)
	}
}

func TestMockProcessor(t *testing.T) {
	// Test mock processor behavior
	mock := &MockProcessor{
		config: Config{
			Type:   "mock",
			Target: "mock-target",
		},
	}

	// Test GetType
	if mock.GetType() != "mock" {
		t.Errorf("Expected GetType to return 'mock', got '%s'", mock.GetType())
	}

	// Test GetConfig
	config := mock.GetConfig()
	if config.Type != "mock" || config.Target != "mock-target" {
		t.Errorf("GetConfig returned unexpected values: %+v", config)
	}

	// Test Process
	response := &models.ModelResponse{
		Content: "test content",
	}
	err := mock.Process(response, "test-template")
	if err != nil {
		t.Errorf("Process returned unexpected error: %v", err)
	}
	if !mock.processCalled {
		t.Error("Process was not called")
	}

	// Test Validate
	err = mock.Validate()
	if err != nil {
		t.Errorf("Validate returned unexpected error: %v", err)
	}
	if !mock.validateCalled {
		t.Error("Validate was not called")
	}
}
