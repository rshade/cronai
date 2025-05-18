package processor

import (
	"testing"
	"time"

	"github.com/rshade/cronai/internal/models"
)

// MockProcessor is a mock implementation of the Processor interface for testing
type MockProcessor struct {
	config         ProcessorConfig
	processError   error
	validateError  error
	processCalled  bool
	validateCalled bool
}

func (m *MockProcessor) Process(response *models.ModelResponse, templateName string) error {
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

func (m *MockProcessor) GetConfig() ProcessorConfig {
	return m.config
}

func TestProcessorInterface(t *testing.T) {
	// Test that each processor implements the interface
	processors := []struct {
		name      string
		processor Processor
	}{
		{
			name: "EmailProcessor",
			processor: &EmailProcessor{
				config: ProcessorConfig{Type: "email", Target: "test@example.com"},
			},
		},
		{
			name: "SlackProcessor",
			processor: &SlackProcessor{
				config: ProcessorConfig{Type: "slack", Target: "test-channel"},
			},
		},
		{
			name: "WebhookProcessor",
			processor: &WebhookProcessor{
				config: ProcessorConfig{Type: "webhook", Target: "test-webhook"},
			},
		},
		{
			name: "FileProcessor",
			processor: &FileProcessor{
				config: ProcessorConfig{Type: "file"},
			},
		},
		{
			name: "ConsoleProcessor",
			processor: &ConsoleProcessor{
				config: ProcessorConfig{Type: "console"},
			},
		},
	}

	for _, p := range processors {
		t.Run(p.name, func(t *testing.T) {
			// Test GetType
			if p.processor.GetType() == "" {
				t.Errorf("GetType() returned empty string for %s", p.name)
			}

			// Test GetConfig
			config := p.processor.GetConfig()
			if config.Type != p.processor.GetType() {
				t.Errorf("GetConfig().Type (%s) doesn't match GetType() (%s) for %s",
					config.Type, p.processor.GetType(), p.name)
			}

			// Test that Process method exists (compilation would fail if not)
			response := &models.ModelResponse{
				Content:     "test",
				Model:       "test-model",
				PromptName:  "test-prompt",
				Timestamp:   time.Now(),
				ExecutionID: "test-exec",
			}
			_ = p.processor.Process(response, "test-template")

			// Test that Validate method exists
			_ = p.processor.Validate()
		})
	}
}

func TestMockProcessor(t *testing.T) {
	// Test mock processor behavior
	mock := &MockProcessor{
		config: ProcessorConfig{
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
