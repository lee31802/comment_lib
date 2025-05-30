package main

import (
	"github.com/lee31802/comment_lib/gweb"
	"github.com/lee31802/comment_lib/gzero/server"
	"github.com/lee31802/comment_lib/logkit"
	"google.golang.org/grpc"
)

var (
	cmd = &gweb.Command{
		Name: "test",
		PreRun: func(router gweb.Router) error {
			logkit.InitByCfg()
			logkit.Error("success2222")
			router.GET("/", func() string { return "OK" })
			return nil
		},
		PreStop: func() error {
			return logkit.Sync()
		},
	}
	gz = &server.Command{
		RegisterServer: func(grpcServer *grpc.Server) {

		},
	}
)

func main() {
	// log： 1.调整了输出样式更加直观 2.自动打印context里边携带的requestid和traceid3. 结合lumberjack做了日志切割等4.用到了池化思想
	//defer logkit.Sync()
	//type request struct {
	//	A int
	//}
	//req := request{
	//	A: 1,
	//}
	//traceCtx := logkit.NewContextWith(context.Background(), trace.FieldRequestID("22"), trace.FieldTraceID("333"))
	//err := fmt.Errorf("lxlerr")
	////logger := logkit.FromContext(traceCtx).With(logkit.Err(err))
	//logkit.FromContext(traceCtx).Error("AreaModule.BatchGetLocationGroupsByDistrict rpc failed", logkit.Any("req", req), logkit.Err(err))
	//logkit.FromContext(traceCtx).Info("infolog")
	//
	// client:1.继承了requeset会自动调用validate等方法 2.自动绑定request 3.自动生成一个可以复现的curl命令 4.自动生成接口文档
	// 5.query携带了_show_request_id，那么就会返回requestid 6.提供了pprof 7.提供了捕捉panic的选项
	logkit.Error("success2222")
	cmd.Execute()
	//gz.Execute()
}
