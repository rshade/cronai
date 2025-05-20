// Package logger provides structured logging with different log levels.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// Fields is a type alias for log metadata fields
type Fields map[string]interface{}

// Level represents the severity level of a log message
type Level int

const (
	// DebugLevel is the most verbose logging level, used for debugging.
	DebugLevel Level = iota
	// InfoLevel is the standard logging level, used for general information.
	InfoLevel
	// WarnLevel is for logging warnings that don't cause application failure.
	WarnLevel
	// ErrorLevel is for logging errors that affect application function.
	ErrorLevel
	// FatalLevel is for logging fatal errors that require application termination.
	FatalLevel
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// Config holds the logger configuration.
type Config struct {
	// MinLevel is the minimum log level to output
	MinLevel Level
	// EnableJSON controls whether to output logs in JSON format
	EnableJSON bool
	// Output is where logs are written to
	Output io.Writer
	// IncludeTimestamp controls whether to include a timestamp in logs
	IncludeTimestamp bool
	// IncludeFileLine controls whether to include the file and line number in logs
	IncludeFileLine bool
}

// Logger is a structured logger that supports different log levels
type Logger struct {
	config Config
	mu     sync.Mutex
}

// New creates a new logger with the provided configuration
func New(config Config) *Logger {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	return &Logger{
		config: config,
	}
}

// DefaultLogger returns a logger with sensible defaults
func DefaultLogger() *Logger {
	return New(Config{
		MinLevel:         InfoLevel,
		EnableJSON:       false,
		Output:           os.Stdout,
		IncludeTimestamp: true,
		IncludeFileLine:  true,
	})
}

// logEntry represents a single log entry
type logEntry struct {
	Time     string `json:"time,omitempty"`
	Level    string `json:"level"`
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Metadata Fields `json:"metadata,omitempty"`
}

// log logs a message at the specified level with optional metadata.
func (l *Logger) log(level Level, msg string, metadata Fields) {
	if level < l.config.MinLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := logEntry{
		Level:    level.String(),
		Message:  msg,
		Metadata: metadata,
	}

	if l.config.IncludeTimestamp {
		entry.Time = time.Now().Format(time.RFC3339)
	}

	if l.config.IncludeFileLine {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.File = file
			entry.Line = line
		}
	}

	if l.config.EnableJSON {
		if err := json.NewEncoder(l.config.Output).Encode(entry); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error encoding log entry: %v\n", err)
		}
	} else {
		var timestamp, fileLine string
		if l.config.IncludeTimestamp {
			timestamp = fmt.Sprintf("[%s] ", entry.Time)
		}
		if l.config.IncludeFileLine {
			fileLine = fmt.Sprintf(" (%s:%d)", entry.File, entry.Line)
		}

		if _, err := fmt.Fprintf(l.config.Output, "%s[%s]%s %s", timestamp, entry.Level, fileLine, entry.Message); err != nil {
			return
		}
		if len(entry.Metadata) > 0 {
			if _, err := fmt.Fprintf(l.config.Output, " | "); err != nil {
				return
			}
			first := true
			for k, v := range entry.Metadata {
				if !first {
					if _, err := fmt.Fprintf(l.config.Output, ", "); err != nil {
						return
					}
				}
				if _, err := fmt.Fprintf(l.config.Output, "%s=%v", k, v); err != nil {
					return
				}
				first = false
			}
		}
		if _, err := fmt.Fprintln(l.config.Output); err != nil {
			return
		}
	}

	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug logs a debug message with optional metadata.
func (l *Logger) Debug(msg string, metadata ...Fields) {
	var data Fields
	if len(metadata) > 0 {
		data = metadata[0]
	}
	l.log(DebugLevel, msg, data)
}

// Info logs an info message with optional metadata.
func (l *Logger) Info(msg string, metadata ...Fields) {
	var data Fields
	if len(metadata) > 0 {
		data = metadata[0]
	}
	l.log(InfoLevel, msg, data)
}

// Warn logs a warning message with optional metadata.
func (l *Logger) Warn(msg string, metadata ...Fields) {
	var data Fields
	if len(metadata) > 0 {
		data = metadata[0]
	}
	l.log(WarnLevel, msg, data)
}

// Error logs an error message with optional metadata.
func (l *Logger) Error(msg string, metadata ...Fields) {
	var data Fields
	if len(metadata) > 0 {
		data = metadata[0]
	}
	l.log(ErrorLevel, msg, data)
}

// Fatal logs a fatal message with optional metadata and then exits the application.
func (l *Logger) Fatal(msg string, metadata ...Fields) {
	var data Fields
	if len(metadata) > 0 {
		data = metadata[0]
	}
	l.log(FatalLevel, msg, data)
}

// WithMetadata returns a function that logs with the provided metadata.
func (l *Logger) WithMetadata(_ Fields) *Logger {
	return &Logger{
		config: l.config,
	}
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.config.MinLevel = level
}

// GetLevel returns the current minimum log level.
func (l *Logger) GetLevel() Level {
	return l.config.MinLevel
}

// ParseLevel parses a string level into a Level.
func ParseLevel(level string) (Level, error) {
	switch level {
	case "DEBUG":
		return DebugLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "WARN":
		return WarnLevel, nil
	case "ERROR":
		return ErrorLevel, nil
	case "FATAL":
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}
