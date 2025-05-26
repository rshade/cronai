// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rshade/cronai/internal/logger"
)

// mockTaskProcessor for testing
type mockTaskProcessor struct {
	mu         sync.Mutex
	processed  []*TaskMessage
	failNext   bool
	alwaysFail bool
	delay      time.Duration
}

func (m *mockTaskProcessor) Process(_ context.Context, task *TaskMessage) error {
	// optional processing delay
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.failNext {
		m.failNext = false
		return fmt.Errorf("mock processing error")
	}

	if m.alwaysFail {
		return fmt.Errorf("mock processing error")
	}

	m.processed = append(m.processed, task)
	return nil
}

func (m *mockTaskProcessor) getProcessedCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.processed)
}

// testConsumer for testing
type testConsumer struct {
	name          string
	connected     bool
	messages      chan *Message
	errors        chan error
	ackCount      int
	rejectCount   int
	rejectRequeue bool
	mu            sync.Mutex
}

func newTestConsumer(name string) *testConsumer {
	return &testConsumer{
		name:     name,
		messages: make(chan *Message, 10),
		errors:   make(chan error, 10),
	}
}

func (t *testConsumer) Connect(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connected = true
	return nil
}

func (t *testConsumer) Disconnect(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connected = false
	close(t.messages)
	close(t.errors)
	return nil
}

func (t *testConsumer) Consume(_ context.Context) (<-chan *Message, <-chan error) {
	return t.messages, t.errors
}

func (t *testConsumer) Acknowledge(_ context.Context, _ *Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ackCount++
	return nil
}

func (t *testConsumer) Reject(_ context.Context, _ *Message, requeue bool) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rejectCount++
	t.rejectRequeue = requeue
	return nil
}

func (t *testConsumer) Name() string {
	return t.name
}

func (t *testConsumer) Validate() error {
	return nil
}

func (t *testConsumer) getAckCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.ackCount
}

func (t *testConsumer) getRejectCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.rejectCount
}

func TestDefaultCoordinator_AddConsumer(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	tests := []struct {
		name     string
		consName string
		consumer Consumer
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "successful add",
			consName: "test-consumer",
			consumer: newTestConsumer("test"),
			wantErr:  false,
		},
		{
			name:     "empty name",
			consName: "",
			consumer: newTestConsumer("test"),
			wantErr:  true,
			errMsg:   "consumer name cannot be empty",
		},
		{
			name:     "nil consumer",
			consName: "test",
			consumer: nil,
			wantErr:  true,
			errMsg:   "consumer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := coordinator.AddConsumer(tt.consName, tt.consumer)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}

	// Test duplicate registration
	consumer := newTestConsumer("duplicate")
	if err := coordinator.AddConsumer("dup", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}
	err := coordinator.AddConsumer("dup", consumer)
	if err == nil {
		t.Errorf("expected error for duplicate consumer")
	}
}

func TestDefaultCoordinator_RemoveConsumer(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add a consumer first
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Remove existing consumer
	err := coordinator.RemoveConsumer("test-consumer")
	if err != nil {
		t.Errorf("unexpected error removing consumer: %v", err)
	}

	// Try to remove non-existent consumer
	err = coordinator.RemoveConsumer("non-existent")
	if err == nil {
		t.Errorf("expected error for non-existent consumer")
	}
}

func TestDefaultCoordinator_GetConsumer(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add a consumer
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Get existing consumer
	retrieved, exists := coordinator.GetConsumer("test-consumer")
	if !exists {
		t.Errorf("expected consumer to exist")
	}
	if retrieved != consumer {
		t.Errorf("retrieved different consumer instance")
	}

	// Get non-existent consumer
	_, exists = coordinator.GetConsumer("non-existent")
	if exists {
		t.Errorf("expected consumer not to exist")
	}
}

func TestDefaultCoordinator_ListConsumers(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Initially empty
	list := coordinator.ListConsumers()
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d consumers", len(list))
	}

	// Add consumers
	names := []string{"consumer1", "consumer2", "consumer3"}
	for _, name := range names {
		consumer := newTestConsumer(name)
		if err := coordinator.AddConsumer(name, consumer); err != nil {
			t.Fatalf("failed to add consumer %s: %v", name, err)
		}
	}

	// Check list
	list = coordinator.ListConsumers()
	if len(list) != len(names) {
		t.Errorf("expected %d consumers, got %d", len(names), len(list))
	}

	// Verify all names are present
	nameMap := make(map[string]bool)
	for _, name := range list {
		nameMap[name] = true
	}
	for _, expectedName := range names {
		if !nameMap[expectedName] {
			t.Errorf("expected consumer %s not found", expectedName)
		}
	}
}

