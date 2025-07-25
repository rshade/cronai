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
- **Slack Integration**: The Slack processor supports both webhook URLs and OAuth token methods for maximum flexibility
  - Dual authentication: Uses `SLACK_TOKEN` (OAuth) or `SLACK_WEBHOOK_URL` (webhook)
  - Automatic method selection based on available environment variables
  - JSON payload validation with proper error handling
  - Automatic monitoring template detection for alert-style formatting
  - Comprehensive error handling with proper request/response lifecycle management

### Claude GitHub Action Integration

The project uses a sophisticated multi-persona Claude GitHub Assistant system:

- **Persona System**: Three distinct personas activated by GitHub labels:
  - `claude-reviewer`: Code reviewer persona for detailed technical review (.github/claude/code-reviewer.md)
  - `claude-engineer`: Software engineer persona for implementation guidance (.github/claude/software-engineer.md)
  - `claude-assistant`: General assistant persona for questions and help (.github/claude/default.md)
- **Label-Based Selection**: The workflow automatically selects the appropriate persona based on issue/PR labels
- **System Prompts**: Stored in `.github/claude/` directory as separate markdown files for version control
- **Workflow Integration**: The GitHub Action workflow reads prompt files and passes them to the `system_prompt` parameter
- **Issue Templates**: Three specialized templates that automatically add the correct persona labels
- **Documentation**: Comprehensive setup guide in `.github/claude/README.md`

#### Claude GitHub Action Configuration Patterns

- Use conditional logic in workflows to select different configurations based on labels
- Add explicit permissions to workflows for security (`contents: read`, `issues: write`, `pull-requests: write`)
- Set conversation limits (`max_turns: "10"`) to control API costs
- Check out repository before reading local files in workflows
- Pass project context via `claude_env` environment variables
- The `anthropics/claude-code-action@beta` doesn't natively read prompt files - content must be read and passed as string

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

- **For CI/PR checks: run individual lint targets to avoid timeout issues**:
  - `make vet` - Run go vet on all packages at once
  - `make vet-pkg` - Run go vet package by package (slower but avoids timeouts)
  - `make lint-fmt` - Check Go formatting
  - `make lint-golangci` - Run golangci-lint on all packages
  - `make lint-golangci-pkg` - Run golangci-lint package by package (slower but avoids timeouts)
  - `make lint-markdown` - Run markdownlint
- **Alternative: Use scripts directly if make times out**:
  - `./scripts/lint-fmt.sh` - Check Go formatting
  - `./scripts/lint-vet.sh` - Run go vet
  - `./scripts/lint-vet.sh pkg` - Run go vet package by package
  - `./scripts/lint-golangci.sh` - Run golangci-lint
  - `./scripts/lint-golangci.sh pkg` - Run golangci-lint package by package
  - `./scripts/lint-all.sh` - Run all linters
  - `./scripts/test-pkg.sh` - Run tests package by package
- **For local development: run `make lint-fix-all` to automatically fix all issues including go vet**
- **Individual fix targets (to avoid timeouts)**:
  - `make lint-fix-fmt` - Fix Go formatting
  - `make lint-fix-golangci` - Fix golangci-lint issues
  - `make lint-fix-markdown` - Fix markdown issues
- **Testing targets**:
  - `make test` - Run all tests at once
  - `make test-pkg` - Run tests package by package (avoids timeouts)
  - `make test-verbose` - Run tests with verbose output
  - `make test-coverage` - Run tests with coverage report
- **Quick linting without go vet: run `make lint` or `make lint-fix`**
- **Run all linting checks: run `make lint-all` (includes go vet + lint)**
- **View all available commands: run `make help`**
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

### Testing Patterns

#### OAuth/External API Testing Pattern

When testing functions that call external APIs, create a testable wrapper function that accepts the API URL as a parameter:

```go
// Main function calls production API
func (s *SlackProcessor) sendViaOAuth(token string, payload []byte) error {
    return s.sendViaOAuthWithURL(token, payload, "https://slack.com/api/chat.postMessage")
}

// Testable wrapper accepts custom URL for mocking
func (s *SlackProcessor) sendViaOAuthWithURL(token string, payload []byte, apiURL string) error {
    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
    // ... rest of implementation
}
```

