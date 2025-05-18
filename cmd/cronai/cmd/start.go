package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/cron"
	"github.com/spf13/cobra"
)

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
	Example: `  # Start with default config file
  cronai start

  # Start with custom config
  cronai start --config=/etc/cronai/production.config

  # Run in background (systemd)
  sudo systemctl start cronai`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file path
		configPath := cfgFile
		if configPath == "" {
			configPath = "./cronai.config"
		}

		fmt.Printf("Starting CronAI service with config: %s\n", configPath)

		// Start the cron service
		if err := cron.StartService(configPath); err != nil {
			fmt.Printf("Error starting cron service: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
