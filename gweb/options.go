package gweb

import (
	"github.com/lee31802/comment_lib/conf"
	"log"
)

type Jaeger struct {
	Enable       bool    `mapstructure:"enable"`
	SamplingRate float64 `mapstructure:"sampling_rate"`
}

type Options struct {
	// Default is current work path.
	Address string `mapstructure:"address"`
	// Default is true.
	Recovery bool `mapstructure:"recovery"`
	// Default is false.
	Pprof bool `mapstructure:"pprof"`
	// Default is "/".
	RootPath      string
	UploadMetrics bool `mapstructure:"upload_metrics"`
	// trace
	Jaeger Jaeger `mapstructure:"jaeger"`
}

func newOptions() *Options {
	return &Options{
		Recovery: true,
		RootPath: "/",
	}
}

func (opts *Options) updateFromConfig(cfg *conf.Configuration) {
	err := cfg.UnmarshalKey("ginweb", &opts)
	if err != nil {
		log.Printf("unmarshal ginweb config err: %v", err)
	}
	if !cfg.IsSet("ginweb.upload_metrics") {
		opts.UploadMetrics = true
	}
}

// Option defines a function to modify client.
type Option func(*Options)
