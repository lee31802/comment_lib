package gredis

import (
	"github.com/lee31802/comment_lib/conf"
	redis "github.com/redis/go-redis/v9"
)

// 分片
type RingOptions struct {
	appPath string
	redis.RingOptions
}

// 集群
type ClusterOptions struct {
	appPath string
	redis.ClusterOptions
}

// 哨兵
type FailoverClientOptions struct {
	appPath string
	redis.FailoverOptions
}

type DefaultClientOptions struct {
	appPath string
	redis.Options
}

const (
	DefaultMaxRetries   = 5
	DefaultPoolSize     = 100
	DefaultDialTimeout  = 500
	DefaultReadTimeout  = 100
	DefaultWriteTimeout = 100
	DefaultPoolTimeout  = 1000
	DefaultIdleTimeout  = 60000
)

func newDefaultOptions() *DefaultClientOptions {
	return &DefaultClientOptions{}
}

// Option defines a function to modify client.
type Option func(*DefaultClientOptions)

func (opts *DefaultClientOptions) updateFromConfig(cfg *conf.Configuration) {
	//zClient := zrpc.RpcClientConf{}
	//err := cfg.UnmarshalKey(opts.serviceName, &zClient)
	//if err != nil {
	//	log.Printf("unmarshal gozero client config err: %v", err)
	//}
	//opts.Client = zClient
}

// WithAppPath sets application path.
func WithAppPath(appPath string) Option {
	return func(opts *DefaultClientOptions) {
		opts.appPath = appPath
	}
}
