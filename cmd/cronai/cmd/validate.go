package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rshade/cronai/internal/processor/template"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate response template files",
	Long:  `Validate response template files for syntax errors. This command helps identify template syntax errors early before they cause issues in response processing.`,
	Example: `  # Validate a single template file
  cronai validate --file=templates/email_report.tmpl

  # Validate all templates in a directory
  cronai validate --dir=templates/`,
	Run: func(cmd *cobra.Command, _ []string) {
		fileFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("Error getting file flag: %v\n", err)
			return
		}

		dirFlag, err := cmd.Flags().GetString("dir")
		if err != nil {
			fmt.Printf("Error getting dir flag: %v\n", err)
			return
		}

		if fileFlag == "" && dirFlag == "" {
			fmt.Println("Please specify either --file or --dir")
			return
		}

		if fileFlag != "" {
			validateTemplateFile(fileFlag)
		} else if dirFlag != "" {
			validateTemplateDir(dirFlag)
		}
	},
}

// validateTemplateFile validates a single template file
func validateTemplateFile(filePath string) {
	// Only validate .tmpl files
	if !strings.HasSuffix(filePath, ".tmpl") {
		fmt.Printf("Not a template file (must end with .tmpl): %s\n", filePath)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	manager := template.GetManager()
	// Use filename as template name for validation
	templateName := filepath.Base(filePath)

	err = manager.Validate(templateName, string(content))
	if err != nil {
		fmt.Printf("❌ Invalid template in %s: %v\n", filePath, err)
	} else {
		fmt.Printf("✅ Template %s is valid\n", filePath)
	}
}

// validateTemplateDir validates all template files in a directory
func validateTemplateDir(dirPath string) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.tmpl"))
	if err != nil {
		fmt.Printf("Error finding template files: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No template files found in directory")
		return
	}

	allValid := true
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			allValid = false
			continue
		}

		manager := template.GetManager()
		// Use filename as template name for validation
		templateName := filepath.Base(file)

		err = manager.Validate(templateName, string(content))
		if err != nil {
			fmt.Printf("❌ Invalid template in %s: %v\n", file, err)
			allValid = false
		} else {
			fmt.Printf("✅ Template %s is valid\n", file)
		}
	}

	if allValid {
		fmt.Println("\nAll templates are valid")
	} else {
		fmt.Println("\nSome templates have errors")
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().String("file", "", "Path to a template file to validate")
	validateCmd.Flags().String("dir", "", "Path to a directory of templates to validate")
}
