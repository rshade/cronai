// Package memory provides an in-memory consumer implementation for testing and development.
package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/queue"
)

// Consumer implements queue.Consumer for in-memory message processing
type Consumer struct {
	config           *queue.ConsumerConfig
	messageQueue     chan *queue.Message
	mu               sync.RWMutex
	connected        bool
	logger           *logger.Logger
	acknowledgedMsgs map[string]bool
	rejectedMsgs     map[string]bool
	done             chan struct{}
	wg               sync.WaitGroup
}

// NewConsumer creates a new in-memory consumer
func NewConsumer(config *queue.ConsumerConfig) (queue.Consumer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Consumer{
		config:           config,
		messageQueue:     make(chan *queue.Message, 100), // Buffered channel
		logger:           logger.DefaultLogger(),
		acknowledgedMsgs: make(map[string]bool),
		rejectedMsgs:     make(map[string]bool),
		done:             make(chan struct{}),
	}, nil
}

// Name returns the name of this consumer
func (c *Consumer) Name() string {
	return "memory"
}

// Validate checks if the configuration is valid
func (c *Consumer) Validate() error {
	if c.config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if c.config.Queue == "" {
		return fmt.Errorf("queue name is required")
	}

	return nil
}

// Connect establishes connection (no-op for memory consumer)
func (c *Consumer) Connect(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Connecting to memory queue", logger.Fields{"queue": c.config.Queue})
	c.connected = true

	// Create a new done channel if it was closed
	select {
	case <-c.done:
		// Channel was closed, create a new one
		c.done = make(chan struct{})
	default:
		// Channel is still open, nothing to do
	}

	c.logger.Info("Connected to memory queue successfully", logger.Fields{"queue": c.config.Queue})
	return nil
}

// Disconnect closes the connection (no-op for memory consumer)
func (c *Consumer) Disconnect(_ context.Context) error {
	c.mu.Lock()
	c.logger.Info("Disconnecting from memory queue")
	c.connected = false

	// Signal shutdown to all goroutines
	select {
	case <-c.done:
		// Already closed
	default:
		close(c.done)
	}
	c.mu.Unlock()

	// Wait for all goroutines to finish (outside the lock to avoid deadlock)
	c.wg.Wait()

	// Now safe to close the message queue
	c.mu.Lock()
	close(c.messageQueue)
	c.messageQueue = make(chan *queue.Message, 100)
	c.mu.Unlock()

	c.logger.Info("Disconnected from memory queue")
	return nil
}

// Consume starts consuming messages from the memory queue
func (c *Consumer) Consume(ctx context.Context) (<-chan *queue.Message, <-chan error) {
	messages := make(chan *queue.Message)
	errors := make(chan error)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer close(messages)
		defer close(errors)

		c.logger.Info("Started consuming messages", logger.Fields{"queue": c.config.Queue})

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("Consumer stopping due to context cancellation")
				return

			case <-c.done:
				c.logger.Info("Consumer stopping due to disconnect")
				return

			case msg, ok := <-c.messageQueue:
				if !ok {
					c.logger.Info("Message queue closed")
					return
				}

				if msg != nil {
					select {
					case messages <- msg:
						// Message sent successfully
					case <-ctx.Done():
						return
					case <-c.done:
						return
					}
				}
			}
		}
	}()

	return messages, errors
}

// Acknowledge marks a message as successfully processed
func (c *Consumer) Acknowledge(_ context.Context, message *queue.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected to memory queue")
	}

	c.acknowledgedMsgs[message.ID] = true
	c.logger.Debug("Message acknowledged", logger.Fields{"messageId": message.ID})
	return nil
}

// Reject marks a message as failed and optionally requeues it
func (c *Consumer) Reject(_ context.Context, message *queue.Message, requeue bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected to memory queue")
	}

	c.rejectedMsgs[message.ID] = true

	if requeue {
		// For memory consumer, we can add the message back to the queue
		// In a real scenario, this would be handled by the queue system
		go func() {
			time.Sleep(100 * time.Millisecond) // Small delay for requeue
			select {
			case c.messageQueue <- message:
				c.logger.Debug("Message requeued", logger.Fields{"messageId": message.ID})
			default:
				c.logger.Warn("Failed to requeue message - queue full", logger.Fields{"messageId": message.ID})
			}
		}()
	}

	c.logger.Debug("Message rejected", logger.Fields{
		"messageId": message.ID,
		"requeue":   requeue,
	})
	return nil
}

// AddMessage adds a message to the memory queue (for testing)
func (c *Consumer) AddMessage(msg *queue.Message) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("not connected to memory queue")
	}

	select {
	case <-c.done:
		return fmt.Errorf("consumer is shutting down")
	case c.messageQueue <- msg:
		c.logger.Debug("Message added to queue", logger.Fields{"messageId": msg.ID})
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// GetAcknowledgedMessages returns the IDs of acknowledged messages (for testing)
func (c *Consumer) GetAcknowledgedMessages() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var messages []string
	for id := range c.acknowledgedMsgs {
		messages = append(messages, id)
	}
	return messages
}

// GetRejectedMessages returns the IDs of rejected messages (for testing)
func (c *Consumer) GetRejectedMessages() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var messages []string
	for id := range c.rejectedMsgs {
		messages = append(messages, id)
	}
	return messages
}
