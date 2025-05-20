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
	templateName  string
	varsString    string
	modelParams   string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a single AI task immediately",
	Long: `Execute a single AI task immediately for testing or one-off operations.

This command allows you to run any prompt with your chosen AI model and processor
without scheduling. Perfect for testing new prompts or running ad-hoc tasks.

Features:
  • Variable substitution in prompts
  • Conditional logic with Go templates
  • Special variables like {{CURRENT_DATE}}
  • Model parameter customization
  • Response templating

Prompt Template Syntax:
  {{.Variables.name}}              - Variable substitution
  {{if eq .condition "value"}}...{{end}} - Conditional logic
  {{include "header.md"}}          - Include other files`,
	Example: `  # Basic execution
  cronai run --model=openai --prompt=daily_summary --processor=file

  # With variables
  cronai run --model=claude --prompt=report --processor=email \
    --vars="department=Sales,period=Q1"

  # With model parameters
  cronai run --model=gemini --prompt=creative --processor=file \
    --model-params="temperature=0.9,max_tokens=2000"

  # With special variables and template
  cronai run --model=openai --prompt=status --processor=slack \
    --vars="date={{CURRENT_DATE}}" --template=alert`,
	Run: func(_ *cobra.Command, _ []string) {
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

		if modelParams != "" {
			fmt.Printf("Model parameters: %s\n", modelParams)
		}

		// Load the prompt with variables if provided
		var promptContent string
		var err error

		if len(variables) > 0 {
			// Load the prompt with variables, which will use template processing if needed
			fmt.Println("Loading prompt with variables and template processing if applicable...")
			promptContent, err = prompt.LoadPromptWithVariables(promptName, variables)
		} else {
			// Load the prompt without variables
			fmt.Println("Loading prompt without variables...")
			promptContent, err = prompt.LoadPrompt(promptName)
		}

		if err != nil {
			// Provide more specific error messages for template-related errors
			if strings.Contains(err.Error(), "template") {
				fmt.Printf("Error processing prompt template: %v\n", err)
				fmt.Println("Please check the template syntax in your prompt file.")
			} else {
				fmt.Printf("Error loading prompt: %v\n", err)
			}
			return
		}

		// Execute the model with model parameters
		response, err := models.ExecuteModel(modelName, promptContent, variables, modelParams)
		if err != nil {
			fmt.Printf("Error executing model: %v\n", err)
			return
		}

		// Process the response
		err = processor.ProcessResponse(processorName, response, templateName)
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
	runCmd.Flags().StringVar(&templateName, "template", "", "Optional template name to use for formatting the response")
	runCmd.Flags().StringVar(&varsString, "vars", "", "Variables in format key1=value1,key2=value2")
	runCmd.Flags().StringVar(&modelParams, "model-params", "", "Model parameters in format temperature=0.7,max_tokens=1024")

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
