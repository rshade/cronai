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
