package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ModelConfig defines common configuration parameters for AI models
type ModelConfig struct {
	// Common parameters
	Temperature    float64 // Controls randomness (0.0-1.0)
	MaxTokens      int     // Maximum response length
	TopP           float64 // Nucleus sampling parameter (0.0-1.0)
	FrequencyPenalty float64 // Penalize frequent tokens (-2.0 to 2.0)
	PresencePenalty  float64 // Penalize new tokens based on presence (-2.0 to 2.0)
	
	// Model-specific configurations
	OpenAIConfig   *OpenAIConfig
	ClaudeConfig   *ClaudeConfig
	GeminiConfig   *GeminiConfig
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	Model         string // GPT model to use (e.g., "gpt-4", "gpt-3.5-turbo")
	SystemMessage string // System message for chat completions
}

// ClaudeConfig holds Anthropic Claude-specific configuration
type ClaudeConfig struct {
	Model         string // Claude model to use (e.g., "claude-3-opus-20240229", "claude-3-sonnet-20240229")
	SystemMessage string // System message for Claude
}

// GeminiConfig holds Google Gemini-specific configuration
type GeminiConfig struct {
	Model         string // Gemini model to use (e.g., "gemini-pro", "gemini-1.5-pro")
	SafetySettings map[string]string // Safety settings for Gemini
}

// DefaultModelConfig returns default configuration values
func DefaultModelConfig() *ModelConfig {
	return &ModelConfig{
		Temperature:     0.7,
		MaxTokens:       1024,
		TopP:            1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		OpenAIConfig:    &OpenAIConfig{
			Model:         "gpt-3.5-turbo",
			SystemMessage: "You are a helpful assistant.",
		},
		ClaudeConfig:    &ClaudeConfig{
			Model:         "claude-3-sonnet-20240229",
			SystemMessage: "You are a helpful assistant.",
		},
		GeminiConfig:    &GeminiConfig{
			Model:         "gemini-pro",
			SafetySettings: make(map[string]string),
		},
	}
}

// LoadFromEnvironment loads model configuration from environment variables
func (mc *ModelConfig) LoadFromEnvironment() {
	// Common parameters
	if temp, err := strconv.ParseFloat(os.Getenv("MODEL_TEMPERATURE"), 64); err == nil && temp >= 0 && temp <= 1 {
		mc.Temperature = temp
	}
	if tokens, err := strconv.Atoi(os.Getenv("MODEL_MAX_TOKENS")); err == nil && tokens > 0 {
		mc.MaxTokens = tokens
	}
	if topP, err := strconv.ParseFloat(os.Getenv("MODEL_TOP_P"), 64); err == nil && topP >= 0 && topP <= 1 {
		mc.TopP = topP
	}
	if freqP, err := strconv.ParseFloat(os.Getenv("MODEL_FREQUENCY_PENALTY"), 64); err == nil {
		mc.FrequencyPenalty = freqP
	}
	if presP, err := strconv.ParseFloat(os.Getenv("MODEL_PRESENCE_PENALTY"), 64); err == nil {
		mc.PresencePenalty = presP
	}

	// OpenAI specific
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		mc.OpenAIConfig.Model = model
	}
	if sysMsg := os.Getenv("OPENAI_SYSTEM_MESSAGE"); sysMsg != "" {
		mc.OpenAIConfig.SystemMessage = sysMsg
	}
	// Additional OpenAI params can be added here

	// Claude specific
	if model := os.Getenv("CLAUDE_MODEL"); model != "" {
		mc.ClaudeConfig.Model = model
	}
	if sysMsg := os.Getenv("CLAUDE_SYSTEM_MESSAGE"); sysMsg != "" {
		mc.ClaudeConfig.SystemMessage = sysMsg
	}
	// Additional Claude params can be added here

	// Gemini specific
	if model := os.Getenv("GEMINI_MODEL"); model != "" {
		mc.GeminiConfig.Model = model
	}

	// Load safety settings from environment if provided
	// Format: GEMINI_SAFETY_SETTINGS=category1=level1,category2=level2
	if safetySettings := os.Getenv("GEMINI_SAFETY_SETTINGS"); safetySettings != "" {
		if mc.GeminiConfig.SafetySettings == nil {
			mc.GeminiConfig.SafetySettings = make(map[string]string)
		}

		for _, setting := range strings.Split(safetySettings, ",") {
			parts := strings.SplitN(setting, "=", 2)
			if len(parts) == 2 {
				category := strings.TrimSpace(parts[0])
				level := strings.TrimSpace(parts[1])
				mc.GeminiConfig.SafetySettings[category] = level
			}
		}
	}
}

