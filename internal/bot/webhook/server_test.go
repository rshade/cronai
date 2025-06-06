// Package webhook provides HTTP server functionality for receiving GitHub webhook events in bot mode.
package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockRouter implements the Router interface for testing
type mockRouter struct {
	lastEventType string
	lastPayload   []byte
	err           error
}

func (m *mockRouter) Route(eventType string, payload []byte) error {
	m.lastEventType = eventType
	m.lastPayload = payload
	return m.err
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		wantPort string
	}{
		{
			name:     "default port",
			cfg:      Config{},
			wantPort: "8080",
		},
		{
			name:     "custom port",
			cfg:      Config{Port: "9090"},
			wantPort: "9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.cfg)
			if s.port != tt.wantPort {
				t.Errorf("New() port = %v, want %v", s.port, tt.wantPort)
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	s := New(Config{})

	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   map[string]string
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"status": "healthy", "mode": "bot"},
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			s.handleHealth(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("handleHealth() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantBody != nil {
				var body map[string]string
				if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				for k, v := range tt.wantBody {
					if body[k] != v {
						t.Errorf("handleHealth() body[%s] = %v, want %v", k, body[k], v)
					}
				}
			}
		})
	}
}

func TestHandleWebhook(t *testing.T) {
	router := &mockRouter{}
	s := New(Config{Router: router})

	tests := []struct {
		name         string
		method       string
		headers      map[string]string
		body         string
		wantStatus   int
		wantRouted   bool
		wantResponse map[string]string
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Missing event type",
			method:     http.MethodPost,
			headers:    map[string]string{},
			body:       `{"test": "data"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "Valid webhook",
			method: http.MethodPost,
			headers: map[string]string{
				"X-GitHub-Event": "push",
			},
			body:         `{"ref": "refs/heads/main"}`,
			wantStatus:   http.StatusOK,
			wantRouted:   true,
			wantResponse: map[string]string{"status": "accepted", "event": "push"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.lastEventType = ""
			router.lastPayload = nil

			req := httptest.NewRequest(tt.method, "/webhook", bytes.NewBufferString(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			s.handleWebhook(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("handleWebhook() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantRouted && router.lastEventType == "" {
				t.Error("handleWebhook() expected event to be routed")
			}

			if !tt.wantRouted && router.lastEventType != "" {
				t.Error("handleWebhook() unexpected event routing")
			}

			if tt.wantResponse != nil {
				var body map[string]string
				if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				for k, v := range tt.wantResponse {
					if body[k] != v {
						t.Errorf("handleWebhook() body[%s] = %v, want %v", k, body[k], v)
					}
				}
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	secret := "test-secret"
	s := New(Config{Secret: secret})

	payload := []byte(`{"test": "data"}`)

	// Generate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		payload   []byte
		signature string
		want      bool
	}{
		{
			name:      "valid signature",
			payload:   payload,
			signature: validSignature,
			want:      true,
		},
		{
			name:      "invalid signature",
			payload:   payload,
			signature: "sha256=invalid",
			want:      false,
		},
		{
			name:      "missing prefix",
			payload:   payload,
			signature: "invalid",
			want:      false,
		},
		{
			name:      "empty signature",
			payload:   payload,
			signature: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.verifySignature(tt.payload, tt.signature); got != tt.want {
				t.Errorf("verifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleWebhookWithSignature(t *testing.T) {
	secret := "test-secret"
	router := &mockRouter{}
	s := New(Config{Secret: secret, Router: router})

	payload := []byte(`{"test": "data"}`)

	// Generate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name       string
		signature  string
		wantStatus int
		wantRouted bool
	}{
		{
			name:       "valid signature",
			signature:  validSignature,
			wantStatus: http.StatusOK,
			wantRouted: true,
		},
		{
			name:       "invalid signature",
			signature:  "sha256=invalid",
			wantStatus: http.StatusUnauthorized,
			wantRouted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router.lastEventType = ""

			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(payload))
			req.Header.Set("X-GitHub-Event", "push")
			req.Header.Set("X-Hub-Signature-256", tt.signature)
			w := httptest.NewRecorder()

			s.handleWebhook(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("handleWebhook() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantRouted && router.lastEventType == "" {
				t.Error("handleWebhook() expected event to be routed")
			}

			if !tt.wantRouted && router.lastEventType != "" {
				t.Error("handleWebhook() unexpected event routing")
			}
		})
	}
}

func TestServerStartStop(t *testing.T) {
	// Create a test server that doesn't actually start the HTTP server
	// This tests the lifecycle without the complexities of actual network operations
	s := New(Config{Port: "0"}) // Use port 0 for automatic port assignment

	// Test that we can create the server
	if s == nil {
		t.Fatal("New() returned nil")
	}

	// Test Stop on unstarted server (should not error)
	if err := s.Stop(); err != nil {
		t.Errorf("Stop() on unstarted server error = %v", err)
	}

	// For actual start/stop testing with network operations,
	// we would need more complex test infrastructure.
	// This basic test ensures the server can be created and stopped safely.
}

func TestServerLifecycle(t *testing.T) {
	// Find a free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find free port: %v", err)
	}
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("Failed to get TCP address")
	}
	port := fmt.Sprintf("%d", addr.Port)
	_ = listener.Close() //nolint:errcheck // Just finding a free port

	s := New(Config{Port: port})

	// Start server in background
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	// Wait for server to be ready by polling the health endpoint
	var healthCheckErr error
	for i := 0; i < 50; i++ { // Try for up to 5 seconds
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
		if err == nil {
			_ = resp.Body.Close() //nolint:errcheck // Test cleanup
			if resp.StatusCode == http.StatusOK {
				healthCheckErr = nil
				break
			}
		}
		healthCheckErr = err
	}

	if healthCheckErr != nil {
		t.Fatalf("Server failed to start: %v", healthCheckErr)
	}

	// Stop the server
	if err := s.Stop(); err != nil {
		t.Errorf("Stop() error = %v", err)
	}

	// Wait for server to shut down
	select {
	case err := <-serverErr:
		if err != nil {
			t.Errorf("Server error on shutdown: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Error("Server did not shut down within timeout")
	}
}

func TestRouterError(t *testing.T) {
	router := &mockRouter{err: fmt.Errorf("routing error")}
	s := New(Config{Router: router})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString(`{"test": "data"}`))
	req.Header.Set("X-GitHub-Event", "push")
	w := httptest.NewRecorder()

	s.handleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("handleWebhook() status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
