// Package models provides implementations for different AI model clients
package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/rshade/cronai/pkg/config"
)

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
		Model:       modelName,
		MaxTokens:   maxTokens,
		Temperature: anthropic.Float(c.config.Temperature),
		TopP:        anthropic.Float(c.config.TopP),
		System:      []anthropic.TextBlockParam{{Text: systemMessage}},
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{
						OfRequestTextBlock: &anthropic.TextBlockParam{
							Text: promptContent,
						},
					},
				},
			},
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
		return c.config.ClaudeConfig.Model
	}
	// Default to a reasonable model if not specified
	return "claude-3-sonnet-20240229"
}

// getSystemMessage returns the Claude system message to use
func (c *ClaudeClient) getSystemMessage() string {
	if c.config != nil && c.config.ClaudeConfig != nil && c.config.ClaudeConfig.SystemMessage != "" {
		return c.config.ClaudeConfig.SystemMessage
	}
	// Default to a standard system message if not specified
	return "You are a helpful assistant."
}
