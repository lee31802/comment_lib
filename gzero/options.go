package gzero

import (
	"github.com/lee31802/comment_lib/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"log"
)

type Options struct {
	AppPath            string
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	zrpc.RpcServerConf
	zrpc.RpcClientConf
}

func newOptions() *Options {
	return &Options{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name: "demo",
			},
			ListenOn: "0.0.0.0:8081",
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   "demo",
			},
		},
		RpcClientConf: zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   "demo",
			},
		},
	}
}

func (opts *Options) updateFromConfig(cfg *conf.Configuration) {
	err := cfg.UnmarshalKey("gozero", &opts)
	if err != nil {
		log.Printf("unmarshal gozero config err: %v", err)
	}
}

// Option defines a function to modify options.
type Option func(*Options)

// WithAppPath sets application path.
func WithAppPath(appPath string) Option {
	return func(opts *Options) {
		opts.AppPath = appPath
	}
}

// WithUnaryInterceptors sets unary interceptors.
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(opts *Options) {
		opts.UnaryInterceptors = interceptors
	}
}

// WithUnaryInterceptors sets unary interceptors.
func WithStreamServerInterceptor(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(opts *Options) {
		opts.StreamInterceptors = interceptors
	}
}
