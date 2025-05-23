# CronAI Template Examples

This directory contains example templates for CronAI response processors. For comprehensive documentation on the template system, see [Response Processor Templating](../docs/response-templating.md).

The `library/` subdirectory contains reusable template components that can be used with template inheritance and composition patterns.

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
- `{{.Metadata}}` - Additional metadata about the execution

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

```text
# Format: timestamp model prompt response_processor [template] [variables]
0 9 1 * * claude report_template email-team@company.com monthly_report reportType=Monthly,date={{CURRENT_DATE}}
```text

You can also validate templates using the `validate` command:

```bash
# Validate a single template
./cronai validate --file templates/monthly_report_html.tmpl

# Validate all templates in a directory
./cronai validate --dir templates/
```text

## Template Inheritance and Composition

CronAI supports advanced template inheritance and composition patterns:

```text
# Template inheritance
{{extends "base_template"}}
{{block "section" .}}Override content{{end}}

# Template inclusion
{{include "component_template"}}

# Template embedding
{{template "partial_template" .}}
```text

For complete documentation on these features, see [Template Inheritance and Composition](../docs/template-inheritance.md).

## Template Functions

The template system provides many utility functions:

- Variable handling: `hasVar`, `getVar`
- String operations: `eq`, `ne`, `contains`, `hasPrefix`, `hasSuffix`, `upper`, `lower`, `title`, `trim`
- Numeric comparison: `lt`, `le`, `gt`, `ge`
- Boolean operations: `not`, `and`, `or`
- Date operations: `now`, `formatDate`, `addDays`
- JSON utilities: `marshalJSON`

See the full documentation for details and examples.
