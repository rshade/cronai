// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file contains configuration parsing logic for queue tasks, handling the conversion
// from configuration lines to structured task definitions.
package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/kballard/go-shellquote"
)

// Task represents a task configuration from the queue section
type Task struct {
	Name       string                 // Queue consumer name
	Type       string                 // Queue type (rabbitmq, sqs, etc.)
	Connection string                 // Connection string
	Queue      string                 // Queue name or topic
	Options    map[string]interface{} // Provider-specific options
	RetryLimit int                    // Maximum retry attempts
	RetryDelay time.Duration          // Delay between retries
}

// ParseQueueConfig parses a queue configuration line
// Format: queue <name> <type> <connection> <queue> [options]
// Example: queue main-queue rabbitmq amqp://localhost:5672 tasks retry_limit=3,retry_delay=5s
func ParseQueueConfig(line string) (*Task, error) {
	// Skip empty lines and comments
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	// Check if this is a queue configuration line
	if !strings.HasPrefix(line, "queue ") {
		return nil, nil
	}

	// keep quoted segments intact so spaces inside quotes survive
	parts, err := shellquote.Split(line)
	if err != nil {
		return nil, fmt.Errorf("invalid queue config: %w", err)
	}
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid queue format: need at least 4 fields (name, type, connection, queue)")
	}

	// Skip the "queue" prefix
	task := &Task{
		Name:       parts[1],
		Type:       parts[2],
		Connection: parts[3],
		Queue:      parts[4],
		Options:    nil,
		RetryLimit: 3,               // Default retry limit
		RetryDelay: 5 * time.Second, // Default retry delay
	}

	// Parse optional configuration
	if len(parts) > 5 {
		optionsStr := strings.Join(parts[5:], " ")
		options, err := parseQueueOptions(optionsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid options format: %w", err)
		}

		// Extract standard options
		if retryLimit, ok := options["retry_limit"]; ok {
			if limit, ok := retryLimit.(int); ok {
				task.RetryLimit = limit
			}
		}

		if retryDelay, ok := options["retry_delay"]; ok {
			if delay, ok := retryDelay.(time.Duration); ok {
				task.RetryDelay = delay
			}
		}

		// Store all options for provider-specific use
		task.Options = options
	}

	return task, nil
}

// parseQueueOptions parses queue options string
func parseQueueOptions(optionsStr string) (map[string]interface{}, error) {
	options := make(map[string]interface{})

	// Split by comma
	pairs := strings.Split(optionsStr, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// Split key=value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid option format: %s", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Parse special values
		switch key {
		case "retry_limit":
			// Parse as integer
			var limit int
			if _, err := fmt.Sscanf(value, "%d", &limit); err != nil {
				return nil, fmt.Errorf("invalid retry_limit: %s", value)
			}
			options[key] = limit

		case "retry_delay":
			// Parse as duration
			duration, err := time.ParseDuration(value)
			if err != nil {
				return nil, fmt.Errorf("invalid retry_delay: %s", value)
			}
			options[key] = duration

		default:
			// Store as string for provider-specific parsing
			options[key] = value
		}
	}

	return options, nil
}

// CreateConsumerConfig creates a ConsumerConfig from a Task
func CreateConsumerConfig(task *Task) *ConsumerConfig {
	return &ConsumerConfig{
		Type:       task.Type,
		Connection: task.Connection,
		Queue:      task.Queue,
		Options:    task.Options,
		RetryLimit: task.RetryLimit,
		RetryDelay: task.RetryDelay,
	}
}

// IsQueueConfig checks if a configuration line is a queue configuration
func IsQueueConfig(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "queue ")
}
