# Integration Tests for CronAI

This directory contains integration tests for the CronAI application, focusing on testing the full end-to-end pipeline from prompt loading through model execution to response processing.

## Test Coverage

The integration tests cover:

1. **OpenAI Model Integration**
   - Model initialization and parameter handling
   - Prompt execution
   - Error handling and recovery
   - Rate limit handling

2. **File Processor**
   - File output creation
   - Directory creation if not exists
   - File permissions handling
   - Template-based formatting

3. **GitHub Processor**
   - GitHub API integration
   - Issue comment creation
   - Authentication handling
   - Template-based formatting

## Running the Tests

### Basic Test Run

To run all integration tests in mock mode (no real API calls):

```bash
go test -v ./internal/integration/
```

### With Real API Keys

To run tests with real API keys:

```bash
# Set your API keys
export OPENAI_API_KEY=your_openai_key
export GITHUB_TOKEN=your_github_token

# Run tests with real APIs
export RUN_INTEGRATION_TESTS=1
go test -v ./internal/integration/
```

### Continuous Integration

The integration tests run automatically in GitHub Actions:

- In pull requests: Tests run in mock mode (no real API calls)
- On the main branch: Tests run with real API calls using GitHub secrets
- Manual runs: Can be triggered with or without real API calls

The GitHub Actions workflow can be found at `.github/workflows/integration-tests.yml`.

## Test Configuration

The integration tests use the following environment variables:

- `GO_TEST=1`: Enables test mode (mocks external API calls)
- `RUN_INTEGRATION_TESTS=1`: Enables tests that use real API calls
- `OPENAI_API_KEY`: API key for OpenAI
- `GITHUB_TOKEN`: API token for GitHub
- `LOGS_DIRECTORY`: Directory for file processor outputs

## Test Files

- `pipeline_test.go`: Tests the full end-to-end pipeline from prompt to processor

## Adding New Tests

When adding new integration tests:

1. Focus on testing the integration between components, not individual functionality
2. Use temporary directories for test files
3. Mock external dependencies in CI environments
4. Include both happy path and error handling tests
5. Clean up resources after tests complete
