package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Jaeger struct {
	Enable       bool    `mapstructure:"enable"`
	SamplingRate float64 `mapstructure:"sampling_rate"`
}

type Options struct {
	// Default is current work path.
	AppPath string
	Address string `mapstructure:"address"`
	// Default is true.
	Recovery bool `mapstructure:"recovery"`
	// Default is false.
	Pprof bool `mapstructure:"pprof"`
	// Default is "/".
	RootPath      string
	UploadMetrics bool `mapstructure:"upload_metrics"`
	Middlewares   []gin.HandlerFunc
	Plugins       []Plugin
	Engine        *gin.Engine
	// trace
	Jaeger Jaeger `mapstructure:"jaeger"`
}

func newOptions() *Options {
	return &Options{
		Recovery: true,
		Engine:   gin.New(),
		RootPath: "/",
	}
}

func (opts *Options) updateFromConfig(cfg *Configuration) {
	err := cfg.UnmarshalKey("ginweb", &opts)
	if err != nil {
		log.Printf("unmarshal ginweb config err: %v", err)
	}
	if !cfg.IsSet("ginweb.upload_metrics") {
		opts.UploadMetrics = true
	}
}

// Option defines a function to modify options.
type Option func(*Options)

// WithAppApth sets application path.
func WithAppApth(appPath string) Option {
	return func(opts *Options) {
		opts.AppPath = appPath
	}
}

// WithEngine sets custom gin engine.
func WithEngine(engine *gin.Engine) Option {
	return func(opts *Options) {
		opts.Engine = engine
	}
}

// WithRootPath sets api root path.
func WithRootPath(rootPath string) Option {
	return func(opts *Options) {
		opts.RootPath = rootPath
	}
}

// WithMiddlewares sets global middlewares.
func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(opts *Options) {
		opts.Middlewares = middlewares
	}
}

// WithPlugins adapts custom plugins.
func WithPlugins(plugins ...Plugin) Option {
	return func(opts *Options) {
		opts.Plugins = plugins
	}
}
