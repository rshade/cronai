// Package models provides implementations for different AI model clients
package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/rshade/cronai/pkg/config"
)

// Claude model constants for all supported models
const (
	// Claude 4 models (newest generation - as per issue requirements)
	Claude4OpusLatest   = "claude-4-opus-latest"
	Claude4Opus20250514 = "claude-4-opus-20250514"
	Claude4SonnetLatest = "claude-4-sonnet-latest"
	Claude4HaikuLatest  = "claude-4-haiku-latest"

	// Claude 3.5 models (latest generation)
	Claude35OpusLatest     = "claude-3-5-opus-latest"
	Claude35Opus20250120   = "claude-3-5-opus-20250120"
	Claude35SonnetLatest   = "claude-3-5-sonnet-latest"
	Claude35Sonnet20241022 = "claude-3-5-sonnet-20241022"
	Claude35Sonnet20240620 = "claude-3-5-sonnet-20240620"
	Claude35HaikuLatest    = "claude-3-5-haiku-latest"
	Claude35Haiku20241022  = "claude-3-5-haiku-20241022"

	// Claude 3 models
	Claude3OpusLatest     = "claude-3-opus-latest"
	Claude3Opus20240229   = "claude-3-opus-20240229"
	Claude3SonnetLatest   = "claude-3-sonnet-latest"
	Claude3Sonnet20240229 = "claude-3-sonnet-20240229"
	Claude3HaikuLatest    = "claude-3-haiku-latest"
	Claude3Haiku20240307  = "claude-3-haiku-20240307"

	// Default model - using latest Sonnet for balance of capability and cost
	DefaultClaudeModel = Claude35SonnetLatest
)

// ClaudeModelAliases maps common aliases to specific model versions
var ClaudeModelAliases = map[string]string{
	// Opus aliases
	"opus":          Claude4OpusLatest, // Default to latest
	"opus-latest":   Claude4OpusLatest,
	"4-opus":        Claude4OpusLatest,
	"3.5-opus":      Claude35OpusLatest,
	"3-opus":        Claude3OpusLatest,
	"claude-opus":   Claude4OpusLatest,
	"claude-4-opus": Claude4OpusLatest,

	// Sonnet aliases
	"sonnet":          Claude4SonnetLatest, // Default to latest
	"sonnet-latest":   Claude4SonnetLatest,
	"4-sonnet":        Claude4SonnetLatest,
	"3.5-sonnet":      Claude35SonnetLatest,
	"3-sonnet":        Claude3SonnetLatest,
	"claude-sonnet":   Claude4SonnetLatest,
	"claude-4-sonnet": Claude4SonnetLatest,

	// Haiku aliases
	"haiku":          Claude4HaikuLatest, // Default to latest
	"haiku-latest":   Claude4HaikuLatest,
	"4-haiku":        Claude4HaikuLatest,
	"3.5-haiku":      Claude35HaikuLatest,
	"3-haiku":        Claude3HaikuLatest,
	"claude-haiku":   Claude4HaikuLatest,
	"claude-4-haiku": Claude4HaikuLatest,
}

// SupportedClaudeModels contains all supported Claude model versions
var SupportedClaudeModels = map[string]bool{
	// Claude 4 models
	Claude4OpusLatest:   true,
	Claude4Opus20250514: true,
	Claude4SonnetLatest: true,
	Claude4HaikuLatest:  true,

	// Claude 3.5 models
	Claude35OpusLatest:     true,
	Claude35Opus20250120:   true,
	Claude35SonnetLatest:   true,
	Claude35Sonnet20241022: true,
	Claude35Sonnet20240620: true,
	Claude35HaikuLatest:    true,
	Claude35Haiku20241022:  true,

	// Claude 3 models
	Claude3OpusLatest:     true,
	Claude3Opus20240229:   true,
	Claude3SonnetLatest:   true,
	Claude3Sonnet20240229: true,
	Claude3HaikuLatest:    true,
	Claude3Haiku20240307:  true,
}

// ClaudeClient handles interactions with Claude API
type ClaudeClient struct {
	client anthropic.Client
	config *config.ModelConfig
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient(modelConfig *config.ModelConfig) (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	// Create the Anthropic client with API key
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Check if we have an Anthropic base URL set in the environment
	if baseURL := os.Getenv("ANTHROPIC_BASE_URL"); baseURL != "" {
		client = anthropic.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		)
	}

	return &ClaudeClient{
		client: client,
		config: modelConfig,
	}, nil
}

