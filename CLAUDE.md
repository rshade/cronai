# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Implementation Patterns

### Model Configuration System

The model configuration system follows these patterns:

- Provides consistent configuration across multiple AI models (OpenAI, Claude, Gemini)
- Offers common parameters (temperature, max_tokens) applicable to all models
- Supports model-specific parameters through dot notation (e.g., `openai.system_message`, `gemini.safety_setting`)
- Configuration can be loaded from environment variables, command-line arguments, or config files
- Implements a clear parameter precedence: task-specific > environment > defaults
- Uses official SDKs for all model integrations (go-openai, anthropic-sdk-go, generative-ai-go)
- Validates all parameter values before use

### Templating System

The templating system follows these patterns:

- Uses Go's built-in `text/template` package
- Singleton pattern for the template manager
- Default templates are registered at initialization
- SafeExecute method with fallback mechanism for templates
- Naming conventions for templates based on processor type and purpose

### Processor System

The processor system follows these patterns:

- **Standardized API**: All processors implement the `Processor` interface
- **Registry Pattern**: Global registry for processor factories
- **Factory Pattern**: Processors are created via factory functions
- **Configuration Management**: Processors use `ProcessorConfig` for standardized configuration
- **Environment Variable Naming**: Consistent naming scheme for environment variables
  - Base variables: `PROCESSOR_OPTION` (e.g., `SLACK_TOKEN`, `SMTP_SERVER`, `GITHUB_TOKEN`)
  - Type-specific variables: `PROCESSOR_OPTION_TYPE` (e.g., `WEBHOOK_URL_MONITORING`)
- **Validation**: All processors implement a `Validate()` method to check configuration
- **Template Integration**: Processors integrate with the templating system for output formatting
- **Error Handling**: Consistent error wrapping using the internal errors package
- **GitHub Integration**: The GitHub processor uses the google/go-github library for API calls

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

- **For CI/PR checks: run `make vet` and `make lint` separately (to avoid timeout issues)**
- **For local development: run `make lint-fix-all` to automatically fix all issues including go vet**
- **Quick linting without go vet: run `make lint` or `make lint-fix`**
- **Run all linting checks: run `make lint-all` (includes go vet + lint)**
- Fix all linting issues before submitting pull requests
- The linting process checks:
  - Code formatting with `gofmt`
  - Static analysis with `go vet` (separated into `make vet`)
  - Comprehensive linting with `golangci-lint`
  - Markdown linting with `markdownlint`
  - Revive linter with the following rules:
    - `atomic`: Avoid shadowing variables
    - `blank-imports`: Avoid blank imports
    - `context-as-argument`: Pass context as first parameter
    - `context-keys-type`: Use value types in context keys
    - `dot-imports`: Avoid dot imports
    - `empty-block`: Avoid empty blocks
    - `error-return`: Return values on error
    - `error-strings`: Error strings should not be capitalized
    - `error-naming`: Error variables should be prefixed with "err" or "Err"
    - `exported`: Exported function/variable comments
    - `if-return`: Avoid redundant if-then-return statements
    - `increment-decrement`: Use i++ and i-- instead of i += 1 and i -= 1
    - `var-naming`: Variable naming conventions
    - `var-declaration`: Reduce variable scope
    - `package-comments`: All packages should have comments
    - `range`: Simplify range expressions
    - `receiver-naming`: Receiver names should be consistent
    - `time-naming`: Avoid using time.Now in variable names
    - `unexported-return`: Unexported return values
    - `indent-error-flow`: Error handling indent
    - `errorf`: Use errors.New() instead of fmt.Errorf() when appropriate
    - `empty-lines`: Control empty lines between declarations
    - `superfluous-else`: Avoid superfluous else statements
    - `unreachable-code`: Check for unreachable code
    - `redefines-builtin-id`: Don't redefine built-in identifiers
    - `waitgroup-by-value`: Don't pass sync.WaitGroup by value
- Adhere to Go idiomatic patterns
- Use named return values where it improves readability
- Use switch statements instead of long if-else chains
- Follow strict function signature requirements

### Go Package Documentation