Test the wrapper function using `httptest.NewServer`:

```go
func TestSlackProcessor_sendViaOAuth(t *testing.T) {
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify headers, method, etc.
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
    }))
    defer mockServer.Close()
    
    err := processor.sendViaOAuthWithURL(token, payload, mockServer.URL)
    // assertions...
}
```

This pattern allows comprehensive testing without requiring real API keys while keeping production code clean.

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
- **CodeRabbit Configuration**: When creating new files or directories, check if they need to be added to `.coderabbit.yaml` configuration to ensure proper code review coverage

### Common Issues & Solutions

#### Test and Build Timeouts
- **Issue**: `make test`, `make lint`, or other make commands timeout after 2 minutes
- **Root Cause**: Go downloading dependencies during test execution
- **Solutions**:
  1. Run `go mod download` first to cache dependencies
  2. Use package-specific targets: `make test-pkg`, `make lint-pkg` 
  3. Run scripts directly: `./scripts/test-pkg.sh`, `./scripts/lint-all.sh`
  4. Test individual packages: `go test ./internal/processor`

#### Conventional Commit Check Failures
- **Issue**: Historical commits don't follow conventional commit format
- **Solution**: Requires maintainer intervention - cannot be fixed by modifying files
- **Note**: Only affects existing commits in PR history, not new commits

#### Markdown Linting Errors
- **Issue**: Multiple markdown linting failures (MD022, MD026, MD032)
- **Solution**: Often fixable with single edit by:
  - Removing trailing punctuation from headings (`:` → ``)
  - Adding blank lines before/after headings and lists
  - Example: Fix `#### Heading:` + `- list item` → `#### Heading` + `\n` + `- list item`

#### Coverage Failures
- **Issue**: CodeCov patch/project failures due to low test coverage
- **Focus Areas for Quick Wins**:
  1. Error handling methods (`Error()`, `Unwrap()`)
  2. Edge cases and fallback scenarios  
  3. Input validation functions
  4. Testable wrapper functions for external APIs
- **Areas to Avoid**: Actual API execution methods (require real credentials)

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

## GitHub CLI Commands

### Label Management
- `gh label create "LABEL_NAME" --description "Description" --color "HEX_COLOR" --repo OWNER/REPO` - Create GitHub labels with descriptions and colors
- `gh label list --repo OWNER/REPO` - List all repository labels

### Secrets Management
- `gh api repos/OWNER/REPO/actions/secrets --method GET` - Check existing GitHub repository secrets
- `gh secret set SECRET_NAME --repo OWNER/REPO` - Add repository secrets via CLI (prompts for secure input)

### Repository Information
- `gh repo view OWNER/REPO --json defaultBranchRef,name,owner,url` - Get repository details in JSON format
- `gh auth status` - Check GitHub authentication status and token scopes

### PR and CI Monitoring
- `gh pr checks` - Check PR CI status and failures for current branch
- `gh pr checks 160` - Check CI status for specific PR number
- `gh run view RUN_ID` - Get detailed GitHub Actions run information and logs
- `gh run view RUN_ID --log-failed` - Get logs for failed jobs only

### CodeRabbit Tools
- `~/bin/coderabbit-fix PR_NUMBER --ai-format` - Generate AI-formatted CodeRabbit issue analysis with detailed fix instructions
- `~/bin/coderabbit-fix PR_NUMBER --dry-run` - Show what would be changed without making actual changes
- `~/bin/coderabbit-fix PR_NUMBER --prioritize` - Group issues by priority for systematic fixing

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

## Claude Models

All claude models come from this site: <https://docs.anthropic.com/en/docs/about-claude/models/overview>

## Roadmap Overview

The project has a defined roadmap divided into milestones:

### v0.0.2 - Claude Support & Queue/Bot MVP (Due: June 15, 2025)

