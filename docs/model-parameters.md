# Model Parameters Configuration

CronAI supports model-specific parameters that allow you to fine-tune AI model behavior for each prompt. This document explains how to configure and use these parameters.

## Supported Parameters

### Common Parameters

The following parameters are supported across all models:

| Parameter          | Type   | Range        | Description                                        |
|--------------------|--------|-------------|----------------------------------------------------|
| temperature        | float  | 0.0 - 1.0   | Controls response randomness (higher = more random) |
| max_tokens         | int    | > 0         | Maximum number of tokens to generate                |
| top_p              | float  | 0.0 - 1.0   | Nucleus sampling parameter                         |
| frequency_penalty  | float  | -2.0 - 2.0  | Penalize frequent tokens                           |
| presence_penalty   | float  | -2.0 - 2.0  | Penalize new tokens based on presence              |
| model              | string | -           | Specific model version to use                      |
| system_message     | string | -           | System message for the model                       |

### Model-Specific Parameters

Each model can also be configured with specific parameters using the prefix notation `model_name.parameter`:

#### OpenAI

| Parameter          | Type   | Description                                     |
|--------------------|--------|-------------------------------------------------|
| openai.model       | string | Specific OpenAI model to use                    |
| openai.system_message | string | System message specific to OpenAI            |

#### Claude

| Parameter           | Type   | Description                                    |
|---------------------|--------|------------------------------------------------|
| claude.model        | string | Specific Claude model to use                   |
| claude.system_message | string | System message specific to Claude            |

#### Gemini

| Parameter           | Type   | Description                                    |
|---------------------|--------|------------------------------------------------|
| gemini.model        | string | Specific Gemini model to use                   |
| gemini.safety_setting | string | Safety settings in format "category=level"   |

## Model-Specific Default Values

### OpenAI
- **Default Model**: `gpt-3.5-turbo`
- **System Message**: `You are a helpful assistant.`
- **Supported Models**: 
  - `gpt-3.5-turbo` - Fast and cost-effective for most tasks
  - `gpt-3.5-turbo-16k` - Extended context length for 3.5
  - `gpt-4` - Strong reasoning and instruction following
  - `gpt-4o` - Optimized version with improved speed and 16k context
  - `gpt-4-turbo` - Advanced capabilities with 128k context
  - `gpt-4-32k` - Extended context length for GPT-4
  - `gpt-4.1` - Latest model with improved coding capabilities

### Claude
- **Default Model**: `claude-3-sonnet-20240229`
- **System Message**: `You are a helpful assistant.`
- **Supported Models**: 
  - `claude-3-opus-20240229` - Most powerful Claude model for complex tasks
  - `claude-3-sonnet-20240229` - Balanced performance and speed
  - `claude-3-haiku-20240307` - Fast and economical
  - `claude-3.5-sonnet` - Enhanced reasoning and capabilities
  - `claude-3.7-sonnet` - Latest model with advanced multimodal capabilities

### Gemini
- **Default Model**: `gemini-pro`
- **Supported Models**: 
  - `gemini-pro` - Original Gemini model
  - `gemini-1.5-pro` - Enhanced capabilities
  - `gemini-1.5-flash` - Optimized for speed
  - `gemini-2.5-pro` - Latest model with enhanced reasoning
  - `gemini-2.5-pro-latest` - Most current version with 2M token context
- **Safety Categories**: 
  - `harassment` - Harmful content targeting identity and/or protected attributes
  - `hate_speech`, `hate` - Content that is rude, disrespectful, or profane
  - `sexually_explicit`, `sexual` - Content meant to arouse sexual excitement
  - `dangerous_content`, `dangerous` - Promotes, facilitates, or encourages harmful acts
- **Safety Levels**:
  - `block_none`, `none` - Block no content based on safety
  - `block_low`, `low` - Block only very harmful content
  - `block_medium`, `medium` - Block moderately harmful content (default)
  - `block_high`, `high` - Block content that has a low risk of harm
  - `block`, `block_all` - Block all potentially harmful content

## SDK Implementation

CronAI uses official client SDKs for all supported AI models:

- **OpenAI**: Uses the official `github.com/sashabaranov/go-openai` SDK
- **Claude**: Uses the official `github.com/anthropics/anthropic-sdk-go` SDK
- **Gemini**: Uses the official `github.com/google/generative-ai-go` SDK

These SDKs ensure reliable connectivity to the respective AI services and proper implementation of all model parameters.

## Configuration Methods

You can configure model parameters in three ways, listed in order of precedence:

1. **Task-specific parameters in the config file**
2. **Environment variables**
3. **Default values**

### 1. Task-specific Configuration

In the `cronai.config` file, you can specify model parameters after variables using the prefix `model_params:`:

