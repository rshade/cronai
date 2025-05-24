# cronai

AI agent for scheduled prompt execution - Your automated AI assistant.

## Overview

CronAI is an intelligent agent that schedules and executes AI model prompts automatically. It acts as your personal AI automation system, running tasks on schedule and delivering results through your preferred channels.

## MVP Features

The current MVP release includes:

- ✅ Cron-style scheduling for automated execution
- ✅ Support for multiple AI models:
  - OpenAI (gpt-3.5-turbo, gpt-4)
  - Claude (claude-3-sonnet, claude-3-opus)
  - Gemini
- ✅ Customizable prompts stored as markdown files
- ✅ Response processing options:
  - File output
  - GitHub (issues and comments)
  - Console output
- ✅ Variable substitution in prompts
- ✅ Systemd service for deployment

### Planned Post-MVP Features (Coming Soon)

The following features are in development and will be available in future releases:

- Email processor integration
- Slack processor integration
- Webhook processor integration
- Enhanced templating capabilities
- Web UI for prompt management
- Bot mode for event-driven webhook handling (stub available via `--mode bot`)
- Queue mode for distributed task execution (stub available via `--mode queue`)

See [limitations-and-improvements.md](docs/limitations-and-improvements.md) for a detailed breakdown of current limitations and planned improvements.

## Installation

```bash
# Install directly
go install github.com/rshade/cronai/cmd/cronai@latest

# Or clone and build
git clone https://github.com/rshade/cronai.git
cd cronai
go build -o cronai ./cmd/cronai
```text

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

## Environment Variables

Create a `.env` file with your API keys and configuration:

```text
# Model API Keys (at least one is required)
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_claude_key
GOOGLE_API_KEY=your_gemini_key

# GitHub configuration
GITHUB_TOKEN=your_github_token
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

# Future operation modes (coming soon)
cronai start --mode bot    # Event-driven webhook handler (planned)
cronai start --mode queue  # Job queue processor (planned)

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

As of v0.0.2, CronAI supports the `--mode` flag to prepare for future operation modes:

- **cron** (default): Traditional scheduled task execution using cron syntax
- **bot** (coming soon): Event-driven webhook handler for real-time responses
- **queue** (coming soon): Job queue processor for distributed task execution

The `--mode` flag establishes the CLI interface early, allowing users to prepare for future features without breaking changes.text

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
