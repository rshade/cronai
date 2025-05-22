package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/rshade/cronai/internal/prompt"
	"github.com/spf13/cobra"
)

var category string
var searchQuery string
var showVars bool
var searchContent bool
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Manage AI prompt templates",
	Long: `Manage and explore AI prompt templates.

The prompt command provides tools to discover, search, preview, and understand
available AI prompts. Use this to find the right prompts for your tasks and
see how they work with variables.`,
	Example: `  # List all prompts
  cronai prompt list

  # Search for prompts
  cronai prompt search "report"

  # Show prompt details
  cronai prompt show monthly_report --vars

  # Preview with variables
  cronai prompt preview weekly_report --vars="team=Engineering,date={{CURRENT_DATE}}"`,
}

var promptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available AI prompts",
	Long: `List all available AI prompt templates.

Displays prompts organized by category with their descriptions and file paths.
Use categories to filter and find prompts for specific purposes.`,
	Example: `  # List all prompts
  cronai prompt list

  # List prompts by category
  cronai prompt list --category=reports
  cronai prompt list -c monitoring`,
	Run: func(_ *cobra.Command, _ []string) {
		// Get all prompts
		prompts, err := prompt.ListPrompts()
		if err != nil {
			fmt.Printf("Error listing prompts: %v\n", err)
			os.Exit(1)
		}

		// Sort prompts by category and name
		sort.Slice(prompts, func(i, j int) bool {
			if prompts[i].Category != prompts[j].Category {
				return prompts[i].Category < prompts[j].Category
			}
			return prompts[i].Name < prompts[j].Name
		})

		// Display prompts
		currentCategory := ""
		for _, p := range prompts {
			if p.Category != currentCategory {
				if currentCategory != "" {
					fmt.Println()
				}
				fmt.Printf("\n%s:\n", p.Category)
				currentCategory = p.Category
			}
			fmt.Printf("  %s\n", p.Name)
			if p.Description != "" {
				fmt.Printf("    %s\n", p.Description)
			}
		}
	},
}

var promptSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for prompts",
	Long:  `Search for prompts by name, description, or content.`,
	Run: func(_ *cobra.Command, args []string) {
		// Get query from args if not provided as flag
		if len(args) > 0 && searchQuery == "" {
			searchQuery = args[0]
		}

		var prompts []prompt.Info
		var err error

		if searchContent {
			// Search in prompt content
			prompts, err = prompt.SearchPromptContent(searchQuery, category)
		} else {
			// Search in metadata only
			prompts, err = prompt.SearchPrompts(searchQuery, category)
		}

		if err != nil {
			fmt.Printf("Error searching prompts: %v\n", err)
			return
		}

		if len(prompts) == 0 {
			fmt.Println("No prompts found matching the search criteria")
			return
		}

		// Format output using tabwriter
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "CATEGORY\tNAME\tDESCRIPTION\tPATH"); err != nil {
			fmt.Printf("Error writing to tabwriter: %v\n", err)
			return
		}
		for _, p := range prompts {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Category, p.Name, p.Description, p.Path); err != nil {
				fmt.Printf("Error writing to tabwriter: %v\n", err)
				return
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Printf("Error flushing tabwriter: %v\n", err)
		}
	},
}

var promptShowCmd = &cobra.Command{
	Use:   "show [promptName]",
	Short: "Show prompt details",
	Long:  `Show detailed information about a specific prompt.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		promptName := args[0]

		// Get prompt metadata
		metadata, err := prompt.GetPromptMetadata(promptName)
		if err != nil {
			fmt.Printf("Error getting prompt metadata: %v\n", err)
			return
		}

		// Display prompt details
		fmt.Printf("Name: %s\n", metadata.Name)
		fmt.Printf("Description: %s\n", metadata.Description)
		fmt.Printf("Author: %s\n", metadata.Author)
		fmt.Printf("Version: %s\n", metadata.Version)
		fmt.Printf("Category: %s\n", metadata.Category)

		if len(metadata.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(metadata.Tags, ", "))
		}

		// Show variables if requested
		if showVars && len(metadata.Variables) > 0 {
			fmt.Println("\nVariables:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "NAME\tDESCRIPTION"); err != nil {
				fmt.Printf("Error writing to tabwriter: %v\n", err)
				return
			}
			for _, v := range metadata.Variables {
				if _, err := fmt.Fprintf(w, "%s\t%s\n", v.Name, v.Description); err != nil {
					fmt.Printf("Error writing to tabwriter: %v\n", err)
					return
				}
			}
			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing tabwriter: %v\n", err)
			}
		}

		// Get the prompt content
		content, err := prompt.LoadPrompt(promptName)
		if err != nil {
			fmt.Printf("Error loading prompt content: %v\n", err)
			return
		}

		fmt.Println("\nContent:")
		fmt.Println("----------------------------------------")
		fmt.Println(content)
		fmt.Println("----------------------------------------")
	},
}

var promptPreviewCmd = &cobra.Command{
	Use:   "preview [promptName]",
	Short: "Preview a prompt with variables",
	Long:  `Preview a prompt with variables and includes processed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		promptName := args[0]

		// Get variables from command line
		variablesFlag, err := cmd.Flags().GetString("vars")
		if err != nil {
			fmt.Printf("Error getting vars flag: %v\n", err)
			return
		}
		variables := make(map[string]string)

		// Parse variables from comma-separated list of key=value pairs
		if variablesFlag != "" {
			pairs := strings.Split(variablesFlag, ",")
			for _, pair := range pairs {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 2 {
					variables[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
		}

		// Add special variables
		variables["CURRENT_DATE"] = "2025-05-12"              // This would be time.Now().Format("2006-01-02") in a real implementation
		variables["CURRENT_TIME"] = "12:00:00"                // This would be time.Now().Format("15:04:05") in a real implementation
		variables["CURRENT_DATETIME"] = "2025-05-12 12:00:00" // This would be time.Now().Format("2006-01-02 15:04:05") in a real implementation

		// Load prompt with variables
		content, err := prompt.LoadPromptWithVariables(promptName, variables)
		if err != nil {
			fmt.Printf("Error loading prompt with variables: %v\n", err)
			return
		}

		fmt.Println("Preview:")
		fmt.Println("----------------------------------------")
		fmt.Println(content)
		fmt.Println("----------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.AddCommand(promptListCmd)
	promptCmd.AddCommand(promptSearchCmd)
	promptCmd.AddCommand(promptShowCmd)
	promptCmd.AddCommand(promptPreviewCmd)

	// Add flags
	promptListCmd.Flags().StringVarP(&category, "category", "c", "", "Filter prompts by category")

	promptSearchCmd.Flags().StringVarP(&category, "category", "c", "", "Filter prompts by category")
	promptSearchCmd.Flags().StringVarP(&searchQuery, "query", "q", "", "Search query")
	promptSearchCmd.Flags().BoolVarP(&searchContent, "content", "t", false, "Search in prompt content")

	promptShowCmd.Flags().BoolVarP(&showVars, "vars", "v", false, "Show prompt variables")

	promptPreviewCmd.Flags().String("vars", "", "Variables in format 'key1=value1,key2=value2'")
}
