package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/rshade/cronai/internal/prompt"
)

func main() {
	// Define flags
	queryFlag := flag.String("query", "", "Search query")
	categoryFlag := flag.String("category", "", "Filter by category")
	contentFlag := flag.Bool("content", false, "Search in prompt content")
	flag.Parse()

	// If no query provided, try to get it from args
	query := *queryFlag
	if query == "" && flag.NArg() > 0 {
		query = strings.Join(flag.Args(), " ")
	}

	// If still no query, show usage
	if query == "" && *categoryFlag == "" {
		fmt.Println("Usage: search_prompt [options] [query]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Search for prompts
	var prompts []prompt.PromptInfo
	var err error

	if *contentFlag {
		// Search in prompt content
		prompts, err = prompt.SearchPromptContent(query, *categoryFlag)
	} else {
		// Search in metadata only
		prompts, err = prompt.SearchPrompts(query, *categoryFlag)
	}

	if err != nil {
		fmt.Printf("Error searching prompts: %v\n", err)
		os.Exit(1)
	}

	// Display results
	if len(prompts) == 0 {
		fmt.Println("No prompts found matching the search criteria")
		os.Exit(0)
	}

	// Format output using tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "CATEGORY\tNAME\tDESCRIPTION\tPATH"); err != nil {
		fmt.Printf("Error writing to tabwriter: %v\n", err)
		os.Exit(1)
	}
	for _, p := range prompts {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Category, p.Name, p.Description, p.Path); err != nil {
			fmt.Printf("Error writing to tabwriter: %v\n", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Printf("Error flushing tabwriter: %v\n", err)
		os.Exit(1)
	}
}
