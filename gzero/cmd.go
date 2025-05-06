package gzero

import (
	"google.golang.org/grpc"
)

type Command struct {
	Name         string
	AppPath      string
	RegisTerFunc func(grpcServer *grpc.Server)
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
