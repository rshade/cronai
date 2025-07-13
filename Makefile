.PHONY: all help build test test-pkg test-verbose test-coverage coverage coverage-report \
	clean clean_branches \
	vet vet-pkg lint lint-all lint-fix lint-fix-all \
	lint-fmt lint-golangci lint-golangci-pkg lint-markdown \
	lint-fix-fmt lint-fix-golangci lint-fix-markdown \
	format run run-task list changelog release setup install

# Default target
all: build

# Help target to show available commands
help:
	@echo "Available targets:"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run all tests"
	@echo "  make test-pkg       - Run tests package by package (avoids timeouts)"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make coverage       - Show coverage percentage"
	@echo "  make vet            - Run go vet on all packages"
	@echo "  make vet-pkg        - Run go vet package by package (avoids timeouts)"
	@echo "  make lint-fmt       - Check Go formatting"
	@echo "  make lint-golangci  - Run golangci-lint"
	@echo "  make lint-golangci-pkg - Run golangci-lint package by package (avoids timeouts)"
	@echo "  make lint-markdown  - Run markdownlint"
	@echo "  make lint           - Run all linters"
	@echo "  make lint-fix-fmt   - Fix Go formatting"
	@echo "  make lint-fix-golangci - Fix golangci-lint issues"
	@echo "  make lint-fix-markdown - Fix markdown issues"
	@echo "  make lint-fix       - Fix all linting issues"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install the application"
	@echo "  make run            - Run the application"
	@echo "  make setup          - Setup development environment"

# Get version information
VERSION ?= $(shell git describe --tags --always --dirty --match='v*' 2>/dev/null || echo "v0.0.0-dev")
BUILD_VERSION := $(shell echo $(VERSION) | sed 's/^v//')
LDFLAGS := -s -w -X github.com/rshade/cronai/cmd/cronai/cmd.Version=$(BUILD_VERSION)

# Build the application
build:
	go build -ldflags="$(LDFLAGS)" -o cronai ./cmd/cronai

# Run tests
test:
	go test ./...

# Run tests package by package (useful for debugging and avoiding timeouts)
test-pkg:
	@./scripts/test-pkg.sh

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate coverage percentage
coverage: test-coverage
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}'

# Clean build artifacts
clean:
	rm -f cronai coverage.out

# Clean up local branches that no longer have remote tracking branches
clean_branches:
	@echo "Fetching latest changes and pruning deleted remote branches..."
	@git fetch --prune
	@echo "Cleaning up local branches that no longer exist on remote..."
	@git branch -vv | grep ': gone]' | awk '{print $$1}' | while read branch; do \
		echo "Deleting branch: $$branch"; \
		git branch -d "$$branch" 2>/dev/null || \
		(echo "Warning: Branch $$branch has unmerged changes. Use 'git branch -D $$branch' to force delete." >&2); \
	done
	@echo "Branch cleanup complete."

# View coverage report in browser (for local development)
coverage-report: test-coverage
	go tool cover -html=coverage.out

# Get all Go packages
GO_PACKAGES := $(shell go list ./... 2>/dev/null || echo "")

# Run go vet (separated for CI timeout management)
vet:
	@echo "Running go vet..."
	go vet ./...

# Run go vet package by package (slower but avoids timeouts)
vet-pkg:
	@./scripts/lint-vet.sh pkg

# Individual lint targets for better control and avoiding timeouts
lint-fmt:
	@./scripts/lint-fmt.sh

lint-golangci:
	@./scripts/lint-golangci.sh

# Run golangci-lint package by package (slower but avoids timeouts)
lint-golangci-pkg:
	@./scripts/lint-golangci.sh pkg

lint-markdown:
	@echo "Running markdownlint..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore-path .markdownlintignore; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

# Run linter (strict mode for CI, excludes go vet)
# This runs all linters sequentially - use individual targets to avoid timeouts
lint: lint-fmt lint-golangci lint-markdown

# Run all linting including go vet (for convenience)
lint-all: vet lint

# Individual fix targets for better control
lint-fix-fmt:
	@echo "Running gofmt with fix..."
	@gofmt -w .

lint-fix-golangci:
	@echo "Running golangci-lint with fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Note: go vet should be run separately with 'make vet'"; \
	fi

lint-fix-markdown:
	@echo "Running markdownlint with fix..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore node_modules --fix; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

# Run linter with automatic fixes (for local development)
# This runs all fixers sequentially - use individual targets to avoid timeouts
lint-fix: lint-fix-fmt lint-fix-golangci lint-fix-markdown

# Run all linting with fixes including go vet (for convenience)
lint-fix-all: vet lint-fix

# Format all Go files
format:
	@echo "Formatting all Go files..."
	go fmt ./...

# Run the application
run:
	go run ./cmd/cronai/main.go start

# Install the application
install:
	go install -ldflags="$(LDFLAGS)" ./cmd/cronai

# Run a specific task immediately
run-task:
	@if [ -z "$(MODEL)" ] || [ -z "$(PROMPT)" ] || [ -z "$(PROCESSOR)" ]; then \
		echo "Usage: make run-task MODEL=claude PROMPT=product_manager PROCESSOR=slack-pm-channel"; \
		exit 1; \
	fi
	go run ./cmd/cronai/main.go run --model $(MODEL) --prompt $(PROMPT) --processor $(PROCESSOR)

# List tasks from config
list:
	go run ./cmd/cronai/main.go list

# Generate changelog from conventional commits
changelog:
	@echo "Generating changelog..."
	@./scripts/generate_changelog.sh $(FROM) $(TO)
	@echo "Changelog generated at CHANGELOG.md"

# Release the application using goreleaser
release:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "v0.0.0-dev" ]; then \
		echo "Error: VERSION must be set to a valid version tag for releases"; \
		echo "Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	goreleaser release --clean

# Setup development environment
setup:
	go mod download
	go mod tidy
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please edit it with your API keys."; \
	fi
	@echo "Installing markdownlint..."
	@npm install -g markdownlint-cli || echo "Warning: Failed to install markdownlint-cli. Make sure npm is installed."