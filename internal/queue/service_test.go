// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.coordinator)
	assert.NotNil(t, service.taskProcessor)
	assert.NotNil(t, service.logger)
}

func TestService_loadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		expected *ServiceConfig
	}{
		{
			name: "memory queue from env",
			envVars: map[string]string{
				"QUEUE_TYPE": "memory",
				"QUEUE_NAME": "test-queue",
			},
			wantErr: false,
			expected: &ServiceConfig{
				Consumers: []ConsumerConfig{
					{
						Type:       "memory",
						Connection: "memory://localhost",
						Queue:      "test-queue",
					},
				},
			},
		},
		{
			name: "rabbitmq queue from env",
			envVars: map[string]string{
				"QUEUE_TYPE":       "rabbitmq",
				"QUEUE_CONNECTION": "amqp://guest:guest@localhost:5672/",
				"QUEUE_NAME":       "cronai-tasks",
			},
			wantErr: false,
			expected: &ServiceConfig{
				Consumers: []ConsumerConfig{
					{
						Type:       "rabbitmq",
						Connection: "amqp://guest:guest@localhost:5672/",
						Queue:      "cronai-tasks",
					},
				},
			},
		},
		{
			name:    "no config - should create default",
			envVars: map[string]string{},
			wantErr: false,
			expected: &ServiceConfig{
				Consumers: []ConsumerConfig{
					{
						Type:       "memory",
						Connection: "memory://localhost",
						Queue:      "cronai-tasks",
					},
				},
			},
		},
		{
			name: "unsupported queue type",
			envVars: map[string]string{
				"QUEUE_TYPE": "unsupported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("failed to set environment variable %s: %v", key, err)
				}
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					if err := os.Unsetenv(key); err != nil {
						t.Logf("failed to unset environment variable %s: %v", key, err)
					}
				}
			}()

			service := NewService()
			config, err := service.loadConfig("")

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, config)
			assert.Len(t, config.Consumers, len(tt.expected.Consumers))

			if len(config.Consumers) > 0 {
				consumer := config.Consumers[0]
				expected := tt.expected.Consumers[0]

				assert.Equal(t, expected.Type, consumer.Type)
				assert.Equal(t, expected.Connection, consumer.Connection)
				assert.Equal(t, expected.Queue, consumer.Queue)
				assert.Equal(t, 3, consumer.RetryLimit) // Default
				assert.NotZero(t, consumer.RetryDelay)  // Default
			}
		})
	}
}

func TestGetSupportedQueueTypes(t *testing.T) {
	// Initially, no types should be registered
	types := GetSupportedQueueTypes()
	initialCount := len(types)

	// Register a test type
	testFactory := func(_ *ConsumerConfig) (Consumer, error) {
		return nil, nil
	}
	err := Register("test-type", testFactory)
	require.NoError(t, err)

	// Check that the new type is included
	types = GetSupportedQueueTypes()
	assert.Greater(t, len(types), initialCount)
	assert.Contains(t, types, "test-type")
}

func TestValidateQueueConfig(t *testing.T) {
	tests := []struct {
		name       string
		queueType  string
		connection string
		queueName  string
		wantErr    bool
	}{
		{
			name:       "valid config",
			queueType:  "memory",
			connection: "memory://localhost",
			queueName:  "test-queue",
			wantErr:    false,
		},
		{
			name:       "empty queue type",
			queueType:  "",
			connection: "memory://localhost",
			queueName:  "test-queue",
			wantErr:    true,
		},
		{
			name:       "empty queue name",
			queueType:  "memory",
			connection: "memory://localhost",
			queueName:  "",
			wantErr:    true,
		},
		{
			name:       "whitespace queue type",
			queueType:  "  ",
			connection: "memory://localhost",
			queueName:  "test-queue",
			wantErr:    true,
		},
		{
			name:       "whitespace queue name",
			queueType:  "memory",
			connection: "memory://localhost",
			queueName:  "  ",
			wantErr:    true,
		},
	}

	// Register a memory type for testing
	testFactory := func(_ *ConsumerConfig) (Consumer, error) {
		return nil, nil
	}
	if err := Register("memory", testFactory); err != nil {
		t.Fatalf("failed to register memory consumer: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQueueConfig(tt.queueType, tt.connection, tt.queueName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
