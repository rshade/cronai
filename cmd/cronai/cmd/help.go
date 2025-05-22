// Package cmd implements the command-line interface for CronAI.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Display help and examples for CronAI commands",
	Long: `Display detailed help information and usage examples for CronAI commands.
	
CronAI is an AI agent that schedules and runs AI model prompts automatically.
It supports multiple AI models and can process responses through various channels.`,
	Run: showHelp,
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

func showHelp(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		// Show general help
		fmt.Print(`CronAI - AI Agent for Scheduled Prompt Execution

CronAI is your automated AI assistant that:
- Schedules AI prompts to run on a cron schedule
- Supports multiple AI models (OpenAI, Claude, Gemini)
- Processes responses through various channels (email, Slack, webhooks, files)
- Handles dynamic variables and conditional logic in prompts
- Provides templating for response formatting

Available Commands:
  start       Start the CronAI service with scheduled tasks
  run         Execute a single AI task immediately
  list        Display all scheduled tasks
  prompt      Manage and explore prompt templates
  validate    Check template files for syntax errors
  help        Show this help information

Configuration:
  CronAI uses a configuration file (cronai.config) with entries like:
  0 8 * * * openai daily_report email-report@company.com

  Each line follows: schedule model prompt processor [variables]

Examples:
  # Start the service
  cronai start --config=./cronai.config

  # Run a single task
  cronai run --model=claude --prompt=product_review --processor=file --vars="product=Widget,category=Hardware"

  # List scheduled tasks
  cronai list

  # Search for prompts
  cronai prompt search "report"

Use "cronai help [command]" for more information about a command.
`)
		return
	}

	// Show command-specific help
	switch args[0] {
	case "start":
		fmt.Print(`Start Command - Launch CronAI Service

The start command launches the CronAI service, which reads your configuration
file and schedules AI tasks to run automatically.

Usage:
  cronai start [flags]

Flags:
  --config string   Path to configuration file (default "./cronai.config")

Configuration File Format:
  # Schedule format: minute hour day month weekday
  0 8 * * * openai morning_briefing email-team@company.com
  0 */4 * * * claude system_check slack-ops-channel
  30 9 * * 1 gemini weekly_report file reportType=Weekly,format=markdown

Examples:
  # Start with default config
  cronai start

  # Start with custom config
  cronai start --config=/etc/cronai/production.config

  # Run in background with systemd
  sudo systemctl start cronai

The service will run continuously, executing tasks according to the schedule.
Use Ctrl+C to stop the service.
`)

	case "run":
		fmt.Print(`Run Command - Execute Single AI Task

The run command executes a single AI task immediately without scheduling.
Perfect for testing prompts or running one-off tasks.

Usage:
  cronai run --model=MODEL --prompt=PROMPT --processor=PROCESSOR [flags]

Required Flags:
  --model string      AI model to use (openai, claude, gemini)
  --prompt string     Name of prompt file in cron_prompts directory
  --processor string  Response processor (email, slack, webhook, file)

Optional Flags:
  --vars string         Variables in format "key1=value1,key2=value2"
  --template string     Response template name for formatting
  --model-params string Model parameters (temperature=0.7,max_tokens=1024)

Special Variables:
  {{CURRENT_DATE}}     Replaced with current date (YYYY-MM-DD)
  {{CURRENT_TIME}}     Replaced with current time (HH:MM:SS)
  {{CURRENT_DATETIME}} Replaced with date and time

Examples:
  # Basic execution
  cronai run --model=openai --prompt=daily_summary --processor=file

  # With variables and template
  cronai run --model=claude --prompt=report_template --processor=email \
    --vars="department=Engineering,period=Q1" \
    --template=quarterly_report

  # With model parameters
  cronai run --model=gemini --prompt=creative_writing --processor=file \
    --model-params="temperature=0.9,max_tokens=2000"

  # With special variables
  cronai run --model=openai --prompt=status_check --processor=slack \
    --vars="date={{CURRENT_DATE}},system=production"
`)

	case "list":
		fmt.Print(`List Command - Display Scheduled Tasks

The list command shows all scheduled tasks from your configuration file.

Usage:
  cronai list [flags]

Flags:
  --config string   Path to configuration file (default "./cronai.config")

Output Format:
  The command displays tasks in the format:
  [index]. [schedule] [model] [prompt] [processor]

Examples:
  # List tasks from default config
  cronai list

  # List tasks from custom config
  cronai list --config=/etc/cronai/production.config

Example Output:
  Listing tasks from config: ./cronai.config
  Scheduled tasks:
  1. 0 8 * * * openai morning_briefing email-team@company.com
  2. 0 */4 * * * claude system_check slack-ops-channel
  3. 30 9 * * 1 gemini weekly_report file
`)

	case "prompt":
		fmt.Print(`Prompt Command - Manage AI Prompts

The prompt command helps you explore, search, and preview AI prompt templates.

Usage:
  cronai prompt [subcommand] [flags]

Subcommands:
  list      List all available prompts
  search    Search for prompts by name or content
  show      Display detailed information about a prompt
  preview   Preview a prompt with processed variables

Examples:
  # List all prompts
  cronai prompt list

  # List prompts by category
  cronai prompt list --category=reports

  # Search for prompts
  cronai prompt search "weekly"
  cronai prompt search --query="system" --content

  # Show prompt details
  cronai prompt show monthly_report --vars

  # Preview with variables
  cronai prompt preview template_name --vars="key1=value1,key2=value2"

Prompt Structure:
  Prompts can include:
  - Basic text content
  - Variables: {{.Variables.name}}
  - Conditionals: {{if eq .Variables.env "prod"}}...{{end}}
  - Includes: {{include "common_header.md"}}
  - Special vars: {{.CURRENT_DATE}}, {{.CURRENT_TIME}}
`)

	case "validate":
		fmt.Print(`Validate Command - Check Template Syntax

The validate command checks template files for syntax errors before use.

Usage:
  cronai validate [flags]

Flags:
  --file string, -f   Validate a specific template file
  --dir string, -d    Validate all template files in a directory

Examples:
  # Validate single file
  cronai validate --file=templates/email_report.tmpl

  # Validate directory
  cronai validate --dir=templates/

Output:
  ✅ Template templates/email_report.tmpl is valid
  ❌ Invalid template in templates/broken.tmpl: template: line 5: unexpected "}"

This command helps ensure your response templates are correctly formatted
before they're used in production.
`)

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		fmt.Println("Use 'cronai help' to see available commands")
	}
}
