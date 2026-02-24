// Package logger provides structured logging for the Lab CMS application.
// It uses the standard library's log package with custom formatting.
// In production, logs are output as JSON; in development, they're human-readable.
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity of a log entry
type LogLevel int

const (
	// DebugLevel is the most verbose, includes detailed information
	DebugLevel LogLevel = iota
	// InfoLevel includes general operational messages
	InfoLevel
	// WarnLevel includes warning messages
	WarnLevel
	// ErrorLevel includes error messages only
	ErrorLevel
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}

// ParseLogLevel converts a string to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// Logger provides structured logging with fields
type Logger struct {
	level     LogLevel
	isJSON    bool
	output    *log.Logger
	fields    map[string]interface{}
	fieldsMu  sync.RWMutex
	requestID string
	userID    int64
}

var (
	// globalLogger is the global logger instance
	globalLogger *Logger
	mu           sync.RWMutex
)

// Init initializes the global logger with the specified configuration
func Init(level string, isProduction bool) {
	mu.Lock()
	defer mu.Unlock()

	logLevel := ParseLogLevel(level)

	globalLogger = &Logger{
		level:  logLevel,
		isJSON: isProduction,
		output: log.New(os.Stdout, "", 0),
		fields: make(map[string]interface{}),
	}
}

// L returns the global logger instance
// If not initialized, it creates a default logger with info level
func L() *Logger {
	mu.RLock()
	l := globalLogger
	mu.RUnlock()

	if l == nil {
		// Double-check after acquiring write lock
		mu.Lock()
		defer mu.Unlock()
		if globalLogger == nil {
			// Initialize without calling Init() to avoid deadlock
			globalLogger = &Logger{
				level:  InfoLevel,
				isJSON: false,
				output: log.New(os.Stdout, "", 0),
				fields: make(map[string]interface{}),
			}
		}
		l = globalLogger
	}

	return l
}

// WithRequestID returns a new logger with the request ID set
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.clone()
	newLogger.requestID = requestID
	return newLogger
}

// WithUserID returns a new logger with the user ID set
func (l *Logger) WithUserID(userID int64) *Logger {
	newLogger := l.clone()
	newLogger.userID = userID
	return newLogger
}

// WithField returns a new logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.clone()
	newLogger.fieldsMu.Lock()
	newLogger.fields[key] = value
	newLogger.fieldsMu.Unlock()
	return newLogger
}

// WithFields returns a new logger with multiple additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := l.clone()
	newLogger.fieldsMu.Lock()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	newLogger.fieldsMu.Unlock()
	return newLogger
}

// clone creates a copy of the logger
func (l *Logger) clone() *Logger {
	l.fieldsMu.RLock()
	defer l.fieldsMu.RUnlock()

	newFields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		newFields[k] = v
	}

	return &Logger{
		level:     l.level,
		isJSON:    l.isJSON,
		output:    l.output,
		fields:    newFields,
		requestID: l.requestID,
		userID:    l.userID,
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg, nil)
}

// Debugf logs a formatted message at debug level
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...), nil)
}

// Info logs a message at info level
func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg, nil)
}

// Infof logs a formatted message at info level
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...), nil)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg, nil)
}

// Warnf logs a formatted message at warn level
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...), nil)
}

// Error logs a message at error level
func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg, nil)
}

// Errorf logs a formatted message at error level
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
}

// Fatal logs a message at error level and exits the program
func (l *Logger) Fatal(msg string) {
	l.log(ErrorLevel, msg, nil)
	os.Exit(1)
}

// Fatalf logs a formatted message at error level and exits the program
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...), nil)
	os.Exit(1)
}

// log writes a log entry with the given level and message
func (l *Logger) log(level LogLevel, msg string, err error) {
	if level < l.level {
		return
	}

	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   msg,
	}

	// Add request ID if set
	if l.requestID != "" {
		entry.RequestID = l.requestID
	}

	// Add user ID if set
	if l.userID > 0 {
		entry.UserID = l.userID
	}

	// Add error if present
	if err != nil {
		entry.Error = err.Error()
	}

	// Add custom fields
	l.fieldsMu.RLock()
	if len(l.fields) > 0 {
		entry.Fields = make(map[string]interface{}, len(l.fields))
		for k, v := range l.fields {
			entry.Fields[k] = v
		}
	}
	l.fieldsMu.RUnlock()

	// Output the log entry
	if l.isJSON {
		output, _ := json.Marshal(entry)
		l.output.Println(string(output))
	} else {
		l.output.Println(entry.String())
	}
}

// logEntry represents a structured log entry
type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    int64                  `json:"user_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// String returns a human-readable string representation
func (e logEntry) String() string {
	var parts []string

	// Format: [TIMESTAMP] [LEVEL] message
	parts = append(parts, fmt.Sprintf("[%s]", e.Timestamp))
	parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(e.Level)))

	// Add request ID if present
	if e.RequestID != "" {
		parts = append(parts, fmt.Sprintf("[req:%s]", e.RequestID))
	}

	// Add user ID if present
	if e.UserID > 0 {
		parts = append(parts, fmt.Sprintf("[user:%d]", e.UserID))
	}

	// Add message
	parts = append(parts, e.Message)

	// Add error if present
	if e.Error != "" {
		parts = append(parts, fmt.Sprintf("| error: %s", e.Error))
	}

	// Add fields
	if len(e.Fields) > 0 {
		var fieldParts []string
		for k, v := range e.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("| %s", strings.Join(fieldParts, " ")))
	}

	return strings.Join(parts, " ")
}

// IsLevelEnabled returns true if the given level is enabled
func (l *Logger) IsLevelEnabled(level LogLevel) bool {
	return level >= l.level
}
