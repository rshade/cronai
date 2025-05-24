#!/bin/bash
# Run tests package by package

echo "Running tests package by package..."
failed=0
for pkg in $(go list ./... 2>/dev/null); do
    echo "Testing $pkg..."
    if ! go test "$pkg"; then
        failed=1
    fi
done
exit $failed