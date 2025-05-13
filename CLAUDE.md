# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Implementation Patterns

### Templating System

The templating system follows these patterns:
- Uses Go's built-in `text/template` package
- Singleton pattern for the template manager
- Default templates are registered at initialization
- SafeExecute method with fallback mechanism for templates
- Naming conventions for templates based on processor type and purpose

## Coding Standards and Robustness

### Error Handling

- All error returns must be properly checked and handled
- Use `fmt.Errorf("error message: %w", err)` for error wrapping
- Defer functions should handle errors using anonymous functions
- Use multi-error approach for preserving both original and secondary errors (e.g., file close errors)
- For critical configuration errors, use early failure strategies (e.g., panic with clear messages)
- File operations must have proper error checking, especially for:
  - Close operations on files
  - Remove operations on temporary files

### Code Quality and Linting

- **ALWAYS run `make lint` before committing any code changes**
- **ALWAYS run `make lint` after every code change to ensure compliance with coding standards**
- Fix all linting issues before submitting pull requests
- The linting process checks:
  - Code formatting with `gofmt`
  - Static analysis with `go vet`
  - Comprehensive linting with `golangci-lint`
- Adhere to Go idiomatic patterns
- Use named return values where it improves readability
- Use switch statements instead of long if-else chains
- Follow strict function signature requirements

### Testing Requirements

- All new functions should have corresponding tests
- Test coverage should be maintained or improved with each PR
- Mock external dependencies in tests
- Use table-driven tests when appropriate

### Git and CI/CD Practices

