---
id: prompt-management
title: Prompt Management
sidebar_label: Prompt Management
---

CronAI includes a simple file-based prompt management system that helps you organize and use prompts efficiently.

## Directory Structure

Prompts are stored as markdown files in the `cron_prompts/` directory:

```text
cron_prompts/
├── README.md           # Documentation for prompt structure
├── monitoring/         # Prompts for monitoring purposes
├── reports/            # Prompts for report generation
├── system/             # Prompts for system operations
└── [other_categories]/ # Custom categories
```

## Prompt Files

Prompts are standard markdown files with a `.md` extension. During the MVP phase, prompts are simple text files that can contain variables:

```markdown
# System Health Check

Analyze the following system metrics and provide recommendations:

- CPU Usage: {{cpu_usage}}%
- Memory Usage: {{memory_usage}}%
- Disk Usage: {{disk_usage}}%

Please provide:
1. Assessment of current system health
2. Potential issues identified
3. Recommended actions
```

## Variables in Prompts

Variables in prompts use the `{{variable_name}}` syntax. CronAI automatically provides the following special variables:

- `{{CURRENT_DATE}}`: Current date in YYYY-MM-DD format
- `{{CURRENT_TIME}}`: Current time in HH:MM:SS format
- `{{CURRENT_DATETIME}}`: Current date and time in YYYY-MM-DD HH:MM:SS format

Custom variables can be provided in the configuration file or command line using a comma-separated list of key=value pairs.

## CLI Commands

CronAI provides several commands to help you manage your prompts:

### List Prompts

List all available prompts, optionally filtered by category:

```bash
# List all prompts
cronai prompt list

# List prompts in a specific category
cronai prompt list --category system
```

### Search Prompts

Search for prompts by name or content:

```bash
# Search by name or description
cronai prompt search "health check"

# Search in a specific category
cronai prompt search --query "monitoring" --category system

# Search in prompt content
cronai prompt search --content --query "CPU usage"
```

### Show Prompt Details

Show detailed information about a specific prompt:

```bash
cronai prompt show system/system_health
```

### Preview Prompt

Preview a prompt with variables substituted:

```bash
cronai prompt preview system/system_health --vars "cpu_usage=85,memory_usage=70,disk_usage=50"
```

## Using Prompts in CronAI Configuration

Reference prompts in your cronai.config file using either the full path or category/name format:

```text
# Using a prompt from a category
0 8 * * * openai system/system_health file-/var/log/cronai/health.log

# Using a prompt with variables
0 9 * * 1 claude reports/weekly_report github-issue:owner/repo date={{CURRENT_DATE}},team=Engineering
```

## Example Prompt Files

### Basic Prompt

```markdown
# Daily Status Report

Generate a daily status report for {{project}} on {{CURRENT_DATE}}.

1. Current status overview
2. Progress since yesterday
3. Planned tasks for today
4. Blocking issues, if any
```

### Prompt with System Metrics

```markdown
# System Health Check

CPU Usage: {{cpu_usage}}%
Memory Usage: {{memory_usage}}%
Disk Space: {{disk_usage}}%

Please analyze these metrics and provide:
1. Current system health assessment
2. Potential issues or warning signs
3. Recommended actions
```

## Post-MVP Features

The following prompt management features are planned for future releases:

- Prompt metadata (YAML frontmatter)
- Template inheritance and composition
- Includes for reusing common prompt components
- Conditional logic in prompts
- Advanced variable validation and defaults

For more information on these upcoming features, see the project roadmap.
