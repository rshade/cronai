# Prompt Management

CronAI includes a robust file-based prompt management system that helps you organize, discover, and reuse prompts efficiently.

## Directory Structure

Prompts are organized in the following directory structure:

```
cron_prompts/
├── README.md           # Documentation for prompt structure
├── monitoring/         # Prompts for monitoring purposes
├── reports/            # Prompts for report generation
├── system/             # Prompts for system operations
├── templates/          # Reusable prompt templates
└── [other_categories]/ # Custom categories
```

## Prompt Metadata

Each prompt file can include an optional YAML metadata section at the beginning of the file, enclosed by triple dashes (`---`):

```markdown
---
name: System Health Check
description: Analyzes system health metrics and provides recommendations
author: CronAI Team
version: 1.0
category: system
tags: health, monitoring, metrics, analysis
extends: templates/base_system_check
variables:
  - name: cpu_usage
    description: Current CPU usage percentage
  - name: memory_usage
    description: Current memory usage percentage
---

# Actual Prompt Content

The content of your prompt goes here...
```

The `extends` field is optional and allows the prompt to inherit from another template. See [Template Inheritance and Composition](template-inheritance.md) for more details.

## Prompt Composition

CronAI supports prompt composition through includes. You can include content from other prompt files using the `{{include}}` directive:

```markdown
{{include "templates/common_header.md"}}

# Main Content

Your specific prompt content goes here.

{{include "templates/common_footer.md"}}
```

For more advanced template reuse patterns, CronAI now also supports template inheritance with the `{{extends}}` and `{{block}}` directives. See [Template Inheritance and Composition](template-inheritance.md) for complete documentation.

## CLI Commands

CronAI provides several commands to help you manage your prompts:

### List Prompts

List all available prompts, optionally filtered by category:

```bash
cronai prompt list
cronai prompt list --category system
```

### Search Prompts

Search for prompts by name, description, or content:

```bash
cronai prompt search "health check"
cronai prompt search --query "monitoring" --category system
cronai prompt search --content --query "CPU usage"
```

### Show Prompt Details

Show detailed information about a specific prompt:

```bash
cronai prompt show system/system_health
cronai prompt show --vars system/system_health
```

### Preview Prompt

Preview a prompt with variables substituted:

```bash
cronai prompt preview system/system_health --vars "cpu_usage=85,memory_usage=70"
```

## Using Prompts in CronAI Configuration

Reference prompts in your cronai.config file using either the full path or category/name format:

```
# Using a prompt from a category
0 8 * * * claude system/system_health slack-alerts

# Using a prompt with variables
0 9 * * 1 claude reports/weekly_report email-team@company.com date={{CURRENT_DATE}},team=Engineering
```

## Variables in Prompts

Variables in prompts use the `{{variable_name}}` syntax. CronAI automatically provides the following special variables:

- `{{CURRENT_DATE}}`: Current date in YYYY-MM-DD format
- `{{CURRENT_TIME}}`: Current time in HH:MM:SS format
- `{{CURRENT_DATETIME}}`: Current date and time in YYYY-MM-DD HH:MM:SS format

Custom variables can be defined in the prompt metadata and provided in the configuration file or command line using a comma-separated list of key=value pairs.