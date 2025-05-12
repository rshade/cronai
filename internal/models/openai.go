package models

import (
	"context"
	"fmt"
	"os"

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

	client := openai.NewClient(apiKey)
	return &OpenAIClient{
		client: client,
		config: modelConfig,
	}, nil
}

// Execute sends a prompt to OpenAI and returns the model response
func (c *OpenAIClient) Execute(promptContent string) (*ModelResponse, error) {
	// Create the chat completion request
	req := openai.ChatCompletionRequest{
		Model:       c.getModelName(),
		Temperature: float32(c.config.Temperature),
		MaxTokens:   c.config.MaxTokens,
		TopP:        float32(c.config.TopP),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.config.OpenAIConfig.SystemMessage,
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
	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("openai API error: %w", err)
	}

	// Extract the response text
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &ModelResponse{
		Content: resp.Choices[0].Message.Content,
		Model:   resp.Model,
	}, nil
}

// getModelName returns the OpenAI model name to use
func (c *OpenAIClient) getModelName() string {
	if c.config != nil && c.config.OpenAIConfig != nil && c.config.OpenAIConfig.Model != "" {
		return c.config.OpenAIConfig.Model
	}
	// Default to a reasonable model if not specified
	return "gpt-3.5-turbo"
}