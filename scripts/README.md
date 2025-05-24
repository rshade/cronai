# Scripts

This directory contains various scripts used for development, testing, and deployment.

## Development Scripts

### Build

```bash
./build.sh
```

This script builds the project for the current platform.

### Test

```bash
./test.sh
```

This script runs all tests in the project.

### Lint

```bash
./lint.sh
```

This script runs all linters on the project.

### Format

```bash
./format.sh
```

This script formats all code in the project.

### Clean

```bash
./clean.sh
```

This script cleans all build artifacts.

### Generate

```bash
./generate.sh
```

This script generates code from templates and protobuf definitions.

## Common Issues

### Make Commands

- Make commands are timing out
- Make commands are failing with permission errors
- Make commands are not finding dependencies

### Solutions

1. Ensure you have all required dependencies installed
2. Check that you have the correct permissions
3. Try running the commands with sudo if needed
4. Check the Makefile for any custom requirements

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
