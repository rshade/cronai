// Package webhook provides HTTP server functionality for receiving GitHub webhook events in bot mode.
package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rshade/cronai/internal/logger"
)

// Server represents the webhook server
type Server struct {
	port        string
	secret      string
	router      Router
	httpServer  *http.Server
	logger      *logger.Logger
	rateLimiter RateLimiter
}

// RateLimiter interface for rate limiting
type RateLimiter interface {
	Allow() bool
}

// Router defines the interface for event routing
type Router interface {
	Route(eventType string, payload []byte) error
}

// Config holds the server configuration
type Config struct {
	Port        string
	Secret      string
	Router      Router
	RateLimiter RateLimiter
}

// New creates a new webhook server
func New(cfg Config) *Server {
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return &Server{
		port:        cfg.Port,
		secret:      cfg.Secret,
		router:      cfg.Router,
		logger:      logger.GetLogger(),
		rateLimiter: cfg.RateLimiter,
	}
}

// Start starts the webhook server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/webhook", s.handleWebhook)

	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Channel to notify when server has shut down
	done := make(chan error, 1)

	// Start server in goroutine
	go func() {
		s.logger.Info("Starting webhook server", logger.Fields{"port": s.port})
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- err
		} else {
			done <- nil // Signal successful shutdown
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case <-stop:
		s.logger.Info("Received interrupt signal, shutting down...")
		return s.gracefulShutdown()
	case err := <-done:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
		s.logger.Info("Server stopped")
		return nil
	}
}

// gracefulShutdown performs a graceful shutdown of the server
func (s *Server) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	s.logger.Info("Server gracefully stopped")
	return nil
}

// Stop stops the webhook server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.gracefulShutdown()
	}
	return nil
}

// loggingMiddleware logs all incoming requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.logger.Debug("Incoming request", logger.Fields{"method": r.Method, "path": r.URL.Path})
		next.ServeHTTP(w, r)
		s.logger.Debug("Request completed", logger.Fields{"duration": time.Since(start)})
	})
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"mode":   "bot",
	}); err != nil {
		s.logger.Error("Failed to encode health response", logger.Fields{"error": err})
	}
}

// handleWebhook handles incoming webhook requests
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Apply rate limiting if configured
	if s.rateLimiter != nil && !s.rateLimiter.Allow() {
		s.logger.Warn("Webhook request rate limited", logger.Fields{"remote_addr": r.RemoteAddr})
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("Failed to read request body", logger.Fields{"error": err})
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			s.logger.Error("Failed to close request body", logger.Fields{"error": closeErr})
		}
	}()

	// Verify signature if secret is configured
	if s.secret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !s.verifySignature(body, signature) {
			s.logger.Warn("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Get event type
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		s.logger.Warn("Missing X-GitHub-Event header")
		http.Error(w, "Missing event type", http.StatusBadRequest)
		return
	}

	s.logger.Info("Received event", logger.Fields{"eventType": eventType})

	// Route the event
	if s.router != nil {
		if err := s.router.Route(eventType, body); err != nil {
			s.logger.Error("Failed to route event", logger.Fields{"error": err})
			http.Error(w, "Failed to process event", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"event":  eventType,
	}); err != nil {
		s.logger.Error("Failed to encode webhook response", logger.Fields{"error": err})
	}
}

// verifySignature verifies the GitHub webhook signature
func (s *Server) verifySignature(payload []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	expected := signature[7:] // Remove "sha256=" prefix
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(payload)
	computed := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(computed))
}
