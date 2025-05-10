package client

import (
	"github.com/zeromicro/go-zero/zrpc"
)

func InitClient(options ...Option) zrpc.Client {
	cli := newGoZeroClient(options...)
	return cli.initClient()
}
