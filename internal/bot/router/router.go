// Package router provides event routing functionality for processing GitHub webhook events.
package router

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rshade/cronai/internal/logger"
)

// EventHandler processes a specific type of GitHub event
type EventHandler interface {
	Handle(event Event) error
}

// Event represents a parsed GitHub webhook event
type Event struct {
	Type    string          // Event type from X-GitHub-Event header
	Action  string          // Action field from payload (if applicable)
	Payload json.RawMessage // Raw JSON payload
}

// Router routes GitHub webhook events to appropriate handlers
type Router struct {
	handlers map[string]EventHandler
	filters  []EventFilter
	logger   *logger.Logger
	mu       sync.RWMutex
}

// EventFilter filters events before routing
type EventFilter func(event Event) bool

// New creates a new event router
func New() *Router {
	return &Router{
		handlers: make(map[string]EventHandler),
		filters:  []EventFilter{},
		logger:   logger.DefaultLogger(),
	}
}

// RegisterHandler registers a handler for a specific event type
func (r *Router) RegisterHandler(eventType string, handler EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[eventType] = handler
	r.logger.Debug("Registered handler for event type", logger.Fields{"eventType": eventType})
}

// AddFilter adds an event filter
func (r *Router) AddFilter(filter EventFilter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.filters = append(r.filters, filter)
}

// Route processes an incoming webhook event
func (r *Router) Route(eventType string, payload []byte) error {
	event := Event{
		Type:    eventType,
		Payload: json.RawMessage(payload),
	}

	// Extract action from payload if present
	var basePayload struct {
		Action string `json:"action,omitempty"`
	}
	if err := json.Unmarshal(payload, &basePayload); err == nil {
		event.Action = basePayload.Action
	}

	// Apply filters
	r.mu.RLock()
	filters := r.filters
	r.mu.RUnlock()

	for _, filter := range filters {
		if !filter(event) {
			r.logger.Debug("Event filtered out", logger.Fields{"type": event.Type, "action": event.Action})
			return nil
		}
	}

	// Find and execute handler
	r.mu.RLock()
	handler, exists := r.handlers[eventType]
	r.mu.RUnlock()

	if !exists {
		r.logger.Warn("No handler registered for event type", logger.Fields{"eventType": eventType})
		return nil // Not an error - just no handler registered
	}

	r.logger.Info("Routing event to handler", logger.Fields{"eventType": eventType, "action": event.Action})

	if err := handler.Handle(event); err != nil {
		return fmt.Errorf("handler error for %s event: %w", eventType, err)
	}

	return nil
}

// GetRegisteredTypes returns a list of event types with registered handlers
func (r *Router) GetRegisteredTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		types = append(types, t)
	}
	return types
}

// DefaultFilters returns commonly used event filters
func DefaultFilters() []EventFilter {
	return []EventFilter{
		// Filter out bot events
		BotEventFilter(),
	}
}

// BotEventFilter filters out events from bot users
func BotEventFilter() EventFilter {
	return func(event Event) bool {
		// Parse sender information
		var payload struct {
			Sender struct {
				Type string `json:"type"`
			} `json:"sender"`
		}

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			// If we can't parse, allow the event through
			return true
		}

		// Filter out bot events
		return payload.Sender.Type != "Bot"
	}
}
