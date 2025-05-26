// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"fmt"
	"strings"
	"sync"
)

// Registry manages queue consumer factories
type Registry struct {
	mu        sync.RWMutex
	factories map[string]ConsumerFactory
}

// globalRegistry is the singleton registry instance
var globalRegistry = &Registry{
	factories: make(map[string]ConsumerFactory),
}

// Register adds a new consumer factory to the registry
func Register(queueType string, factory ConsumerFactory) error {
	return globalRegistry.Register(queueType, factory)
}

// Get retrieves a consumer factory from the registry
func Get(queueType string) (ConsumerFactory, error) {
	return globalRegistry.Get(queueType)
}

// List returns all registered queue types
func List() []string {
	return globalRegistry.List()
}

// Register adds a new consumer factory to this registry
func (r *Registry) Register(queueType string, factory ConsumerFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	queueType = strings.TrimSpace(queueType)
	if queueType == "" {
		return fmt.Errorf("queue type cannot be empty")
	}

	if factory == nil {
		return fmt.Errorf("factory cannot be nil")
	}

	if _, exists := r.factories[queueType]; exists {
		return fmt.Errorf("queue type %s already registered", queueType)
	}

	r.factories[queueType] = factory
	return nil
}

// Get retrieves a consumer factory from this registry
func (r *Registry) Get(queueType string) (ConsumerFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[queueType]
	if !exists {
		return nil, fmt.Errorf("queue type %s not registered", queueType)
	}

	return factory, nil
}

// List returns all registered queue types in this registry
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for queueType := range r.factories {
		types = append(types, queueType)
	}

	return types
}

// CreateConsumer creates a new consumer instance using the registered factory
func CreateConsumer(config *ConsumerConfig) (Consumer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if strings.TrimSpace(config.Type) == "" {
		return nil, fmt.Errorf("consumer type cannot be empty")
	}

	if strings.TrimSpace(config.Connection) == "" {
		return nil, fmt.Errorf("consumer connection cannot be empty")
	}

	if strings.TrimSpace(config.Queue) == "" {
		return nil, fmt.Errorf("consumer queue name cannot be empty")
	}

	factory, err := Get(config.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get factory for queue type %s: %w", config.Type, err)
	}

	consumer, err := factory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Validate the consumer configuration
	if err := consumer.Validate(); err != nil {
		return nil, fmt.Errorf("consumer validation failed: %w", err)
	}

	return consumer, nil
}
