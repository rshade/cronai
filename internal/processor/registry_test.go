package processor

import (
	"fmt"
	"testing"

	"github.com/rshade/cronai/internal/models"
)

func TestRegistry(t *testing.T) {
	// Create a new registry for testing
	registry := &Registry{
		factories: make(map[string]Factory),
	}

	// Test RegisterFactory
	t.Run("RegisterFactory", func(t *testing.T) {
		// Test with valid processor
		registry.RegisterFactory("test", func(config Config) (Processor, error) {
			return &MockProcessor{config: config}, nil
		})

		// Verify the factory was registered
		if _, exists := registry.factories["test"]; !exists {
			t.Error("Factory was not registered")
		}
	})

	// Test CreateProcessor
	t.Run("CreateProcessor", func(t *testing.T) {
		// Register a test processor
		registry.RegisterFactory("test", func(config Config) (Processor, error) {
			return &MockProcessor{config: config}, nil
		})

		// Test creating registered processor
		config := Config{
			Type:   "test",
			Target: "test-target",
		}
		processor, err := registry.CreateProcessor("test", config)
		if err != nil {
			t.Errorf("CreateProcessor failed: %v", err)
		}
		if processor == nil {
			t.Error("Expected processor, got nil")
		}

		// Test creating unregistered processor
		config.Type = "unregistered"
		_, err = registry.CreateProcessor("unregistered", config)
		if err == nil {
			t.Error("Expected error for unregistered processor type")
		}

		// Test processor that fails creation
		registry.RegisterFactory("failing", func(_ Config) (Processor, error) {
			return nil, fmt.Errorf("creation failed")
		})
		config.Type = "failing"
		_, err = registry.CreateProcessor("failing", config)
		if err == nil {
			t.Error("Expected error for failing processor creation")
		}

		// Test processor that fails validation
		registry.RegisterFactory("invalid", func(config Config) (Processor, error) {
			return &MockProcessor{
				config:        config,
				validateError: fmt.Errorf("validation failed"),
			}, nil
		})
		config.Type = "invalid"
		_, err = registry.CreateProcessor("invalid", config)
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
	expectedTypes := []string{"email", "slack", "webhook", "file", "console", "github"}

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
		factories: make(map[string]Factory),
	}

	// Register a mock processor factory
	mockFactory := func(config Config) (Processor, error) {
		return &MockProcessor{
			config: config,
		}, nil
	}

	registry.RegisterFactory("mock", mockFactory)

	// Create a processor using the factory
	config := Config{
		Type:   "mock",
		Target: "test-target",
		Options: map[string]interface{}{
			"option1": "value1",
		},
	}

	processor, err := registry.CreateProcessor("mock", config)
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
