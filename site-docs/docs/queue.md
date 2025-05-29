---
sidebar_position: 5
---

# Queue Integration

## Overview

CronAI's queue integration feature (available in v0.0.2+) provides core infrastructure for consuming tasks from message queues. This enables dynamic task distribution, real-time processing, and seamless integration with external systems.

## Architecture

The queue system is built with a plugin architecture that allows for easy extension with new queue providers. The core components include:

- **Consumer Interface**: Standard interface that all queue providers must implement
- **Plugin Registry**: Dynamic registration system for queue providers
- **Message Parser**: Handles both minimal and comprehensive message formats
- **Task Processor**: Bridges queue messages with the existing model execution system
- **Coordinator**: Manages multiple queue consumers concurrently
- **Retry Mechanisms**: Configurable retry policies for failed messages

## Configuration

Queue consumers are configured alongside cron tasks in the `cronai.config` file.

### Syntax

```text
queue <name> <type> <connection> <queue> [options]
```

### Parameters

- **name**: Unique identifier for the queue consumer (e.g., `main-queue`, `priority-tasks`)
- **type**: Queue provider type (e.g., `rabbitmq`, `sqs`, `servicebus`, `pubsub`)
- **connection**: Connection string or URL for the queue service
- **queue**: Queue name, topic, or subscription identifier
- **options**: Comma-separated key-value pairs for configuration

### Standard Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| retry_limit | int | 3 | Maximum number of retry attempts |
| retry_delay | duration | 5s | Delay between retry attempts |

### Provider-Specific Options

Each queue provider may support additional options. These will be documented with each provider implementation.

## Message Formats

The queue system supports two message formats to accommodate different use cases:

### Comprehensive Format

The comprehensive format includes all task details in the message body:

```json
{
  "model": "openai",
  "prompt": "weekly_report",
  "processor": "email-team@example.com",
  "variables": {
    "week": "2024-W01",
    "project": "CronAI",
    "department": "Engineering"
  },
  "is_inline": false
}
```

Fields:

- **model** (required): AI model to use (`openai`, `claude`, `gemini`)
- **prompt** (required): Prompt file name or inline prompt content
- **processor** (required): Response processor configuration
- **variables** (optional): Key-value pairs for variable substitution
- **is_inline** (optional): Whether `prompt` contains inline content (default: false)

### Minimal Format

The minimal format includes only variables in the message body, with task configuration provided via message attributes/metadata:

```json
{
  "variables": {
    "date": "2024-01-01",
    "environment": "production",
    "severity": "high"
  }
}
```

Required message attributes:

- **model**: AI model to use
- **prompt**: Prompt file name
- **processor**: Response processor configuration

## Inline Prompts

Queue messages can include inline prompts for dynamic content generation:

```json
{
  "model": "claude",
  "prompt": "Analyze the following metrics and provide recommendations:\n\nCPU Usage: {{cpu_usage}}%\nMemory Usage: {{memory_usage}}%\nDisk Usage: {{disk_usage}}%\n\nFocus on optimization strategies.",
  "processor": "slack-ops-alerts",
  "variables": {
    "cpu_usage": "85",
    "memory_usage": "72",
    "disk_usage": "45"
  },
  "is_inline": true
}
```

When `is_inline` is true, the `prompt` field contains the actual prompt content rather than a file reference.

## Retry Policies

The queue system includes built-in retry mechanisms for handling transient failures:

### Exponential Backoff

Default retry policy with exponentially increasing delays:

- First retry: 1 second
- Second retry: 2 seconds
- Third retry: 4 seconds
- Maximum delay capped at 30 seconds

### Linear Retry

Fixed delay between retry attempts:

```text
queue tasks sqs https://sqs.region.amazonaws.com/account/queue queue-name retry_delay=10s
```

### No Retry

Disable retries for specific consumers:

```text
queue critical rabbitmq amqp://localhost:5672 critical-tasks retry_limit=0
```

## Error Handling

The queue system provides comprehensive error handling:

