# Model Parameters Configuration (MVP)

CronAI supports model-specific parameters that allow you to fine-tune AI model behavior for each prompt. This document explains how to configure and use these parameters in the MVP release.

## Supported Parameters

### Common Parameters for MVP

The following parameters are supported across all models in the MVP:

| Parameter          | Type   | Range        | Description                                        |
|--------------------|--------|-------------|-------------------------------------------------|
| temperature        | float  | 0.0 - 1.0   | Controls response randomness (higher = more random) |
| max_tokens         | int    | > 0         | Maximum number of tokens to generate                |
| model              | string | -           | Specific model version to use                      |

### Model-Specific Parameters

Each model can also be configured with specific parameters using the prefix notation `model_name.parameter`:

#### OpenAI

| Parameter          | Type   | Description                                     |
|--------------------|--------|-------------------------------------------------|
| openai.model       | string | Specific OpenAI model to use                    |

#### Claude

| Parameter           | Type   | Description                                    |
|---------------------|--------|------------------------------------------------|
| claude.model        | string | Specific Claude model to use                   |

#### Gemini

| Parameter           | Type   | Description                                    |
|---------------------|--------|------------------------------------------------|
| gemini.model        | string | Specific Gemini model to use                   |

## Supported Processors in MVP

**Important Note:** For the MVP release, only a subset of response processors is fully implemented:

| Processor | Status | Description |
|-----------|--------|-------------|
| console   | ✅ Available | Outputs response to the console |
| file      | ✅ Available | Writes response to a file |
| github    | ✅ Available | Creates or updates GitHub issues/comments |
| email     | ⏳ Planned | Email delivery (coming post-MVP) |
| slack     | ⏳ Planned | Slack messaging (coming post-MVP) |
| webhook   | ⏳ Planned | HTTP webhook integration (coming post-MVP) |

When configuring tasks in your `cronai.config` file, please use only the available processors for the MVP release.

## Model-Specific Default Values

### OpenAI Default Settings

- **Default Model**: `gpt-3.5-turbo`
- **Supported Models**:
  - `gpt-3.5-turbo` - Fast and cost-effective for most tasks
  - `gpt-4` - Strong reasoning and instruction following

### Claude Default Settings (v0.0.2+)

- **Default Model**: `claude-3-5-sonnet-latest`
- **Supported Models**:
  
  **Claude 4 Models (Latest - v0.0.2+)**:
  - `claude-4-opus-latest` / `claude-4-opus-20250514` - Most capable model for complex tasks
  - `claude-4-sonnet-latest` - Balanced performance and cost
  - `claude-4-haiku-latest` - Fastest and most efficient
  
  **Claude 3.5 Models (v0.0.2+)**:
  - `claude-3-5-opus-latest` / `claude-3-5-opus-20250120` - Most capable 3.5 model
  - `claude-3-5-sonnet-latest` / `claude-3-5-sonnet-20241022` / `claude-3-5-sonnet-20240620` - Balanced 3.5 model
  - `claude-3-5-haiku-latest` / `claude-3-5-haiku-20241022` - Fast 3.5 model
  
  **Claude 3 Models (v0.0.2+)**:
  - `claude-3-opus-latest` / `claude-3-opus-20240229` - Most powerful Claude 3 model
  - `claude-3-sonnet-latest` / `claude-3-sonnet-20240229` - Balanced performance and speed
  - `claude-3-haiku-latest` / `claude-3-haiku-20240307` - Fast and economical

- **Model Aliases (v0.0.2+)**:
  - `opus` → `claude-4-opus-latest`
  - `sonnet` → `claude-4-sonnet-latest`
  - `haiku` → `claude-4-haiku-latest`
  - `3.5-opus` → `claude-3-5-opus-latest`
  - `3.5-sonnet` → `claude-3-5-sonnet-latest`
  - `3.5-haiku` → `claude-3-5-haiku-latest`
  - `3-opus` → `claude-3-opus-latest`
  - `3-sonnet` → `claude-3-sonnet-latest`
  - `3-haiku` → `claude-3-haiku-latest`

