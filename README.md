# cronai

AI agent for scheduled prompt execution - Your automated AI assistant.

## Overview

CronAI is an intelligent agent that schedules and executes AI model prompts automatically. It acts as your personal AI automation system, running tasks on schedule and delivering results through your preferred channels.

## MVP Features

The current MVP release includes:

- ✅ Cron-style scheduling for automated execution
- ✅ Support for multiple AI models:
  - OpenAI (gpt-3.5-turbo, gpt-4)
  - Claude 4 (opus, sonnet, haiku) - Available in v0.0.2+
  - Claude 3.5 (opus, sonnet, haiku) - Available in v0.0.2+
  - Claude 3 (opus, sonnet, haiku) - Available in v0.0.2+
  - Gemini
- ✅ Customizable prompts stored as markdown files
- ✅ Response processing options:
  - File output
  - GitHub (issues and comments)
  - Microsoft Teams webhooks - Available in v0.0.2+
  - Console output
- ✅ Variable substitution in prompts
- ✅ Systemd service for deployment
- ✅ Queue mode for distributed task execution via RabbitMQ and in-memory queues - Available in v0.0.2+
- ✅ Bot mode for event-driven webhook handling - Available in v0.0.2+

### Planned Post-MVP Features (Coming Soon)

The following features are in development and will be available in future releases:

- Email processor integration
- Slack processor integration
- Generic webhook processor integration
- Enhanced templating capabilities
- Web UI for prompt management

See [limitations-and-improvements.md](docs/limitations-and-improvements.md) for a detailed breakdown of current limitations and planned improvements.

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

```text
timestamp model prompt response_processor [variables] [model_params:...]
```text

- **timestamp**: Standard cron format (minute hour day-of-month month day-of-week)
- **model**: AI model to use (openai, claude, gemini)
- **prompt**: Name of prompt file in cron_prompts directory (with or without .md extension)
- **response_processor**: How to process the response:
  - `file-path/to/output.txt`: Save to file
  - `github-issue:owner/repo`: Create GitHub issue
  - `github-comment:owner/repo#123`: Add comment to GitHub issue
  - `teams-channel`: Send to Microsoft Teams webhook (Available in v0.0.2+)
  - `console`: Display in console
- **variables** (optional): Variables to replace in the prompt file, in the format `key1=value1,key2=value2,...`
- **model_params** (optional): Model-specific parameters in the format `model_params:param1=value1,param2=value2,...`

### Example Configuration

```text
# Run daily at 8 AM using OpenAI, saving to file
0 8 * * * openai product_manager file-/var/log/cronai/product_manager.log

# Run weekly on Monday at 9 AM using Claude, creating GitHub issue
0 9 * * 1 claude weekly_report github-issue:your-org/your-repo

# Run daily health check with variables
0 6 * * * openai system_check file-/var/log/cronai/health.log system=production,check_level=detailed

# Send monitoring alerts to Microsoft Teams (v0.0.2+)
0 */4 * * * claude system_health teams-monitoring

# Run with custom model parameters (temperature and specific model version)
0 9 * * 1 openai weekly_report file-/var/log/cronai/report.log model_params:temperature=0.5,model=gpt-4
```text

See [cronai.config.example](cronai.config.example), [cronai.config.variables.example](cronai.config.variables.example), and [cronai.config.model-params.example](cronai.config.model-params.example) for more examples.

## Prompt Management

CronAI uses a file-based prompt management system to help you organize and use prompts efficiently.

### Prompt Structure

Prompts are stored as markdown files in the `cron_prompts/` directory:

```text
cron_prompts/
├── monitoring/     # Monitoring-related prompts
├── reports/        # Report generation prompts
├── system/         # System operations prompts
└── [custom]/       # Your custom categories
```text

### Variables in Prompts

You can use variables in prompt files with the syntax `{{variable_name}}`. These variables will be replaced with values from the configuration:

```markdown
# Report for {{CURRENT_DATE}}

## Overview

This is a report for {{project}} created on {{CURRENT_DATE}}.

## Details

Environment: {{environment}}
Team: {{team}}
```text

