# CronAI Scripts

This directory contains utility scripts for development tasks. These scripts can be used directly if the Makefile commands timeout or if you prefer a more direct approach.

## Linting Scripts

- **lint-fmt.sh** - Check Go code formatting

  ```bash
  ./scripts/lint-fmt.sh
  ```

- **lint-vet.sh** - Run go vet

  ```bash
  ./scripts/lint-vet.sh      # Run on all packages
  ./scripts/lint-vet.sh pkg  # Run package by package (avoids timeouts)
  ```

- **lint-golangci.sh** - Run golangci-lint

  ```bash
  ./scripts/lint-golangci.sh      # Run on all packages
  ./scripts/lint-golangci.sh pkg  # Run directory by directory (avoids timeouts)
  ```

- **lint-all.sh** - Run all linters

  ```bash
  ./scripts/lint-all.sh
  ```

## Testing Scripts

- **test-pkg.sh** - Run tests package by package

  ```bash
  ./scripts/test-pkg.sh
  ```

## Other Scripts

- **generate_changelog.sh** - Generate changelog from conventional commits

  ```bash
  ./scripts/generate_changelog.sh [FROM] [TO]
  ```

## Usage

All scripts are executable and can be run directly:

```bash
cd /path/to/cronai
./scripts/lint-fmt.sh
```

These scripts are particularly useful when:

- Make commands are timing out
- You want more control over the execution
- You're debugging specific issues
- You're running in CI/CD environments with timeout constraints