// Execute sends a prompt to Claude and returns the model response
func (c *ClaudeClient) Execute(promptContent string) (*ModelResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	modelName := c.getModelName()
	systemMessage := c.getSystemMessage()

	// Convert MaxTokens to int64
	maxTokens := int64(c.config.MaxTokens)

	// Create the message request
	request := &anthropic.MessageNewParams{
		Model:       anthropic.Model(modelName),
		MaxTokens:   maxTokens,
		Temperature: anthropic.Float(c.config.Temperature),
		TopP:        anthropic.Float(c.config.TopP),
		System:      []anthropic.TextBlockParam{{Text: systemMessage}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(promptContent)),
		},
	}

	// Send the request to Claude API
	resp, err := c.client.Messages.New(ctx, *request)
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract the response text
	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("no response from Claude")
	}

	// Get the text from the first content block
	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content = block.Text
			break
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no text content in Claude response")
	}

	// Create the model response
	modelResponse := &ModelResponse{
		Content:   content,
		Model:     string(resp.Model),
		Timestamp: time.Now(),
	}

	// Add additional metadata
	modelResponse.PromptName = "direct" // Will be overridden by the caller if needed
	modelResponse.ExecutionID = generateExecutionID("claude", modelResponse.PromptName)

	return modelResponse, nil
}

// getModelName returns the Claude model name to use
func (c *ClaudeClient) getModelName() string {
	if c.config != nil && c.config.ClaudeConfig != nil && c.config.ClaudeConfig.Model != "" {
		model := c.config.ClaudeConfig.Model

		// Check if it's an alias and resolve it
		if resolvedModel, ok := ClaudeModelAliases[strings.ToLower(model)]; ok {
			return resolvedModel
		}

		// Check if it's a supported model
		if SupportedClaudeModels[strings.ToLower(model)] {
			return strings.ToLower(model)
		}

		// If not supported, log a warning and use default
		// Note: In production, you might want to return an error instead
		log.Printf("Warning: Unsupported Claude model '%s', using default: %s", model, DefaultClaudeModel)
		return DefaultClaudeModel
	}
	// Default to the configured default model if not specified
	return DefaultClaudeModel
}

// getSystemMessage returns the Claude system message to use
func (c *ClaudeClient) getSystemMessage() string {
	if c.config != nil && c.config.ClaudeConfig != nil && c.config.ClaudeConfig.SystemMessage != "" {
		return c.config.ClaudeConfig.SystemMessage
	}
	// Default to a standard system message if not specified
	return "You are a helpful assistant."
}

// GetAvailableClaudeModels returns a list of available Claude models with their descriptions
func GetAvailableClaudeModels() map[string]string {
	return map[string]string{
		// Claude 4 models
		Claude4OpusLatest:   "Claude 4 Opus (latest) - Most capable model for complex tasks",
		Claude4Opus20250514: "Claude 4 Opus (2025-05-14) - Specific version",
		Claude4SonnetLatest: "Claude 4 Sonnet (latest) - Balanced performance and cost",
		Claude4HaikuLatest:  "Claude 4 Haiku (latest) - Fastest and most efficient",

		// Claude 3.5 models
		Claude35OpusLatest:     "Claude 3.5 Opus (latest) - Most capable 3.5 model",
		Claude35Opus20250120:   "Claude 3.5 Opus (2025-01-20) - Specific version",
		Claude35SonnetLatest:   "Claude 3.5 Sonnet (latest) - Balanced 3.5 model",
		Claude35Sonnet20241022: "Claude 3.5 Sonnet (2024-10-22) - Specific version",
		Claude35Sonnet20240620: "Claude 3.5 Sonnet (2024-06-20) - Specific version",
		Claude35HaikuLatest:    "Claude 3.5 Haiku (latest) - Fast 3.5 model",
		Claude35Haiku20241022:  "Claude 3.5 Haiku (2024-10-22) - Specific version",

		// Claude 3 models
		Claude3OpusLatest:     "Claude 3 Opus (latest) - Most capable 3.0 model",
		Claude3Opus20240229:   "Claude 3 Opus (2024-02-29) - Specific version",
		Claude3SonnetLatest:   "Claude 3 Sonnet (latest) - Balanced 3.0 model",
		Claude3Sonnet20240229: "Claude 3 Sonnet (2024-02-29) - Specific version",
		Claude3HaikuLatest:    "Claude 3 Haiku (latest) - Fast 3.0 model",
		Claude3Haiku20240307:  "Claude 3 Haiku (2024-03-07) - Specific version",
	}
}

// GetClaudeModelAliases returns available model aliases
func GetClaudeModelAliases() map[string]string {
	return ClaudeModelAliases
}