- **ALWAYS add package comments to every Go file**:
  - Every package should have at least one file with a package comment
  - For consistency, add a package comment to every file
  - Use the standard format: `// Package packagename provides ...`
  - Example: `// Package cmd implements the command line interface for CronAI.`
  - Package comments should be concise but descriptive of the package's purpose
  - Package comments must appear directly before the package clause
- Package comments are required to pass the `package-comments` linting rule
- When creating new files, always start with an appropriate package comment
- When modifying existing files without package comments, add them

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
  - **ALWAYS validate PR_MESSAGE.md with commitlint before committing**:
    - Ensure commit messages follow all commitlint rules
    - Keep body lines under 100 characters to avoid `body-max-line-length` errors
    - Use line breaks to split long descriptions into multiple shorter lines
    - If using multi-paragraph bodies, ensure proper line breaks and formatting
    - Run `cat PR_MESSAGE.md | npx commitlint` to validate before committing
    - Fix any commitlint errors before proceeding with the commit
- Make PR titles follow the same conventional commit format
- Keep PR sizes manageable (ideally under 300 lines of changes)
- Update documentation when changing functionality
  - When adding or modifying features, always update the README.md to keep it in sync
  - Ensure all examples in documentation reflect the latest capabilities
- Use GitHub Actions for CI/CD pipelines
- Ensure CI passes before merging PRs

### Changelog Management

- Changelogs are generated automatically from conventional commit messages
- Use `make changelog` to generate a CHANGELOG.md file for releases
  - By default, it generates changes since the last tag (or initial commit if no tags)
  - Specify custom range with `make changelog FROM=<tag/commit> TO=<tag/commit>`
- The generated changelog categorizes commits by type:
  - Features (feat)
  - Bug Fixes (fix)
  - Performance Improvements (perf)
  - Code Refactoring (refactor)
  - Documentation (docs)
  - Tests (test)
  - Build System (build)
  - Continuous Integration (ci)
  - Chores (chore)
- Always generate a changelog before creating a new release
- Include the changelog in release notes when creating GitHub releases

### Release Management

The project uses GitHub Actions for automated releases through GoReleaser:

#### Workflows

- **Build Workflow** (`.github/workflows/build.yml`):
  - Triggers on pushes to main branch
  - Runs tests with coverage
  - Builds the application
  - Uploads build artifacts
  - Does NOT create releases

- **Release Workflow** (`.github/workflows/release.yml`):
  - Triggers on any tag matching `v*` pattern
  - Runs tests before release
  - Generates changelog using git-chglog
  - Creates GitHub releases using GoReleaser
  - Supports both stable and prerelease versions

#### Release Types

The release system automatically handles different release types based on tag naming:

- **Stable Releases**: Tags like `v1.0.0`, `v2.1.3`
  - Creates normal GitHub releases
  - Marked as "Latest" release

- **Prerelease Versions**: Tags containing `-alpha`, `-beta`, `-rc`, etc.
  - Examples: `v0.0.2-beta`, `v1.0.0-rc1`, `v2.0.0-alpha.1`
  - Automatically marked as "Pre-release" on GitHub
  - Not marked as "Latest" release

#### Creating Releases

To create a release:

1. Ensure all changes are committed and pushed to main
2. Create and push a tag:

   ```bash
   # For stable release
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   
   # For beta release
   git tag -a v0.0.2-beta -m "Release v0.0.2-beta"
   git push origin v0.0.2-beta
   ```

3. GitHub Actions will automatically:
   - Run tests
   - Generate changelog
   - Build binaries for multiple platforms
   - Create GitHub release with assets
   - Mark prerelease appropriately

#### GoReleaser Configuration

The `.goreleaser.yml` file configures:

- Multi-platform builds (Linux, Windows, macOS)
- Architecture support (amd64, arm64)
- Archive naming and contents
- Automatic prerelease detection (`prerelease: auto`)
- Changelog generation and filtering

## Project Overview

CronAI is a Go utility designed to run AI model prompts on a cron-type schedule. It allows scheduled execution of AI prompts and automatic processing of responses through various channels (email, Slack, webhooks, file output).

## Repository Structure

