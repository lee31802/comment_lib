package client

import (
	"github.com/lee31802/comment_lib/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
)

type Options struct {
	appPath     string
	serviceName string
	Client      zrpc.RpcClientConf
}

func newOptions() *Options {
	return &Options{
		Client: zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   "demo.rpc",
			},
		},
	}
}

// Option defines a function to modify client.
type Option func(*Options)

func (opts *Options) updateFromConfig(cfg *conf.Configuration) {
	zClient := zrpc.RpcClientConf{}
	err := cfg.UnmarshalKey(opts.serviceName, &zClient)
	if err != nil {
		log.Printf("unmarshal gozero client config err: %v", err)
	}
	opts.Client = zClient
}

// WithAppPath sets application path.
func WithAppPath(appPath string) Option {
	return func(opts *Options) {
		opts.appPath = appPath
	}
}

// WithServiceName sets service name.
func WithServiceName(serviceName string) Option {
	return func(opts *Options) {
		opts.serviceName = serviceName
	}
}
