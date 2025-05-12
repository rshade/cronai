package cmd

import (
	"fmt"

	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
	"github.com/rshade/cronai/internal/prompt"
	"github.com/spf13/cobra"
)

var (
	modelName     string
	promptName    string
	processorName string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a single task immediately",
	Long:  `Run a single task immediately without scheduling.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running task with model: %s, prompt: %s, processor: %s\n", modelName, promptName, processorName)

		// Load the prompt
		promptContent, err := prompt.LoadPrompt(promptName)
		if err != nil {
			fmt.Printf("Error loading prompt: %v\n", err)
			return
		}

		// Execute the model
		response, err := models.ExecuteModel(modelName, promptContent)
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

	runCmd.MarkFlagRequired("model")
	runCmd.MarkFlagRequired("prompt")
	runCmd.MarkFlagRequired("processor")
}
