package gweb

import "github.com/gin-gonic/gin"

// Structure to manage groups for commands
type Group struct {
	ID    string
	Title string
}

type Command struct {
	Name    string
	AppPath string
	// 可以自定义engine
	Engine  *gin.Engine
	Plugins []Plugin
	// 全局的
	Middlewares []gin.HandlerFunc
	// PreRun: children of this command will not inherit.
	PreRun func(router Router) error
	// PostRun: run after the Run command.
	PostRun  func() error
	PreStop  func() error
	PostStop func() error
	Modules  []Module
}

func (cmd Command) Execute() error {
	//if cmd.Middlewares != nil {
	//	Configure(WithMiddlewares(cmd.Middlewares...))
	//}
	return gw.Run(cmd)
}
