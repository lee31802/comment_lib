package main

import (
	"context"
	"fmt"
	"github.com/lee31802/comment_lib/logkit"
)

func init() {
	logkit.Init(
		logkit.Level("debug"),
		logkit.Path("log"),
		logkit.MaxSize(1024),
		logkit.MaxBackups(1024),
		logkit.MaxAge(1),
		logkit.EnableCaller(true),
		logkit.EnableConsole(true),
		logkit.ErrorAsync(false),
	)
}

func main() {
	defer logkit.Sync()
	type request struct {
		A int
	}
	req := request{
		A: 1,
	}
	traceCtx := logkit.NewContextWith(context.Background(), logkit.FieldRequestID("22"), logkit.FieldTraceID("333"))
	err := fmt.Errorf("lxlerr")
	//logger := logkit.FromContext(traceCtx).With(logkit.Err(err))
	logkit.FromContext(traceCtx).Error("AreaModule.BatchGetLocationGroupsByDistrict rpc failed", logkit.Any("req", req), logkit.Err(err))
	logkit.FromContext(traceCtx).Info("infolog")
}
