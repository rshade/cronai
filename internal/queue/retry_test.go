// Package queue provides the core infrastructure for message queue integration in CronAI.
package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestExponentialBackoffRetryPolicy_ShouldRetry(t *testing.T) {
	policy := NewExponentialBackoffRetryPolicy(3, 1*time.Second, 30*time.Second)

	tests := []struct {
		name     string
		message  *Message
		err      error
		expected bool
	}{
		{
			name:     "nil message",
			message:  nil,
			err:      fmt.Errorf("some error"),
			expected: false,
		},
		{
			name: "first retry",
			message: &Message{
				ID:         "test-1",
				RetryCount: 0,
			},
			err:      fmt.Errorf("temporary error"),
			expected: true,
		},
		{
			name: "second retry",
			message: &Message{
				ID:         "test-2",
				RetryCount: 1,
			},
			err:      fmt.Errorf("temporary error"),
			expected: true,
		},
		{
			name: "third retry",
			message: &Message{
				ID:         "test-3",
				RetryCount: 2,
			},
			err:      fmt.Errorf("temporary error"),
			expected: true,
		},
		{
			name: "exceeded max retries",
			message: &Message{
				ID:         "test-4",
				RetryCount: 3,
			},
			err:      fmt.Errorf("temporary error"),
			expected: false,
		},
		{
			name: "no error",
			message: &Message{
				ID:         "test-5",
				RetryCount: 0,
			},
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.ShouldRetry(tt.message, tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExponentialBackoffRetryPolicy_NextRetryDelay(t *testing.T) {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second
	policy := NewExponentialBackoffRetryPolicy(3, baseDelay, maxDelay)

	tests := []struct {
		name     string
		message  *Message
		expected time.Duration
	}{
		{
			name:     "nil message",
			message:  nil,
			expected: baseDelay,
		},
		{
			name: "first retry",
			message: &Message{
				ID:         "test-1",
				RetryCount: 0,
			},
			expected: 1 * time.Second,
		},
		{
			name: "second retry",
			message: &Message{
				ID:         "test-2",
				RetryCount: 1,
			},
			expected: 2 * time.Second,
		},
		{
			name: "third retry",
			message: &Message{
				ID:         "test-3",
				RetryCount: 2,
			},
			expected: 4 * time.Second,
		},
		{
			name: "fourth retry",
			message: &Message{
				ID:         "test-4",
				RetryCount: 3,
			},
			expected: 8 * time.Second,
		},
		{
			name: "exceeds max delay",
			message: &Message{
				ID:         "test-5",
				RetryCount: 10, // 2^10 = 1024 seconds > 30 seconds
			},
			expected: maxDelay,
		},
		{
			name: "high retry count hits max delay",
			message: &Message{
				ID:         "test-6",
				RetryCount: 20, // Very high, ensures we hit the cap
			},
			expected: maxDelay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.NextRetryDelay(tt.message)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExponentialBackoffRetryPolicy_MaxRetries(t *testing.T) {
	maxRetries := 5
	policy := NewExponentialBackoffRetryPolicy(maxRetries, 1*time.Second, 30*time.Second)

	if policy.MaxRetries() != maxRetries {
		t.Errorf("expected max retries %d, got %d", maxRetries, policy.MaxRetries())
	}
}

func TestLinearRetryPolicy_ShouldRetry(t *testing.T) {
	policy := NewLinearRetryPolicy(2, 5*time.Second)

	tests := []struct {
		name     string
		message  *Message
		err      error
		expected bool
	}{
		{
			name:     "nil message",
			message:  nil,
			err:      fmt.Errorf("error"),
			expected: false,
		},
		{
			name: "first retry",
			message: &Message{
				ID:         "test-1",
				RetryCount: 0,
			},
			err:      fmt.Errorf("error"),
			expected: true,
		},
		{
			name: "second retry",
			message: &Message{
				ID:         "test-2",
				RetryCount: 1,
			},
			err:      fmt.Errorf("error"),
			expected: true,
		},
		{
			name: "exceeded max retries",
			message: &Message{
				ID:         "test-3",
				RetryCount: 2,
			},
			err:      fmt.Errorf("error"),
			expected: false,
		},
		{
			name: "no error",
			message: &Message{
				ID:         "test-4",
				RetryCount: 0,
			},
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := policy.ShouldRetry(tt.message, tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLinearRetryPolicy_NextRetryDelay(t *testing.T) {
	delay := 5 * time.Second
	policy := NewLinearRetryPolicy(2, delay)

	// Test with various messages - should always return the same delay
	messages := []*Message{
		nil,
		{ID: "test-1", RetryCount: 0},
		{ID: "test-2", RetryCount: 1},
		{ID: "test-3", RetryCount: 5},
	}

	for i, msg := range messages {
		t.Run(fmt.Sprintf("message_%d", i), func(t *testing.T) {
			result := policy.NextRetryDelay(msg)
			if result != delay {
				t.Errorf("expected %v, got %v", delay, result)
			}
		})
	}
}

func TestLinearRetryPolicy_MaxRetries(t *testing.T) {
	maxRetries := 3
	policy := NewLinearRetryPolicy(maxRetries, 5*time.Second)

	if policy.MaxRetries() != maxRetries {
		t.Errorf("expected max retries %d, got %d", maxRetries, policy.MaxRetries())
	}
}

func TestNoRetryPolicy(t *testing.T) {
	policy := NewNoRetryPolicy()

	// Test ShouldRetry - should always return false
	tests := []struct {
		name    string
		message *Message
		err     error
	}{
		{
			name:    "nil message",
			message: nil,
			err:     fmt.Errorf("error"),
		},
		{
			name:    "with message and error",
			message: &Message{ID: "test-1"},
			err:     fmt.Errorf("error"),
		},
		{
			name:    "with message no error",
			message: &Message{ID: "test-2"},
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if policy.ShouldRetry(tt.message, tt.err) {
				t.Errorf("NoRetryPolicy.ShouldRetry should always return false")
			}
		})
	}

	// Test NextRetryDelay - should always return 0
	if delay := policy.NextRetryDelay(&Message{ID: "test"}); delay != 0 {
		t.Errorf("expected 0 delay, got %v", delay)
	}

	// Test MaxRetries - should always return 0
	if maxRetries := policy.MaxRetries(); maxRetries != 0 {
		t.Errorf("expected 0 max retries, got %d", maxRetries)
	}
}

// TestExponentialBackoffRetryPolicy_LargeDurations tests that the exponential backoff
// calculation doesn't overflow with large durations
func TestExponentialBackoffRetryPolicy_LargeDurations(t *testing.T) {
	// Test with very large base delay and max delay
	largeBaseDelay := 24 * time.Hour
	largeMaxDelay := 30 * 24 * time.Hour // 30 days

	policy := NewExponentialBackoffRetryPolicy(10, largeBaseDelay, largeMaxDelay)

	tests := []struct {
		name       string
		retryCount int
		maxDelay   time.Duration
	}{
		{
			name:       "first retry with large base",
			retryCount: 1,
			maxDelay:   largeMaxDelay,
		},
		{
			name:       "multiple retries should hit max",
			retryCount: 5,
			maxDelay:   largeMaxDelay,
		},
		{
			name:       "many retries should still hit max",
			retryCount: 20,
			maxDelay:   largeMaxDelay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &Message{
				ID:         "test",
				RetryCount: tt.retryCount,
			}

			delay := policy.NextRetryDelay(message)

			// Verify no overflow occurred - delay should be valid
			if delay < 0 {
				t.Errorf("delay overflowed: %v", delay)
			}

			// Verify delay doesn't exceed max
			if delay > tt.maxDelay {
				t.Errorf("delay %v exceeds max delay %v", delay, tt.maxDelay)
			}

			// For high retry counts, should hit max delay
			if tt.retryCount > 3 && delay != largeMaxDelay {
				t.Errorf("expected max delay for retry count %d, got %v", tt.retryCount, delay)
			}
		})
	}
}

// TestExponentialBackoffRetryPolicy_PreciseCalculation tests the precise calculation
// of exponential backoff delays
func TestExponentialBackoffRetryPolicy_PreciseCalculation(t *testing.T) {
	baseDelay := 100 * time.Millisecond
	maxDelay := 10 * time.Second
	backoffFactor := 2.0

	policy := &ExponentialBackoffRetryPolicy{
		maxRetries:    5,
		baseDelay:     baseDelay,
		maxDelay:      maxDelay,
		backoffFactor: backoffFactor,
	}

	tests := []struct {
		retryCount    int
		expectedDelay time.Duration
	}{
		{0, 100 * time.Millisecond},  // 100ms
		{1, 200 * time.Millisecond},  // 100ms * 2^1
		{2, 400 * time.Millisecond},  // 100ms * 2^2
		{3, 800 * time.Millisecond},  // 100ms * 2^3
		{4, 1600 * time.Millisecond}, // 100ms * 2^4
		{5, 3200 * time.Millisecond}, // 100ms * 2^5
		{6, 6400 * time.Millisecond}, // 100ms * 2^6
		{7, 10 * time.Second},        // Would be 12.8s, but capped at max
		{8, 10 * time.Second},        // Capped at max
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("retry_%d", tt.retryCount), func(t *testing.T) {
			message := &Message{
				ID:         "test",
				RetryCount: tt.retryCount,
			}

			delay := policy.NextRetryDelay(message)

			if delay != tt.expectedDelay {
				t.Errorf("retry %d: expected %v, got %v", tt.retryCount, tt.expectedDelay, delay)
			}
		})
	}
}
