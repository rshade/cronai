package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/spf13/cobra"
)

var (
	modelName     string
	promptName    string
	processorName string
	varsString    string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a single task immediately",
	Long:  `Run a single task immediately without scheduling.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse variables if provided
		variables := make(map[string]string)
		if varsString != "" {
			for _, varPair := range strings.Split(varsString, ",") {
				keyValue := strings.SplitN(varPair, "=", 2)
				if len(keyValue) == 2 {
					key := strings.TrimSpace(keyValue[0])
					value := strings.TrimSpace(keyValue[1])

					// Handle special variables
					switch value {
					case "{{CURRENT_DATE}}":
						value = time.Now().Format("2006-01-02")
					case "{{CURRENT_TIME}}":
						value = time.Now().Format("15:04:05")
					case "{{CURRENT_DATETIME}}":
						value = time.Now().Format("2006-01-02 15:04:05")
					}

					variables[key] = value
				}
			}
		}

		fmt.Printf("Running task with model: %s, prompt: %s, processor: %s\n", modelName, promptName, processorName)
		if len(variables) > 0 {
			fmt.Println("Variables:")
			for k, v := range variables {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}

		// Load the prompt with variables if provided
		var promptContent string
		var err error

		if len(variables) > 0 {
			promptContent, err = prompt.LoadPromptWithVariables(promptName, variables)
		} else {
			promptContent, err = prompt.LoadPrompt(promptName)
		}

		if err != nil {
			fmt.Printf("Error loading prompt: %v\n", err)
			return
		}

		// Execute the model
		response, err := models.ExecuteModel(modelName, promptContent, variables)
		if err != nil {
			fmt.Printf("Error executing model: %v\n", err)
			return
		}

		// Process the response
		err = processor.ProcessResponse(processorName, response)
		if err != nil {
			fmt.Printf("Error processing response: %v\n", err)
			return
		}

		fmt.Println("Task completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&modelName, "model", "", "AI model to use (openai, claude, gemini)")
	runCmd.Flags().StringVar(&promptName, "prompt", "", "Name of prompt file in cron_prompts directory")
	runCmd.Flags().StringVar(&processorName, "processor", "", "Response processor to use")
	runCmd.Flags().StringVar(&varsString, "vars", "", "Variables in format key1=value1,key2=value2")

	// Fail fast if we can't mark flags as required - this indicates a serious configuration issue
	markFlagRequiredOrFail(runCmd, "model")
	markFlagRequiredOrFail(runCmd, "prompt")
	markFlagRequiredOrFail(runCmd, "processor")
}

// markFlagRequiredOrFail marks a flag as required and fails early if there's an issue
func markFlagRequiredOrFail(cmd *cobra.Command, flagName string) {
	if err := cmd.MarkFlagRequired(flagName); err != nil {
		// This is a fatal error that indicates a serious misconfiguration
		// Use panic to ensure immediate termination with a clear error message
		panic(fmt.Sprintf("Critical configuration error: unable to mark %s flag as required: %v", flagName, err))
	}
}
