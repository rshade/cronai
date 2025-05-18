# CronAI Prompt Library

This directory contains organized prompt files for use with CronAI.

## Directory Structure

- `monitoring/` - Prompts related to system and service monitoring
- `system/` - Prompts for system health checks and diagnostics
- `reports/` - Prompts for generating various reports
- `templates/` - Reusable prompt templates with variables and includes

## Prompt File Format

Each prompt file should include a metadata header in the following format:

```markdown
---
name: Prompt Name
description: Brief description of what this prompt does
author: Author Name
version: 1.0
category: monitoring|system|report|template
tags: tag1, tag2, tag3
variables:
  - name: variable1
    description: Description of variable1
  - name: variable2 
    description: Description of variable2
---

# Main Prompt Content

The actual prompt content starts here...
```text

## Usage

Prompts can be referenced in your cronai configuration using the relative path from this directory:

```text
0 8 * * * claude monitoring/monitoring_check slack-alerts
```text

or for templates with variables:

```text
0 9 * * * claude templates/report_template email-team@example.com reportType=Weekly,date={{CURRENT_DATE}}
```text

## Includes and Composition

Prompts can include content from other prompt files using the include syntax:

```markdown
{{include "templates/common_header.md"}}

Your prompt content here...

{{include "templates/common_footer.md"}}
```text
