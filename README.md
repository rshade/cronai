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
timestamp model prompt response_processor [variables]
```

- **timestamp**: Standard cron format (minute hour day-of-month month day-of-week)
- **model**: AI model to use (openai, claude, gemini)
- **prompt**: Name of prompt file in cron_prompts directory (with or without .md extension)
- **response_processor**: How to process the response (email, slack, webhook, file)
- **variables** (optional): Variables to replace in the prompt file, in the format `key1=value1,key2=value2,...`

### Example Configuration

```
# Run daily at 8 AM using Claude, sending to slack
0 8 * * * claude product_manager slack-pm-channel

# Run every Monday at 9 AM using OpenAI, sending to email
0 9 * * 1 openai weekly_report email-team@company.com

# Run monthly report with variables on the 1st of each month
0 9 1 * * claude report_template email-execs@company.com reportType=Monthly,date={{CURRENT_DATE}},project=CronAI
```

See [cronai.config.example](cronai.config.example) and [cronai.config.variables.example](cronai.config.variables.example) for more examples.

## Prompt Files

Store your prompt files as markdown in the `cron_prompts` directory. Example:

```markdown
# Product Manager Daily Task List

As a product manager, please generate a prioritized task list for today.

Include the following:
1. Top 3 urgent items to address
2. Customer feedback that needs immediate attention
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

## Response Processors

CronAI supports various response processors:

- **Slack**: `slack-channelname` - Send the response to a Slack channel
- **Email**: `email-address@example.com` - Send the response via email
- **Webhook**: `webhook-monitoring` - Send the response to a webhook
- **File**: `log-to-file` - Save the response to a file in the logs directory

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

# Run a task with variables
cronai run --model claude --prompt report_template --processor email-execs@company.com --vars "reportType=Weekly,date=2025-05-11,project=CronAI"

# List all scheduled tasks
cronai list
```

## Running as a systemd Service

The application can be run as a systemd service for automatic startup and management. See [docs/systemd.md](docs/systemd.md) for detailed setup instructions.

## Development

CronAI is designed to be extended with additional models and processors. See the [CLAUDE.md](CLAUDE.md) file for development information.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.