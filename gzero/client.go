package gzero

import (
	"github.com/zeromicro/go-zero/zrpc"
)

var DefaultClient zrpc.Client

func (g *goZero) InitClient() zrpc.Client {
	DefaultClient = zrpc.MustNewClient(g.opts.RpcClientConf)
	return DefaultClient
}
