#!/bin/bash
# Run golangci-lint

if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

echo "Running golangci-lint..."
if [ "$1" == "pkg" ]; then
    # Run directory by directory
    echo "Running golangci-lint directory by directory..."
    failed=0
    for dir in $(find . -name "*.go" -not -path "./vendor/*" -not -path "./node_modules/*" | xargs -n1 dirname | sort -u); do
        echo "Linting $dir..."
        if ! golangci-lint run "$dir/..."; then
            failed=1
        fi
    done
    exit $failed
else
    # Run all at once
    golangci-lint run
fi