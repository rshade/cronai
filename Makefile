.PHONY: all build test test-coverage coverage-report clean lint lint-fix format run changelog

# Default target
all: build

# Build the application
build:
	go build -o cronai ./cmd/cronai

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Clean build artifacts
clean:
	rm -f cronai coverage.out

# View coverage report in browser (for local development)
coverage-report: test-coverage
	go tool cover -html=coverage.out

# Run linter (strict mode for CI)
lint:
	@echo "Running gofmt check..."
	@! gofmt -d . 2>&1 | grep -q '^' || (echo "Code not formatted. Run 'make lint-fix' to fix."; exit 1)
	@echo "Running go vet..."
	go vet ./...
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi
	@echo "Running markdownlint..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore node_modules; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

# Run linter with automatic fixes (for local development)
lint-fix:
	@echo "Running gofmt with fix..."
	@gofmt -w .
	@echo "Running go vet..."
	go vet ./...
	@echo "Running golangci-lint with fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi
	@echo "Running markdownlint with fix..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint . --ignore node_modules --fix; \
	else \
		echo "markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

# Format all Go files
format:
	@echo "Formatting all Go files..."
	go fmt ./...

# Run the application
run:
	go run ./cmd/cronai/main.go start

# Install the application
install:
	go install ./cmd/cronai

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