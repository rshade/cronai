.PHONY: all build test clean lint run

# Default target
all: build

# Build the application
build:
	go build -o cronai ./cmd/cronai

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f cronai

# Run linter
lint:
	@echo "Running gofmt check..."
	@test -z "$$(gofmt -l .)" || (echo "The following files need formatting with gofmt:"; gofmt -l . && exit 1)
	@echo "Running go vet..."
	go vet ./...
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

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

# Setup development environment
setup:
	go mod download
	go mod tidy
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please edit it with your API keys."; \
	fi