package logkit

import (
	"github.com/lee31802/comment_lib/logkit/encoder"
	"github.com/lee31802/comment_lib/util"
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
	appPath       string
	path          string
	maxSize       int
	maxBackups    int
	maxAge        int
	level         string
	bufferSize    int
	channelSize   int
	enableConsole bool
	enableCaller  bool
	errorAsync    bool // error log using async writer

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

func AppPath(path string) Option {
	return func(o *Options) {
		o.appPath = util.GetWorkDir()
	}
}
func Path(p string) Option {
	return func(o *Options) {
		o.path = p
	}
}

func MaxSize(m int) Option {
	return func(o *Options) {
		o.maxSize = m
	}
}

func MaxBackups(m int) Option {
	return func(o *Options) {
		o.maxBackups = m
	}
}

func MaxAge(m int) Option {
	return func(o *Options) {
		o.maxAge = m
	}
}

func Level(l string) Option {
	return func(o *Options) {
		o.level = l
	}
}

func BufferSize(b int) Option {
	return func(o *Options) {
		o.bufferSize = b
	}
}

func ChannelSize(c int) Option {
	return func(o *Options) {
		o.channelSize = c
	}
}

func EnableConsole(e bool) Option {
	return func(o *Options) {
		o.enableConsole = e
	}
}

func EnableCaller(e bool) Option {
	return func(o *Options) {
		o.enableCaller = e
	}
}

func ErrorAsync(e bool) Option {
	return func(o *Options) {
		o.errorAsync = e
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
