package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/cron"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	Long:  `List all scheduled tasks from the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file path
		configPath := cfgFile
		if configPath == "" {
			configPath = "./cronai.config"
		}

		fmt.Printf("Listing tasks from config: %s\n", configPath)

		// List all tasks
		tasks, err := cron.ListTasks(configPath)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

		// Display tasks
		if len(tasks) == 0 {
			fmt.Println("No tasks found in configuration file")
			return
		}

		fmt.Println("Scheduled tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s %s %s %s\n", i+1, task.Schedule, task.Model, task.Prompt, task.Processor)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