Special variables that are automatically populated:

- `{{CURRENT_DATE}}`: Current date in YYYY-MM-DD format
- `{{CURRENT_TIME}}`: Current time in HH:MM:SS format
- `{{CURRENT_DATETIME}}`: Current date and time in YYYY-MM-DD HH:MM:SS format

### Managing Prompts

CronAI includes CLI commands for prompt management:

```bash
# List all prompts
cronai prompt list

# Search for prompts
cronai prompt search "health check"

# Show prompt details
cronai prompt show system/system_health

# Preview a prompt with variables
cronai prompt preview system/system_health --vars "cpu_usage=85,memory_usage=70"
```text

## Model Parameters

You can configure model-specific parameters to fine-tune AI model behavior. Basic supported parameters:

| Parameter          | Type   | Range        | Description                                        |
|--------------------|--------|-------------|-------------------------------------------------|
| temperature        | float  | 0.0 - 1.0   | Controls response randomness (higher = more random) |
| max_tokens         | int    | > 0         | Maximum number of tokens to generate                |
| model              | string | -           | Specific model version to use                      |

### Default Model Versions

- **OpenAI**: `gpt-3.5-turbo`
- **Claude**: `claude-3-sonnet-20240229`
- **Gemini**: `gemini-pro`

### Configuration Methods

You can configure model parameters in three ways (in order of precedence):

1. **Task-specific parameters in the config file**:

   ```text
   0 8 * * * claude product_manager file-output.txt model_params:temperature=0.8,model=claude-3-opus-20240229
   ```text

2. **Environment variables**:

   ```text
   MODEL_TEMPERATURE=0.7
   MODEL_MAX_TOKENS=2048
   OPENAI_MODEL=gpt-4
   CLAUDE_MODEL=claude-3-opus-20240229
   ```text

3. **Command line parameters** (with the `run` command):

   ```bash
   cronai run --model openai --prompt weekly_report --processor file-report.txt --model-params "temperature=0.5,model=gpt-4"
   ```text

## Response Processors

CronAI currently supports these response processors in the MVP:

- **File**: `file-path/to/file.txt` - Save the response to a file
- **GitHub**: `github-issue:owner/repo` or `github-comment:owner/repo#123` - Create issues or comments
- **Microsoft Teams**: `teams-channel` - Send to Teams webhook (v0.0.2+)
- **Console**: `console` - Display the response in the console

### GitHub Processor

The GitHub processor allows you to create issues and add comments to existing issues.

#### Format

```text
github-action:owner/repo
```text

Where `action` can be:

- `issue` - Create a new issue
- `comment` - Add a comment to an existing issue (format: `comment:owner/repo#issue_number`)

#### Examples

```text
# Create a GitHub issue
0 9 * * 1 claude weekly_report github-issue:myorg/myrepo

# Add a comment to issue #123  
0 10 * * * claude issue_analysis github-comment:myorg/myrepo#123
```text

### Microsoft Teams Processor (v0.0.2+)

The Teams processor sends formatted messages to Microsoft Teams channels via webhooks.

#### Format

```text
teams-channel_identifier
```text

Where `channel_identifier` is an optional identifier for your Teams channel (e.g., `general`, `monitoring`, `alerts`).

#### Configuration

Set up your Teams webhook URL using one of these environment variables:
- `TEAMS_WEBHOOK_URL` - Primary Teams webhook URL
- `WEBHOOK_URL_TEAMS` - Alternative configuration
- `WEBHOOK_URL_<CHANNEL>` - Channel-specific URLs (e.g., `WEBHOOK_URL_MONITORING`)

#### Examples

```text
# Send daily reports to Teams
0 9 * * * claude daily_report teams-general

# Send monitoring alerts to a specific Teams channel
*/30 * * * * openai system_monitor teams-monitoring

# Send critical alerts with custom formatting
0 * * * * claude critical_check teams-alerts
```text

The Teams processor uses Microsoft's MessageCard format with:
- Themed color coding (blue for general, red for alerts)
- Structured sections with activity titles
- Facts display for metadata (model, timestamp, variables)
- Markdown support in message content
- Automatic 25KB message size validation

