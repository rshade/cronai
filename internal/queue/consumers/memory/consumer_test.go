// Package memory provides an in-memory consumer implementation for testing and development.
package memory

import (
	"context"
	"testing"
	"time"

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
				Type:  "memory",
				Queue: "test-queue",
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
				assert.Equal(t, "memory", consumer.Name())
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
				Type:  "memory",
				Queue: "test-queue",
			},
			wantErr: false,
		},
		{
			name: "empty queue name",
			config: &queue.ConsumerConfig{
				Type:  "memory",
				Queue: "",
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

func TestConsumer_ConnectDisconnect(t *testing.T) {
	config := &queue.ConsumerConfig{
		Type:  "memory",
		Queue: "test-queue",
	}

	consumer, err := NewConsumer(config)
	require.NoError(t, err)

	ctx := context.Background()

	// Test connect
	err = consumer.Connect(ctx)
	assert.NoError(t, err)

	// Test disconnect
	err = consumer.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestConsumer_ConsumeAndProcess(t *testing.T) {
	config := &queue.ConsumerConfig{
		Type:  "memory",
		Queue: "test-queue",
	}

	consumer, err := NewConsumer(config)
	require.NoError(t, err)

	memConsumer, ok := consumer.(*Consumer)
	require.True(t, ok, "consumer should be of type *Consumer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect
	err = consumer.Connect(ctx)
	require.NoError(t, err)

	// Start consuming
	messages, errors := consumer.Consume(ctx)

	// Add a test message
	testMsg := &queue.Message{
		ID:         "test-msg-1",
		Body:       []byte(`{"model":"openai","prompt":"test","processor":"console"}`),
		Attributes: map[string]string{"source": "test"},
		ReceivedAt: time.Now(),
	}

	err = memConsumer.AddMessage(testMsg)
	require.NoError(t, err)

	// Receive the message
	select {
	case msg := <-messages:
		assert.Equal(t, testMsg.ID, msg.ID)
		assert.Equal(t, testMsg.Body, msg.Body)

		// Test acknowledge
		err = consumer.Acknowledge(ctx, msg)
		assert.NoError(t, err)

		// Check it was acknowledged
		ackMsgs := memConsumer.GetAcknowledgedMessages()
		assert.Contains(t, ackMsgs, msg.ID)

	case err := <-errors:
		t.Fatalf("Unexpected error: %v", err)

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	// Test reject with requeue
	testMsg2 := &queue.Message{
		ID:         "test-msg-2",
		Body:       []byte(`{"model":"claude","prompt":"test2","processor":"console"}`),
		Attributes: map[string]string{"source": "test"},
		ReceivedAt: time.Now(),
	}

	err = memConsumer.AddMessage(testMsg2)
	require.NoError(t, err)

	select {
	case msg := <-messages:
		// Test reject
		err = consumer.Reject(ctx, msg, false)
		assert.NoError(t, err)

		// Check it was rejected
		rejMsgs := memConsumer.GetRejectedMessages()
		assert.Contains(t, rejMsgs, msg.ID)

	case err := <-errors:
		t.Fatalf("Unexpected error: %v", err)

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	// Disconnect
	err = consumer.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestConsumer_RejectWithRequeue(t *testing.T) {
	config := &queue.ConsumerConfig{
		Type:  "memory",
		Queue: "test-queue",
	}

	consumer, err := NewConsumer(config)
	require.NoError(t, err)

	memConsumer, ok := consumer.(*Consumer)
	require.True(t, ok, "consumer should be of type *Consumer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect
	err = consumer.Connect(ctx)
	require.NoError(t, err)

	// Start consuming
	messages, errors := consumer.Consume(ctx)

	// Add a test message
	testMsg := &queue.Message{
		ID:         "test-msg-requeue",
		Body:       []byte(`{"model":"openai","prompt":"test","processor":"console"}`),
		Attributes: map[string]string{"source": "test"},
		ReceivedAt: time.Now(),
	}

	err = memConsumer.AddMessage(testMsg)
	require.NoError(t, err)

	// Receive the message
	select {
	case msg := <-messages:
		// Test reject with requeue
		err = consumer.Reject(ctx, msg, true)
		assert.NoError(t, err)

		// Wait a bit for requeue to happen
		time.Sleep(200 * time.Millisecond)

		// Should receive the requeued message
		select {
		case requeuedMsg := <-messages:
			assert.Equal(t, msg.ID, requeuedMsg.ID)
			// Acknowledge the requeued message
			err = consumer.Acknowledge(ctx, requeuedMsg)
			assert.NoError(t, err)

		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for requeued message")
		}

	case err := <-errors:
		t.Fatalf("Unexpected error: %v", err)

	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for initial message")
	}

	// Disconnect
	err = consumer.Disconnect(ctx)
	assert.NoError(t, err)
}
