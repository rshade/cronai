# cronai

A command-line utility to run AI model prompts on a cron-type schedule.

## Overview

CronAI allows you to schedule AI prompts to run at specified intervals and process the responses automatically. It supports multiple AI models and response processors.

## Features

- Cron-style scheduling
- Support for multiple AI models (OpenAI, Claude, Gemini)
- Customizable prompts stored as markdown files
- Flexible response processing options (email, Slack, webhooks, file output)
- Can be run as a systemd service

## Installation

```bash
# Install directly
go install github.com/rshade/cronai/cmd/cronai@latest

# Or clone and build
git clone https://github.com/rshade/cronai.git
cd cronai
go build -o cronai ./cmd/cronai
```

## Configuration

Create a configuration file called `cronai.config` with your scheduled tasks.

### Format

```
timestamp model prompt response_processor [template] [variables] [model_params:...]
```

- **timestamp**: Standard cron format (minute hour day-of-month month day-of-week)
- **model**: AI model to use (openai, claude, gemini)
- **prompt**: Name of prompt file in cron_prompts directory (with or without .md extension)
- **response_processor**: How to process the response (email, slack, webhook, file)
- **template** (optional): Name of template to use for formatting the response
- **variables** (optional): Variables to replace in the prompt file, in the format `key1=value1,key2=value2,...`
- **model_params** (optional): Model-specific parameters in the format `model_params:param1=value1,param2=value2,...`

### Example Configuration

```
# Run daily at 8 AM using Claude, sending to slack
0 8 * * * claude product_manager slack-pm-channel

# Run every Monday at 9 AM using OpenAI, sending to email
0 9 * * 1 openai weekly_report email-team@company.com

# Run monthly report with custom template and variables on the 1st of each month
0 9 1 * * claude report_template email-execs@company.com monthly_report reportType=Monthly,date={{CURRENT_DATE}},project=CronAI

# Run with custom model parameters (temperature and specific model version)
0 9 * * 1 openai weekly_report email-team@company.com model_params:temperature=0.5,model=gpt-4
```

See [cronai.config.example](cronai.config.example), [cronai.config.variables.example](cronai.config.variables.example), and [cronai.config.model-params.example](cronai.config.model-params.example) for more examples.

## Prompt Management

CronAI includes a robust file-based prompt management system to help you organize, discover, and reuse prompts efficiently.

### Prompt Structure and Organization

Prompts are organized in a category-based directory structure:

```
cron_prompts/
├── monitoring/     # Monitoring-related prompts
├── reports/        # Report generation prompts
├── system/         # System operations prompts
├── templates/      # Reusable templates
└── [custom]/       # Your custom categories
```

### Prompt Files

Each prompt is stored as a markdown file and can include an optional metadata section:

```markdown
---
name: System Health Check
description: Analyzes system health metrics
author: CronAI Team
version: 1.0
category: system
tags: health, monitoring, metrics
variables:
  - name: cpu_usage
    description: Current CPU usage percentage
  - name: memory_usage
    description: Current memory usage percentage
---

# System Health Check

Analyze the following system metrics:
- CPU Usage: {{cpu_usage}}%
- Memory Usage: {{memory_usage}}%
...
```

### Variables in Prompts

You can use variables in prompt files with the syntax `{{variable_name}}`. These variables will be replaced with values from the configuration:

```markdown
# {{reportType}} Report for {{date}}

## Overview

This is an automatically generated {{reportType}} report for {{project}} created on {{date}}.

## Team Details

Team: {{team}}
```

Special variables that are automatically populated:
- `{{CURRENT_DATE}}`: Current date in YYYY-MM-DD format
- `{{CURRENT_TIME}}`: Current time in HH:MM:SS format
- `{{CURRENT_DATETIME}}`: Current date and time in YYYY-MM-DD HH:MM:SS format

### Conditional Logic in Prompts

CronAI supports conditional blocks in prompts, allowing for dynamic content based on variables:

```markdown
{{if eq .Variables.environment "production"}}
## Production Environment
This is a production system. Be cautious with recommendations.
{{else if eq .Variables.environment "staging"}}
## Staging Environment
This is a staging system. Testing is allowed.
{{else}}
## Development Environment
This is a development system. Feel free to experiment.
{{end}}

{{if hasVar .Variables "includeMetrics"}}
## Metrics Analysis
Detailed metrics included as requested.
{{end}}

{{if gt (getVar .Variables "errorCount" "0") "5"}}
High error count detected: {{.Variables.errorCount}}
{{else}}
Low or no errors detected.
{{end}}
```

Available conditional features:
- If-else branching based on variable values
- Variable existence checks with `hasVar`
- Default values with `getVar`
- String comparisons: eq (equals), ne (not equals), contains, hasPrefix, hasSuffix
- Numeric comparisons: lt (less than), gt (greater than), le (less than or equal), ge (greater than or equal)
- Logical operators: and, or, not
- Nested conditionals for complex logic

