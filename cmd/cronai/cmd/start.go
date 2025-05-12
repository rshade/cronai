package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/cron"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the CronAI service",
	Long:  `Start the CronAI service with the specified configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file path
		configPath := cfgFile
		if configPath == "" {
			configPath = "./cronai.config"
		}

		fmt.Printf("Starting CronAI service with config: %s\n", configPath)

		// Start the cron service
		cron.StartService(configPath)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
