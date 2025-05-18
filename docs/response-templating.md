# Response Processor Templating (MVP)

## MVP Status: In Development

> ⚠️ **Note**: The full response templating system described in this document is currently in development and will be available in a future update. For the MVP release, CronAI provides basic output formatting without custom templating. This document describes how templating will work in the upcoming releases.

## Current MVP Capabilities

In the current MVP release, CronAI provides:

- Basic output formatting for file processor (raw output)
- Basic GitHub issue and comment creation with default formatting
- Console output with simple formatting

## Future Templating System Overview

The future templating system will be built on Go's `text/template` package and will allow for:

- Consistent formatting across different output channels
- Custom templates defined by users
- Layout customization for different processor types
- Conditional content based on response attributes
- Template inheritance and composition

### Planned Template Types

Each processor type will support specific template formats:

- **Email**: Subject, HTML content, and plain text fallback templates
- **Slack**: Message content, blocks, and attachments templates
- **Webhook**: JSON payload templates
- **File**: Content templates with optional metadata
- **GitHub**: Issue title, issue body, and comment templates

### Template Variables

Templates will have access to various data attributes:

```text
{{ .Content }}     - The AI model's response text
{{ .Model }}       - The model name (e.g., "openai", "claude")
{{ .PromptName }}  - Name of the prompt that was used
{{ .Timestamp }}   - When the response was generated
{{ .ExecutionID }} - Unique execution identifier
{{ .Variables }}   - Map of variables used in the prompt
```text

### Template Location

Templates will be stored in the `templates/` directory:

```text
templates/
├── email/
│   ├── default_subject.tmpl
│   ├── default_html.tmpl
│   └── default_text.tmpl
├── slack/
│   └── default.tmpl
├── webhook/
│   └── default.tmpl
├── file/
│   └── default.tmpl
├── github/
│   ├── default_issue.tmpl
│   └── default_comment.tmpl
└── library/
    ├── header.tmpl
    ├── footer.tmpl
    └── common.tmpl
```text

## Future Implementation Plan

The full templating system will include:

1. Default templates for all processor types
2. Custom template registration and management
3. Template inheritance for component reuse
4. Conditional logic in templates
5. Helper functions for common formatting tasks
6. Multi-format output for email and Slack

## Current MVP Usage

For the MVP release, processors use the following default behavior:

### File Processor

- Writes the raw AI response to the specified file
- No templating is applied in the MVP

### GitHub Processor

The GitHub processor uses built-in JSON templates for formatting:

**For Issues (`github-issue:owner/repo`):**

- Title format: `[PromptName] - [Date]`
- Body includes: Model info, timestamp, execution ID, variables (if provided), and formatted content
- Automatically adds labels: `["auto-generated", "cronai"]`
- Example format:

  ```text
  github-issue:myorg/myrepo
  ```text

**For Comments (`github-comment:owner/repo#issue_number`):**

- Adds a formatted comment with model info and AI response
- Includes metadata: Model name, timestamp, and prompt name
- Example format:

  ```text
  github-comment:myorg/myrepo#42
  ```text

**For Pull Requests (`github-pr:owner/repo`):**

- Similar formatting to issues
- Requires `head_branch` variable
- Optional `base_branch` variable (defaults to "main")
- Note: In MVP, PR creation is logged rather than executed
- Example with required variables:

  ```text
  github-pr:myorg/myrepo head_branch=feature/auto-update

### Console Processor

- Displays the AI response with minimal formatting
- No custom templating is available in the MVP

## Example (Coming in Future Release)

Example of how templating will work in the future release:

```text
# Use custom template for file output
0 9 * * 1 openai weekly_report file-/var/log/cronai/report.md report_template

# Use custom template for GitHub issue
0 8 * * * claude system_health github-issue:owner/repo system_alert
```text

Stay tuned for updates as we implement the full templating system in upcoming releases.
