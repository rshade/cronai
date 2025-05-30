// Package consumers provides initialization for all queue consumer implementations.
package consumers

import (
	"github.com/rshade/cronai/internal/queue"
	"github.com/rshade/cronai/internal/queue/consumers/memory"
	"github.com/rshade/cronai/internal/queue/consumers/rabbitmq"
)

// RegisterAll registers all available queue consumers with the registry
func RegisterAll() error {
	// Register memory consumer
	if err := queue.Register("memory", memory.NewConsumer); err != nil {
		return err
	}

	// Register RabbitMQ consumer
	if err := queue.Register("rabbitmq", rabbitmq.NewConsumer); err != nil {
		return err
	}

	return nil
}
