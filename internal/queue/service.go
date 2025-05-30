// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file implements the queue service that manages queue consumers and coordinates
// message processing for the queue operation mode.
package queue

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rshade/cronai/internal/logger"
)

// Service manages the queue operation mode
type Service struct {
	coordinator   Coordinator
	taskProcessor TaskProcessor
	config        *ServiceConfig
	logger        *logger.Logger
}

// ServiceConfig holds configuration for the queue service
type ServiceConfig struct {
	Consumers []ConsumerConfig `json:"consumers"`
}

// NewService creates a new queue service
func NewService() *Service {
	taskProcessor := NewTaskProcessor()
	coordinator := NewCoordinator(taskProcessor)

	return &Service{
		coordinator:   coordinator,
		taskProcessor: taskProcessor,
		logger:        logger.DefaultLogger(),
	}
}

// StartService starts the queue service with the given configuration
func StartService(configPath string) error {
	service := NewService()
	return service.Start(configPath)
}

// Start starts the service with configuration from file
func (s *Service) Start(configPath string) error {
	s.logger.Info("Starting CronAI Queue Service", logger.Fields{"config": configPath})

	// Load configuration
	config, err := s.loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	s.config = config

	// Create and register consumers
	for i, consumerConfig := range config.Consumers {
		consumerName := fmt.Sprintf("consumer-%d", i+1)
		if consumerConfig.Queue != "" {
			consumerName = fmt.Sprintf("%s-%s", consumerConfig.Type, consumerConfig.Queue)
		}

		s.logger.Info("Creating consumer", logger.Fields{
			"name":  consumerName,
			"type":  consumerConfig.Type,
			"queue": consumerConfig.Queue,
		})

		// Create consumer using registry
		consumer, err := CreateConsumer(&consumerConfig)
		if err != nil {
			return fmt.Errorf("failed to create consumer %s: %w", consumerName, err)
		}

		// Add consumer to coordinator
		if err := s.coordinator.AddConsumer(consumerName, consumer); err != nil {
			return fmt.Errorf("failed to add consumer %s: %w", consumerName, err)
		}
	}

	// Start the coordinator
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.coordinator.Start(ctx); err != nil {
		return fmt.Errorf("failed to start coordinator: %w", err)
	}

	s.logger.Info("Queue service started successfully", logger.Fields{
		"consumers": len(config.Consumers),
	})

	// Wait for interrupt signal
	s.waitForShutdown(cancel)

	// Graceful shutdown
	s.logger.Info("Shutting down queue service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := s.coordinator.Stop(shutdownCtx); err != nil {
		s.logger.Error("Error during shutdown", logger.Fields{"error": err.Error()})
		return err
	}

	s.logger.Info("Queue service stopped successfully")
	return nil
}

// loadConfig loads queue configuration from file or environment
func (s *Service) loadConfig(_ string) (*ServiceConfig, error) {
	// For now, we'll use environment variables for configuration
	// In the future, this could be extended to support config files

	config := &ServiceConfig{
		Consumers: []ConsumerConfig{},
	}

	// Check for environment-based configuration
	if queueType := os.Getenv("QUEUE_TYPE"); queueType != "" {
		consumerConfig := ConsumerConfig{
			Type:       queueType,
			Connection: os.Getenv("QUEUE_CONNECTION"),
			Queue:      os.Getenv("QUEUE_NAME"),
			Options:    make(map[string]interface{}),
			RetryLimit: 3,
			RetryDelay: 5 * time.Second,
		}

		// Parse retry configuration
		if retryLimit := os.Getenv("QUEUE_RETRY_LIMIT"); retryLimit != "" {
			if limit, err := time.ParseDuration(retryLimit); err == nil {
				consumerConfig.RetryLimit = int(limit.Seconds())
			}
		}

		if retryDelay := os.Getenv("QUEUE_RETRY_DELAY"); retryDelay != "" {
			if delay, err := time.ParseDuration(retryDelay); err == nil {
				consumerConfig.RetryDelay = delay
			}
		}

		// Validate required fields
		if consumerConfig.Type == "" {
			return nil, fmt.Errorf("QUEUE_TYPE environment variable is required")
		}

		// Set defaults based on queue type
		switch consumerConfig.Type {
		case "memory":
			if consumerConfig.Queue == "" {
				consumerConfig.Queue = "cronai-tasks"
			}
			consumerConfig.Connection = "memory://localhost"

		case "rabbitmq":
			if consumerConfig.Connection == "" {
				consumerConfig.Connection = "amqp://guest:guest@localhost:5672/"
			}
			if consumerConfig.Queue == "" {
				consumerConfig.Queue = "cronai-tasks"
			}

		default:
			return nil, fmt.Errorf("unsupported queue type: %s", consumerConfig.Type)
		}

		config.Consumers = append(config.Consumers, consumerConfig)
	}

	// If no consumers configured, create a default memory consumer for demonstration
	if len(config.Consumers) == 0 {
		s.logger.Info("No queue configuration found, creating default memory consumer for demonstration")
		config.Consumers = append(config.Consumers, ConsumerConfig{
			Type:       "memory",
			Connection: "memory://localhost",
			Queue:      "cronai-tasks",
			Options:    make(map[string]interface{}),
			RetryLimit: 3,
			RetryDelay: 5 * time.Second,
		})
	}

	return config, nil
}

// waitForShutdown waits for interrupt signals
func (s *Service) waitForShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	s.logger.Info("Received shutdown signal", logger.Fields{"signal": sig.String()})
	cancel()
}

// GetSupportedQueueTypes returns a list of supported queue types
func GetSupportedQueueTypes() []string {
	// This will be populated by the registered consumers
	return List()
}

// ValidateQueueConfig validates queue configuration
func ValidateQueueConfig(queueType, _, queueName string) error {
	if strings.TrimSpace(queueType) == "" {
		return fmt.Errorf("queue type cannot be empty")
	}

	if strings.TrimSpace(queueName) == "" {
		return fmt.Errorf("queue name cannot be empty")
	}

	// Check if queue type is supported
	supportedTypes := GetSupportedQueueTypes()
	for _, supportedType := range supportedTypes {
		if supportedType == queueType {
			return nil // Found supported type
		}
	}

	return fmt.Errorf("unsupported queue type '%s'. Supported types: %s",
		queueType, strings.Join(supportedTypes, ", "))
}