For a full guide to conditional syntax and examples, see [docs/conditional-templates.md](docs/conditional-templates.md).

### Prompt Composition

You can include content from other prompt files using the include directive:

```markdown
{{include "templates/common_header.md"}}

# Main Content

Your specific prompt content goes here.

{{include "templates/common_footer.md"}}
```

### Managing Prompts

CronAI includes CLI commands for prompt management:

```bash
# List all prompts
cronai prompt list

# List prompts in a specific category
cronai prompt list --category system

# Search for prompts
cronai prompt search "health check"
cronai prompt search --content --query "CPU"

# Show prompt details
cronai prompt show system/system_health --vars

# Preview a prompt with variables
cronai prompt preview system/system_health --vars "cpu_usage=85,memory_usage=70"
```

For more details, see [docs/prompt-management.md](docs/prompt-management.md).

## Model Parameters

You can configure model-specific parameters to fine-tune AI model behavior. The supported parameters include:

| Parameter          | Type   | Range        | Description                                        |
|--------------------|--------|-------------|----------------------------------------------------|
| temperature        | float  | 0.0 - 1.0   | Controls response randomness (higher = more random) |
| max_tokens         | int    | > 0         | Maximum number of tokens to generate                |
| top_p              | float  | 0.0 - 1.0   | Nucleus sampling parameter                         |
| frequency_penalty  | float  | -2.0 - 2.0  | Penalize frequent tokens                           |
| presence_penalty   | float  | -2.0 - 2.0  | Penalize new tokens based on presence              |
| model              | string | -           | Specific model version to use                      |
| system_message     | string | -           | System message for the model                       |

### Default Model Versions

- **OpenAI**: `gpt-3.5-turbo`
- **Claude**: `claude-3-sonnet-20240229`
- **Gemini**: `gemini-pro`

### Configuration Methods

You can configure model parameters in three ways (in order of precedence):

1. **Task-specific parameters in the config file**:
   ```
   0 8 * * * claude product_manager slack-pm-channel model_params:temperature=0.8,model=claude-3-opus-20240229
   ```

2. **Environment variables**:
   ```
   MODEL_TEMPERATURE=0.7
   MODEL_MAX_TOKENS=2048
   OPENAI_MODEL=gpt-4
   CLAUDE_MODEL=claude-3-opus-20240229
   ```

3. **Command line parameters** (with the `run` command):
   ```bash
   cronai run --model openai --prompt weekly_report --processor email-team@company.com --model-params "temperature=0.5,model=gpt-4"
   ```

For more details, see [docs/model-parameters.md](docs/model-parameters.md).

## Response Processors

CronAI supports various response processors:

- **Slack**: `slack-channelname` - Send the response to a Slack channel
- **Email**: `email-address@example.com` - Send the response via email
- **Webhook**: `webhook-monitoring` - Send the response to a webhook
- **File**: `log-to-file` or `file` - Save the response to a file in the logs directory

### Response Templating

CronAI includes a powerful templating system for formatting responses. This allows you to:

- Apply consistent formatting across different output channels
- Create custom templates for different response types
- Include model metadata and execution context in the output
- Use conditional logic to customize output format based on content

The template system uses Go's `text/template` syntax and supports different template formats for each processor type. You can specify a template in your configuration:

```
# Format: timestamp model prompt response_processor [template] [variables]
0 9 1 * * claude report_template email-team@company.com monthly_report reportType=Monthly
```

In this example, `monthly_report` is the template name, which will look for the appropriate template files based on the processor type (e.g., `monthly_report_subject.tmpl`, `monthly_report_html.tmpl`, etc. for email).

For full details on the templating system, see [docs/response-templating.md](docs/response-templating.md).

## Environment Variables

Create a `.env` file with your API keys and configuration (see [.env.example](.env.example)):

```
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_claude_key
GOOGLE_API_KEY=your_gemini_key

# Response processor configs
SMTP_SERVER=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user
SMTP_PASSWORD=pass
SLACK_TOKEN=your_slack_token
WEBHOOK_URL=https://example.com/webhook
```

## Usage

```bash
# Start the service with default config file
cronai start

# Specify a custom config file
cronai start --config /path/to/config

# Run a single task immediately
cronai run --model claude --prompt product_manager --processor slack-pm-channel

# Run a task with template and variables
cronai run --model claude --prompt report_template --processor email-execs@company.com --template monthly_report --vars "reportType=Weekly,date=2025-05-11,project=CronAI"

# List all scheduled tasks
cronai list

# Manage prompts
cronai prompt list
cronai prompt search "monitoring"
cronai prompt show system/system_health
cronai prompt preview reports/monthly_report --vars "month=May,year=2025,team=Engineering"
```

## Running as a systemd Service

The application can be run as a systemd service for automatic startup and management. See [docs/systemd.md](docs/systemd.md) for detailed setup instructions.

## Development

CronAI is designed to be extended with additional models and processors. See the [CLAUDE.md](CLAUDE.md) file for development information.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.