// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file contains the coordinator implementation that manages multiple queue consumers
// and handles message processing, retries, and error handling.
package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rshade/cronai/internal/logger"
)

// Default logger for the queue package
var log = logger.DefaultLogger()

// SetLogger allows overriding the default logger
func SetLogger(l *logger.Logger) {
	log = l
}

// DefaultCoordinator implements the Coordinator interface
type DefaultCoordinator struct {
	consumers     map[string]Consumer
	mu            sync.RWMutex
	wg            sync.WaitGroup
	cancel        context.CancelFunc
	parser        MessageParser
	retryPolicy   RetryPolicy
	taskProcessor TaskProcessor
}

// TaskProcessor processes parsed task messages
type TaskProcessor interface {
	Process(ctx context.Context, task *TaskMessage) error
}

// NewCoordinator creates a new coordinator instance
func NewCoordinator(processor TaskProcessor, opts ...CoordinatorOption) Coordinator {
	c := &DefaultCoordinator{
		consumers:     make(map[string]Consumer),
		parser:        NewMessageParser(),
		retryPolicy:   NewExponentialBackoffRetryPolicy(3, 1*time.Second, 30*time.Second),
		taskProcessor: processor,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// CoordinatorOption is a functional option for configuring the coordinator
type CoordinatorOption func(*DefaultCoordinator)

// WithParser sets a custom message parser
func WithParser(parser MessageParser) CoordinatorOption {
	return func(c *DefaultCoordinator) {
		c.parser = parser
	}
}

// WithRetryPolicy sets a custom retry policy
func WithRetryPolicy(policy RetryPolicy) CoordinatorOption {
	return func(c *DefaultCoordinator) {
		c.retryPolicy = policy
	}
}

// AddConsumer adds a new consumer to the coordinator
func (c *DefaultCoordinator) AddConsumer(name string, consumer Consumer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "" {
		return fmt.Errorf("consumer name cannot be empty")
	}

	if consumer == nil {
		return fmt.Errorf("consumer cannot be nil")
	}

	if _, exists := c.consumers[name]; exists {
		return fmt.Errorf("consumer %s already registered", name)
	}

	c.consumers[name] = consumer
	log.Info("Added consumer", logger.Fields{"name": name, "type": consumer.Name()})

	return nil
}

// RemoveConsumer removes a consumer from the coordinator
func (c *DefaultCoordinator) RemoveConsumer(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	consumer, exists := c.consumers[name]
	if !exists {
		return fmt.Errorf("consumer %s not found", name)
	}

	delete(c.consumers, name)
	log.Info("Removed consumer", logger.Fields{"name": name, "type": consumer.Name()})

	return nil
}

// Start begins processing messages from all consumers
func (c *DefaultCoordinator) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.consumers) == 0 {
		return fmt.Errorf("no consumers registered")
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	// Start each consumer
	for name, consumer := range c.consumers {
		c.wg.Add(1)
		go c.runConsumer(ctx, name, consumer)
	}

	log.Info("Coordinator started", logger.Fields{"consumers": len(c.consumers)})

	return nil
}

// Stop gracefully shuts down all consumers
func (c *DefaultCoordinator) Stop(ctx context.Context) error {
	log.Info("Stopping coordinator")

	// Cancel the context to signal consumers to stop
	if c.cancel != nil {
		c.cancel()
	}

	// Wait for all consumers to finish with timeout
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("All consumers stopped gracefully")
		return nil
	case <-ctx.Done():
		log.Warn("Timeout waiting for consumers to stop")
		return fmt.Errorf("timeout waiting for consumers to stop")
	}
}

// GetConsumer retrieves a specific consumer by name
func (c *DefaultCoordinator) GetConsumer(name string) (Consumer, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	consumer, exists := c.consumers[name]
	return consumer, exists
}

// ListConsumers returns names of all registered consumers
func (c *DefaultCoordinator) ListConsumers() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.consumers))
	for name := range c.consumers {
		names = append(names, name)
	}

	return names
}

// runConsumer runs a single consumer
func (c *DefaultCoordinator) runConsumer(ctx context.Context, name string, consumer Consumer) {
	defer c.wg.Done()

	log.Info("Starting consumer", logger.Fields{"name": name})

	// Connect to the queue
	if err := consumer.Connect(ctx); err != nil {
		log.Error("Failed to connect consumer", logger.Fields{"name": name, "error": err.Error()})
		return
	}

	defer func() {
		if err := consumer.Disconnect(ctx); err != nil {
			log.Error("Failed to disconnect consumer", logger.Fields{"name": name, "error": err.Error()})
		}
	}()

	// Start consuming messages
	messages, errors := consumer.Consume(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Info("Consumer stopped", logger.Fields{"name": name})
			return

		case err, ok := <-errors:
			if !ok {
				log.Info("Error channel closed", logger.Fields{"name": name})
				return
			}
			if err != nil {
				log.Error("Consumer error", logger.Fields{"name": name, "error": err.Error()})
			}

		case msg, ok := <-messages:
			if !ok {
				log.Info("Message channel closed", logger.Fields{"name": name})
				return
			}
			if msg != nil {
				c.processMessage(ctx, name, consumer, msg)
			}
		}
	}
}

// processMessage processes a single message
func (c *DefaultCoordinator) processMessage(ctx context.Context, consumerName string, consumer Consumer, msg *Message) {
	log.Debug("Processing message", logger.Fields{"consumer": consumerName, "messageId": msg.ID})

	// Parse the message
	task, err := c.parser.Parse(msg)
	if err != nil {
		log.Error("Failed to parse message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
		// Parsing errors are not retryable
		if err := consumer.Reject(ctx, msg, false); err != nil {
			log.Error("Failed to reject message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
		}
		return
	}

	// Validate the task
	if err := c.parser.Validate(task); err != nil {
		log.Error("Invalid task message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
		// Validation errors are not retryable
		if err := consumer.Reject(ctx, msg, false); err != nil {
			log.Error("Failed to reject message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
		}
		return
	}

	// Process the task
	err = c.taskProcessor.Process(ctx, task)
	if err != nil {
		log.Error("Failed to process task", logger.Fields{"messageId": msg.ID, "error": err.Error()})

		// Check if we should retry
		if c.retryPolicy.ShouldRetry(msg, err) {
			msg.RetryCount++
			delay := c.retryPolicy.NextRetryDelay(msg)
			log.Info("Retrying message", logger.Fields{"messageId": msg.ID, "retryCount": msg.RetryCount, "delay": delay.String()})

			// Reject with requeue
			if err := consumer.Reject(ctx, msg, true); err != nil {
				log.Error("Failed to requeue message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
			}
		} else {
			// Max retries exceeded or non-retryable error
			log.Error("Message processing failed permanently", logger.Fields{"messageId": msg.ID, "retryCount": msg.RetryCount})
			if err := consumer.Reject(ctx, msg, false); err != nil {
				log.Error("Failed to reject message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
			}
		}
		return
	}

	// Acknowledge successful processing
	if err := consumer.Acknowledge(ctx, msg); err != nil {
		log.Error("Failed to acknowledge message", logger.Fields{"messageId": msg.ID, "error": err.Error()})
	} else {
		log.Info("Message processed successfully", logger.Fields{"messageId": msg.ID})
	}
}
