package logkit

import (
	"golog/comm_lib/logkit/lumberjack"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(opts ...Option) (*LogkitLogger, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	err := level.UnmarshalText([]byte(options.level)) //传入 "info" 时，会转换为 zapcore.InfoLevel
	if err != nil {
		return nil, err
	}

	infoFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:    options.path + "/info.log",
		MaxSize:     options.maxSize,
		MaxBackups:  options.maxBackups,
		MaxAge:      options.maxAge,
		BufferSize:  options.bufferSize,
		ChannelSize: options.channelSize,
		AsyncWrite:  true,
	})
	errFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:    options.path + "/error.log",
		MaxSize:     options.maxSize,
		MaxBackups:  options.maxBackups,
		MaxAge:      options.maxAge,
		BufferSize:  options.bufferSize,
		ChannelSize: options.channelSize,
		AsyncWrite:  options.errorAsync,
	})
	debugFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:    options.path + "/debug.log",
		MaxSize:     options.maxSize,
		MaxBackups:  options.maxBackups,
		MaxAge:      options.maxAge,
		BufferSize:  options.bufferSize,
		ChannelSize: options.channelSize,
		AsyncWrite:  true,
	})

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	liveEncoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	devEncoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	zapCores := []zapcore.Core{
		zapcore.NewCore(options.encoderBuilder(liveEncoder), errFileHandler, NewLevelEnabler(&level, zapcore.ErrorLevel)),
		zapcore.NewCore(options.encoderBuilder(liveEncoder), infoFileHandler, NewLevelEnabler(&level, zapcore.InfoLevel)),
		zapcore.NewCore(options.encoderBuilder(devEncoder), debugFileHandler, NewLevelEnabler(&level, zapcore.DebugLevel)),
	}

	if options.enableConsole {
		//zapcore.Lock 用于确保在多线程环境下对标准输出的安全访问
		consoleHandler := zapcore.Lock(os.Stdout)
		zapCores = append(zapCores, zapcore.NewCore(zapcore.NewConsoleEncoder(devEncoder), consoleHandler, zapcore.DebugLevel))
	}
	// create options with priority for our opts
	defaultZapOptions := []zap.Option{}
	if options.enableCaller {
		defaultZapOptions = append(
			defaultZapOptions,
			zap.AddCaller(),
			zap.AddCallerSkip(1), //避免日志内部的调用
		)
	}

	core := zapcore.NewTee(
		zapCores...,
	)

	logger := &LogkitLogger{
		zap.New(NewErrorsExtractCore(core), defaultZapOptions...),
	}

	return logger, err
}

type LogkitLogger struct {
	*zap.Logger
}

func (wrapper *LogkitLogger) Error(msg string, fields ...Field) {
	wrapper.Logger.Error(msg, append(fields, extractFields(fields)...)...)
}

func (wrapper *LogkitLogger) Info(msg string, fields ...Field) {
	wrapper.Logger.Info(msg, append(fields, extractFields(fields)...)...)
}

func (wrapper *LogkitLogger) Warn(msg string, fields ...Field) {
	wrapper.Logger.Warn(msg, append(fields, extractFields(fields)...)...)
}

func (wrapper *LogkitLogger) Debug(msg string, fields ...Field) {
	wrapper.Logger.Debug(msg, append(fields, extractFields(fields)...)...)
}

func (wrapper *LogkitLogger) With(fields ...Field) *LogkitLogger {
	return &LogkitLogger{wrapper.Logger.With(fields...)}
}

func (wrapper *LogkitLogger) WithFields(fields map[string]interface{}) *LogkitLogger {
	fieldList := make([]Field, 0, len(fields))
	for field, value := range fields {
		fieldList = append(fieldList, Reflect(field, value))
	}
	return wrapper.With(fieldList...)
}

func (wrapper *LogkitLogger) WithField(key string, value interface{}) *LogkitLogger {
	return wrapper.With(Reflect(key, value))
}

func (wrapper *LogkitLogger) WithError(err error) *LogkitLogger {
	return wrapper.With(Err(err))
}
