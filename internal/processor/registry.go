package processor

import (
	"fmt"
	"sync"

	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
)

// Registry manages available processors
type Registry struct {
	factories map[string]ProcessorFactory
	mu        sync.RWMutex
}

// global registry instance
var (
	registry *Registry
	once     sync.Once
)

// GetRegistry returns the global processor registry
func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{
			factories: make(map[string]ProcessorFactory),
		}
		// Register default processors
		registry.RegisterDefaults()
	})
	return registry
}

// RegisterProcessor registers a new processor factory
func (r *Registry) RegisterProcessor(processorType string, factory ProcessorFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if processorType == "" {
		return errors.Wrap(errors.CategoryValidation, fmt.Errorf("processor type cannot be empty"),
			"invalid processor registration")
	}

	if factory == nil {
		return errors.Wrap(errors.CategoryValidation, fmt.Errorf("processor factory cannot be nil"),
			"invalid processor registration")
	}

	r.factories[processorType] = factory
	log.Debug("Registered processor", logger.Fields{
		"type": processorType,
	})
	return nil
}

// CreateProcessor creates a new processor instance from configuration
func (r *Registry) CreateProcessor(config ProcessorConfig) (Processor, error) {
	r.mu.RLock()
	factory, exists := r.factories[config.Type]
	r.mu.RUnlock()

	if !exists {
		return nil, errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("unknown processor type: %s", config.Type),
			"processor not registered")
	}

	processor, err := factory(config)
	if err != nil {
		return nil, errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("failed to create processor: %w", err),
			"processor creation failed")
	}

	// Validate the processor before returning
	if err := processor.Validate(); err != nil {
		return nil, errors.Wrap(errors.CategoryValidation, err,
			"processor validation failed")
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
	_ = r.RegisterProcessor("email", func(config ProcessorConfig) (Processor, error) {
		return NewEmailProcessor(config)
	})

	// Register slack processor
	_ = r.RegisterProcessor("slack", func(config ProcessorConfig) (Processor, error) {
		return NewSlackProcessor(config)
	})

	// Register webhook processor
	_ = r.RegisterProcessor("webhook", func(config ProcessorConfig) (Processor, error) {
		return NewWebhookProcessor(config)
	})

	// Register file processor
	_ = r.RegisterProcessor("file", func(config ProcessorConfig) (Processor, error) {
		return NewFileProcessor(config)
	})

	// Register console processor
	_ = r.RegisterProcessor("console", func(config ProcessorConfig) (Processor, error) {
		return NewConsoleProcessor(config)
	})
}
