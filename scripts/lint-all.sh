#!/bin/bash
# Run all linters

echo "=== Running all linters ==="
echo

# Run format check
echo "1. Checking Go formatting..."
./scripts/lint-fmt.sh
echo

# Run go vet
echo "2. Running go vet..."
go vet ./...
echo

# Run golangci-lint if available
echo "3. Running golangci-lint..."
if command -v golangci-lint >/dev/null 2>&1; then
    golangci-lint run
else
    echo "golangci-lint not installed. Skipping..."
fi
echo

# Run markdownlint if available
echo "4. Running markdownlint..."
if command -v markdownlint >/dev/null 2>&1; then
    markdownlint . --ignore node_modules
else
    echo "markdownlint not installed. Skipping..."
fi

echo "=== Linting complete ==="