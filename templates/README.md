# CronAI Template Examples

This directory contains example templates for CronAI response processors.

## Template Files

Templates use Go's standard `text/template` package syntax. Within a template, you have access to the following variables:

- `{{.Content}}` - The content of the AI model response
- `{{.Model}}` - The name of the model (openai, claude, gemini)
- `{{.PromptName}}` - The name of the prompt used
- `{{.Timestamp}}` - The time when the response was generated
  - Use format function: `{{.Timestamp.Format "2006-01-02 15:04:05"}}`
- `{{.Variables}}` - Map of variables used in the prompt
  - Access specific variable: `{{.Variables.key_name}}`
- `{{.ExecutionID}}` - Unique identifier for this execution

## Naming Conventions

Templates should follow these naming conventions:

- Email templates:
  - `name_subject.tmpl`: Email subject line
  - `name_html.tmpl`: HTML email body
  - `name_text.tmpl`: Plain text email body
- Slack templates:
  - `name.tmpl`: Slack Block Kit JSON payload
- Webhook templates:
  - `name.tmpl`: Webhook JSON payload
- File templates:
  - `name_filename.tmpl`: Template for the output filename
  - `name_content.tmpl`: Template for the file content

## Using Templates

To use a template, specify its name (without the `.tmpl` extension) in your cronai configuration:

```
# Format: timestamp model prompt response_processor [template] [variables]
0 9 1 * * claude report_template email-team@company.com monthly_report reportType=Monthly,date={{CURRENT_DATE}}
```

You can also validate templates using the `validate` command:

```bash
# Validate a single template
./cronai validate --file templates/monthly_report_html.tmpl

# Validate all templates in a directory
./cronai validate --dir templates/
```