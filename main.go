package main

import (
	"context"
	"fmt"
	"github.com/lee31802/comment_lib/logkit"
	"github.com/lee31802/comment_lib/trace"
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
	// log： 1.调整了输出样式更加直观 2.自动打印context里边携带的requestid和traceid3. 结合lumberjack做了日志切割等
	defer logkit.Sync()
	type request struct {
		A int
	}
	req := request{
		A: 1,
	}
	traceCtx := logkit.NewContextWith(context.Background(), trace.FieldRequestID("22"), trace.FieldTraceID("333"))
	err := fmt.Errorf("lxlerr")
	//logger := logkit.FromContext(traceCtx).With(logkit.Err(err))
	logkit.FromContext(traceCtx).Error("AreaModule.BatchGetLocationGroupsByDistrict rpc failed", logkit.Any("req", req), logkit.Err(err))
	logkit.FromContext(traceCtx).Info("infolog")
	//
	// server:1.继承了requeset会自动调用validate等方法 2.自动绑定request 3.自动生成一个可以复现的curl命令 4.自动生成接口文档

}
