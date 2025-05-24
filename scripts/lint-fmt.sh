#!/bin/bash
# Check Go formatting

echo "Running gofmt check..."
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
    echo "Code not formatted. The following files need formatting:"
    echo "$unformatted"
    echo "Run 'gofmt -w .' to fix."
    exit 1
else
    echo "All Go files are properly formatted."
fi