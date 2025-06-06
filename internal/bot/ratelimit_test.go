// Package bot provides the main service for running CronAI in bot mode.
package bot

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, time.Second)

	if rl.maxTokens != 10 {
		t.Errorf("NewRateLimiter() maxTokens = %v, want 10", rl.maxTokens)
	}

	if rl.tokens != 10 {
		t.Errorf("NewRateLimiter() initial tokens = %v, want 10", rl.tokens)
	}

	if rl.refillRate != time.Second {
		t.Errorf("NewRateLimiter() refillRate = %v, want %v", rl.refillRate, time.Second)
	}
}

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)

	// First two requests should be allowed
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}

	if !rl.Allow() {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if rl.Allow() {
		t.Error("Third request should be denied")
	}

	// Wait for refill and try again
	time.Sleep(150 * time.Millisecond)

	if !rl.Allow() {
		t.Error("Request after refill should be allowed")
	}
}

func TestRateLimiterRefill(t *testing.T) {
	rl := NewRateLimiter(5, 50*time.Millisecond)

	// Exhaust all tokens
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Should be denied now
	if rl.Allow() {
		t.Error("Request should be denied when tokens exhausted")
	}

	// Wait for multiple refill periods
	time.Sleep(120 * time.Millisecond) // Should refill 2 tokens

	// Should allow 2 more requests
	if !rl.Allow() {
		t.Error("First request after refill should be allowed")
	}

	if !rl.Allow() {
		t.Error("Second request after refill should be allowed")
	}

	// Third should be denied
	if rl.Allow() {
		t.Error("Third request after partial refill should be denied")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute) // Very slow refill for testing
	middleware := RateLimitMiddleware(rl)

	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck // Test helper, error not important
		w.Write([]byte("success"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(handler)

	// First request should succeed
	req1 := httptest.NewRequest("GET", "/test", nil)
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request status = %v, want %v", w1.Code, http.StatusOK)
	}

	if w1.Body.String() != "success" {
		t.Errorf("First request body = %v, want 'success'", w1.Body.String())
	}

	// Second request should be rate limited
	req2 := httptest.NewRequest("GET", "/test", nil)
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request status = %v, want %v", w2.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimiterMaxTokens(t *testing.T) {
	rl := NewRateLimiter(3, 50*time.Millisecond)

	// Exhaust tokens
	for i := 0; i < 3; i++ {
		rl.Allow()
	}

	// Wait for a long time (should refill more than max)
	time.Sleep(200 * time.Millisecond) // 4 refill periods

	// Should only have max tokens available
	allowedCount := 0
	for i := 0; i < 5; i++ {
		if rl.Allow() {
			allowedCount++
		}
	}

	if allowedCount != 3 {
		t.Errorf("After long wait, allowed %d requests, want 3", allowedCount)
	}
}

func TestRateLimiterConcurrency(t *testing.T) {
	// Use a very long refill rate to prevent refills during the test
	rl := NewRateLimiter(100, time.Hour)

	// Run concurrent requests
	results := make(chan bool, 200)

	for i := 0; i < 200; i++ {
		go func() {
			results <- rl.Allow()
		}()
	}

	// Collect results
	allowed := 0
	denied := 0

	for i := 0; i < 200; i++ {
		if <-results {
			allowed++
		} else {
			denied++
		}
	}

	// Should have allowed exactly 100 requests (initial tokens)
	if allowed != 100 {
		t.Errorf("Concurrent test allowed %d requests, want 100", allowed)
	}

	if denied != 100 {
		t.Errorf("Concurrent test denied %d requests, want 100", denied)
	}
}
