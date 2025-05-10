package server

import (
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/lee31802/comment_lib/conf"
	"github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/gzero"
	"github.com/lee31802/comment_lib/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"os/signal"
	"path"
	"sync"
)

// Stopper is callback invoked before ginweb has stopped.
type Stopper func() error

type goZero struct {
	appPath  string
	injector inject.Injector
	environ  *env.Environ
	config   *conf.Configuration
	opts     *Options
	mu       sync.RWMutex

	stopChan chan bool

	errChan chan error

	whenStops []Stopper
}

var gs *goZero

// shortcuts
var (
	Config *conf.Configuration
	Env    *env.Environ
	Opts   *Options
)

func init() {
	gs = newGoZeroServer()
	Config = gs.config
	Env = gs.environ
	Opts = gs.opts
	logx.Disable()
}

func (g *goZero) initConfig(cmd Command) {
	appPath := cmd.AppPath
	if appPath == "" {
		appPath = util.GetWorkDir()
	}
	g.appPath = appPath
	*g.config = *initConfiguration(appPath, g.environ)
	g.opts.updateFromConfig(g.config)
}

func (g *goZero) initBeforeRun(cmd Command) {
	g.initConfig(cmd)
	g.initPlugins(cmd.Plugins)
	g.registerSignals()
}

func (g *goZero) initPlugins(plugins []Plugin) {
	for _, plugin := range plugins {
		plugin.Install(g.injector, func(s Stopper) {
			g.whenStops = append(g.whenStops, s)
		})
	}
}

// newGoZeroServer returns a newGoZeroServer application instance with given config.
func newGoZeroServer(options ...Option) *goZero {
	// Default config
	opts := newOptions()
	for _, setter := range options {
		setter(opts)
	}
	appPath := util.GetWorkDir()
	defaultEnv := env.DefaultEnviron()
	config := conf.NewConfiguration()
	server := &goZero{
		injector: inject.New(),
		appPath:  appPath,
		opts:     opts,
		config:   config,
		environ:  defaultEnv,
		stopChan: make(chan bool),
		errChan:  make(chan error),
	}
	return server
}
func initConfiguration(appPath string, env *env.Environ) *conf.Configuration {
	config := conf.NewConfiguration()
	configName := fmt.Sprintf("config_%v.yml", env.Env)
	configPath := path.Join(appPath, "conf", configName)
	defaultConfigPath := path.Join(appPath, "conf", "config.yml")
	for _, path := range []string{configPath, defaultConfigPath} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			gzero.DebugPrint("load app config error: %v", err.Error())
		} else {
			gzero.DebugPrint("load app config: %v", path)
			config.Apply(path)
			break
		}
	}
	return config
}

func (g *goZero) registerSignals() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		g.Stop()
	}()
}

// Execute starts listening and serving HTTP requests.
func (g *goZero) Run(cmd Command) error {
	g.initBeforeRun(cmd)
	if cmd.PreRun != nil {
		if err := cmd.PreRun(); err != nil {
			gzero.DebugPrint("service PreRun() error: %v", err.Error())
			return err
		}
	}
	s := &zrpc.RpcServer{}
	go func() {
		s = zrpc.MustNewServer(g.opts.Server, cmd.RegisterServer)
		s.AddUnaryInterceptors(cmd.UnaryInterceptors...)
		s.AddStreamInterceptors(cmd.StreamInterceptors...)
		s.Start()
	}()
	if cmd.PostRun != nil {
		if err := cmd.PostRun(); err != nil {
			g.errChan <- err
			gzero.DebugPrint("service AfterStart() error: %v", err.Error())
			return err
		}
	}
	var retErr error
	select {
	case <-g.stopChan:
		gzero.DebugPrint("receive stop signal")
		break
	case err := <-g.errChan:
		gzero.DebugPrint("error: %v", err.Error())
		retErr = err
		break
	}
	if cmd.PreStop != nil {
		if err := cmd.PreStop(); err != nil {
			gzero.DebugPrint("service BeforeStop() error: %v", err.Error())
			return err
		}
	}
	for _, callback := range g.whenStops {
		err := callback()
		if err != nil {
			gzero.DebugPrint("callback error: %v", err.Error())
		}
	}
	s.Stop()
	if cmd.PostStop != nil {
		if err := cmd.PostStop(); err != nil {
			gzero.DebugPrint("service AfterStop() error: %v", err.Error())
			return err
		}
	}
	return retErr
}

// Stop terminates the application.
func (g *goZero) Stop() {
	go func() { g.stopChan <- true }()
}
