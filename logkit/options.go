package logkit

import (
	"github.com/lee31802/comment_lib/logkit/encoder"
	"go.uber.org/zap/zapcore"
)

const (
	LOG_LEVEL          = "INFO"
	LOG_PATH           = "log"
	LOG_MAX_SIZE       = 100
	LOG_MAX_BACKUPS    = 10
	LOG_MAX_AGE        = 3
	LOG_BUFFER_SIZE    = 4 * 1024 * 1024
	LOG_CHANNEL_SIZE   = 32 * 1024
	LOG_ENABLE_CONSOLE = false
	LOG_ENABLE_CALLER  = false
	LOG_ERROR_ASYNC    = true
)

var defaultOptions = []Option{
	Level(LOG_LEVEL),
	Path(LOG_PATH),
	MaxSize(LOG_MAX_SIZE),
	MaxBackups(LOG_MAX_BACKUPS),
	MaxAge(LOG_MAX_AGE),
	BufferSize(LOG_BUFFER_SIZE),
	ChannelSize(LOG_CHANNEL_SIZE),
	EnableConsole(LOG_ENABLE_CONSOLE),
	EnableCaller(LOG_ENABLE_CALLER),
	ErrorAsync(LOG_ERROR_ASYNC),
	Encoder(encoder.NewEncoder),
}

func newOptions() *Options {
	options := &Options{}
	for _, opt := range defaultOptions {
		opt(options)
	}
	return options
}

type Options struct {
	configPath    string
	Path          string `mapstructure:"Path"`
	MaxSize       int    `mapstructure:"MaxSize"`
	MaxBackups    int    `mapstructure:"MaxBackups"`
	MaxAge        int    `mapstructure:"MaxAge"`
	Level         string `mapstructure:"Level"`
	BufferSize    int    `mapstructure:"BufferSize"`
	ChannelSize   int    `mapstructure:"ChannelSize"`
	EnableConsole bool   `mapstructure:"EnableConsole"`
	EnableCaller  bool   `mapstructure:"EnableCaller"`
	ErrorAsync    bool   `mapstructure:"ErrorAsync"` // error log using async writer

	encoderBuilder encoderBuilder
}

var logOpts Options

type Option func(*Options)
type encoderBuilder func(cfg zapcore.EncoderConfig) zapcore.Encoder

func Configure(options ...Option) {
	for _, setter := range options {
		setter(&logOpts)
	}
}

func ConfigPath(path string) Option {
	return func(o *Options) {
		o.configPath = path
	}
}
func Path(p string) Option {
	return func(o *Options) {
		if p == "" {
			return
		}
		o.Path = p
	}
}

func MaxSize(m int) Option {
	return func(o *Options) {
		if m == 0 {
			return
		}
		o.MaxSize = m
	}
}

func MaxBackups(m int) Option {
	return func(o *Options) {
		if m == 0 {
			return
		}
		o.MaxBackups = m
	}
}

func MaxAge(m int) Option {
	return func(o *Options) {
		if m == 0 {
			return
		}
		o.MaxAge = m
	}
}

func Level(l string) Option {
	return func(o *Options) {
		if l == "" {
			return
		}
		o.Level = l
	}
}

func BufferSize(b int) Option {
	return func(o *Options) {
		if b == 0 {
			return
		}
		o.BufferSize = b
	}
}

func ChannelSize(c int) Option {
	return func(o *Options) {
		if c == 0 {
			return
		}
		o.ChannelSize = c
	}
}

func EnableConsole(e bool) Option {
	return func(o *Options) {
		o.EnableConsole = e
	}
}

func EnableCaller(e bool) Option {
	return func(o *Options) {
		o.EnableCaller = e
	}
}

func ErrorAsync(e bool) Option {
	return func(o *Options) {
		o.ErrorAsync = e
	}
}

func Encoder(eb encoderBuilder) Option {
	return func(o *Options) {
		o.encoderBuilder = eb
	}
}

func JsonEncoder() Option {
	return func(o *Options) {
		o.encoderBuilder = zapcore.NewJSONEncoder
	}
}
