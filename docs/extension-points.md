# CronAI Extension Points

> **Note**: This document is a placeholder for more detailed extension documentation. Comprehensive extension guides are planned for post-MVP releases.

This document outlines the main extension points in CronAI, providing a high-level overview of how to extend the system with new capabilities.

## Current Status

The extension points documentation is in development and will be expanded in post-MVP releases. This placeholder covers the main extension points in the current architecture.

## Extension Categories

CronAI is designed to be extended in the following key areas:

1. **Model Integrations**: Adding support for new AI models
2. **Response Processors**: Creating new ways to process and deliver AI responses
3. **CLI Commands**: Extending the command-line interface
4. **Templating**: Creating custom templates for response formatting

## 1. Adding New AI Models

CronAI supports OpenAI, Claude, and Gemini in the MVP. To add a new model:

1. Implement the `ModelClient` interface in a new file in `internal/models/`
2. Register the model in the `defaultCreateModelClient` function
3. Add model-specific configuration to the `pkg/config` package

Basic implementation pattern:

```go
// NewCustomModelClient creates a client for a custom AI model
func NewCustomModelClient(modelConfig *config.ModelConfig) (*CustomModelClient, error) {
    // Initialize client with API key from environment
    // Implement model-specific configuration
    return &CustomModelClient{
        client: client,
        config: modelConfig,
    }, nil
}

// Execute sends a prompt to the custom model and returns the response
func (c *CustomModelClient) Execute(promptContent string) (*ModelResponse, error) {
    // Implement the API call to your model
    // Transform the response into the standard ModelResponse format
    return &ModelResponse{
        Content:   response.Text,
        Model:     "custom-model",
        Timestamp: time.Now(),
    }, nil
}
```text

## 2. Adding New Response Processors

Response processors handle the output from AI models. The MVP includes File, GitHub, and Console processors. To add a new processor:

1. Create a new file in `internal/processor/` implementing the `Processor` interface
2. Register the processor in `internal/processor/registry.go`
3. Add processor-specific configuration and environment variables

Basic implementation pattern:

```go
// NewCustomProcessor creates a new processor
func NewCustomProcessor(config ProcessorConfig) (Processor, error) {
    return &CustomProcessor{
        config: config,
    }, nil
}

// Process handles the model response
func (p *CustomProcessor) Process(response *models.ModelResponse, templateName string) error {
    // Process the response according to your requirements
    // Use templates if appropriate
    return nil
}

// Validate checks if the processor is properly configured
func (p *CustomProcessor) Validate() error {
    // Validate configuration
    return nil
}

// GetType returns the processor type identifier
func (p *CustomProcessor) GetType() string {
    return "custom"
}

// GetConfig returns the processor configuration
func (p *CustomProcessor) GetConfig() ProcessorConfig {
    return p.config
}
```text

Don't forget to register your processor:

```go
func init() {
    // Register in registry.go
    registry.RegisterProcessor("custom", NewCustomProcessor)
}
```text

## 3. Adding CLI Commands

CronAI uses the Cobra framework for CLI commands. To add a new command:

1. Create a new file in `cmd/cronai/cmd/` for your command
2. Implement the command logic
3. Register the command in `cmd/cronai/cmd/root.go`

Basic implementation pattern:

```go
// customCmd represents the custom command
var customCmd = &cobra.Command{
    Use:   "custom",
    Short: "A brief description of your command",
    Long:  `A longer description of your command`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implement command logic here
    },
}

func init() {
    rootCmd.AddCommand(customCmd)
    // Define flags for your command
    customCmd.Flags().StringVarP(&option, "option", "o", "", "Description of option")
}
```text

## 4. Adding Custom Templates

The templating system allows customization of response formats. To add custom templates:

1. Create template files in the `templates/` directory
2. Register your templates with the template manager
3. Use them in processors or other components

Basic implementation pattern:

```go
// Register a custom template
template.GetInstance().RegisterTemplate("custom_template", `
Subject: {{.Metadata.subject}}
Content: {{.Content}}
Execution: {{.ExecutionID}}
Timestamp: {{.Metadata.timestamp}}
`)

// Use the template in a processor
func (p *CustomProcessor) Process(response *models.ModelResponse, templateName string) error {
    // Use template name if provided, or default to processor-specific template
    tmplName := templateName
    if tmplName == "" {
        tmplName = "custom_template"
    }
    
    // Create template data
    tmplData := template.TemplateData{
        Content:     response.Content,
        Model:       response.Model,
        Timestamp:   response.Timestamp,
        PromptName:  response.PromptName,
        Variables:   response.Variables,
        ExecutionID: response.ExecutionID,
        Metadata:    make(map[string]string),
    }
    
    // Add metadata
    tmplData.Metadata["subject"] = "Custom Subject"
    
    // Execute template
    content, err := template.GetInstance().Execute(tmplName, tmplData)
    if err != nil {
        return err
    }
    
    // Use the formatted content
    // ...
    
    return nil
}
```text

## Future Extension Points

In post-MVP releases, additional extension points are planned:

1. **Web UI Extensions**: Adding custom components to the web interface
2. **Authentication Providers**: Integrating with custom authentication systems
3. **Storage Backends**: Supporting different storage options for responses
4. **Conditional Logic**: Extending the conditional execution system
5. **Metrics and Monitoring**: Custom monitoring integrations

## Best Practices

When extending CronAI, follow these best practices:

1. **Follow Project Patterns**: Match the existing code style and patterns
2. **Error Handling**: Use the error wrapping pattern with clear messages
3. **Configuration**: Support both environment variables and config file options
4. **Validation**: Validate all inputs and configurations
5. **Testing**: Write comprehensive tests for your extensions
6. **Documentation**: Update relevant documentation to include your extension

## Getting Help

If you need assistance with extending CronAI:

- Check the existing code for examples
- Review the [CONTRIBUTING.md](../CONTRIBUTING.md) file
- Look at the tests for usage patterns
- Submit questions as GitHub issues

Detailed extension guides will be available in future releases as the API matures.
