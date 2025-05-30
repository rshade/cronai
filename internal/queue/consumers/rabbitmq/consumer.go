// Package rabbitmq provides a RabbitMQ consumer implementation for the CronAI queue system.
package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/queue"
)

// Consumer implements queue.Consumer for RabbitMQ
type Consumer struct {
	config     *queue.ConsumerConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	mu         sync.RWMutex
	logger     *logger.Logger
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(config *queue.ConsumerConfig) (queue.Consumer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Consumer{
		config: config,
		logger: logger.DefaultLogger(),
	}, nil
}

// Name returns the name of this consumer
func (c *Consumer) Name() string {
	return "rabbitmq"
}

// Validate checks if the configuration is valid
func (c *Consumer) Validate() error {
	if c.config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if c.config.Connection == "" {
		return fmt.Errorf("connection string is required")
	}

	if c.config.Queue == "" {
		return fmt.Errorf("queue name is required")
	}

	return nil
}

// Connect establishes connection to RabbitMQ
func (c *Consumer) Connect(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Connecting to RabbitMQ", logger.Fields{"queue": c.config.Queue})

	// Connect to RabbitMQ
	conn, err := amqp.Dial(c.config.Connection)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	c.connection = conn

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		if closeErr := c.connection.Close(); closeErr != nil {
			c.logger.Error("failed to close connection after channel error", logger.Fields{"error": closeErr})
		}
		return fmt.Errorf("failed to open channel: %w", err)
	}
	c.channel = ch

	// Declare the queue (ensure it exists)
	queue, err := ch.QueueDeclare(
		c.config.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		if chErr := c.channel.Close(); chErr != nil {
			c.logger.Error("failed to close channel after queue declare error", logger.Fields{"error": chErr})
		}
		if connErr := c.connection.Close(); connErr != nil {
			c.logger.Error("failed to close connection after queue declare error", logger.Fields{"error": connErr})
		}
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	c.queue = queue

	// Set QoS to process one message at a time
	if err := ch.Qos(1, 0, false); err != nil {
		if chErr := c.channel.Close(); chErr != nil {
			c.logger.Error("failed to close channel after QoS error", logger.Fields{"error": chErr})
		}
		if connErr := c.connection.Close(); connErr != nil {
			c.logger.Error("failed to close connection after QoS error", logger.Fields{"error": connErr})
		}
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	c.logger.Info("Connected to RabbitMQ successfully", logger.Fields{
		"queue":    c.config.Queue,
		"messages": queue.Messages,
	})

	return nil
}

// Disconnect closes the connection to RabbitMQ
func (c *Consumer) Disconnect(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Disconnecting from RabbitMQ")

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			c.logger.Warn("Error closing channel", logger.Fields{"error": err.Error()})
		}
		c.channel = nil
	}

	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			c.logger.Warn("Error closing connection", logger.Fields{"error": err.Error()})
		}
		c.connection = nil
	}

	c.logger.Info("Disconnected from RabbitMQ")
	return nil
}

// Consume starts consuming messages from the queue
func (c *Consumer) Consume(ctx context.Context) (<-chan *queue.Message, <-chan error) {
	messages := make(chan *queue.Message)
	errors := make(chan error)

	go func() {
		defer close(messages)
		defer close(errors)

		c.mu.RLock()
		channel := c.channel
		queueName := c.config.Queue
		c.mu.RUnlock()

		if channel == nil {
			errors <- fmt.Errorf("not connected to RabbitMQ")
			return
		}

		// Start consuming messages
		deliveries, err := channel.Consume(
			queueName, // queue
			"",        // consumer
			false,     // auto-ack (we'll ack manually)
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			errors <- fmt.Errorf("failed to register consumer: %w", err)
			return
		}

		c.logger.Info("Started consuming messages", logger.Fields{"queue": queueName})

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("Consumer stopping due to context cancellation")
				return

			case delivery, ok := <-deliveries:
				if !ok {
					c.logger.Info("Delivery channel closed")
					return
				}

				// Convert AMQP delivery to queue message
				msg := &queue.Message{
					ID:          delivery.MessageId,
					Body:        delivery.Body,
					Attributes:  make(map[string]string),
					ReceivedAt:  time.Now(),
					QueueSource: queueName,
				}

				// If no message ID, generate one
				if msg.ID == "" {
					msg.ID = fmt.Sprintf("amqp-%d", time.Now().UnixNano())
				}

				// Copy headers to attributes
				for key, value := range delivery.Headers {
					if str, ok := value.(string); ok {
						msg.Attributes[key] = str
					}
				}

				// Store the delivery tag for acknowledgment
				msg.Attributes["_delivery_tag"] = fmt.Sprintf("%d", delivery.DeliveryTag)

				select {
				case messages <- msg:
					// Message sent successfully
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return messages, errors
}

// Acknowledge marks a message as successfully processed
func (c *Consumer) Acknowledge(_ context.Context, message *queue.Message) error {
	c.mu.RLock()
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Get the delivery tag from message attributes
	tagStr, exists := message.Attributes["_delivery_tag"]
	if !exists {
		return fmt.Errorf("delivery tag not found in message attributes")
	}

	var deliveryTag uint64
	if _, err := fmt.Sscanf(tagStr, "%d", &deliveryTag); err != nil {
		return fmt.Errorf("invalid delivery tag: %w", err)
	}

	if err := channel.Ack(deliveryTag, false); err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	c.logger.Debug("Message acknowledged", logger.Fields{"messageId": message.ID})
	return nil
}

// Reject marks a message as failed and optionally requeues it
func (c *Consumer) Reject(_ context.Context, message *queue.Message, requeue bool) error {
	c.mu.RLock()
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Get the delivery tag from message attributes
	tagStr, exists := message.Attributes["_delivery_tag"]
	if !exists {
		return fmt.Errorf("delivery tag not found in message attributes")
	}

	var deliveryTag uint64
	if _, err := fmt.Sscanf(tagStr, "%d", &deliveryTag); err != nil {
		return fmt.Errorf("invalid delivery tag: %w", err)
	}

	if err := channel.Nack(deliveryTag, false, requeue); err != nil {
		return fmt.Errorf("failed to reject message: %w", err)
	}

	c.logger.Debug("Message rejected", logger.Fields{
		"messageId": message.ID,
		"requeue":   requeue,
	})
	return nil
}