1. **Parse Errors**: Invalid message format results in immediate rejection without retry
2. **Validation Errors**: Invalid task configuration results in rejection without retry
3. **Processing Errors**: Transient failures trigger retry based on policy
4. **Connection Errors**: Logged and may trigger consumer restart

## Integration with Existing Features

Queue tasks integrate seamlessly with existing CronAI features:

- **Prompt Management**: File-based prompts work identically for queue and cron tasks
- **Variable Substitution**: Same variable syntax and special variables
- **Response Processors**: All processors available for both queue and cron tasks
- **Model Parameters**: Can be specified in message attributes or variables

## Examples

### Basic Queue Configuration

```text
# Single queue consumer
queue main rabbitmq amqp://guest:guest@localhost:5672 tasks
```

### Multiple Queue Consumers

```text
# High-priority tasks with no retry
queue priority sqs https://sqs.us-east-1.amazonaws.com/123/priority priority-queue retry_limit=0

# Standard tasks with custom retry
queue standard sqs https://sqs.us-east-1.amazonaws.com/123/standard standard-queue retry_limit=5,retry_delay=30s

# Batch processing queue
queue batch rabbitmq amqp://localhost:5672 batch-tasks retry_limit=10,retry_delay=1m
```

### Mixed Configuration

```text
# Cron tasks
0 8 * * * openai daily_summary file-/var/log/summary.log
0 */4 * * * claude system_check console

# Queue consumers
queue realtime rabbitmq amqp://localhost:5672 realtime-tasks retry_delay=1s
queue batch servicebus Endpoint=sb://namespace.servicebus.windows.net/ batch-topic
```

## Security Considerations

When implementing queue consumers:

1. **Connection Security**: Use encrypted connections (AMQPS, HTTPS) in production
2. **Authentication**: Use strong authentication for queue services
3. **Message Validation**: All messages are validated before processing
4. **Secret Management**: Never include API keys or secrets in messages
5. **Access Control**: Limit queue access to authorized services only

## Performance Considerations

1. **Concurrent Consumers**: Each queue consumer runs in its own goroutine
2. **Message Batching**: Future providers may support batch message processing
3. **Connection Pooling**: Providers should implement connection pooling
4. **Resource Limits**: Configure appropriate retry limits to prevent resource exhaustion

## Monitoring and Observability

The queue system includes comprehensive logging:

- Consumer lifecycle events (start, stop, connect, disconnect)
- Message processing (received, processed, acknowledged, rejected)
- Error conditions with context
- Performance metrics (processing duration)

## Future Enhancements

Planned improvements for the queue system:

1. **Provider Implementations**: RabbitMQ, AWS SQS, Azure Service Bus, Google Pub/Sub
2. **Dead Letter Queues**: Automatic handling of permanently failed messages
3. **Message Batching**: Process multiple messages in a single AI call
4. **Priority Queues**: Support for message priorities
5. **Metrics Export**: Prometheus/OpenTelemetry integration
6. **Circuit Breakers**: Automatic failure detection and recovery
7. **Message Encryption**: End-to-end encryption for sensitive prompts

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Verify connection string format
   - Check network connectivity
   - Confirm authentication credentials

2. **Message Rejected**
   - Check message format (JSON syntax)
   - Verify required fields are present
   - Ensure the model and processor are valid

3. **High Retry Rate**
   - Check API key validity
   - Monitor API rate limits
   - Review error logs for specific failures

### Debug Logging

Enable debug logging for detailed queue operation information:

```bash
LOG_LEVEL=debug cronai start
```

## API Reference

### Consumer Interface

```go
type Consumer interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Consume(ctx context.Context) (<-chan *Message, <-chan error)
    Acknowledge(ctx context.Context, message *Message) error
    Reject(ctx context.Context, message *Message, requeue bool) error
    Name() string
    Validate() error
}
```

### Message Structure

```go
type Message struct {
    ID          string
    Body        []byte
    Attributes  map[string]string
    ReceivedAt  time.Time
    RetryCount  int
    QueueSource string
}
```

### Task Message

```go
type TaskMessage struct {
    Model     string
    Prompt    string
    Processor string
    Variables map[string]string
    IsInline  bool
}
```