- **ALWAYS use conventional commit format for ALL commits**:
  - Format: `<type>(<scope>): <description>`
  - Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert
  - Scope: Optional component name (e.g., prompt, cron, models)
  - Description: Concise present-tense summary
  - Examples:
    - `feat(prompt): implement file-based prompt management`
    - `fix(cron): resolve issue with scheduler timing`
    - `docs: update API documentation for prompt commands`
    - `refactor(models): improve claude client implementation`
  - For breaking changes, add `BREAKING CHANGE:` in the footer
  - See [Conventional Commits specification](https://www.conventionalcommits.org/) for details
- Commits that do not follow this format will be automatically rejected by CI
- **NEVER commit directly - always update PR_MESSAGE.md** with your commit message:
  - PR_MESSAGE.md should use the same conventional commit format
  - Update PR_MESSAGE.md whenever you have changes ready to commit
  - This allows for review of commit messages before actual commits are created
- Make PR titles follow the same conventional commit format
- Keep PR sizes manageable (ideally under 300 lines of changes)
- Update documentation when changing functionality
  - When adding or modifying features, always update the README.md to keep it in sync
  - Ensure all examples in documentation reflect the latest capabilities
- Use GitHub Actions for CI/CD pipelines
- Ensure CI passes before merging PRs

## Project Overview

CronAI is a Go utility designed to run AI model prompts on a cron-type schedule. It allows scheduled execution of AI prompts and automatic processing of responses through various channels (email, Slack, webhooks, file output).

## Repository Structure

This project follows standard Go project structure:

```
cronai/
├── .github/               # GitHub configuration
│   └── workflows/         # GitHub Actions workflows
│       ├── build.yml      # Build workflow
│       ├── commit-check.yml # Conventional commit checker
│       ├── pr-check.yml   # PR validation workflow
│       └── todo.yml       # TODO to Issue converter
├── cmd/cronai/            # Main application entry point
│   ├── main.go            # Entry point
│   └── cmd/               # CLI commands (using Cobra)
│       ├── root.go        # Root command
│       ├── start.go       # Start the service
│       ├── run.go         # Run a single task
│       └── list.go        # List scheduled tasks
├── internal/              # Private application code
│   ├── cron/              # Cron scheduling functionality
│   │   ├── service.go     # Core scheduling service
│   │   └── service_test.go # Tests for cron service
│   ├── models/            # AI model integrations (OpenAI, Claude, Gemini)
│   │   └── models.go      # Model execution logic
│   ├── processor/         # Response processors
│   │   └── processor.go   # Response processing logic
│   └── prompt/            # Prompt loading
│       ├── loader.go      # Prompt loading logic
│       └── loader_test.go # Tests for prompt loader
├── pkg/                   # Public packages
│   └── config/            # Configuration loading
├── cron_prompts/          # Directory for markdown prompt files
│   ├── product_manager.md # Example prompt
│   ├── report_template.md # Example prompt with variables
│   ├── weekly_report.md   # Example prompt
│   ├── system_health.md   # Example prompt
│   └── monitoring_check.md # Example prompt
├── docs/                  # Documentation
│   └── systemd.md         # Systemd service setup guide
├── .commitlintrc.js       # Commitlint configuration
├── .goreleaser.yml        # GoReleaser configuration
├── cronai.config.example  # Example configuration file
├── cronai.config.variables.example # Example with variables
├── cronai.service         # Systemd service file
├── .env.example           # Example environment variables
├── Makefile               # Build and development commands
└── CLAUDE.md              # This file
```

## Configuration Format

The configuration file uses the following format:
```
timestamp model prompt response_processor [variables]
```

Where:
- **timestamp**: Standard cron format (minute hour day-of-month month day-of-week)
- **model**: AI model to use (openai, claude, gemini)
- **prompt**: Name of prompt file in cron_prompts directory (with or without .md extension)
- **response_processor**: How to process the response (email, slack, webhook, file)
- **variables** (optional): Variables to replace in the prompt file, in the format `key1=value1,key2=value2,...`

Examples:
```
# Basic configuration without variables
0 8 * * * claude product_manager slack-pm-channel

# Configuration with variables
0 9 1 * * claude report_template email-team@company.com reportType=Monthly,date={{CURRENT_DATE}},project=CronAI
```

### Special Variables

Special variables that are automatically populated:
- `{{CURRENT_DATE}}`: Current date in YYYY-MM-DD format
- `{{CURRENT_TIME}}`: Current time in HH:MM:SS format
- `{{CURRENT_DATETIME}}`: Current date and time in YYYY-MM-DD HH:MM:SS format

## Development Commands

CronAI includes a Makefile to simplify common development tasks.

### Setup and Dependencies

```bash
# Initialize module (already done)
go mod init github.com/rshade/cronai

# Setup the development environment
make setup

# Download dependencies
go mod download

# Add required dependencies
go get github.com/spf13/cobra            # CLI framework
go get github.com/joho/godotenv          # Environment variable loading
go get github.com/robfig/cron/v3         # Cron scheduling

# Tidy the go.mod file
go mod tidy
```

### Build and Run

```bash
# Build the project
make build
# or
go build -o cronai ./cmd/cronai

# Run the application
make run
# or
./cronai start

# Install the application
make install
# or
go install ./cmd/cronai
```

### Testing

```bash
# Run all tests
make test
# or
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/models

# Run a specific test
go test -run TestParseConfig ./internal/cron
```

### Linting and Code Quality

```bash
# Run linter
make lint
# or
go vet ./...
golangci-lint run
```

## Roadmap Overview

The project has a defined roadmap divided into four milestones:

### Q2 2025 - MVP Release
- ✅ Basic variable support in prompts (#5)
- Model-specific configuration support (#6)
- Response processor templating (#7)
- Model fallback mechanism (#8)
- Comprehensive logging and error handling (#9)

### Q3 2025 - Enhanced Usability
- Basic web UI (#10)
- Conditional logic in prompt templates (#11)
- Prompt testing tool (#12)
- Claude 3 model support (#13)

### Q4 2025 - Integration & Scale
- External API for integration (#14)
- Performance metrics and analytics (#15)
- Distributed task execution (#16)
- CI/CD platform integrations (#17)

### Q1 2026 - Enterprise Features
- Role-based access control (#18)
- Audit logging and compliance (#19)
- Prompt library management (#20)
- Cost tracking and budget management (#21)

## Key Epics

The development is organized around four major epics:

### Epic: Enhanced Templating System (#1)
This epic covers the implementation of an advanced templating system for CronAI prompts, allowing for more dynamic and flexible prompts with variables, conditional logic, and template inheritance.

### Epic: Advanced Model Support (#2)
This epic focuses on expanding the AI model support in CronAI, including additional models, model versioning capabilities, model fallbacks, and model-specific configurations.

### Epic: Web UI and Management Console (#3)
This epic encompasses the development of a web-based user interface for CronAI, making it easier to configure, monitor, and manage scheduled tasks without editing configuration files manually.

### Epic: Enterprise Integration (#4)
This epic focuses on features necessary for enterprise adoption, including advanced security, monitoring, integration with corporate systems, and scalability improvements.

## Key Components

### Cron Parser
Located in `internal/cron` - Parses the configuration file and sets up the scheduled tasks.

### Model Adapters
Located in `internal/models` - Implements adapters for different AI models (OpenAI, Claude, Gemini).

### Prompt Loader
Located in `internal/prompt` - Loads and prepares prompt files for submission to models.
- `LoadPrompt` - Loads a prompt file from the cron_prompts directory
- `LoadPromptWithVariables` - Loads a prompt file and replaces variables with provided values
- `ApplyVariables` - Replaces variable placeholders in format {{variable_name}} with their values

### Response Processors
Located in `internal/processor` - Processes model responses (sending to email, Slack, webhooks, or saving to file).

## Environment Variables

The application uses a `.env` file for configuration with the following variables:
- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Claude API key
- `GOOGLE_API_KEY` - Gemini API key
- Various processor-specific configuration variables (SMTP, Slack tokens, etc.)

## CLI Commands

The application uses Cobra for CLI commands:
- `cronai start` - Start the service with the configuration file
- `cronai run` - Run a single task immediately
  - `--model`: AI model to use
  - `--prompt`: Name of prompt file
  - `--processor`: Response processor to use
  - `--vars`: Variables for the prompt in format "key1=value1,key2=value2"
- `cronai list` - List all scheduled tasks