package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/inconshreveable/log15"
)

type Logger struct {
	ctx    context.Context
	fields map[string]interface{}
}

func (logger *Logger) WithField(name string, value interface{}) *Logger {
	logger.fields[name] = value
	return logger
}

func (logger *Logger) WithError(err error) *Logger {
	return logger.WithField("error", err)
}

func (logger *Logger) buildOutput(message string) string {
	output := make([]string, 0)
	for name, value := range logger.fields {
		output = append(output, fmt.Sprintf("[%s: %s]", name, value))
	}
	return fmt.Sprintf("%s %s", message, strings.Join(output, " "))
}

func (logger *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log15.Error(logger.buildOutput(message))
}

func (logger *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log15.Info(logger.buildOutput(message))
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{ctx, make(map[string]interface{})}
}
