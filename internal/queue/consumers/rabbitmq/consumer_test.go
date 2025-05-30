// Package rabbitmq provides a RabbitMQ consumer implementation for the CronAI queue system.
package rabbitmq

import (
	"testing"

	"github.com/rshade/cronai/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsumer(t *testing.T) {
	tests := []struct {
		name    string
		config  *queue.ConsumerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &queue.ConsumerConfig{
				Type:       "rabbitmq",
				Connection: "amqp://guest:guest@localhost:5672/",
				Queue:      "test-queue",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer, err := NewConsumer(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, consumer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, consumer)
				assert.Equal(t, "rabbitmq", consumer.Name())
			}
		})
	}
}

func TestConsumer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *queue.ConsumerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &queue.ConsumerConfig{
				Type:       "rabbitmq",
				Connection: "amqp://guest:guest@localhost:5672/",
				Queue:      "test-queue",
			},
			wantErr: false,
		},
		{
			name: "empty connection",
			config: &queue.ConsumerConfig{
				Type:       "rabbitmq",
				Connection: "",
				Queue:      "test-queue",
			},
			wantErr: true,
		},
		{
			name: "empty queue name",
			config: &queue.ConsumerConfig{
				Type:       "rabbitmq",
				Connection: "amqp://guest:guest@localhost:5672/",
				Queue:      "",
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := &Consumer{config: tt.config}
			err := consumer.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsumer_Name(t *testing.T) {
	config := &queue.ConsumerConfig{
		Type:       "rabbitmq",
		Connection: "amqp://guest:guest@localhost:5672/",
		Queue:      "test-queue",
	}

	consumer, err := NewConsumer(config)
	require.NoError(t, err)

	assert.Equal(t, "rabbitmq", consumer.Name())
}

// Note: These tests would require a running RabbitMQ instance to test actual connections.
// For CI/CD purposes, we focus on unit tests for validation and creation logic.
// Integration tests with actual RabbitMQ should be in a separate test suite.

func TestConsumer_ConfigurationHandling(t *testing.T) {
	config := &queue.ConsumerConfig{
		Type:       "rabbitmq",
		Connection: "amqp://guest:guest@localhost:5672/",
		Queue:      "test-queue",
		Options: map[string]interface{}{
			"durable":    true,
			"autoDelete": false,
		},
	}

	consumer, err := NewConsumer(config)
	require.NoError(t, err)

	rabbitConsumer, ok := consumer.(*Consumer)
	require.True(t, ok, "consumer should be of type *Consumer")
	assert.Equal(t, config, rabbitConsumer.config)
	assert.NotNil(t, rabbitConsumer.logger)
}
