package server

import (
	"google.golang.org/grpc"
)

type Command struct {
	AppPath            string
	Plugins            []Plugin
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	RegisterServer     func(grpcServer *grpc.Server)
	// PreRun: children of this command will not inherit.
	PreRun func() error
	// PostRun: run after the Run command.
	PostRun  func() error
	PreStop  func() error
	PostStop func() error
}

func (cmd Command) Execute() error {
	return gs.Run(cmd)
}