### Gemini Default Settings

- **Default Model**: `gemini-pro`
- **Supported Models**:
  - `gemini-pro` - Original Gemini model
  
## SDK Implementation

CronAI uses official client SDKs for all supported AI models:

- **OpenAI**: Uses the official `github.com/sashabaranov/go-openai` SDK
- **Claude**: Uses the official `github.com/anthropics/anthropic-sdk-go` SDK
- **Gemini**: Uses the official `github.com/google/generative-ai-go` SDK

## Configuration Methods

You can configure model parameters in three ways, listed in order of precedence:

1. **Task-specific parameters in the config file**
2. **Environment variables**
3. **Default values**

### 1. Task-specific Configuration

In the `cronai.config` file, you can specify model parameters using the prefix `model_params:`:

```text
# Format: timestamp model prompt response_processor [variables] [model_params:...]
0 8 * * * claude product_manager file-output.txt model_params:temperature=0.8,model=claude-3-opus-20240229
```text

You can also include both variables and model parameters:

```text
0 9 * * 1 openai report_template github-issue:owner/repo reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,model=gpt-4
```text

#### Using Model-Specific Parameters

For model-specific configuration, use the prefix notation:

```text
# Use OpenAI-specific parameters
0 9 * * 1 openai report_template file-output.txt model_params:openai.model=gpt-4

# Use Claude-specific parameters
0 8 * * * claude product_manager file-output.txt model_params:claude.model=claude-3-opus-20240229

# Use Claude 4 models (v0.0.2+)
0 8 * * * claude product_manager file-output.txt model_params:claude.model=opus
0 8 * * * claude product_manager file-output.txt model_params:claude.model=claude-4-opus-latest

# Use Claude 3.5 models (v0.0.2+)
0 8 * * * claude product_manager file-output.txt model_params:claude.model=3.5-sonnet

# Use Gemini-specific parameters
*/15 * * * * gemini system_health file-output.txt model_params:gemini.model=gemini-pro
```text

### 2. Environment Variables

You can set global defaults for all tasks using environment variables:

```bash
# Common parameters
MODEL_TEMPERATURE=0.7
MODEL_MAX_TOKENS=2048

# Model-specific parameters
OPENAI_MODEL=gpt-4
CLAUDE_MODEL=claude-3-opus-20240229
GEMINI_MODEL=gemini-pro
```text

### 3. Command Line Parameters

When using the `run` command, you can specify model parameters with the `--model-params` flag:

```bash
# Using common parameters
cronai run --model openai --prompt weekly_report --processor file-output.txt --model-params "temperature=0.5,max_tokens=4000,model=gpt-4"

# Using model-specific parameters
cronai run --model gemini --prompt system_health --processor file-output.txt --model-params "gemini.model=gemini-pro"
```text

## Examples

### Low Temperature for Consistent Output

```text
# Run system health check with very precise (low temperature) settings
*/15 * * * * claude system_health file-health.log cluster=Primary model_params:temperature=0.1,max_tokens=1000
```text

### Specific Model Version

```text
# Run weekly with OpenAI using a specific model
0 9 * * 1 openai report_template file-report.log model_params:openai.model=gpt-4
```text

### GitHub Processor Example

```text
# Create a GitHub issue with the weekly report
0 9 * * 1 claude weekly_report github-owner/repo reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.7
```text

### Console Output Example

```text
# Output system health check to console (useful for testing)
*/30 * * * * gemini system_health console model_params:temperature=0.3,max_tokens=2000
```text

### Integration with Variables

Model parameters can be used alongside variables:

```text
0 9 * * 1 openai report_template file-report.log reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,openai.model=gpt-4
```text

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

## Post-MVP Features

The following features are planned for post-MVP releases:

- Advanced parameters like top_p, frequency_penalty, and presence_penalty
- System message customization
- Model fallback mechanism for automatic model switching on failure
- Safety setting configurations for Gemini
- Detailed error handling with retry mechanisms
