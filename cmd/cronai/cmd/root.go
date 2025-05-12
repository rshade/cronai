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
	Short: "Run AI model prompts on a cron schedule",
	Long: `CronAI is a command-line utility to run AI model prompts on a cron-type schedule.

It supports multiple AI models (OpenAI, Claude, Gemini) and various response
processors (email, Slack, webhooks, file output).
`,
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

	// If a config file is provided, use it
	if cfgFile != "" {
		// Config handling will be implemented in pkg/config
	}
}