func TestDefaultCoordinator_StartStop(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Try to start with no consumers
	ctx := context.Background()
	err := coordinator.Start(ctx)
	if err == nil {
		t.Errorf("expected error starting with no consumers")
	}

	// Add a consumer
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start coordinator
	err = coordinator.Start(ctx)
	if err != nil {
		t.Errorf("unexpected error starting: %v", err)
	}

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}
}

func TestDefaultCoordinator_ProcessMessage(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add a consumer
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start coordinator
	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start coordinator: %v", err)
	}

	// Send a valid message
	validMsg := &Message{
		ID: "test-msg-1",
		Body: []byte(`{
			"model": "openai",
			"prompt": "test_prompt",
			"processor": "console"
		}`),
	}

	consumer.messages <- validMsg

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Check if message was processed
	if processor.getProcessedCount() != 1 {
		t.Errorf("expected 1 processed message, got %d", processor.getProcessedCount())
	}

	// Check if message was acknowledged
	if consumer.getAckCount() != 1 {
		t.Errorf("expected 1 acknowledged message, got %d", consumer.getAckCount())
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}
}

func TestDefaultCoordinator_InvalidMessage(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add a consumer
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start coordinator
	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start coordinator: %v", err)
	}

	// Send an invalid message
	invalidMsg := &Message{
		ID:   "test-msg-invalid",
		Body: []byte(`{invalid json`),
	}

	consumer.messages <- invalidMsg

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Check if message was rejected
	if consumer.getRejectCount() != 1 {
		t.Errorf("expected 1 rejected message, got %d", consumer.getRejectCount())
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}
}

func TestDefaultCoordinator_RetryPolicy(t *testing.T) {
	processor := &mockTaskProcessor{alwaysFail: true}
	coordinator := NewCoordinator(processor)

	// Add a consumer
	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start coordinator
	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start coordinator: %v", err)
	}

	// Send a message that will fail first time
	msg := &Message{
		ID: "test-retry",
		Body: []byte(`{
			"model": "openai",
			"prompt": "test_prompt",
			"processor": "console"
		}`),
		RetryCount: 0,
	}

	consumer.messages <- msg

	// Wait for processing and retry
	time.Sleep(300 * time.Millisecond)

	// Should have been rejected with requeue
	if consumer.getRejectCount() != 1 {
		t.Errorf("expected 1 rejection, got %d", consumer.getRejectCount())
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}
}

func TestDefaultCoordinator_Options(t *testing.T) {
	processor := &mockTaskProcessor{}
	customParser := NewMessageParser()
	customRetry := NewNoRetryPolicy()

	coordinator, ok := NewCoordinator(processor,
		WithParser(customParser),
		WithRetryPolicy(customRetry),
	).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Verify options were applied
	if coordinator.parser != customParser {
		t.Errorf("custom parser not set")
	}
	if coordinator.retryPolicy != customRetry {
		t.Errorf("custom retry policy not set")
	}
}

// testConsumerWithCloseableChannels is a consumer that can have its channels closed safely
type testConsumerWithCloseableChannels struct {
	name      string
	connected bool
	messages  chan *Message
	errors    chan error
	closed    bool
	mu        sync.Mutex
}

