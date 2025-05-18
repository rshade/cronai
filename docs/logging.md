# Logging in CronAI

CronAI implements structured logging with configurable log levels to help with troubleshooting and monitoring.

## Log Levels

The following log levels are supported, in order of increasing severity:

- **DEBUG**: Detailed information, typically only useful when troubleshooting issues
- **INFO**: General information about the normal operation of the application
- **WARN**: Warnings that don't affect application function but should be addressed
- **ERROR**: Errors that affect application function but don't cause termination
- **FATAL**: Fatal errors that require application termination

## Configuring Log Level

The log level can be configured through the `LOG_LEVEL` environment variable:

```bash
export LOG_LEVEL=DEBUG
./cronai start
```

Valid values for `LOG_LEVEL` are: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`.

If not specified, the default log level is `INFO`.

## Structured Logging

CronAI uses structured logging to provide context for log messages. Each log message includes:

- Timestamp in RFC3339 format
- Log level
- Message
- File and line number (for debugging)
- Additional metadata specific to the log message

Example log output:

```
[2025-05-18T14:30:45Z] [INFO] (service.go:40) Starting CronAI service | config_path=/etc/cronai/cronai.config
[2025-05-18T14:30:45Z] [INFO] (service.go:207) Successfully parsed configuration file | path=/etc/cronai/cronai.config, task_count=3
[2025-05-18T14:30:45Z] [INFO] (service.go:70) Scheduled task | task_index=0, schedule=0 9 * * 1-5, model=claude, prompt=weekly_report, processor=email-team@company.com
```

## JSON Logging

For integration with log management systems, CronAI supports JSON-formatted logs. To enable JSON logging, set the `LOG_FORMAT` environment variable to `JSON`:

```bash
export LOG_FORMAT=JSON
./cronai start
```

Example JSON log output:

```json
{"time":"2025-05-18T14:30:45Z","level":"INFO","message":"Starting CronAI service","file":"service.go","line":40,"metadata":{"config_path":"/etc/cronai/cronai.config"}}
{"time":"2025-05-18T14:30:45Z","level":"INFO","message":"Successfully parsed configuration file","file":"service.go","line":207,"metadata":{"path":"/etc/cronai/cronai.config","task_count":3}}
```

## Error Handling

CronAI implements categorized error handling through the `errors` package. Errors are categorized as:

- **CONFIGURATION**: Errors related to configuration files and parameters
- **VALIDATION**: Errors related to input validation
- **EXTERNAL**: Errors from external services (APIs, etc.)
- **SYSTEM**: System-level errors (file I/O, etc.)
- **APPLICATION**: Application-level errors

Error logs include the error category and additional context information to aid in troubleshooting.

## Log File

By default, logs are written to STDOUT. To direct logs to a file, use the `LOG_FILE` environment variable:

```bash
export LOG_FILE=/var/log/cronai.log
./cronai start
```

If not specified, logs are written to STDOUT.

## Troubleshooting

For troubleshooting issues, set the log level to DEBUG:

```bash
export LOG_LEVEL=DEBUG
./cronai start
```

This will provide detailed logs of all operations, including prompt loading, model execution, and response processing.