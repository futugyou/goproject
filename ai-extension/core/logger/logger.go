/*
Based on  https://github.com/dapr/kit/blob/main/logger/logger.go
*/

package logger

import (
	"context"
	"io"
	"strings"
)

const (
	LogTypeLog     = "log"
	LogTypeRequest = "request"
)

type logContextKeyType struct{}

var logContextKey = logContextKeyType{}

type LogLevel string

const (
	TraceLevel     LogLevel = "trace"
	DebugLevel     LogLevel = "debug"
	InfoLevel      LogLevel = "info"
	WarnLevel      LogLevel = "warn"
	ErrorLevel     LogLevel = "error"
	FatalLevel     LogLevel = "fatal"
	UndefinedLevel LogLevel = "undefined"
)

type Logger interface {
	EnableJSONOutput(enabled bool)
	SetAppID(id string)
	SetOutputLevel(outputLevel LogLevel)
	SetOutput(dst io.Writer)
	IsOutputLevelEnabled(level LogLevel) bool
	WithLogType(logType string) Logger
	WithFields(fields map[string]any) Logger
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func ToLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}

	return UndefinedLevel
}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logContextKey, logger)
}
