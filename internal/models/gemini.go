package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/rshade/cronai/pkg/config"
	"google.golang.org/api/option"
)

// GeminiClient handles interactions with Google's Gemini API
type GeminiClient struct {
	client *genai.Client
	config *config.ModelConfig
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(modelConfig *config.ModelConfig) (*GeminiClient, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		config: modelConfig,
	}, nil
}

// Execute sends a prompt to Gemini and returns the model response
func (c *GeminiClient) Execute(promptContent string) (*ModelResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Before we clean up resources
	defer func() {
		if err := c.client.Close(); err != nil {
			// Log the error but don't override the main error
			fmt.Printf("Warning: Failed to close Gemini client: %v\n", err)
		}
	}()

	modelName := c.getModelName()

	// Create the generative model with the specified model name
	model := c.client.GenerativeModel(modelName)

	// Set the common parameters
	model.SetTemperature(float32(c.config.Temperature))
	model.SetTopP(float32(c.config.TopP))
	model.SetMaxOutputTokens(int32(c.config.MaxTokens))

	// Apply any safety settings if configured
	if len(c.config.GeminiConfig.SafetySettings) > 0 {
		var safetySettings []*genai.SafetySetting
		for category, level := range c.config.GeminiConfig.SafetySettings {
			// Parse the category string to genai.HarmCategory
			harmCategory, err := parseHarmCategory(category)
			if err != nil {
				// Log the error but continue with other settings
				fmt.Printf("Warning: Invalid safety category %s: %v\n", category, err)
				continue
			}

			// Parse the level string to genai.HarmBlockThreshold
			harmLevel, err := parseHarmLevel(level)
			if err != nil {
				// Log the error but continue with other settings
				fmt.Printf("Warning: Invalid safety level %s: %v\n", level, err)
				continue
			}

			safetySettings = append(safetySettings, &genai.SafetySetting{
				Category:  harmCategory,
				Threshold: harmLevel,
			})
		}
		model.SafetySettings = safetySettings
	}

	// Generate content from the prompt
	resp, err := model.GenerateContent(ctx, genai.Text(promptContent))
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract the response text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Extract text from the response
	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content += string(text)
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no text content in Gemini response")
	}

	// Create the model response
	modelResponse := &ModelResponse{
		Content:   content,
		Model:     modelName,
		Timestamp: time.Now(),
	}

	// Add additional metadata
	modelResponse.PromptName = "direct" // Will be overridden by the caller if needed
	modelResponse.ExecutionID = generateExecutionID("gemini", modelResponse.PromptName)

	return modelResponse, nil
}

// getModelName returns the Gemini model name to use
func (c *GeminiClient) getModelName() string {
	if c.config != nil && c.config.GeminiConfig != nil && c.config.GeminiConfig.Model != "" {
		return c.config.GeminiConfig.Model
	}
	// Default to a reasonable model if not specified
	return "gemini-pro"
}

// parseHarmCategory converts a string to a genai.HarmCategory
func parseHarmCategory(category string) (genai.HarmCategory, error) {
	switch category {
	case "harassment":
		return genai.HarmCategoryHarassment, nil
	case "hate_speech", "hate":
		return genai.HarmCategoryHateSpeech, nil
	case "sexually_explicit", "sexual":
		return genai.HarmCategorySexuallyExplicit, nil
	case "dangerous_content", "dangerous":
		return genai.HarmCategoryDangerousContent, nil
	case "derogatory":
		return genai.HarmCategoryDerogatory, nil
	case "toxicity", "toxic":
		return genai.HarmCategoryToxicity, nil
	case "violence":
		return genai.HarmCategoryViolence, nil
	case "medical":
		return genai.HarmCategoryMedical, nil
	default:
		return genai.HarmCategoryUnspecified, fmt.Errorf("unknown harm category: %s", category)
	}
}

// parseHarmLevel converts a string to a genai.HarmBlockThreshold
func parseHarmLevel(level string) (genai.HarmBlockThreshold, error) {
	switch level {
	case "block_none", "none":
		return genai.HarmBlockNone, nil
	case "block_low", "low":
		return genai.HarmBlockLowAndAbove, nil
	case "block_medium", "medium":
		return genai.HarmBlockMediumAndAbove, nil
	case "block_high", "high":
		return genai.HarmBlockOnlyHigh, nil
	case "block", "block_all":
		return genai.HarmBlockUnspecified, nil
	default:
		return genai.HarmBlockUnspecified, fmt.Errorf("unknown harm level: %s", level)
	}
}
