package processor

import (
	"fmt"
	"sync"
)

// Registry manages available processors
type Registry struct {
	factories map[string]Factory
	mu        sync.RWMutex
}

// global registry instance
var (
	registry *Registry
	once     sync.Once
)

// GetRegistry returns the singleton registry instance
func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{
			factories: make(map[string]Factory),
		}
		// Register default processors
		registry.RegisterFactory("console", NewConsoleProcessor)
		registry.RegisterFactory("file", NewFileProcessor)
		registry.RegisterFactory("email", NewEmailProcessor)
		registry.RegisterFactory("slack", NewSlackProcessor)
		registry.RegisterFactory("webhook", NewWebhookProcessor)
		registry.RegisterFactory("github", NewGitHubProcessor)
	})
	return registry
}

// RegisterFactory registers a new processor factory
func (r *Registry) RegisterFactory(processorType string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[processorType] = factory
}

// RegisterProcessor registers a new processor factory
func (r *Registry) RegisterProcessor(processorType string, factory Factory) {
	r.factories[processorType] = factory
}

// CreateProcessor creates a new processor instance
func (r *Registry) CreateProcessor(processorType string, config Config) (Processor, error) {
	r.mu.RLock()
	factory, exists := r.factories[processorType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown processor type: %s", processorType)
	}

	processor, err := factory(config)
	if err != nil {
		return nil, err
	}

	// Validate the processor
	if err := processor.Validate(); err != nil {
		return nil, fmt.Errorf("processor validation failed: %w", err)
	}

	return processor, nil
}

// GetProcessorTypes returns all registered processor types
func (r *Registry) GetProcessorTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}

// RegisterDefaults registers all default processors
func (r *Registry) RegisterDefaults() {
	// Register email processor
	r.RegisterFactory("email", NewEmailProcessor)

	// Register slack processor
	r.RegisterFactory("slack", NewSlackProcessor)

	// Register webhook processor
	r.RegisterFactory("webhook", NewWebhookProcessor)

	// Register file processor
	r.RegisterFactory("file", NewFileProcessor)

	// Register console processor
	r.RegisterFactory("console", NewConsoleProcessor)

	// Register github processor
	r.RegisterFactory("github", NewGitHubProcessor)
}

// GetProcessorFunc returns a function that creates a processor
func GetProcessorFunc(config Config) func() (Processor, error) {
	return func() (Processor, error) {
		return CreateProcessor(config.Type, config)
	}
}

// GetProcessor creates a processor from the given config
func GetProcessor(config Config) (Processor, error) {
	return CreateProcessor(config.Type, config)
}

// CreateProcessor creates a new processor using the global registry
func CreateProcessor(processorType string, config Config) (Processor, error) {
	return GetRegistry().CreateProcessor(processorType, config)
}
