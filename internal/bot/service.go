// Package bot provides the main service for running CronAI in bot mode.
package bot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/bot/router"
	"github.com/rshade/cronai/internal/bot/webhook"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
)

// MockModelClient is a basic model client for testing/development
type MockModelClient struct {
	modelName string
}

// Execute implements the ModelClient interface
func (m *MockModelClient) Execute(promptContent string) (*models.ModelResponse, error) {
	return &models.ModelResponse{
		Content: fmt.Sprintf("Mock response from %s model for prompt: %s", m.modelName, promptContent[:minInt(50, len(promptContent))]),
		Model:   m.modelName,
	}, nil
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Service represents the bot mode service
type Service struct {
	server    *webhook.Server
	router    *router.Router
	model     models.ModelClient
	processor processor.Processor
	logger    *logger.Logger
}

// Config holds the bot service configuration
type Config struct {
	Port      string
	Secret    string
	Model     string
	Processor string
}

// NewService creates a new bot service
func NewService(cfg Config) (*Service, error) {
	log := logger.GetLogger()

	// Create model client
	var modelClient models.ModelClient

	// For bot mode, we'll use a simple default model if none specified
	modelName := cfg.Model
	if modelName == "" {
		modelName = "openai" // Default to OpenAI
	}

	// Note: In production, this would need proper model initialization
	// For now, we'll create a mock model for testing
	modelClient = &MockModelClient{modelName: modelName}

	// Create processor (can be nil for testing)
	var proc processor.Processor
	if cfg.Processor != "" {
		// Note: In production, this would create the actual processor
		// For now, we'll use nil processor which handlers can handle
		proc = nil
	}

	// Create router
	r := router.New()

	// Add default filters
	for _, filter := range router.DefaultFilters() {
		r.AddFilter(filter)
	}

	// Register handlers for common events
	r.RegisterHandler("issues", router.NewIssuesHandler(modelClient, proc))
	r.RegisterHandler("pull_request", router.NewPullRequestHandler(modelClient, proc))
	r.RegisterHandler("push", router.NewPushHandler(modelClient, proc))
	r.RegisterHandler("release", router.NewReleaseHandler(modelClient, proc))

	// Create rate limiter (100 requests per minute)
	rateLimiter := NewRateLimiter(100, time.Minute)

	// Create webhook server
	server := webhook.New(webhook.Config{
		Port:        cfg.Port,
		Secret:      cfg.Secret,
		Router:      r,
		RateLimiter: rateLimiter,
	})

	return &Service{
		server:    server,
		router:    r,
		model:     modelClient,
		processor: proc,
		logger:    log,
	}, nil
}

// Start starts the bot service
func (s *Service) Start() error {
	s.logger.Info("Starting bot mode service")

	// Log registered event types
	types := s.router.GetRegisteredTypes()
	s.logger.Info("Registered event handlers", logger.Fields{"handlers": strings.Join(types, ", ")})

	// Start webhook server
	return s.server.Start()
}

// Stop stops the bot service
func (s *Service) Stop() error {
	s.logger.Info("Stopping bot mode service")
	return s.server.Stop()
}

// StartService starts the bot service with the given configuration
func StartService(_ string) error {
	// Get and validate configuration from environment variables
	port := os.Getenv("CRONAI_BOT_PORT")
	if port == "" {
		port = "8080"
	}

	if err := ValidatePort(port); err != nil {
		return fmt.Errorf("invalid port configuration: %w", err)
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if err := ValidateWebhookSecret(secret); err != nil {
		return fmt.Errorf("invalid webhook secret: %w", err)
	}

	// Validate model if specified
	model := os.Getenv("CRONAI_DEFAULT_MODEL")
	if model == "" {
		model = "openai"
	}
	if err := ValidateModel(model); err != nil {
		return fmt.Errorf("invalid default model: %w", err)
	}

	// Validate processor if specified
	processor := os.Getenv("CRONAI_BOT_PROCESSOR")
	if err := ValidateProcessor(processor); err != nil {
		return fmt.Errorf("invalid processor configuration: %w", err)
	}

	// Create and start service
	service, err := NewService(Config{
		Port:      port,
		Secret:    secret,
		Model:     model,
		Processor: processor,
	})
	if err != nil {
		return fmt.Errorf("failed to create bot service: %w", err)
	}

	return service.Start()
}
