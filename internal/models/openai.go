package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rshade/cronai/pkg/config"
	"github.com/sashabaranov/go-openai"
)

// OpenAIClient handles interactions with OpenAI API
type OpenAIClient struct {
	client *openai.Client
	config *config.ModelConfig
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(modelConfig *config.ModelConfig) (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	// Check if we have an OpenAI base URL set in the environment
	baseURL := os.Getenv("OPENAI_BASE_URL")
	var client *openai.Client

	if baseURL != "" {
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = baseURL
		client = openai.NewClientWithConfig(config)
	} else {
		client = openai.NewClient(apiKey)
	}

	return &OpenAIClient{
		client: client,
		config: modelConfig,
	}, nil
}

// Execute sends a prompt to OpenAI and returns the model response
func (c *OpenAIClient) Execute(promptContent string) (*ModelResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Create the chat completion request
	req := openai.ChatCompletionRequest{
		Model:       c.getModelName(),
		Temperature: float32(c.config.Temperature),
		MaxTokens:   c.config.MaxTokens,
		TopP:        float32(c.config.TopP),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.getSystemMessage(),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: promptContent,
			},
		},
		FrequencyPenalty: float32(c.config.FrequencyPenalty),
		PresencePenalty:  float32(c.config.PresencePenalty),
	}

	// Make the API call
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("openai API error: %w", err)
	}

	// Extract the response text
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Create the model response
	modelResponse := &ModelResponse{
		Content:   resp.Choices[0].Message.Content,
		Model:     resp.Model,
		Timestamp: time.Now(),
	}

	// Add additional metadata
	modelResponse.PromptName = "direct" // Will be overridden by the caller if needed
	modelResponse.ExecutionID = generateExecutionID("openai", modelResponse.PromptName)

	return modelResponse, nil
}

// getModelName returns the OpenAI model name to use
func (c *OpenAIClient) getModelName() string {
	if c.config != nil && c.config.OpenAIConfig != nil && c.config.OpenAIConfig.Model != "" {
		return c.config.OpenAIConfig.Model
	}
	// Default to a reasonable model if not specified
	return "gpt-3.5-turbo"
}

// getSystemMessage returns the OpenAI system message to use
func (c *OpenAIClient) getSystemMessage() string {
	if c.config != nil && c.config.OpenAIConfig != nil && c.config.OpenAIConfig.SystemMessage != "" {
		return c.config.OpenAIConfig.SystemMessage
	}
	// Default to a standard system message if not specified
	return "You are a helpful assistant."
}
