// Package consumers provides initialization for all queue consumer implementations.
package consumers

import (
	"testing"

	"github.com/rshade/cronai/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAll(t *testing.T) {
	// Get initial state
	initialTypes := queue.List()

	// Register all consumers
	err := RegisterAll()
	require.NoError(t, err)

	// Check that new consumer types were registered
	types := queue.List()
	assert.Greater(t, len(types), len(initialTypes))

	// Check specific consumer types
	assert.Contains(t, types, "memory")
	assert.Contains(t, types, "rabbitmq")

	// Test that we can get the factories
	memoryFactory, err := queue.Get("memory")
	assert.NoError(t, err)
	assert.NotNil(t, memoryFactory)

	rabbitmqFactory, err := queue.Get("rabbitmq")
	assert.NoError(t, err)
	assert.NotNil(t, rabbitmqFactory)

	// Test creating consumers with the factories
	memoryConfig := &queue.ConsumerConfig{
		Type:  "memory",
		Queue: "test-queue",
	}
	memoryConsumer, err := memoryFactory(memoryConfig)
	assert.NoError(t, err)
	assert.NotNil(t, memoryConsumer)
	assert.Equal(t, "memory", memoryConsumer.Name())

	rabbitmqConfig := &queue.ConsumerConfig{
		Type:       "rabbitmq",
		Connection: "amqp://guest:guest@localhost:5672/",
		Queue:      "test-queue",
	}
	rabbitmqConsumer, err := rabbitmqFactory(rabbitmqConfig)
	assert.NoError(t, err)
	assert.NotNil(t, rabbitmqConsumer)
	assert.Equal(t, "rabbitmq", rabbitmqConsumer.Name())
}