```
# Format: timestamp model prompt.md response_processor [variables] [model_params:...]
0 8 * * * claude product_manager slack-pm-channel model_params:temperature=0.8,model=claude-3-opus-20240229
```

You can also include both variables and model parameters:

```
0 9 * * 1 openai report_template email-team@company.com reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,model=gpt-4
```

#### Using Model-Specific Parameters

For model-specific configuration, use the prefix notation:

```
# Use OpenAI-specific parameters
0 9 * * 1 openai report_template email-team@company.com model_params:openai.model=gpt-4,openai.system_message=You are a business analyst.

# Use Claude-specific parameters
0 8 * * * claude product_manager slack-pm-channel model_params:claude.model=claude-3-opus-20240229,claude.system_message=You are a product manager.

# Use Gemini-specific parameters
*/15 * * * * gemini system_health webhook-monitoring model_params:gemini.model=gemini-1.5-pro,gemini.safety_setting=harmful=block,gemini.safety_setting=harassment=warn
```

### 2. Environment Variables

You can set global defaults for all tasks using environment variables:

```bash
# Common parameters
MODEL_TEMPERATURE=0.7
MODEL_MAX_TOKENS=2048
MODEL_TOP_P=0.9
MODEL_FREQUENCY_PENALTY=0.0
MODEL_PRESENCE_PENALTY=0.0

# Model-specific parameters
OPENAI_MODEL=gpt-4
OPENAI_SYSTEM_MESSAGE="You are a helpful assistant specialized in business analysis."
OPENAI_BASE_URL="https://your-custom-endpoint.com/v1" # Optional custom endpoint

CLAUDE_MODEL=claude-3-opus-20240229
CLAUDE_SYSTEM_MESSAGE="You are a helpful assistant specialized in technical documentation."
ANTHROPIC_BASE_URL="https://your-custom-endpoint.com" # Optional custom endpoint

GEMINI_MODEL=gemini-pro
GEMINI_SAFETY_SETTINGS="harmful=block,harassment=warn"
```

### 3. Command Line Parameters

When using the `run` command, you can specify model parameters with the `--model-params` flag:

```bash
# Using common parameters
cronai run --model openai --prompt weekly_report --processor email-team@company.com --model-params "temperature=0.5,max_tokens=4000,model=gpt-4"

# Using model-specific parameters
cronai run --model gemini --prompt system_health --processor webhook-monitoring --model-params "gemini.model=gemini-1.5-pro,gemini.safety_setting=harmful=block"
```

## Examples

### Low Temperature for Consistent Output

```
# Run system health check with very precise (low temperature) settings
*/15 * * * * claude system_health webhook-monitoring cluster=Primary model_params:temperature=0.1,max_tokens=1000
```

### Custom System Message

```
# Run daily at 10 PM using Claude with custom system message
0 22 * * * claude test_prompt slack-dev-channel model_params:claude.system_message=You are a systems analyst providing clear, actionable insights.
```

### Specific Model Version

```
# Run weekly with OpenAI using a specific model
0 9 * * 1 openai report_template email-team@company.com model_params:openai.model=gpt-4-turbo
```

### Safety Settings for Gemini

```
# Run with Gemini using safety settings
0 9-17 * * 1-5 gemini monitoring_check log-to-file model_params:gemini.safety_setting=harmful=block,gemini.safety_setting=harassment=warn
```

### Mixed Common and Model-Specific Parameters

```
# Use both common and model-specific parameters
0 9 * * 1 openai report_template email-team@company.com model_params:temperature=0.7,max_tokens=2000,openai.model=gpt-4,openai.system_message=You are a report generator.
```

## Integration with Variables

Model parameters can be used alongside variables:

```
0 9 * * 1 openai report_template email-team@company.com reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,openai.model=gpt-4
```

## Advanced Configuration

### Custom API Endpoints

For organizations using proxy services or custom endpoints for AI models, CronAI supports custom base URLs through environment variables:

- **OpenAI**: Set `OPENAI_BASE_URL` environment variable to point to a custom endpoint
- **Claude**: Set `ANTHROPIC_BASE_URL` environment variable to point to a custom endpoint

### Timeout Configuration

All model clients have a default timeout of 120 seconds (2 minutes) for API requests. This can be adjusted by setting the appropriate environment variables:

- **OpenAI**: Set `OPENAI_TIMEOUT` to the desired timeout in seconds
- **Claude**: Set `ANTHROPIC_TIMEOUT` to the desired timeout in seconds 
- **Gemini**: Set `GEMINI_TIMEOUT` to the desired timeout in seconds

### Error Handling

CronAI implements robust error handling for all model clients:

- Network errors are properly captured and reported
- API rate limiting is handled with appropriate backoff
- Token limit errors are clearly reported
- Safety filter blocks are identified and reported

When an error occurs during model execution, CronAI will log detailed information while ensuring that sensitive data like prompt contents are not leaked in logs.