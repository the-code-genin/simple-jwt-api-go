package logger

import (
	"context"

	"go.uber.org/zap"
)

var logger *zap.Logger

const (
	FunctionNameField = "functionName"
	RequestBodyField  = "requestBody"

	loggerfields = "logger.fields"
)

func init() {
	var err error
	logger, err = zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func With(ctx context.Context, fields ...zap.Field) context.Context {
	data := ctx.Value(loggerfields)
	var storedFields = []zap.Field{}
	if data != nil {
		storedFields = data.([]zap.Field)
	}
	storedFields = append(storedFields, fields...)
	return context.WithValue(ctx, loggerfields, storedFields)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	data := ctx.Value(loggerfields)
	var storedFields = []zap.Field{}
	if data != nil {
		storedFields = data.([]zap.Field)
	}
	storedFields = append(storedFields, fields...)

	logger.Error(msg, storedFields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	data := ctx.Value(loggerfields)

	var storedFields = []zap.Field{}
	if data != nil {
		storedFields = data.([]zap.Field)
	}
	storedFields = append(storedFields, fields...)

	logger.Info(msg, storedFields...)
}
