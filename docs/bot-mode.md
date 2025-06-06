# Bot Mode Documentation

Bot mode enables CronAI to act as a GitHub webhook handler, processing events in real-time and generating AI-powered responses.

## Overview

When running in bot mode, CronAI starts an HTTP server that listens for GitHub webhook events. It processes these events through AI models and can trigger various actions based on the responses.

## Configuration

### Environment Variables

- `CRONAI_BOT_PORT`: The port to run the webhook server on (default: 8080)
- `GITHUB_WEBHOOK_SECRET`: Secret for validating GitHub webhook signatures (minimum 8 characters recommended for security)
- `CRONAI_DEFAULT_MODEL`: The AI model to use for processing events (default: openai)
- `CRONAI_BOT_PROCESSOR`: Optional processor to send AI responses to (e.g., `console`, `file:/path/to/log`)
- `CRONAI_RATE_LIMIT_DEFAULT`: Default rate limit for requests per minute (default: 100)
- `CRONAI_RATE_LIMIT_BUCKET_SIZE`: Maximum burst capacity in tokens (default: 100)
- `CRONAI_RATE_LIMIT_REFILL_RATE`: Token refill rate per minute (default: 1)
- `CRONAI_WEBHOOK_SECRET_MIN_LENGTH`: Minimum required length for webhook secret (default: 8)

### Rate Limiting Configuration

Bot mode includes built-in rate limiting using a token bucket algorithm to protect against abuse:

- **Default Rate Limit**: 100 requests per minute
- **Bucket Size**: 100 tokens (maximum burst capacity)
- **Refill Rate**: 1 token per minute (sustained throughput)
- **Algorithm**: Token bucket with thread-safe implementation

The rate limiter is applied at the webhook endpoint level and automatically refills tokens based on elapsed time intervals. This provides protection against webhook spam while allowing legitimate GitHub webhook traffic to flow through.

### Starting Bot Mode

```bash
# Basic start
cronai start --mode bot

# With custom configuration
export CRONAI_BOT_PORT=9090
export GITHUB_WEBHOOK_SECRET=your-secret-here
export CRONAI_DEFAULT_MODEL=claude
export CRONAI_BOT_PROCESSOR=console
cronai start --mode bot --config cronai.config
```

## GitHub Webhook Setup

1. Go to your GitHub repository settings
2. Navigate to Webhooks
3. Click "Add webhook"
4. Configure the webhook:

   - **Payload URL**: `https://your-server:8080/webhook`
   - **Content type**: `application/json`
   - **Secret**: Your `GITHUB_WEBHOOK_SECRET` value
   - **Events**: Select the events you want to receive

## Supported Events

Bot mode currently supports the following GitHub events:

### Issues Events

- `opened`: New issue created
- `closed`: Issue closed
- `reopened`: Issue reopened
- `edited`: Issue title or body edited
- `labeled`: Label added to issue
- `unlabeled`: Label removed from issue

### Pull Request Events

- `opened`: New pull request created
- `closed`: Pull request closed
- `reopened`: Pull request reopened
- `synchronize`: New commits pushed to pull request
- `ready_for_review`: Draft PR marked as ready

### Push Events

Triggered when commits are pushed to any branch

### Release Events

- `created`: New release created
- `published`: Release published
- `edited`: Release edited
- `deleted`: Release deleted

## Event Processing

When an event is received:

1. The webhook signature is verified (if secret is configured)
2. The event is routed to the appropriate handler
3. Bot events are filtered out by default
4. An AI prompt is generated based on the event context
5. The configured AI model processes the prompt
6. The response can be sent to a configured processor

## Example Prompts Generated

### Issue Opened

```text
A GitHub issues event occurred with action 'opened'. Issue #123 'Bug in login system' was opened by johndoe in myorg/myrepo. Issue description: Users are unable to login with valid credentials.

Please provide a brief analysis of this event and suggest any actions that might be needed.
```

### Pull Request Opened

```text
A GitHub pull_request event occurred with action 'opened'. Pull Request #456 'Fix login bug' was opened by janedoe in myorg/myrepo.

Please provide a brief analysis of this event and suggest any actions that might be needed.
```

## Security Considerations

1. **Always use webhook secrets**: Configure `GITHUB_WEBHOOK_SECRET` to validate incoming webhooks
2. **Use HTTPS**: In production, run behind a reverse proxy with SSL
3. **Limit exposed ports**: Only expose the webhook port to GitHub's IP ranges
4. **Monitor logs**: Check logs regularly for suspicious activity
5. **Rate limit monitoring**: Monitor rate-limit logs and adjust limits as needed; consider lowering the default 100 requests/minute for high-security environments

## Health Check

The bot mode server provides a health check endpoint:

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "healthy",
  "mode": "bot"
}
```

## Customization

### Adding Custom Event Handlers

To add support for additional GitHub events, you can extend the router in the source code:

1. Create a new handler implementing the `EventHandler` interface
2. Register it in the bot service initialization
3. Add any custom event processing logic

### Custom Prompt Generation

The current implementation generates basic prompts. For more sophisticated prompt generation:

1. Modify the `generatePrompt` method in `internal/bot/service.go`
2. Add template-based prompt generation
3. Include more context from the event payload

## Troubleshooting

### Webhook Not Receiving Events

- Verify the webhook URL is accessible from the internet
- Check the webhook secret matches between GitHub and your configuration
- Look for errors in the GitHub webhook delivery logs

### Signature Validation Failures

- Ensure the `GITHUB_WEBHOOK_SECRET` environment variable is set correctly
- Check that the secret in GitHub matches exactly (no extra spaces)

### Event Processing Errors

- Check the logs for detailed error messages
- Verify the AI model API keys are configured
- Ensure the model specified in `CRONAI_DEFAULT_MODEL` is valid

## Future Enhancements

- Support for more GitHub event types
- Configurable prompt templates
- Event-specific AI model selection
- Response caching for similar events
- Integration with GitHub Actions
- Support for other webhook providers (GitLab, Bitbucket, etc.)
