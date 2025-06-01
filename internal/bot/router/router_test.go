// Package router provides event routing functionality for processing GitHub webhook events.
package router

import (
	"encoding/json"
	"fmt"
	"testing"
)

// mockHandler implements EventHandler for testing
type mockHandler struct {
	handled []Event
	err     error
}

func (m *mockHandler) Handle(event Event) error {
	m.handled = append(m.handled, event)
	return m.err
}

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.handlers == nil {
		t.Error("New() handlers map is nil")
	}
	if r.filters == nil {
		t.Error("New() filters slice is nil")
	}
}

func TestRegisterHandler(t *testing.T) {
	r := New()
	handler := &mockHandler{}

	r.RegisterHandler("push", handler)

	if len(r.handlers) != 1 {
		t.Errorf("RegisterHandler() handlers count = %d, want 1", len(r.handlers))
	}

	if r.handlers["push"] != handler {
		t.Error("RegisterHandler() handler not registered correctly")
	}
}

func TestGetRegisteredTypes(t *testing.T) {
	r := New()

	// Register multiple handlers
	r.RegisterHandler("push", &mockHandler{})
	r.RegisterHandler("pull_request", &mockHandler{})
	r.RegisterHandler("issues", &mockHandler{})

	types := r.GetRegisteredTypes()

	if len(types) != 3 {
		t.Errorf("GetRegisteredTypes() returned %d types, want 3", len(types))
	}

	// Check all types are present (order not guaranteed)
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	expectedTypes := []string{"push", "pull_request", "issues"}
	for _, expected := range expectedTypes {
		if !typeMap[expected] {
			t.Errorf("GetRegisteredTypes() missing type: %s", expected)
		}
	}
}

func TestRoute(t *testing.T) {
	tests := []struct {
		name        string
		eventType   string
		payload     string
		handlers    map[string]*mockHandler
		filters     []EventFilter
		wantHandled bool
		wantAction  string
		wantErr     bool
	}{
		{
			name:      "successful routing",
			eventType: "push",
			payload:   `{"ref": "refs/heads/main"}`,
			handlers: map[string]*mockHandler{
				"push": {},
			},
			wantHandled: true,
		},
		{
			name:      "routing with action",
			eventType: "issues",
			payload:   `{"action": "opened", "issue": {"number": 1}}`,
			handlers: map[string]*mockHandler{
				"issues": {},
			},
			wantHandled: true,
			wantAction:  "opened",
		},
		{
			name:        "no handler registered",
			eventType:   "release",
			payload:     `{"action": "created"}`,
			handlers:    map[string]*mockHandler{},
			wantHandled: false,
		},
		{
			name:      "handler returns error",
			eventType: "push",
			payload:   `{"ref": "refs/heads/main"}`,
			handlers: map[string]*mockHandler{
				"push": {err: fmt.Errorf("handler error")},
			},
			wantHandled: true,
			wantErr:     true,
		},
		{
			name:      "event filtered out",
			eventType: "push",
			payload:   `{"ref": "refs/heads/main"}`,
			handlers: map[string]*mockHandler{
				"push": {},
			},
			filters: []EventFilter{
				func(_ Event) bool { return false },
			},
			wantHandled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			// Register handlers
			for eventType, handler := range tt.handlers {
				r.RegisterHandler(eventType, handler)
			}

			// Add filters
			for _, filter := range tt.filters {
				r.AddFilter(filter)
			}

			// Route event
			err := r.Route(tt.eventType, []byte(tt.payload))

			if (err != nil) != tt.wantErr {
				t.Errorf("Route() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if handler was called
			handler := tt.handlers[tt.eventType]
			if handler != nil {
				handled := len(handler.handled) > 0
				if handled != tt.wantHandled {
					t.Errorf("Route() handled = %v, want %v", handled, tt.wantHandled)
				}

				if tt.wantHandled && len(handler.handled) > 0 {
					event := handler.handled[0]
					if event.Type != tt.eventType {
						t.Errorf("Route() event.Type = %v, want %v", event.Type, tt.eventType)
					}
					if event.Action != tt.wantAction {
						t.Errorf("Route() event.Action = %v, want %v", event.Action, tt.wantAction)
					}
				}
			}
		})
	}
}

func TestBotEventFilter(t *testing.T) {
	filter := BotEventFilter()

	tests := []struct {
		name    string
		payload string
		want    bool
	}{
		{
			name:    "human user event",
			payload: `{"sender": {"type": "User"}}`,
			want:    true,
		},
		{
			name:    "bot user event",
			payload: `{"sender": {"type": "Bot"}}`,
			want:    false,
		},
		{
			name:    "no sender field",
			payload: `{"action": "created"}`,
			want:    true,
		},
		{
			name:    "invalid json",
			payload: `{invalid}`,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{
				Type:    "test",
				Payload: json.RawMessage(tt.payload),
			}

			if got := filter(event); got != tt.want {
				t.Errorf("BotEventFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddFilter(t *testing.T) {
	r := New()

	// Add multiple filters
	filter1Called := false
	filter2Called := false

	r.AddFilter(func(_ Event) bool {
		filter1Called = true
		return true
	})

	r.AddFilter(func(_ Event) bool {
		filter2Called = true
		return false // This filter blocks the event
	})

	// Register handler
	handler := &mockHandler{}
	r.RegisterHandler("test", handler)

	// Route event
	_ = r.Route("test", []byte(`{}`)) //nolint:errcheck // Test doesn't need error handling

	// Check filters were called
	if !filter1Called {
		t.Error("First filter was not called")
	}
	if !filter2Called {
		t.Error("Second filter was not called")
	}

	// Handler should not be called because filter2 blocked it
	if len(handler.handled) > 0 {
		t.Error("Handler was called despite filter blocking event")
	}
}

func TestConcurrentAccess(_ *testing.T) {
	r := New()

	// Run concurrent operations
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			r.RegisterHandler(fmt.Sprintf("event%d", i), &mockHandler{})
			r.AddFilter(func(_ Event) bool { return true })
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			r.GetRegisteredTypes()
			_ = r.Route("event0", []byte(`{}`)) //nolint:errcheck // Test doesn't need error handling
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}

func TestDefaultFilters(t *testing.T) {
	filters := DefaultFilters()

	if len(filters) == 0 {
		t.Error("DefaultFilters() returned empty slice")
	}

	// Test that bot filter is included
	botEvent := Event{
		Type:    "test",
		Payload: json.RawMessage(`{"sender": {"type": "Bot"}}`),
	}

	// Apply filters
	allowed := true
	for _, filter := range filters {
		if !filter(botEvent) {
			allowed = false
			break
		}
	}

	if allowed {
		t.Error("DefaultFilters() should filter out bot events")
	}
}
