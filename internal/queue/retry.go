// Package queue provides the core infrastructure for message queue integration in CronAI.
// This file implements various retry policies including exponential backoff, linear retry,
// and no-retry strategies for handling transient failures in message processing.
package queue

import (
	"time"
)

// ExponentialBackoffRetryPolicy implements an exponential backoff retry strategy
type ExponentialBackoffRetryPolicy struct {
	maxRetries    int
	baseDelay     time.Duration
	maxDelay      time.Duration
	backoffFactor float64
}

// NewExponentialBackoffRetryPolicy creates a new exponential backoff retry policy
func NewExponentialBackoffRetryPolicy(maxRetries int, baseDelay, maxDelay time.Duration) RetryPolicy {
	// Validate maxRetries is at least 0
	if maxRetries < 0 {
		maxRetries = 0
	}

	// Validate baseDelay is positive, default to time.Second if not
	if baseDelay <= 0 {
		baseDelay = time.Second
	}

	// Validate maxDelay is not less than baseDelay
	if maxDelay < baseDelay {
		maxDelay = baseDelay
	}

	return &ExponentialBackoffRetryPolicy{
		maxRetries:    maxRetries,
		baseDelay:     baseDelay,
		maxDelay:      maxDelay,
		backoffFactor: 2.0,
	}
}

// ShouldRetry determines if a message should be retried
func (p *ExponentialBackoffRetryPolicy) ShouldRetry(message *Message, err error) bool {
	if message == nil {
		return false
	}

	// Don't retry if we've exceeded the maximum retries
	return err != nil && message.RetryCount < p.maxRetries
}

// NextRetryDelay calculates the delay before next retry
func (p *ExponentialBackoffRetryPolicy) NextRetryDelay(message *Message) time.Duration {
	if message == nil {
		return p.baseDelay
	}

	// For zero retries, return base delay
	if message.RetryCount == 0 {
		return p.baseDelay
	}

	// Calculate exponential backoff
	// Start with base delay in milliseconds to avoid overflow
	delayMs := p.baseDelay.Milliseconds()

	// Apply exponential backoff
	for i := 0; i < message.RetryCount; i++ {
		// Check if multiplication would exceed max delay
		maxDelayMs := p.maxDelay.Milliseconds()
		if float64(delayMs) > float64(maxDelayMs)/p.backoffFactor {
			return p.maxDelay
		}

		// Multiply by backoff factor
		delayMs = int64(float64(delayMs) * p.backoffFactor)
	}

	// Convert back to Duration
	delay := time.Duration(delayMs) * time.Millisecond

	// Final check against max delay
	if delay > p.maxDelay {
		return p.maxDelay
	}

	return delay
}

// MaxRetries returns the maximum number of retries
func (p *ExponentialBackoffRetryPolicy) MaxRetries() int {
	return p.maxRetries
}

// LinearRetryPolicy implements a linear retry strategy with fixed delays
type LinearRetryPolicy struct {
	maxRetries int
	delay      time.Duration
}

// NewLinearRetryPolicy creates a new linear retry policy
func NewLinearRetryPolicy(maxRetries int, delay time.Duration) RetryPolicy {
	return &LinearRetryPolicy{
		maxRetries: maxRetries,
		delay:      delay,
	}
}

// ShouldRetry determines if a message should be retried
func (p *LinearRetryPolicy) ShouldRetry(message *Message, err error) bool {
	if message == nil {
		return false
	}

	return err != nil && message.RetryCount < p.maxRetries
}

// NextRetryDelay returns a fixed delay for all retries
func (p *LinearRetryPolicy) NextRetryDelay(_ *Message) time.Duration {
	return p.delay
}

// MaxRetries returns the maximum number of retries
func (p *LinearRetryPolicy) MaxRetries() int {
	return p.maxRetries
}

// NoRetryPolicy implements a policy that never retries
type NoRetryPolicy struct{}

// NewNoRetryPolicy creates a policy that never retries
func NewNoRetryPolicy() RetryPolicy {
	return &NoRetryPolicy{}
}

// ShouldRetry always returns false
func (p *NoRetryPolicy) ShouldRetry(_ *Message, _ error) bool {
	return false
}

// NextRetryDelay returns zero duration
func (p *NoRetryPolicy) NextRetryDelay(_ *Message) time.Duration {
	return 0
}

// MaxRetries returns 0
func (p *NoRetryPolicy) MaxRetries() int {
	return 0
}
