// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file defines the core interfaces and data structures used throughout the queue
// subsystem, including Consumer, Coordinator, MessageParser, and RetryPolicy interfaces.
package queue

import (
	"context"
	"time"
)

// Message represents a task message from the queue
type Message struct {
	ID          string            // Unique message identifier
	Body        []byte            // Raw message body
	Attributes  map[string]string // Message attributes/metadata
	ReceivedAt  time.Time         // When the message was received
	RetryCount  int               // Number of retry attempts
	QueueSource string            // Source queue identifier
}

// TaskMessage represents the parsed task data from a queue message
type TaskMessage struct {
	Model     string            // AI model to use
	Prompt    string            // Prompt file name or inline prompt
	Processor string            // Response processor to use
	Variables map[string]string // Variables for prompt substitution
	IsInline  bool              // Whether prompt is inline or file reference
}

// Consumer is the interface that all queue providers must implement
type Consumer interface {
	// Connect establishes connection to the queue
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the queue
	Disconnect(ctx context.Context) error

	// Consume starts consuming messages from the queue
	// Returns a channel that emits messages and a channel for errors
	Consume(ctx context.Context) (<-chan *Message, <-chan error)

	// Acknowledge marks a message as successfully processed
	Acknowledge(ctx context.Context, message *Message) error

	// Reject marks a message as failed and optionally requeues it
	Reject(ctx context.Context, message *Message, requeue bool) error

	// Name returns the name of the queue provider
	Name() string

	// Validate checks if the configuration is valid
	Validate() error
}

// ConsumerConfig holds configuration for a queue consumer
type ConsumerConfig struct {
	Type       string                 // Queue type (e.g., "rabbitmq", "sqs", "azure-servicebus")
	Connection string                 // Connection string or URL
	Queue      string                 // Queue name or topic
	Options    map[string]interface{} // Provider-specific options
	RetryLimit int                    // Maximum retry attempts
	RetryDelay time.Duration          // Delay between retries
}

// ConsumerFactory is a function that creates a new Consumer instance
type ConsumerFactory func(config *ConsumerConfig) (Consumer, error)

// Coordinator manages multiple queue consumers
type Coordinator interface {
	// AddConsumer adds a new consumer to the coordinator
	AddConsumer(name string, consumer Consumer) error

	// RemoveConsumer removes a consumer from the coordinator
	RemoveConsumer(name string) error

	// Start begins processing messages from all consumers
	Start(ctx context.Context) error

	// Stop gracefully shuts down all consumers
	Stop(ctx context.Context) error

	// GetConsumer retrieves a specific consumer by name
	GetConsumer(name string) (Consumer, bool)

	// ListConsumers returns names of all registered consumers
	ListConsumers() []string
}

// MessageParser parses queue messages into task messages
type MessageParser interface {
	// Parse converts a raw message into a TaskMessage
	Parse(message *Message) (*TaskMessage, error)

	// Validate checks if a TaskMessage is valid
	Validate(task *TaskMessage) error
}

// RetryPolicy defines how failed messages should be retried
type RetryPolicy interface {
	// ShouldRetry determines if a message should be retried
	ShouldRetry(message *Message, err error) bool

	// NextRetryDelay calculates the delay before next retry
	NextRetryDelay(message *Message) time.Duration

	// MaxRetries returns the maximum number of retries
	MaxRetries() int
}