This project follows standard Go project structure:

```text
cronai/
├── .github/               # GitHub configuration
│   └── workflows/         # GitHub Actions workflows
│       ├── build.yml      # Build workflow (runs on main branch)
│       ├── release.yml    # Release workflow (runs on tags)
│       ├── commit-check.yml # Conventional commit checker
│       ├── pr-check.yml   # PR validation workflow
│       ├── integration-tests.yml # Integration tests
│       ├── deploy-docs.yml # Documentation deployment
│       └── todo.yml       # TODO to Issue converter
├── cmd/cronai/            # Main application entry point
│   ├── main.go            # Entry point
│   └── cmd/               # CLI commands (using Cobra)
│       ├── root.go        # Root command
│       ├── start.go       # Start the service
│       ├── run.go         # Run a single task
│       ├── list.go        # List scheduled tasks
│       ├── help.go        # Enhanced help system
│       ├── prompt.go      # Prompt management
│       └── validate.go    # Template validation
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
```text

## Prompt Files Location

All prompt files are stored in the `cron_prompts/` directory. When referencing prompts:
- `product_manager` → `cron_prompts/product_manager.md`
- `weekly_report` → `cron_prompts/weekly_report.md`
- `system_health` → `cron_prompts/system_health.md`
- `monitoring_check` → `cron_prompts/monitoring_check.md`
- `report_template` → `cron_prompts/report_template.md`

The prompt loader automatically looks in the `cron_prompts/` directory, so you only need to specify the filename (with or without .md extension) in configurations.

## Configuration Format

The configuration file uses the following format:

```text
timestamp model prompt response_processor [variables]
```text

Where:

- **timestamp**: Standard cron format (minute hour day-of-month month day-of-week)
- **model**: AI model to use (openai, claude, gemini)
- **prompt**: Name of prompt file in cron_prompts directory (with or without .md extension)
- **response_processor**: How to process the response (email, slack, webhook, file)
- **variables** (optional): Variables to replace in the prompt file, in the format `key1=value1,key2=value2,...`

Examples:

```text
# Basic configuration without variables
0 8 * * * claude product_manager slack-pm-channel

# Configuration with variables
0 9 1 * * claude report_template email-team@company.com reportType=Monthly,date={{CURRENT_DATE}},project=CronAI
```text

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
```text

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
```text

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
```text

### Linting and Code Quality

```bash
# Run linter
make lint
# or
go vet ./...
golangci-lint run
```text

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

## Integration Testing Notes

- **NEVER close GitHub issue #89** - This issue is used for GitHub integration tests and must remain open.
- All GitHub processor tests utilize issue #89 to verify comment functionality.
- Integration tests verify that GitHub comments are actually created by fetching comments before/after processing and checking for new comments with expected content.
- Test comments are kept by default for visibility, but can be automatically cleaned up by setting `CLEANUP_TEST_COMMENTS=1`.
- To run integration tests with real GitHub API calls, set `RUN_INTEGRATION_TESTS=1` and provide a valid `GITHUB_TOKEN`.

## Environment Variables

The application uses a `.env` file for configuration with the following variables:

- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Claude API key
- `GOOGLE_API_KEY` - Gemini API key
- `GITHUB_TOKEN` - GitHub personal access token for GitHub processor integration
- `RUN_INTEGRATION_TESTS` - Set to "1" to run integration tests with real API calls
- `CLEANUP_TEST_COMMENTS` - Set to "1" to automatically delete test comments after verification (default: keep comments)
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

## Documentation Maintenance

When working on CronAI, ensure all documentation stays up to date:

- **README.md**: Keep the main documentation updated with new features and changes
- **CONTRIBUTING.md**: Update the contributing guide when:
  - New development tools or processes are added
  - Processor development process changes
  - New environment variables or configuration patterns are introduced
  - Testing requirements or linting rules change
  - CI/CD processes are modified
- **Other documentation files**:
  - `docs/` directory files should be updated when their respective features change
  - `.env.example` should include all new environment variables
  - Example config files should demonstrate new features

Always ensure documentation reflects the current state of the code and provides accurate guidance for users and contributors.
