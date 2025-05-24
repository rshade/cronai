.PHONY: all build test test-coverage coverage coverage-report clean clean_branches vet lint lint-all lint-fix lint-fix-all format run changelog release

# Default target
all: build

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

# Run go vet (separated for CI timeout management)
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linter (strict mode for CI, excludes go vet)
lint:
	@echo "Running gofmt check..."
	@! gofmt -d . 2>&1 | grep -q '^' || (echo "Code not formatted. Run 'make lint-fix' to fix."; exit 1)
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Note: go vet should be run separately with 'make vet'"; \
	fi
	@echo "Running markdownlint..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore node_modules; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

# Run all linting including go vet (for convenience)
lint-all: vet lint

# Run linter with automatic fixes (for local development)
lint-fix:
	@echo "Running gofmt with fix..."
	@gofmt -w .
	@echo "Running golangci-lint with fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Note: go vet should be run separately with 'make vet'"; \
	fi
	@echo "Running markdownlint with fix..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore node_modules --fix; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

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