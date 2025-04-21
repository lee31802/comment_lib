package logkit

import (
	"context"
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *LogkitLogger
	level         zapcore.Level
)

type key int

// key list
const (
	keyLogger key = iota
)

func init() {
	// init default logger by default options
	// which would not return error, so don't need handle error
	defaultLogger, _ = newLogger(defaultOptions...)
}

// Init will reset logger's options, so need recreate defaultLogger
// maybe Init isn't a good function name
// but keep it in order to maintain compatibility
func Init(opts ...Option) error {
	var err error
	defaultLogger, err = newLogger(opts...)
	return err
}

func GetLogger() *LogkitLogger {
	return defaultLogger
}

func Sync() error {
	return defaultLogger.Sync()
}

func SetLevel(l string) error {
	err := level.UnmarshalText([]byte(l))
	if err != nil {
		return err
	}
	return nil
}

// Fatal outputs a message at fatal level.
func Fatal(msg string, fields ...Field) {
	GetLogger().Logger.Fatal(msg, append(fields, extractFields(fields)...)...)
}

// Error outputs a message at error level.
func Error(msg string, fields ...Field) {
	GetLogger().Logger.Error(msg, append(fields, extractFields(fields)...)...)
}

// Info outputs a message at info level.
func Info(msg string, fields ...Field) {
	GetLogger().Logger.Info(msg, append(fields, extractFields(fields)...)...)
}

// Warn outputs a message at warn level.
func Warn(msg string, fields ...Field) {
	GetLogger().Logger.Warn(msg, append(fields, extractFields(fields)...)...)
}

// Debug outputs a message at debug level.
func Debug(msg string, fields ...Field) {
	GetLogger().Logger.Debug(msg, append(fields, extractFields(fields)...)...)
}

func With(fields ...Field) *LogkitLogger {
	return GetLogger().With(fields...)
}

func NewContext(ctx context.Context, l *LogkitLogger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

func NewContextWith(ctx context.Context, field ...Field) context.Context {
	l := FromContext(ctx)
	l = l.With(field...)
	return NewContext(ctx, l)
}

func FromContext(ctx context.Context) *LogkitLogger {
	l, ok := ctx.Value(keyLogger).(*LogkitLogger)
	if !ok {
		return GetLogger()
	}
	return l
}

func FieldRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

func FieldTraceID(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}