## Environment Variables

Create a `.env` file with your API keys and configuration:

```text
# Model API Keys (at least one is required)
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_claude_key
GOOGLE_API_KEY=your_gemini_key

# GitHub configuration
GITHUB_TOKEN=your_github_token

# Microsoft Teams webhook configuration (v0.0.2+)
TEAMS_WEBHOOK_URL=https://outlook.office.com/webhook/your_webhook_url
# Or use type-specific URLs:
WEBHOOK_URL_TEAMS=https://outlook.office.com/webhook/your_webhook_url
```text

## Usage

### Common Commands

```bash
# Start the service with default config file (cron mode)
cronai start

# Start with explicit operation mode (available since v0.0.2)
cronai start --mode cron

# Specify a custom config file
cronai start --config /path/to/config

# Queue mode for message queue integration (available since v0.0.2)
cronai start --mode queue

# Bot mode for webhook handling (available since v0.0.2)
cronai start --mode bot

# Run a single task immediately
cronai run --model openai --prompt system_health --processor file-health.log

# Run a task with variables
cronai run --model claude --prompt report_template --processor github-issue:myorg/myrepo --vars "reportType=Weekly,date=2025-05-11,project=CronAI"

# List all scheduled tasks
cronai list

# Manage prompts
cronai prompt list
cronai prompt search "monitoring"
cronai prompt show system/system_health
```

### Operation Modes

CronAI supports multiple operation modes via the `--mode` flag:

- **cron** (default): Traditional scheduled task execution using cron syntax
- **queue**: Message queue processor for distributed task execution via RabbitMQ or in-memory queues
- **bot**: Event-driven webhook handler for GitHub events

#### Queue Mode

Queue mode allows CronAI to consume tasks from external message queues. Configure with environment variables:

```bash
# Configure RabbitMQ queue
export QUEUE_TYPE=rabbitmq
export QUEUE_CONNECTION=amqp://guest:guest@localhost:5672/
export QUEUE_NAME=cronai-tasks

# Start queue mode
cronai start --mode queue
```

Send JSON messages to your queue:

```json
{
  "model": "openai",
  "prompt": "system_health",
  "processor": "console",
  "variables": {
    "environment": "production"
  }
}
```

See [docs/queue.md](docs/queue.md) for detailed queue configuration and usage.

#### Bot Mode

Bot mode runs a webhook server that listens for GitHub events and processes them with AI:

```bash
# Configure bot mode
export CRONAI_BOT_PORT=8080
export GITHUB_WEBHOOK_SECRET=your-webhook-secret
export CRONAI_DEFAULT_MODEL=openai
export CRONAI_BOT_PROCESSOR=console  # Optional: where to send AI responses

# Start bot mode
cronai start --mode bot
```

The bot will listen on the configured port for GitHub webhooks at `/webhook`. Supported events:

- Issues (opened, closed, etc.)
- Pull requests (opened, synchronize, etc.)
- Push events
- Release events

See [docs/bot-mode.md](docs/bot-mode.md) for detailed bot configuration and webhook setup.

## Running as a systemd Service

The application can be run as a systemd service for automatic startup and management. See [docs/systemd.md](docs/systemd.md) for detailed setup instructions.

## Documentation

CronAI includes comprehensive documentation in the `docs/` directory:

### User Documentation

- [Installation and Configuration](docs/systemd.md) - Systemd service setup guide
- [Model Parameters](docs/model-parameters.md) - Details on configuring AI models
- [Prompt Management](docs/prompt-management.md) - Working with prompt files

### Developer Documentation

- [Architecture Overview](docs/architecture.md) - System design and components
- [Extension Points](docs/extension-points.md) - How to extend CronAI
- [API Documentation](docs/api.md) - API endpoints (planned for future versions)
- [Limitations and Improvements](docs/limitations-and-improvements.md) - Current limitations and future roadmap

## Development

CronAI is designed to be extended with additional models and processors. See the [CONTRIBUTING.md](CONTRIBUTING.md) file for development guidelines and the [docs/](docs/) directory for technical documentation.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