// ParseModelParams parses model parameters from a comma-separated string
// Format: "temperature=0.8,max_tokens=2048,model=gpt-4"
func ParseModelParams(paramsStr string) (map[string]string, error) {
	params := make(map[string]string)
	
	if paramsStr == "" {
		return params, nil
	}
	
	for _, paramPair := range strings.Split(paramsStr, ",") {
		parts := strings.SplitN(paramPair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid parameter format: %s", paramPair)
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		params[key] = value
	}
	
	return params, nil
}

// UpdateFromParams updates the configuration based on the provided parameters
func (mc *ModelConfig) UpdateFromParams(params map[string]string) error {
	for key, value := range params {
		// Handle common parameters
		switch strings.ToLower(key) {
		case "temperature":
			temp, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid temperature value: %s", value)
			}
			if temp < 0 || temp > 1 {
				return fmt.Errorf("temperature must be between 0 and 1, got: %f", temp)
			}
			mc.Temperature = temp

		case "max_tokens", "maxtokens":
			tokens, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid max_tokens value: %s", value)
			}
			if tokens <= 0 {
				return fmt.Errorf("max_tokens must be positive, got: %d", tokens)
			}
			mc.MaxTokens = tokens

		case "top_p", "topp":
			topP, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid top_p value: %s", value)
			}
			if topP < 0 || topP > 1 {
				return fmt.Errorf("top_p must be between 0 and 1, got: %f", topP)
			}
			mc.TopP = topP

		case "frequency_penalty", "frequencypenalty":
			freqP, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid frequency_penalty value: %s", value)
			}
			if freqP < -2.0 || freqP > 2.0 {
				return fmt.Errorf("frequency_penalty must be between -2.0 and 2.0, got: %f", freqP)
			}
			mc.FrequencyPenalty = freqP

		case "presence_penalty", "presencepenalty":
			presP, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid presence_penalty value: %s", value)
			}
			if presP < -2.0 || presP > 2.0 {
				return fmt.Errorf("presence_penalty must be between -2.0 and 2.0, got: %f", presP)
			}
			mc.PresencePenalty = presP

		case "model":
			// Apply model to all model configs to handle the generic case
			// The actual use will be determined by which client is selected
			// Note: If a model-specific parameter is also provided (e.g. openai.model),
			// it will override this generic setting for that specific model client
			if mc.OpenAIConfig != nil {
				mc.OpenAIConfig.Model = value
			}
			if mc.ClaudeConfig != nil {
				mc.ClaudeConfig.Model = value
			}
			if mc.GeminiConfig != nil {
				mc.GeminiConfig.Model = value
			}

			// Safe logging for debugging if needed
			// log.Printf("Generic model parameter applied to all model configurations")

		case "system_message", "systemmessage":
			// Apply system message to models that support it
			// Note: If a model-specific system message is also provided (e.g. openai.system_message),
			// it will override this generic setting for that specific model client
			if mc.OpenAIConfig != nil {
				mc.OpenAIConfig.SystemMessage = value
			}
			if mc.ClaudeConfig != nil {
				mc.ClaudeConfig.SystemMessage = value
			}

			// Safe logging for debugging if needed
			// log.Printf("Generic system message applied to applicable model configurations")

		// Model-specific parameters with prefixes
		default:
			// Handle prefixed model-specific parameters
			// We'll allow unknown parameters to pass through silently
			// This allows for forward compatibility with future parameters
			_ = mc.handleModelSpecificParam(key, value)
		}
	}

	return nil
}

// handleModelSpecificParam processes model-specific parameters with prefixes
func (mc *ModelConfig) handleModelSpecificParam(key, value string) error {
	// Process parameters with prefixes: openai.*, claude.*, gemini.*
	parts := strings.SplitN(strings.ToLower(key), ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("unrecognized parameter: %s", key)
	}

	modelPrefix := parts[0]
	paramName := parts[1]

	switch modelPrefix {
	case "openai":
		if mc.OpenAIConfig == nil {
			mc.OpenAIConfig = &OpenAIConfig{}
		}
		return mc.handleOpenAIParam(paramName, value)
	case "claude":
		if mc.ClaudeConfig == nil {
			mc.ClaudeConfig = &ClaudeConfig{}
		}
		return mc.handleClaudeParam(paramName, value)
	case "gemini":
		if mc.GeminiConfig == nil {
			mc.GeminiConfig = &GeminiConfig{}
		}
		return mc.handleGeminiParam(paramName, value)
	default:
		return fmt.Errorf("unknown model prefix: %s", modelPrefix)
	}
}

// handleOpenAIParam handles OpenAI-specific parameters
func (mc *ModelConfig) handleOpenAIParam(param, value string) error {
	switch param {
	case "model":
		mc.OpenAIConfig.Model = value
	case "system_message", "systemmessage":
		mc.OpenAIConfig.SystemMessage = value
	default:
		return fmt.Errorf("unknown OpenAI parameter: %s", param)
	}
	return nil
}

// handleClaudeParam handles Claude-specific parameters
func (mc *ModelConfig) handleClaudeParam(param, value string) error {
	switch param {
	case "model":
		mc.ClaudeConfig.Model = value
	case "system_message", "systemmessage":
		mc.ClaudeConfig.SystemMessage = value
	default:
		return fmt.Errorf("unknown Claude parameter: %s", param)
	}
	return nil
}

// handleGeminiParam handles Gemini-specific parameters
func (mc *ModelConfig) handleGeminiParam(param, value string) error {
	switch param {
	case "model":
		mc.GeminiConfig.Model = value
	case "safety_setting", "safetysetting":
		// Parse safety settings in format "category=level"
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid safety setting format (should be category=level): %s", value)
		}
		if mc.GeminiConfig.SafetySettings == nil {
			mc.GeminiConfig.SafetySettings = make(map[string]string)
		}
		mc.GeminiConfig.SafetySettings[parts[0]] = parts[1]
	default:
		return fmt.Errorf("unknown Gemini parameter: %s", param)
	}
	return nil
}

// Validate validates the configuration values
func (mc *ModelConfig) Validate() error {
	if mc.Temperature < 0 || mc.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1, got: %f", mc.Temperature)
	}

	if mc.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be positive, got: %d", mc.MaxTokens)
	}

	if mc.TopP < 0 || mc.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1, got: %f", mc.TopP)
	}

	if mc.FrequencyPenalty < -2.0 || mc.FrequencyPenalty > 2.0 {
		return fmt.Errorf("frequency_penalty must be between -2.0 and 2.0, got: %f", mc.FrequencyPenalty)
	}

	if mc.PresencePenalty < -2.0 || mc.PresencePenalty > 2.0 {
		return fmt.Errorf("presence_penalty must be between -2.0 and 2.0, got: %f", mc.PresencePenalty)
	}

	return nil
}