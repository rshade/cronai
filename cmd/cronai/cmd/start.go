package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/bot"
	"github.com/rshade/cronai/internal/cron"
	"github.com/rshade/cronai/internal/queue"
	"github.com/rshade/cronai/internal/queue/consumers"
	"github.com/spf13/cobra"
)

var operationMode string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CronAI service",
	Long: `Start the CronAI agent service with scheduled AI tasks.

The agent will read your configuration file and continuously execute AI prompts
according to their schedules. It runs as a daemon process until stopped.

Configuration Format:
  schedule model prompt processor [variables]

Examples:
  0 8 * * * openai daily_summary email-report@company.com
  */30 * * * * claude system_monitor slack-ops-channel
  0 9 * * 1 gemini weekly_report file type=weekly,format=pdf`,
	Example: `  # Start with default config file (cron mode)
  cronai start

  # Start with explicit cron mode
  cronai start --mode cron

  # Start with custom config
  cronai start --config=/etc/cronai/production.config

  # Start in bot mode (webhook server)
  cronai start --mode bot

  # Start in queue mode
  cronai start --mode queue

  # Run in background (systemd)
  sudo systemctl start cronai`,
	Run: func(_ *cobra.Command, _ []string) {
		// Validate operation mode
		if err := validateMode(operationMode); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get config file path
		configPath := cfgFile
		if configPath == "" {
			configPath = "./cronai.config"
		}

		fmt.Printf("Starting CronAI service in %s mode with config: %s\n", operationMode, configPath)

		// Start the appropriate service based on mode
		switch operationMode {
		case "cron":
			if err := cron.StartService(configPath); err != nil {
				fmt.Printf("Error starting cron service: %v\n", err)
			}
		case "bot":
			if err := bot.StartService(configPath); err != nil {
				fmt.Printf("Error starting bot service: %v\n", err)
			}
		case "queue":
			// Register all consumer types
			if err := consumers.RegisterAll(); err != nil {
				fmt.Printf("Error registering consumers: %v\n", err)
				return
			}
			if err := queue.StartService(configPath); err != nil {
				fmt.Printf("Error starting queue service: %v\n", err)
			}
		default:
			fmt.Printf("Error: mode '%s' is not implemented\n", operationMode)
		}
	},
}

// validateMode validates the operation mode flag
func validateMode(mode string) error {
	switch mode {
	case "cron", "bot", "queue":
		return nil
	default:
		return fmt.Errorf("invalid mode '%s': must be one of: cron, bot, queue", mode)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&operationMode, "mode", "cron",
		"Operation mode: cron (default), bot, queue\n"+
			"Available modes:\n"+
			"  cron  - Traditional scheduled task execution (default)\n"+
			"  bot   - Event-driven webhook handler for GitHub events\n"+
			"  queue - Job queue processor for external message queues")
}
