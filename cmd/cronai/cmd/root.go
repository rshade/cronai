package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cronai",
	Short: "AI agent for scheduled prompt execution",
	Long: `CronAI - Your Automated AI Assistant

CronAI is an intelligent agent that schedules and executes AI model prompts automatically.
It acts as your personal AI automation system, running tasks on schedule and delivering
results through your preferred channels.

Key Features:
  • Schedule AI prompts using cron syntax
  • Support for multiple AI models (OpenAI, Claude, Gemini)
  • Process responses through email, Slack, webhooks, or files
  • Dynamic variables and conditional logic in prompts
  • Template-based response formatting
  • Model fallback and error handling

Quick Start:
  cronai start                    # Start the service
  cronai run --help              # Run a single task
  cronai prompt list             # Explore available prompts
  cronai help                    # Show detailed help

Configuration Example:
  0 8 * * * openai daily_summary email-team@company.com

For more information, use 'cronai help [command]' or visit the documentation.
`,
	Example: `  # Start the service with default config
  cronai start

  # Run a single task with variables
  cronai run --model=claude --prompt=report --processor=email --vars="type=weekly"

  # List all scheduled tasks
  cronai list

  # Search for available prompts
  cronai prompt search "monitoring"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./cronai.config)")
}

func initConfig() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found or error loading it")
	}

	// Config file will be used in the individual commands
	// The actual config handling is implemented in pkg/config
}
