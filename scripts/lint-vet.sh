#!/bin/bash
# Run go vet

echo "Running go vet..."
if [ "$1" == "pkg" ]; then
    # Run package by package
    echo "Running go vet package by package..."
    failed=0
    for pkg in $(go list ./... 2>/dev/null); do
        echo "Vetting $pkg..."
        if ! go vet "$pkg"; then
            failed=1
        fi
    done
    exit $failed
else
    # Run all at once
    go vet ./...
fi