MVP implementation focusing on:
- Claude model support (#13) - Add support for Anthropic Claude 3 models
- Slack processor (#96) - Enable and test Slack processor functionality  
- Microsoft Teams processor (#54) - Implement Teams webhook processor
- Bot mode foundations (#120, #121, #122) - Basic webhook server and event routing ✅ COMPLETED
- Queue mode infrastructure (#63, #86) - Core queue system and basic message queue support

### Q3 2025 - Enhanced Usability

- Basic web UI (#10)
- Conditional logic in prompt templates (#11)
- Prompt testing tool (#12)
- Advanced bot features (context awareness, state management)

### Q4 2025 - Integration & Scale

- External API for integration (#14)
- Performance metrics and analytics (#15)
- Distributed task execution (#16)
- CI/CD platform integrations (#17)
- Enterprise queue providers (AWS SQS, Azure Service Bus, RabbitMQ)

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
  - `--mode`: Operation mode (cron, bot, queue) - default: cron
  - `--config`: Path to configuration file
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

## LangChainGo Integration Reference

### Overview
LangChainGo (github.com/tmc/langchaingo) is a Go implementation of the LangChain framework for building composable AI applications. It provides valuable patterns and abstractions that could enhance CronAI's architecture.

### Key Architecture Patterns from LangChainGo

#### 1. Model Abstraction
- **Unified Interface**: All LLM providers implement a common `Model` interface with methods:
  - `GenerateContent()`: For multi-modal, chat-like interactions
  - `Call()`: Simplified text generation (being deprecated)
- **Provider Support**: Extensive provider coverage including OpenAI, Anthropic, Google AI, AWS Bedrock, Ollama, Hugging Face, and many others
- **Context-Aware**: All operations accept Go context for cancellation and timeout support

#### 2. Composability Through Chains
- **Chain Interface**: Enables sequential operation composition with:
  - Input/output key definitions
  - Memory management between steps
  - Context propagation
- **Execution Patterns**: Supports both synchronous (`Call()`) and asynchronous (`Apply()`) execution
- **Memory Integration**: Built-in support for conversation memory and state persistence

#### 3. Prompt Management
- **Template System**: Rich prompt templating with:
  - Variable substitution
  - Chat-specific prompt templates
  - Few-shot learning support
  - Example selectors for dynamic prompting
- **Message Abstraction**: Structured handling of chat messages and conversations

### Integration Patterns Relevant to CronAI

#### 1. Model Configuration
- Uses functional options pattern for configuration
- Provider-specific options while maintaining common interface
- Example: `llm, err := openai.New(openai.WithAPIKey(key))`

#### 2. Error Handling
- Consistent error propagation through all layers
- Context-based cancellation support
- Graceful degradation patterns

#### 3. Extensibility
- Clear interface definitions for adding new providers
- Modular design allows selective component usage
- Plugin-like architecture for tools and integrations

### Potential Applications for CronAI

1. **Enhanced Model Support**: Could adopt LangChainGo's model abstraction pattern for more unified provider handling
2. **Chain-Based Workflows**: Implement complex prompt sequences using chain patterns
3. **Memory Integration**: Add conversation history for bot mode
4. **Tool Integration**: Leverage tools pattern for extending CronAI capabilities
5. **Prompt Templates**: Adopt the sophisticated prompt templating system

### Best Practices from LangChainGo

1. **Interface-First Design**: Define clear interfaces before implementation
2. **Context Propagation**: Always pass context through the call stack
3. **Modular Architecture**: Keep components loosely coupled
4. **Provider Abstraction**: Hide provider-specific details behind common interfaces
5. **Functional Options**: Use for flexible, backward-compatible configuration

### Key Differences from CronAI's Current Approach

1. **Abstraction Level**: LangChainGo provides higher-level abstractions (chains, agents) vs CronAI's direct model calls
2. **Prompt Management**: More sophisticated templating vs CronAI's file-based approach
3. **Composability**: Built for complex workflows vs CronAI's single-prompt execution
4. **Memory/State**: Built-in conversation memory vs CronAI's stateless execution

### Recommended Considerations

If integrating LangChainGo patterns into CronAI:
- Start with adopting the model interface pattern for better provider abstraction
- Consider implementing a simplified chain pattern for multi-step workflows
- Evaluate the prompt template system for more dynamic prompt generation
- Look into memory components for bot mode context management
