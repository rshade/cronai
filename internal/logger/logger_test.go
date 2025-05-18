package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{Level(999), "LEVEL(999)"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("expected level %d to be %s, got %s", test.level, test.expected, test.level.String())
		}
	}
}

func TestLoggerBasic(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(Config{
		MinLevel:         DebugLevel,
		EnableJSON:       false,
		Output:           buf,
		IncludeTimestamp: false,
		IncludeFileLine:  false,
	})

	logger.Info("test message")
	if !strings.Contains(buf.String(), "[INFO] test message") {
		t.Errorf("expected log to contain '[INFO] test message', got %s", buf.String())
	}
}

func TestLoggerWithMetadata(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(Config{
		MinLevel:         DebugLevel,
		EnableJSON:       false,
		Output:           buf,
		IncludeTimestamp: false,
		IncludeFileLine:  false,
	})

	metadata := Fields{
		"key1": "value1",
		"key2": 123,
	}

	logger.Info("test message", metadata)
	if !strings.Contains(buf.String(), "key1=value1") || !strings.Contains(buf.String(), "key2=123") {
		t.Errorf("expected log to contain metadata, got %s", buf.String())
	}
}

func TestLoggerJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(Config{
		MinLevel:         DebugLevel,
		EnableJSON:       true,
		Output:           buf,
		IncludeTimestamp: false,
		IncludeFileLine:  false,
	})

	metadata := Fields{
		"key1": "value1",
		"key2": 123,
	}

	logger.Info("test message", metadata)

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("error unmarshaling JSON log: %v", err)
	}

	if entry["level"] != "INFO" {
		t.Errorf("expected level to be INFO, got %v", entry["level"])
	}
	if entry["message"] != "test message" {
		t.Errorf("expected message to be 'test message', got %v", entry["message"])
	}

	meta, ok := entry["metadata"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected metadata to be a map, got %T", entry["metadata"])
	}
	if meta["key1"] != "value1" {
		t.Errorf("expected metadata.key1 to be 'value1', got %v", meta["key1"])
	}
	if val, ok := meta["key2"].(float64); !ok || val != 123 {
		t.Error("Expected key2 to be 123")
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(Config{
		MinLevel:         WarnLevel,
		EnableJSON:       false,
		Output:           buf,
		IncludeTimestamp: false,
		IncludeFileLine:  false,
	})

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	if strings.Contains(buf.String(), "debug message") {
		t.Errorf("debug message should have been filtered out")
	}
	if strings.Contains(buf.String(), "info message") {
		t.Errorf("info message should have been filtered out")
	}
	if !strings.Contains(buf.String(), "warn message") {
		t.Errorf("expected log to contain 'warn message', got %s", buf.String())
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		hasError bool
	}{
		{"DEBUG", DebugLevel, false},
		{"INFO", InfoLevel, false},
		{"WARN", WarnLevel, false},
		{"ERROR", ErrorLevel, false},
		{"FATAL", FatalLevel, false},
		{"UNKNOWN", InfoLevel, true},
	}

	for _, test := range tests {
		level, err := ParseLevel(test.input)
		if test.hasError && err == nil {
			t.Errorf("expected error parsing level %s, got nil", test.input)
		}
		if !test.hasError && err != nil {
			t.Errorf("unexpected error parsing level %s: %v", test.input, err)
		}
		if level != test.expected {
			t.Errorf("expected level %s to parse to %d, got %d", test.input, test.expected, level)
		}
	}
}

func TestSetLevel(t *testing.T) {
	logger := DefaultLogger()
	if logger.GetLevel() != InfoLevel {
		t.Errorf("expected default level to be INFO, got %s", logger.GetLevel())
	}

	logger.SetLevel(ErrorLevel)
	if logger.GetLevel() != ErrorLevel {
		t.Errorf("expected level to be ERROR after SetLevel, got %s", logger.GetLevel())
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := DefaultLogger()
	if logger.config.MinLevel != InfoLevel {
		t.Errorf("expected default level to be INFO, got %s", logger.config.MinLevel)
	}
	if logger.config.EnableJSON != false {
		t.Errorf("expected EnableJSON to be false")
	}
	if logger.config.IncludeTimestamp != true {
		t.Errorf("expected IncludeTimestamp to be true")
	}
	if logger.config.IncludeFileLine != true {
		t.Errorf("expected IncludeFileLine to be true")
	}
}
