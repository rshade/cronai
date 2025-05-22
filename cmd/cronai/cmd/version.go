// Package cmd implements the command line interface for CronAI.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version holds the version information that can be set during build time
var Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CronAI",
	Long:  `Print the version number of CronAI and exit.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("CronAI version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
