// Package integration provides integration tests for the CronAI application.
package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIntegration_APIKeys ensures the required API keys are available
// This test helps diagnose issues with GitHub Actions secrets
func TestIntegration_APIKeys(t *testing.T) {
	// Skip if not running real integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping API key validation - not running real integration tests")
	}

	// Track missing keys for summary
	var missingKeys []string

	// Check if OpenAI API key is set
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		missingKeys = append(missingKeys, "OPENAI_API_KEY")
		t.Log("OPENAI_API_KEY is not set - integration tests will not use real OpenAI API")
	} else {
		// Don't log the actual key, just that it exists and basic validation
		assert.True(t, len(openAIKey) > 10, "OPENAI_API_KEY appears to be valid (non-empty and sufficient length)")

		// Check if it follows expected format
		if len(openAIKey) >= 3 {
			assert.Equal(t, "sk-", openAIKey[0:3], "OPENAI_API_KEY should start with 'sk-'")
		}
	}

	// Check if GitHub token is set
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		missingKeys = append(missingKeys, "GITHUB_TOKEN")
		t.Log("GITHUB_TOKEN is not set - integration tests will not use real GitHub API")
	} else {
		// Don't log the actual token, just that it exists and basic validation
		assert.True(t, len(githubToken) > 10, "GITHUB_TOKEN appears to be valid (non-empty and sufficient length)")
	}

	// If we're missing keys, fail the test with a summary
	if len(missingKeys) > 0 {
		t.Fatalf("Missing required API keys for real integration tests: %v. Tests will run in mock mode.", missingKeys)
	}
}