func (t *testConsumerWithCloseableChannels) Connect(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connected = true
	return nil
}

func (t *testConsumerWithCloseableChannels) Disconnect(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.connected = false
	if !t.closed {
		close(t.messages)
		close(t.errors)
		t.closed = true
	}
	return nil
}

func (t *testConsumerWithCloseableChannels) closeChannels() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.closed {
		close(t.messages)
		close(t.errors)
		t.closed = true
	}
}

func (t *testConsumerWithCloseableChannels) Consume(_ context.Context) (<-chan *Message, <-chan error) {
	return t.messages, t.errors
}

func (t *testConsumerWithCloseableChannels) Acknowledge(_ context.Context, _ *Message) error {
	return nil
}

func (t *testConsumerWithCloseableChannels) Reject(_ context.Context, _ *Message, _ bool) error {
	return nil
}

func (t *testConsumerWithCloseableChannels) Name() string {
	return t.name
}

func (t *testConsumerWithCloseableChannels) Validate() error {
	return nil
}

// TestDefaultCoordinator_ClosedChannels tests that the coordinator properly handles closed channels
func TestDefaultCoordinator_ClosedChannels(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Create a consumer that will close its channels
	consumer := &testConsumerWithCloseableChannels{
		name:     "test-consumer",
		messages: make(chan *Message, 1),
		errors:   make(chan error, 1),
	}

	// Add the consumer
	if err := coordinator.AddConsumer("test-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start the coordinator
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start coordinator: %v", err)
	}

	// Give the consumer a moment to start
	time.Sleep(10 * time.Millisecond)

	// Send a valid message first
	validMessage := `{
		"model": "openai",
		"prompt": "test prompt",
		"processor": "console"
	}`
	consumer.messages <- &Message{
		ID:   "test-message",
		Body: []byte(validMessage),
	}

	// Give time for message processing
	time.Sleep(20 * time.Millisecond)

	// Close the channels to simulate consumer shutdown
	consumer.closeChannels()

	// Give the coordinator time to detect closed channels and exit
	time.Sleep(50 * time.Millisecond)

	// Stop the coordinator
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer stopCancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}

	// Verify that the message was processed before channels closed
	if processor.getProcessedCount() != 1 {
		t.Errorf("expected 1 processed message, got %d", processor.getProcessedCount())
	}
}

// TestDefaultCoordinator_AutoClosingChannels tests coordinator behavior with automatically closing channels
func TestDefaultCoordinator_AutoClosingChannels(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Create a consumer that will automatically close its channels
	consumer := &testConsumerWithCloseableChannels{
		name:     "auto-closing-consumer",
		messages: make(chan *Message, 1),
		errors:   make(chan error, 1),
	}

	// Add the consumer
	if err := coordinator.AddConsumer("auto-closing-consumer", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start the coordinator
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start coordinator: %v", err)
	}

	// Close channels after a short delay to simulate consumer shutdown
	go func() {
		time.Sleep(20 * time.Millisecond)
		consumer.closeChannels()
	}()

	// Wait for the test to complete
	<-ctx.Done()

	// Stop the coordinator
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer stopCancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("unexpected error stopping: %v", err)
	}

	// Test should complete without hanging, which means the coordinator
	// properly detected closed channels and exited the consume loop
}

