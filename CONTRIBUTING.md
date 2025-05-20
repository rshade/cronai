# Contributing to CronAI

Thank you for your interest in contributing to CronAI! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contribution Guidelines](#contribution-guidelines)
- [Adding Processors](#adding-processors)
- [Testing](#testing)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally

   ```bash
   git clone https://github.com/YOUR_USERNAME/cronai.git
   cd cronai
   ```text

3. Add the original repository as an upstream remote

   ```bash
   git remote add upstream https://github.com/rshade/cronai.git
   ```text

4. Create a new branch for your feature or bug fix

   ```bash
   git checkout -b feature/your-feature-name
   ```text

## Development Setup

### Prerequisites

- Go 1.21 or higher
- `make` command
- `golangci-lint` for code linting
- `git-chglog` for changelog generation (optional)
- `commitlint` for commit message validation

### Initial Setup

```bash
# Setup development environment
make setup

# Install dependencies
go mod download

# Install commit message linting (requires Node.js)
npm install

# Run tests to verify setup
make test
```text

### Development Workflow

1. **Write code** following the project's conventions
2. **Write tests** for your changes
3. **Run linter** to ensure code quality

   ```bash
   # For CI and pull requests (strict mode - no automatic fixes)
   make lint
   
   # For local development (automatically fixes issues)
   make lint-fix
   ```text

4. **Run tests** to ensure everything works

   ```bash
   make test
   ```text

5. **Update documentation** if needed (README.md, docs/, etc.)

## Developer Documentation

Before contributing to CronAI, please review these key documentation resources:

- [Architecture Overview](docs/architecture.md) - Understand the system design, component relationships, and data flow
- [Extension Points](docs/extension-points.md) - Learn how to extend CronAI with new features
- [API Documentation](docs/api.md) - API endpoints (planned for future versions)
- [Limitations and Improvements](docs/limitations-and-improvements.md) - Current limitations and future roadmap

## Contribution Guidelines

### Code Quality

- **Always run `make lint` before committing** (or `make lint-fix` for automatic fixes)
- Fix all linting issues reported by `golangci-lint`
- Follow Go idiomatic patterns
- Use meaningful variable and function names
- Add comments for complex logic
- Handle all errors appropriately using our error wrapping pattern
- Use the project's logging system instead of `fmt.Print` statements

### Error Handling

Use the internal errors package for consistent error handling:

```go
import "github.com/rshade/cronai/internal/errors"

// Wrap errors with category and context
if err != nil {
    return errors.Wrap(errors.CategoryConfiguration, err, "failed to load config")
}
```text

Error categories:

- `CategoryConfiguration`: Configuration-related errors
- `CategoryValidation`: Input validation errors
- `CategoryApplication`: General application errors
- `CategoryIO`: I/O related errors

### Logging

Use the internal logger package:

```go
import "github.com/rshade/cronai/internal/logger"

// Get the logger
log := logger.DefaultLogger()

// Log with fields
log.Info("Processing response", logger.Fields{
    "processor": processorName,
    "model":     response.Model,
    "execution": response.ExecutionID,
})

// Log errors
log.Error("Failed to process", logger.Fields{
    "error": err.Error(),
})
```text

## Adding Processors

Processors handle the output from AI models. Follow these steps to add a new processor:

### 1. Create the Processor File

Create a new file in `internal/processor/yourprocessor.go`:

```go
package processor

import (
    "fmt"
    "time"
    
    "github.com/rshade/cronai/internal/errors"
    "github.com/rshade/cronai/internal/logger"
    "github.com/rshade/cronai/internal/models"
    "github.com/rshade/cronai/internal/processor/template"
)

// YourProcessor handles your custom processing
type YourProcessor struct {
    config ProcessorConfig
}

// NewYourProcessor creates a new instance
func NewYourProcessor(config ProcessorConfig) (Processor, error) {
    return &YourProcessor{
        config: config,
    }, nil
}

// Process handles the model response with optional template
func (y *YourProcessor) Process(response *models.ModelResponse, templateName string) error {
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
    
    // Add standard metadata
    tmplData.Metadata["timestamp"] = response.Timestamp.Format(time.RFC3339)
    tmplData.Metadata["date"] = response.Timestamp.Format("2006-01-02")
    tmplData.Metadata["time"] = response.Timestamp.Format("15:04:05")
    tmplData.Metadata["execution_id"] = response.ExecutionID
    tmplData.Metadata["processor"] = y.GetType()
    
    // Process with your logic
    return y.processWithTemplate(tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (y *YourProcessor) Validate() error {
    if y.config.Target == "" {
        return errors.Wrap(errors.CategoryValidation,
            fmt.Errorf("target cannot be empty"),
            "invalid processor configuration")
    }
    
    // Check for required environment variables
    // Example: if os.Getenv("YOUR_API_KEY") == "" { ... }
    
    return nil
}

// GetType returns the processor type identifier
func (y *YourProcessor) GetType() string {
    return "yourprocessor"
}

// GetConfig returns the processor configuration
func (y *YourProcessor) GetConfig() ProcessorConfig {
    return y.config
}

// processWithTemplate implements your processing logic
func (y *YourProcessor) processWithTemplate(data template.TemplateData, templateName string) error {
    // Implement your processor logic here
    return nil
}
```text

### 2. Define Environment Variables

Add your environment constants to `internal/processor/env.go`:

```go
const (
    // YourProcessor environment variables
    EnvYourAPIKey    = "YOUR_API_KEY"
    EnvYourEndpoint  = "YOUR_ENDPOINT"
    
    // Default values
    DefaultYourEndpoint = "https://api.example.com"
)
```text

For dynamic environment variables (supporting multiple configurations):

```go
const (
    // Dynamic environment variables for different types
    EnvYourKeyPrefix = "YOUR_KEY_"  // e.g., YOUR_KEY_PRODUCTION, YOUR_KEY_STAGING
)
```text

### 3. Register the Processor

Add your processor to the registry in `internal/processor/registry.go`:

```go
func init() {
    // ... other registrations
    
    // Register YourProcessor
    registry.RegisterProcessor("yourprocessor", NewYourProcessor)
}
```text

### 4. Update Documentation

1. Add environment variables to `.env.example`:

   ```text
   # YourProcessor Configuration
   YOUR_API_KEY=your-api-key
   YOUR_ENDPOINT=https://api.example.com
   ```text

2. Update README.md with processor usage:

   ```text
   # Using YourProcessor
   0 9 * * * claude daily_report yourprocessor-target
   ```text

### Message Format

All processors receive a `ModelResponse` with the following structure:

```go
type ModelResponse struct {
    Content     string            // The AI model's response text
    Model       string            // Model name (e.g., "openai", "claude", "gemini")
    PromptName  string            // Name of the prompt file used
    Variables   map[string]string // Variables used in the prompt
    Timestamp   time.Time         // When the response was generated
    ExecutionID string            // Unique execution identifier
}
```text

The `TemplateData` structure available in templates:

```go
type TemplateData struct {
    Content     string            // Model response content
    Model       string            // Model name
    Timestamp   time.Time         // Response timestamp
    PromptName  string            // Name of the prompt
    Variables   map[string]string // Custom variables
    ExecutionID string            // Unique execution identifier
    Metadata    map[string]string // Additional metadata
    Parent      interface{}       // Parent template data for inheritance
}
```text

Common metadata fields automatically populated:

- `timestamp`: RFC3339 formatted timestamp
- `date`: Date in YYYY-MM-DD format
- `time`: Time in HH:MM:SS format
- `execution_id`: Unique execution identifier
- `processor`: Processor type name
- `template`: Template name (if specified)

### Environment Variable Patterns

CronAI uses consistent patterns for environment variables:

1. **Base variables**: `PROCESSOR_OPTION`
   - Examples: `SLACK_TOKEN`, `SMTP_SERVER`, `WEBHOOK_URL`

2. **Type-specific variables**: `PROCESSOR_OPTION_TYPE`
   - Examples: `WEBHOOK_URL_MONITORING`, `WEBHOOK_METHOD_ALERTS`
   - Allows different configurations for different use cases

3. **Helper functions**: Use the provided helper functions for dynamic variables:

   ```go
   // Get webhook URL for a specific type
   url := GetWebhookURL("monitoring")  // Checks WEBHOOK_URL_MONITORING first, then WEBHOOK_URL
   
   // Get environment variable with default
   port := GetEnvWithDefault(EnvSMTPPort, DefaultSMTPPort)
   ```text

## Testing

### Writing Tests

- Write unit tests for all new functions
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Follow existing test patterns in the codebase

Example test structure:

```go
func TestYourProcessor_Process(t *testing.T) {
    tests := []struct {
        name     string
        config   ProcessorConfig
        response *models.ModelResponse
        wantErr  bool
    }{
        {
            name: "successful processing",
            config: ProcessorConfig{
                Type:   "yourprocessor",
                Target: "target",
            },
            response: &models.ModelResponse{
                Content: "test content",
                Model:   "test-model",
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            processor, err := NewYourProcessor(tt.config)
            if err != nil {
                t.Fatal(err)
            }
            
            err = processor.Process(tt.response, "")
            if (err != nil) != tt.wantErr {
                t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```text

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/processor/...

# Run a specific test
go test -run TestYourProcessor ./internal/processor
```text

## Code Style

### Go Conventions

- Follow standard Go formatting (enforced by `gofmt`)
- Use meaningful names for packages, types, functions, and variables
- Keep functions small and focused
- Avoid deep nesting
- Return early from functions when possible
- Use interfaces for abstraction
- Document exported types and functions

### Project-Specific Conventions

- Place processor implementations in `internal/processor/`
- Use the error wrapping pattern with categories
- Log operations with appropriate context fields
- Follow existing patterns for configuration and environment variables
- Implement all required interface methods

## Commit Messages

CronAI uses [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.

### Format

```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```text

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Changes that don't affect code meaning
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Changes to build process or auxiliary tools
- `perf`: Performance improvements
- `ci`: CI/CD changes
- `build`: Build system changes
- `revert`: Reverting a previous commit

### Scopes

- `processor`: Changes to processor system
- `cron`: Changes to cron scheduling
- `models`: Changes to AI model integrations
- `prompt`: Changes to prompt handling
- `config`: Changes to configuration
- `template`: Changes to templating system

### Examples

```text
feat(processor): add Discord notification processor
fix(cron): resolve timing issue with overlapping tasks
docs: update processor development guide
refactor(models): improve error handling in Claude client
test(processor): add unit tests for webhook processor
```text

### Validation

Before committing, validate your message:

```bash
# Validate the last commit message
npx commitlint --from HEAD~1

# Validate a commit message file
npx commitlint --from PR_MESSAGE.md
```text

## Pull Request Process

1. **Update your branch** with the latest upstream changes:

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```text

2. **Ensure all tests pass and linting is clean**:

   ```bash
   make lint
   make test
   ```text

3. **Update documentation** if your changes affect user-facing functionality

4. **Create the pull request**:
   - Use a descriptive title following conventional commit format
   - Fill out the PR template completely
   - Reference any related issues
   - Include examples or screenshots if relevant

5. **PR Requirements**:
   - All CI checks must pass
   - Code review approval required
   - Keep PR size manageable (preferably under 300 lines)
   - Respond to review feedback promptly

6. **Update PR_MESSAGE.md** with your commit message for review:

   ```markdown
   feat(processor): add Discord notification processor
   
   Implements a new processor for sending notifications to Discord channels
   using webhooks. Supports both simple text messages and rich embeds with
   customizable templates.
   
   Closes #123
   ```text

### PR Review Checklist

Before requesting review, ensure:

- [ ] Code follows project conventions
- [ ] All tests pass
- [ ] Linting is clean
- [ ] Documentation is updated
- [ ] Commit messages follow conventional commits
- [ ] PR description is complete
- [ ] Related issues are referenced

## Questions or Need Help?

- Check existing issues and pull requests
- Read the project documentation
- Ask questions in GitHub issues
- Refer to the [project roadmap](README.md#roadmap-overview) for planned features

Thank you for contributing to CronAI!
