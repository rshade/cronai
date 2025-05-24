---
id: intro
title: CronAI
sidebar_label: Introduction
slug: /
---

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
- ✅ Implemented response processors:
  - File output - Save responses to files
  - GitHub integration - Create issues and add comments
  - Console output - Display responses in terminal
- ✅ Variable substitution in prompts
- ✅ Systemd service for deployment

### Planned Features (Not Yet Implemented)

The following processors are planned but not yet implemented in the current release:

- ⚠️ Email processor - Currently logs actions only
- ⚠️ Slack processor - Currently logs actions only
- ⚠️ Webhook processor - Currently logs actions only

Additional planned features:

- Enhanced templating capabilities
- Web UI for prompt management
- Bot mode for event-driven webhook handling (stub available via `--mode bot` since v0.0.2)
- Queue mode for distributed task execution (stub available via `--mode queue` since v0.0.2)
- Model fallback mechanisms
- Advanced scheduling options

See [Limitations and Improvements](https://github.com/rshade/cronai/blob/main/docs/limitations-and-improvements.md) for a detailed breakdown of current limitations and planned improvements.

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
```

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
```

See the [Example Configuration Files](https://github.com/rshade/cronai/blob/main/cronai.config.example) in the repository for more examples.

## Usage

### Starting the Service

```bash
# Start with default config file (cron mode)
cronai start

# Start with explicit operation mode (available since v0.0.2)
cronai start --mode cron

# Future operation modes (coming soon)
cronai start --mode bot    # Event-driven webhook handler (planned)
cronai start --mode queue  # Job queue processor (planned)
```

### Operation Modes

As of v0.0.2, CronAI supports the `--mode` flag to prepare for future operation modes:

- **cron** (default): Traditional scheduled task execution using cron syntax
- **bot** (coming soon): Event-driven webhook handler for real-time responses
- **queue** (coming soon): Job queue processor for distributed task execution

The `--mode` flag establishes the CLI interface early, allowing users to prepare for future features without breaking changes.