// TestSetLogger tests the SetLogger function
func TestSetLogger(t *testing.T) {
	// Save original logger
	originalLog := log
	defer func() {
		// Restore original logger after test
		log = originalLog
	}()

	// Create a custom logger with default configuration
	customLogger := logger.DefaultLogger()

	// Set the logger
	SetLogger(customLogger)

	// Verify it was set correctly
	if log != customLogger {
		t.Error("SetLogger did not update the package logger correctly")
	}

	// Test that the logger is actually used by the coordinator
	// This verifies that SetLogger affects the behavior of the queue package
	processor := &mockTaskProcessor{}
	coordinator := NewCoordinator(processor)

	// Add a simple consumer so Start doesn't fail
	consumer := &testConsumer{
		name:     "test-consumer",
		messages: make(chan *Message),
		errors:   make(chan error),
	}

	if err := coordinator.AddConsumer("test", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Basic operations should work with the custom logger
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := coordinator.Start(ctx); err != nil {
		t.Errorf("coordinator.Start failed with custom logger: %v", err)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer stopCancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("coordinator.Stop failed with custom logger: %v", err)
	}
}

// failingConnectConsumer is a consumer that fails to connect
type failingConnectConsumer struct {
	*testConsumer
}

func (f *failingConnectConsumer) Connect(_ context.Context) error {
	return fmt.Errorf("connection failed")
}

// TestDefaultCoordinator_ConsumerConnectError tests consumer connection failures
func TestDefaultCoordinator_ConsumerConnectError(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add consumer that fails to connect
	consumer := &failingConnectConsumer{testConsumer: newTestConsumer("fail-connect")}
	if err := coordinator.AddConsumer("fail", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start should succeed but consumer won't process messages
	ctx := context.Background()
	err := coordinator.Start(ctx)
	if err != nil {
		t.Errorf("start should succeed even if consumer fails to connect: %v", err)
	}

	// Give it time to attempt connection
	time.Sleep(100 * time.Millisecond)

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("failed to stop coordinator: %v", err)
	}

	// Verify no messages were processed
	if processor.getProcessedCount() != 0 {
		t.Errorf("expected 0 processed messages, got %d", processor.getProcessedCount())
	}
}

// failingDisconnectConsumer is a consumer that fails to disconnect
type failingDisconnectConsumer struct {
	*testConsumer
}

func (f *failingDisconnectConsumer) Disconnect(_ context.Context) error {
	return fmt.Errorf("disconnect failed")
}

// TestDefaultCoordinator_ConsumerDisconnectError tests consumer disconnect failures
func TestDefaultCoordinator_ConsumerDisconnectError(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	consumer := &failingDisconnectConsumer{testConsumer: newTestConsumer("fail-disconnect")}
	if err := coordinator.AddConsumer("fail", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Stop coordinator - should log disconnect error but not fail
	stopCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := coordinator.Stop(stopCtx)
	if err != nil {
		t.Errorf("stop should succeed even if disconnect fails: %v", err)
	}
}

// TestDefaultCoordinator_StopTimeout tests timeout during stop
func TestDefaultCoordinator_StopTimeout(t *testing.T) {
	// Create a processor with significant delay
	processor := &mockTaskProcessor{delay: 5 * time.Second}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	// Add consumer
	consumer := newTestConsumer("slow-consumer")
	if err := coordinator.AddConsumer("slow", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	// Start coordinator
	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Send a message to keep it busy
	consumer.messages <- &Message{
		ID:   "slow-msg",
		Body: []byte(`{"model": "openai", "prompt": "test", "processor": "console"}`),
	}

	// Give it time to start processing
	time.Sleep(50 * time.Millisecond)

	// Try to stop with very short timeout
	stopCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	err := coordinator.Stop(stopCtx)
	// Should return timeout error
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got: %v", err)
	}
}

// TestDefaultCoordinator_ConsumerErrorChannel tests error channel handling
func TestDefaultCoordinator_ConsumerErrorChannel(t *testing.T) {
	processor := &mockTaskProcessor{}
	coordinator, ok := NewCoordinator(processor).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	consumer := newTestConsumer("error-consumer")
	if err := coordinator.AddConsumer("error", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Send errors through error channel
	consumer.errors <- fmt.Errorf("consumer error 1")
	consumer.errors <- fmt.Errorf("consumer error 2")

	// Wait for errors to be processed
	time.Sleep(100 * time.Millisecond)

	// Coordinator should continue running despite errors
	// Send a valid message to verify it's still working
	validMsg := &Message{
		ID:   "test-after-error",
		Body: []byte(`{"model": "openai", "prompt": "test", "processor": "console"}`),
	}
	consumer.messages <- validMsg

	time.Sleep(100 * time.Millisecond)

	// Should have processed the message
	if processor.getProcessedCount() != 1 {
		t.Errorf("expected 1 processed message, got %d", processor.getProcessedCount())
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("failed to stop coordinator: %v", err)
	}
}

// mockMessageParser with custom parse and validate functions
type mockMessageParser struct {
	parseFunc    func(*Message) (*TaskMessage, error)
	validateFunc func(*TaskMessage) error
}

func (m *mockMessageParser) Parse(msg *Message) (*TaskMessage, error) {
	if m.parseFunc != nil {
		return m.parseFunc(msg)
	}
	return nil, fmt.Errorf("parse not implemented")
}

func (m *mockMessageParser) Validate(task *TaskMessage) error {
	if m.validateFunc != nil {
		return m.validateFunc(task)
	}
	return nil
}

// TestDefaultCoordinator_ValidationError tests message validation failures
func TestDefaultCoordinator_ValidationError(t *testing.T) {
	processor := &mockTaskProcessor{}

	// Create custom parser that returns invalid task
	mockParser := &mockMessageParser{
		parseFunc: func(_ *Message) (*TaskMessage, error) {
			return &TaskMessage{Model: ""}, nil // Invalid - empty model
		},
		validateFunc: func(task *TaskMessage) error {
			if task.Model == "" {
				return fmt.Errorf("model is required")
			}
			return nil
		},
	}

	coordinator, ok := NewCoordinator(processor, WithParser(mockParser)).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Send message that will fail validation
	consumer.messages <- &Message{
		ID:   "invalid-task",
		Body: []byte(`{"processor": "console"}`),
	}

	time.Sleep(100 * time.Millisecond)

	// Should reject without retry
	if consumer.getRejectCount() != 1 {
		t.Errorf("expected 1 rejection, got %d", consumer.getRejectCount())
	}

	// Verify not requeued
	if consumer.rejectRequeue {
		t.Error("message should not be requeued on validation error")
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("failed to stop coordinator: %v", err)
	}
}

// TestDefaultCoordinator_MaxRetriesExceeded tests when max retries are exceeded
func TestDefaultCoordinator_MaxRetriesExceeded(t *testing.T) {
	// Always fail processor
	processor := &mockTaskProcessor{alwaysFail: true}

	// Policy with 2 max retries
	retryPolicy := NewExponentialBackoffRetryPolicy(2, 10*time.Millisecond, 100*time.Millisecond)
	coordinator, ok := NewCoordinator(processor, WithRetryPolicy(retryPolicy)).(*DefaultCoordinator)
	if !ok {
		t.Fatal("failed to create coordinator")
	}

	consumer := newTestConsumer("test")
	if err := coordinator.AddConsumer("test", consumer); err != nil {
		t.Fatalf("failed to add consumer: %v", err)
	}

	ctx := context.Background()
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Send message that will exceed retries
	msg := &Message{
		ID:         "max-retry-test",
		Body:       []byte(`{"model": "openai", "prompt": "test", "processor": "console"}`),
		RetryCount: 2, // Already at max
	}

	consumer.messages <- msg
	time.Sleep(100 * time.Millisecond)

	// Should reject without requeue after hitting max retries
	if consumer.getRejectCount() != 1 {
		t.Errorf("expected 1 final rejection, got %d", consumer.getRejectCount())
	}

	// Verify not requeued
	if consumer.rejectRequeue {
		t.Error("message should not be requeued after max retries")
	}

	// Stop coordinator
	stopCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := coordinator.Stop(stopCtx); err != nil {
		t.Errorf("failed to stop coordinator: %v", err)
	}
}
