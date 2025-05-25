# Microsoft Teams Integration Testing Guide

This guide explains how to set up a free Microsoft Teams account and configure it for integration testing with CronAI.

## Option 1: Microsoft 365 Developer Program (Recommended)

The Microsoft 365 Developer Program provides a free, renewable 90-day subscription perfect for development and testing.

### Setup Steps

1. **Sign up for the Developer Program**
   - Visit [Microsoft 365 Developer Program](https://developer.microsoft.com/microsoft-365/dev-program)
   - Click "Join now"
   - Sign in with a Microsoft account (or create one)
   - Complete the registration form

2. **Set up your developer subscription**
   - After joining, go to your dashboard
   - Click "Set up E5 subscription"
   - Choose a domain name (e.g., yourcompany.onmicrosoft.com)
   - Create an admin account
   - Wait for provisioning (usually takes 2-5 minutes)

3. **Access Microsoft Teams**
   - Go to [Microsoft Teams](https://teams.microsoft.com)
   - Sign in with your developer account
   - Create a new team for testing

4. **Create an Incoming Webhook**
   - In Teams, go to your test team
   - Click the "..." menu next to the channel name
   - Select "Connectors" → "Incoming Webhook"
   - Click "Configure"
   - Give it a name (e.g., "CronAI Integration Test")
   - Click "Create"
   - **Copy the webhook URL** - you'll need this for testing

### Benefits

- Full Microsoft 365 E5 features
- Auto-renews every 90 days with active development
- Multiple test users and teams
- Full admin capabilities

## Option 2: Teams Free Account

A simpler option with limited features but sufficient for basic webhook testing.

### Free Account Setup Steps

1. **Sign up for Teams Free**
   - Visit [Microsoft Teams Free](https://www.microsoft.com/microsoft-teams/free)
   - Click "Sign up for free"
   - Use your email address to create an account

2. **Create a team and channel**
   - After signing in, create a new team
   - Add a channel for testing

3. **Add Incoming Webhook**
   - Follow the same webhook creation steps as Option 1

### Limitations

- Limited to 100 team members
- Some advanced features unavailable
- No admin center access

## Option 3: Mock Testing with Webhook Services

For unit testing without a real Teams account, use webhook testing services.

### Recommended Services

1. **Webhook.site**
   - Visit [Webhook.site](https://webhook.site)
   - Copy your unique URL
   - Use as `TEAMS_WEBHOOK_URL` for testing
   - View all requests and payloads in real-time

2. **RequestBin**
   - Visit [RequestBin](https://requestbin.com)
   - Create a new bin
   - Use the bin URL for testing

## Setting Up Integration Tests

### 1. Environment Variables

Set the following environment variables:

```bash
# For real Teams webhook
export TEAMS_WEBHOOK_URL="https://outlook.office.com/webhook/YOUR_WEBHOOK_ID"

# Enable integration tests
export RUN_INTEGRATION_TESTS=1

# For mock testing (optional)
export USE_MOCK_HTTP=1
```

### 2. Running Tests

```bash
# Run Teams integration tests only
go test -v ./internal/integration/ -run TestTeamsIntegration

# Run all integration tests
make test-integration
```

### 3. CI/CD Configuration

For GitHub Actions, add the webhook URL as a secret:

1. Go to your repository settings
2. Navigate to Secrets and variables → Actions
3. Add a new secret named `TEAMS_WEBHOOK_URL`
4. Paste your webhook URL

Update your workflow file:

```yaml
- name: Run Integration Tests
  env:
    TEAMS_WEBHOOK_URL: ${{ secrets.TEAMS_WEBHOOK_URL }}
    RUN_INTEGRATION_TESTS: 1
  run: go test -v ./internal/integration/
```

## Teams Webhook Payload Format

CronAI uses the MessageCard format for Teams webhooks:

```json
{
  "@type": "MessageCard",
  "@context": "https://schema.org/extensions",
  "themeColor": "0078D4",
  "summary": "Message from CronAI",
  "sections": [{
    "activityTitle": "CronAI Notification",
    "activitySubtitle": "Automated AI Response",
    "facts": [
      {"name": "Model", "value": "claude-3-haiku"},
      {"name": "Prompt", "value": "daily-report"},
      {"name": "Time", "value": "2024-01-15 10:30:00"}
    ],
    "markdown": true,
    "text": "Your AI-generated content here"
  }]
}
```

## Troubleshooting

### Common Issues

1. **"Invalid webhook URL"**
   - Ensure URL contains `outlook.office.com` or `outlook.office365.com`
   - Check for typos or extra spaces

2. **"400 Bad Request" from Teams**
   - Validate JSON payload format
   - Ensure payload is under 25KB
   - Check for required fields (@type, @context, summary)

3. **"Webhook not found"**
   - Webhook may have been deleted in Teams
   - Regenerate the webhook in Teams settings

4. **Rate Limiting**
   - Teams webhooks have rate limits
   - Add delays between test calls
   - Use mock services for high-volume testing

## Best Practices

1. **Use dedicated test channels**
   - Create specific channels for integration testing
   - Name them clearly (e.g., "cronai-integration-tests")

2. **Implement retry logic**
   - Add exponential backoff for failed requests
   - Handle temporary network issues

3. **Log webhook responses**
   - Capture response codes and error messages
   - Useful for debugging integration issues

4. **Test different scenarios**
   - Large messages (near 25KB limit)
   - Special characters and markdown
   - Different template formats

5. **Clean up test data**
   - Consider message retention in test channels
   - Document test message patterns for easy identification

## Security Considerations

1. **Never commit webhook URLs**
   - Always use environment variables
   - Add `.env` to `.gitignore`

2. **Rotate webhooks regularly**
   - Especially after team member changes
   - Update CI/CD secrets when rotating

3. **Use separate webhooks for each environment**
   - Development, staging, and production
   - Never use production webhooks in tests

4. **Monitor webhook usage**
   - Check Teams audit logs
   - Set up alerts for unusual activity
