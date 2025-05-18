package processor

import (
	"fmt"
	"testing"

	"github.com/rshade/cronai/internal/models"
)

func TestRegistry(t *testing.T) {
	// Create a new registry for testing
	registry := &Registry{
		factories: make(map[string]ProcessorFactory),
	}

	// Test RegisterProcessor
	t.Run("RegisterProcessor", func(t *testing.T) {
		// Test with valid processor
		err := registry.RegisterProcessor("test", func(config ProcessorConfig) (Processor, error) {
			return &MockProcessor{config: config}, nil
		})
		if err != nil {
			t.Errorf("RegisterProcessor failed: %v", err)
		}

		// Test with empty type
		err = registry.RegisterProcessor("", nil)
		if err == nil {
			t.Error("Expected error for empty processor type")
		}

		// Test with nil factory
		err = registry.RegisterProcessor("nil-factory", nil)
		if err == nil {
			t.Error("Expected error for nil factory")
		}
	})

	// Test CreateProcessor
	t.Run("CreateProcessor", func(t *testing.T) {
		// Register a test processor
		err := registry.RegisterProcessor("test", func(config ProcessorConfig) (Processor, error) {
			return &MockProcessor{config: config}, nil
		})
		if err != nil {
			t.Fatalf("Failed to register test processor: %v", err)
		}

		// Test creating registered processor
		config := ProcessorConfig{
			Type:   "test",
			Target: "test-target",
		}
		processor, err := registry.CreateProcessor(config)
		if err != nil {
			t.Errorf("CreateProcessor failed: %v", err)
		}
		if processor == nil {
			t.Error("Expected processor, got nil")
		}

		// Test creating unregistered processor
		config.Type = "unregistered"
		_, err = registry.CreateProcessor(config)
		if err == nil {
			t.Error("Expected error for unregistered processor type")
		}

		// Test processor that fails creation
		err = registry.RegisterProcessor("failing", func(config ProcessorConfig) (Processor, error) {
			return nil, fmt.Errorf("creation failed")
		})
		if err != nil {
			t.Fatalf("Failed to register failing processor: %v", err)
		}
		config.Type = "failing"
		_, err = registry.CreateProcessor(config)
		if err == nil {
			t.Error("Expected error for failing processor creation")
		}

		// Test processor that fails validation
		err = registry.RegisterProcessor("invalid", func(config ProcessorConfig) (Processor, error) {
			return &MockProcessor{
				config:        config,
				validateError: fmt.Errorf("validation failed"),
			}, nil
		})
		if err != nil {
			t.Fatalf("Failed to register invalid processor: %v", err)
		}
		config.Type = "invalid"
		_, err = registry.CreateProcessor(config)
		if err == nil {
			t.Error("Expected error for processor validation failure")
		}
	})

	// Test GetProcessorTypes
	t.Run("GetProcessorTypes", func(t *testing.T) {
		types := registry.GetProcessorTypes()
		expectedTypes := map[string]bool{
			"test":    true,
			"failing": true,
			"invalid": true,
		}

		if len(types) != len(expectedTypes) {
			t.Errorf("Expected %d processor types, got %d", len(expectedTypes), len(types))
		}

		for _, typ := range types {
			if !expectedTypes[typ] {
				t.Errorf("Unexpected processor type: %s", typ)
			}
		}
	})
}

func TestGlobalRegistry(t *testing.T) {
	// Test GetRegistry returns the same instance
	registry1 := GetRegistry()
	registry2 := GetRegistry()

	if registry1 != registry2 {
		t.Error("GetRegistry should return the same instance")
	}

	// Test that default processors are registered
	types := registry1.GetProcessorTypes()
	expectedTypes := []string{"email", "slack", "webhook", "file", "console"}

	for _, expected := range expectedTypes {
		found := false
		for _, actual := range types {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected processor type '%s' not found", expected)
		}
	}
}

func TestRegistryWithMockProcessors(t *testing.T) {
	registry := &Registry{
		factories: make(map[string]ProcessorFactory),
	}

	// Register a mock processor factory
	mockFactory := func(config ProcessorConfig) (Processor, error) {
		return &MockProcessor{
			config: config,
		}, nil
	}

	err := registry.RegisterProcessor("mock", mockFactory)
	if err != nil {
		t.Fatalf("Failed to register mock processor: %v", err)
	}

	// Create a processor using the factory
	config := ProcessorConfig{
		Type:   "mock",
		Target: "test-target",
		Options: map[string]interface{}{
			"option1": "value1",
		},
	}

	processor, err := registry.CreateProcessor(config)
	if err != nil {
		t.Fatalf("Failed to create processor: %v", err)
	}

	// Verify the processor was created correctly
	if processor.GetType() != "mock" {
		t.Errorf("Expected processor type 'mock', got '%s'", processor.GetType())
	}

	createdConfig := processor.GetConfig()
	if createdConfig.Target != "test-target" {
		t.Errorf("Expected target 'test-target', got '%s'", createdConfig.Target)
	}

	// Test processing
	response := &models.ModelResponse{
		Content:     "test content",
		Model:       "test-model",
		PromptName:  "test-prompt",
		ExecutionID: "test-exec",
	}

	err = processor.Process(response, "test-template")
	if err != nil {
		t.Errorf("Process failed: %v", err)
	}
}
