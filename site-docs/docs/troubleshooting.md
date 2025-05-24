---
id: troubleshooting
title: Troubleshooting Guide
sidebar_label: Troubleshooting
---

# Troubleshooting Guide (Coming Soon)

> ⚠️ **Note**: This troubleshooting guide is a placeholder for the upcoming comprehensive guide that will be available in future releases. For the MVP release, we've included some basic troubleshooting tips below.

## Basic Troubleshooting

### Installation Issues

If you encounter issues during installation:

1. Make sure you have Go 1.21 or higher installed

   ```bash
   go version
   ```text

2. Check that your GOPATH is correctly configured

   ```bash
   go env GOPATH
   ```text

3. Try reinstalling from source

   ```bash
   git clone https://github.com/rshade/cronai.git
   cd cronai
   go build -o cronai ./cmd/cronai
   ```text

### Configuration Problems

If your configuration file is not working:

1. Validate your configuration format

   ```text
   timestamp model prompt response_processor [variables] [model_params:...]
   ```text

2. Use the validate command to check your configuration

   ```bash
   cronai validate --config /path/to/cronai.config
   ```text

3. Check that your cron schedule format is valid

   ```bash
   # Standard cron format: minute hour day-of-month month day-of-week
   # Example: 0 9 * * 1 (run at 9 AM every Monday)
   ```text

### API Keys

If you're having issues with model execution:

1. Check that you have set the required API keys in your environment

   ```bash
   # For OpenAI
   export OPENAI_API_KEY=your_openai_key
   
   # For Claude
   export ANTHROPIC_API_KEY=your_anthropic_key
   
   # For Gemini
   export GOOGLE_API_KEY=your_google_key
   ```text

2. Verify that your API keys are valid and have the required permissions

3. For GitHub processors, ensure your GITHUB_TOKEN is set and has appropriate permissions

   ```bash
   export GITHUB_TOKEN=your_github_token
   ```text

### Prompt Issues

If your prompts are not working as expected:

1. Check that your prompt file exists in the cron_prompts directory

   ```bash
   ls -la cron_prompts/your_prompt.md
   ```text

2. Make sure your variable placeholders are correctly formatted

   ```text
   {{variable_name}}
   ```text

3. Preview your prompt with variables to see how it will be rendered

   ```bash
   cronai prompt preview your_prompt --vars "var1=value1,var2=value2"
   ```text

### Running as a Service

If you're having issues with the systemd service:

1. Check the service status

   ```bash
   sudo systemctl status cronai
   ```text

2. View the logs for error messages

   ```bash
   sudo journalctl -u cronai -f
   ```text

3. Verify that your environment file is correctly configured

   ```bash
   cat /etc/cronai/.env
   ```text

## Logging

CronAI includes basic logging capabilities in the MVP. Logs can help identify issues:

```bash
# Set more verbose logging
export LOG_LEVEL=debug

# Run with verbose output
cronai start --config /path/to/config
```text

## Getting More Help

For the MVP release, if you encounter issues not covered in this guide:

1. Check the GitHub repository for known issues: <https://github.com/rshade/cronai/issues>
2. File a new issue with detailed information about your problem

A comprehensive troubleshooting guide with detailed error handling, diagnostics, and solutions will be available in future releases.
