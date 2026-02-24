package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{LogLevel(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"INFO", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"unknown", InfoLevel}, // default
		{"", InfoLevel},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLogLevel(tt.input); got != tt.expected {
				t.Errorf("ParseLogLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  InfoLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	// Debug should not be logged when level is Info
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not be logged when level is Info")
	}

	// Info should be logged
	logger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should be logged")
	}

	buf.Reset()

	// Error should be logged
	logger.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should be logged")
	}
}

func TestLogger_WithRequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	newLogger := logger.WithRequestID("test-request-id")
	newLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test-request-id") {
		t.Errorf("Request ID should be in log output, got: %s", output)
	}
}

func TestLogger_WithUserID(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	newLogger := logger.WithUserID(12345)
	newLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "12345") {
		t.Errorf("User ID should be in log output, got: %s", output)
	}
}

func TestLogger_WithField(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	newLogger := logger.WithField("key", "value")
	newLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("Custom field should be in log output, got: %s", output)
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	fields := map[string]interface{}{
		"field1": "value1",
		"field2": 42,
	}
	newLogger := logger.WithFields(fields)
	newLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "field1=value1") || !strings.Contains(output, "field2=42") {
		t.Errorf("Custom fields should be in log output, got: %s", output)
	}
}

func TestLogger_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: true,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	logger.WithRequestID("req-123").Info("test message")

	// Parse the JSON output
	var entry logEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	if entry.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", entry.Level)
	}
	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", entry.Message)
	}
	if entry.RequestID != "req-123" {
		t.Errorf("Expected request_id 'req-123', got '%s'", entry.RequestID)
	}
}

func TestLogger_Debugf(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  DebugLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	logger.Debugf("formatted %s %d", "string", 42)

	output := buf.String()
	if !strings.Contains(output, "formatted string 42") {
		t.Errorf("Formatted message should be in output, got: %s", output)
	}
}

func TestLogger_Infof(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  InfoLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	logger.Infof("formatted %s", "message")

	output := buf.String()
	if !strings.Contains(output, "formatted message") {
		t.Errorf("Formatted message should be in output, got: %s", output)
	}
}

func TestLogger_Warnf(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  WarnLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	logger.Warnf("warning: %s", "alert")

	output := buf.String()
	if !strings.Contains(output, "warning: alert") {
		t.Errorf("Formatted warning should be in output, got: %s", output)
	}
}

func TestLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  ErrorLevel,
		isJSON: false,
		output: log.New(&buf, "", 0),
		fields: make(map[string]interface{}),
	}

	logger.Errorf("error occurred: %v", "something failed")

	output := buf.String()
	if !strings.Contains(output, "error occurred") {
		t.Errorf("Formatted error should be in output, got: %s", output)
	}
}

func TestGlobalLogger(t *testing.T) {
	// Ensure the mutex is not locked before starting
	mu.Lock()
	// Reset global logger
	globalLogger = nil
	mu.Unlock()

	// Should auto-initialize with default settings
	l := L()
	if l == nil {
		t.Error("Global logger should not be nil")
	}

	if l.level != InfoLevel {
		t.Errorf("Default level should be Info, got %v", l.level)
	}
}

func TestInit(t *testing.T) {
	Init("debug", true)

	if globalLogger == nil {
		t.Fatal("Global logger should be initialized")
	}

	if globalLogger.level != DebugLevel {
		t.Errorf("Level should be Debug, got %v", globalLogger.level)
	}

	if !globalLogger.isJSON {
		t.Error("Should be JSON mode in production")
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	tests := []struct {
		loggerLevel LogLevel
		checkLevel  LogLevel
		expected    bool
	}{
		{InfoLevel, DebugLevel, false},
		{InfoLevel, InfoLevel, true},
		{InfoLevel, WarnLevel, true},
		{InfoLevel, ErrorLevel, true},
		{DebugLevel, DebugLevel, true},
		{ErrorLevel, InfoLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.loggerLevel.String()+"_"+tt.checkLevel.String(), func(t *testing.T) {
			logger := &Logger{
				level: tt.loggerLevel,
			}
			if got := logger.IsLevelEnabled(tt.checkLevel); got != tt.expected {
				t.Errorf("IsLevelEnabled(%v) with logger level %v = %v, want %v",
					tt.checkLevel, tt.loggerLevel, got, tt.expected)
			}
		})
	}
}

func TestLogEntry_String(t *testing.T) {
	entry := logEntry{
		Timestamp: "2026-02-23T10:00:00Z",
		Level:     "info",
		Message:   "test message",
		RequestID: "req-abc",
		UserID:    42,
		Fields: map[string]interface{}{
			"key": "value",
		},
	}

	output := entry.String()

	expectedParts := []string{
		"[2026-02-23T10:00:00Z]",
		"[INFO]",
		"[req:req-abc]",
		"[user:42]",
		"test message",
		"key=value",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected output to contain %q, got: %s", part, output)
		}
	}
}
