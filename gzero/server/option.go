package server

import (
	"github.com/lee31802/comment_lib/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
)

type Options struct {
	Server zrpc.RpcServerConf
}

func newOptions() *Options {
	return &Options{
		Server: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name: "demo.rpc",
			},
			ListenOn: "0.0.0.0:8081",
			Etcd: discov.EtcdConf{
				Hosts: []string{"127.0.0.1:2379"},
				Key:   "demo.rpc",
			},
		},
	}
}

func (opts *Options) updateFromConfig(cfg *conf.Configuration) {
	server := zrpc.RpcServerConf{}
	err := cfg.UnmarshalKey("Server", &server)
	if err != nil {
		log.Printf("unmarshal gozero clientserver config err: %v", err)
	}
	opts.Server = server
}

// Option defines a function to modify client.
type Option func(*Options)
