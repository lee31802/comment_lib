package server

import (
	"fmt"
	"github.com/lee31802/comment_lib/conf"
	"github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/util"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"os/signal"
	"path"
	"sync"
)

// Stopper is callback invoked before ginweb has stopped.
type Stopper func() error

type goZeroServer struct {
	appPath string
	environ *env.Environ
	config  *conf.Configuration
	opts    *Options
	mu      sync.RWMutex

	stopChan chan bool

	errChan chan error

	whenStops []Stopper
}

var gs *goZeroServer

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
}

// Configure default ginweb app options.
func Configure(options ...Option) {
	for _, setter := range options {
		setter(gs.opts)
	}
}

func (g *goZeroServer) initConfig() {
	appPath := g.opts.AppPath
	if appPath == "" {
		appPath = util.GetWorkDir()
	}
	g.appPath = appPath
	*g.config = *initConfiguration(appPath, g.environ)
	g.opts.updateFromConfig(g.config)
}

func (g *goZeroServer) initBeforeRun() {
	g.initConfig()
	g.registerSignals()
}

// newGoZeroServer returns a newGoZeroServer application instance with given config.
func newGoZeroServer(options ...Option) *goZeroServer {
	// Default config
	opts := newOptions()
	for _, setter := range options {
		setter(opts)
	}
	appPath := util.GetWorkDir()
	defaultEnv := env.DefaultEnviron()
	config := conf.NewConfiguration()
	gs := &goZeroServer{
		appPath:  appPath,
		opts:     opts,
		config:   config,
		environ:  defaultEnv,
		stopChan: make(chan bool),
		errChan:  make(chan error),
	}
	return gs
}
func initConfiguration(appPath string, env *env.Environ) *conf.Configuration {
	config := conf.NewConfiguration()
	configName := fmt.Sprintf("config_%v.yml", env.Env)
	configPath := path.Join(appPath, "conf", configName)
	defaultConfigPath := path.Join(appPath, "conf", "config.yml")
	for _, path := range []string{configPath, defaultConfigPath} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			debugPrint("load app config error: %v", err.Error())
		} else {
			debugPrint("load app config: %v", path)
			config.Apply(path)
			break
		}
	}
	return config
}

func (g *goZeroServer) registerSignals() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		g.Stop()
	}()
}

// Execute starts listening and serving HTTP requests.
func (g *goZeroServer) Run(cmd Command) error {

	g.initBeforeRun()
	if cmd.PreRun != nil {
		if err := cmd.PreRun(); err != nil {
			debugPrint("service PreRun() error: %v", err.Error())
			return err
		}
	}
	s := &zrpc.RpcServer{}
	go func() {
		debugPrint("Listening on %v", g.opts.ListenOn)
		s = zrpc.MustNewServer(g.opts.RpcServerConf, cmd.RegisTerFunc)
		s.AddUnaryInterceptors(g.opts.UnaryInterceptors...)
		s.AddStreamInterceptors(g.opts.StreamInterceptors...)
		s.Start()
	}()
	if cmd.PostRun != nil {
		if err := cmd.PostRun(); err != nil {
			g.errChan <- err
			debugPrint("service AfterStart() error: %v", err.Error())
			return err
		}
	}
	var retErr error
	select {
	case <-g.stopChan:
		debugPrint("receive stop signal")
		break
	case err := <-g.errChan:
		debugPrint("error: %v", err.Error())
		retErr = err
		break
	}
	if cmd.PreStop != nil {
		if err := cmd.PreStop(); err != nil {
			debugPrint("service BeforeStop() error: %v", err.Error())
			return err
		}
	}
	for _, callback := range g.whenStops {
		err := callback()
		if err != nil {
			debugPrint("callback error: %v", err.Error())
		}
	}
	s.Stop()
	if cmd.PostStop != nil {
		if err := cmd.PostStop(); err != nil {
			debugPrint("service AfterStop() error: %v", err.Error())
			return err
		}
	}
	return retErr
}

// Stop terminates the application.
func (g *goZeroServer) Stop() {
	go func() { g.stopChan <- true }()
}
