package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/cron"
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

  # Future modes (not yet implemented)
  cronai start --mode bot    # Event-driven webhook handler
  cronai start --mode queue  # Job queue processor

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

		// Start the cron service
		if err := cron.StartService(configPath); err != nil {
			fmt.Printf("Error starting cron service: %v\n", err)
		}
	},
}

// validateMode validates the operation mode flag
func validateMode(mode string) error {
	switch mode {
	case "cron":
		return nil
	case "bot", "queue":
		return fmt.Errorf("mode '%s' is not yet implemented (coming in future releases)", mode)
	default:
		return fmt.Errorf("invalid mode '%s': must be one of: cron, bot, queue", mode)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVar(&operationMode, "mode", "cron",
		"Operation mode: cron (default), bot (future), queue (future)\n"+
			"Available modes:\n"+
			"  cron  - Traditional scheduled task execution (default)\n"+
			"  bot   - Event-driven webhook handler (coming soon)\n"+
			"  queue - Job queue processor (coming soon)")
}
