package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/inconshreveable/log15"
)

// Logger represents a structured logger.
type Logger struct {
	ctx    context.Context
	fields map[string]interface{}
}

// WithField adds a key-value pair to the logger's fields and returns a new Logger instance.
func (logger *Logger) WithField(name string, value interface{}) *Logger {
	logger.fields[name] = value
	return logger
}

// WithError adds an "error" field with the given error value to the logger's fields and returns a new Logger instance.
func (logger *Logger) WithError(err error) *Logger {
	return logger.WithField("error", err)
}

// buildOutput constructs the log output string by appending all fields in a formatted manner.
func (logger *Logger) buildOutput(message string) string {
	output := make([]string, 0)
	for name, value := range logger.fields {
		output = append(output, fmt.Sprintf("[%s: %s]", name, value))
	}
	return fmt.Sprintf("%s %s", message, strings.Join(output, " "))
}

// Error logs an error message with the formatted output and additional fields.
func (logger *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log15.Error(logger.buildOutput(message))
}

// Info logs an informational message with the formatted output and additional fields.
func (logger *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log15.Info(logger.buildOutput(message))
}

// NewLogger creates a new Logger instance with the given context.
func NewLogger(ctx context.Context) *Logger {
	return &Logger{ctx, make(map[string]interface{})}
}
