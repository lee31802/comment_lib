package service

import (
	"github.com/spf13/cobra"
	// 导入其他必要的包
)

// DefaultOptions 定义默认配置选项
func DefaultOptions() []app.Option {
	return []app.Option{
		app.WithPort(8082),
		app.WithUseMysql(true),
		app.WithUseRedis(true),
		app.WithEnableOfProm(true),
	}
}

// RunApp 封装通用的应用运行逻辑
func RunApp(options []app.Option, modules ...app.Module) {
	app.ApplyOption(options...)
	app.Run(modules...)
}

// InitRootCmd 初始化根命令
func InitRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "vip-server",
		PreRun: func(cmd *cobra.Command, args []string) {
			initVipPre()
		},
		Run: func(cmd *cobra.Command, args []string) {
			defaultOptions := DefaultOptions()
			// 这里可以添加更多的通用逻辑
			RunApp(defaultOptions, modules...)
		},
	}
	return rootCmd
}

// initVipPre 通用的预初始化逻辑
func initVipPre() {
	// 实现预初始化逻辑
}
