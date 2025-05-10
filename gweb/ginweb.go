package gweb

import (
	"bytes"
	"context"
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/gin-gonic/gin"
	"github.com/lee31802/comment_lib/conf"
	"github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/gweb/pprof"
	"github.com/lee31802/comment_lib/util"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// Stopper is callback invoked before ginweb has stopped.
type Stopper func() error

// gWeb is the gin application instance, it contains a *gin.Engine instance,
// module maps, configuration settings and environment variables.
type gWeb struct {
	router *router

	injector inject.Injector

	appPath string

	environ *env.Environ

	config *conf.Configuration

	opts *Options

	engine *gin.Engine

	modules map[string]ModuleInfo

	mu sync.RWMutex

	stopChan chan bool

	errChan chan error

	whenStops []Stopper
}

var gw *gWeb

// shortcuts
var (
	Config *conf.Configuration
	Env    *env.Environ
	Opts   *Options
)

func init() {
	gw = newGinWeb()
	Config = gw.config
	Env = gw.environ
	Opts = gw.opts
}

func (g *gWeb) initConfig(cmd Command) {
	appPath := cmd.AppPath
	if appPath == "" {
		appPath = util.GetWorkDir()
	}
	g.appPath = appPath
	*g.config = *initConfiguration(appPath, g.environ)
	g.opts.updateFromConfig(g.config)
}

func (g *gWeb) initBeforeRun(cmd Command) {
	g.initConfig(cmd)
	g.initComponents(cmd)
	g.initPlugins(cmd.Plugins)
	g.registerSignals()
}

// newGinWeb returns a newGinWeb application instance with given config.
func newGinWeb() *gWeb {
	// Default config

	opts := newOptions()
	//for _, setter := range options {
	//	setter(opts)
	//}
	appPath := util.GetWorkDir()
	defaultEnv := env.DefaultEnviron()
	config := conf.NewConfiguration()
	gs := &gWeb{
		injector: inject.New(),
		appPath:  appPath,
		opts:     opts,
		config:   config,
		environ:  defaultEnv,
		modules:  make(map[string]ModuleInfo),
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

func (g *gWeb) initPlugins(plugins []Plugin) {
	for _, plugin := range plugins {
		plugin.Install(g.injector, func(s Stopper) {
			g.whenStops = append(g.whenStops, s)
		})
	}
}

func (g *gWeb) initComponents(cmd Command) {
	// Init engine
	engine := cmd.Engine
	if engine == nil {
		engine = gin.New()
	}
	g.engine = engine
	// Init middlewares
	middlewares := cmd.Middlewares[:]
	if g.opts.Recovery {
		middlewares = append(middlewares, Recovery())
	}
	// Init pprof
	if g.opts.Pprof {
		pprof.Register(g.engine)
	}
	// Init router
	rg := g.engine.Group(g.opts.RootPath)
	g.router = &router{
		injector: g.injector,
		rg:       rg,
	}
	if g.opts.UploadMetrics {
		//p := ginprometheus.NewPrometheus()
		//p.SetGetHandlerNameFunc(g.GetHdlSimpleNameByUrl)
		//rg.Use(p.HandlerFunc(CtxKeyHandlerName))
	}

	if g.opts.Jaeger.Enable {
		//trace.InitJaegerTracer(env.GetService(), g.opts.Jaeger.SamplingRate)
		//middlewares = append(middlewares, jaeger.RequestTracing(pathHanlderMap))
	}

	//if env.SupportPfb() {
	//	middlewares = append(middlewares, pfb.PFB())
	//}

	for _, m := range middlewares {
		rg.Use(m)
	}
	//reporter.Init()
	//go reporter.ListenAndServe()
}

func (g *gWeb) registerSignals() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		g.Stop()
	}()
}

// RegisterModule registers a module to application.
func (g *gWeb) RegisterModule(m Module) error {
	pkgPath := reflect.TypeOf(m).Elem().PkgPath()
	pkgPath = strings.TrimLeft(pkgPath, "_")
	if _, ok := g.modules[pkgPath]; ok {
		return fmt.Errorf("can not register a duplicated module: %v", pkgPath)
	}
	module := ModuleInfo{
		Module:  m,
		AppPath: g.appPath,
		PkgPath: pkgPath,
		Environ: *g.environ,
	}
	debugPrint("register module: %v", pkgPath)
	// g.initModuleConfigs(module)
	g.mu.Lock()
	g.modules[pkgPath] = module
	g.mu.Unlock()
	m.Init(g.router)
	return nil
}

// Engine returns the underlying *gin.Engine instance.
func (g *gWeb) Engine() *gin.Engine {
	return g.engine
}

func (g *gWeb) registerAPIView() {
	debugPrint("register handlers:")
	globalHandlerInfos.prettyPrint(false)
	g.engine.GET("/gwapi", func(c *gin.Context) {
		t := template.New("gwapi")
		t = t.Funcs(template.FuncMap{
			"extractJson": extractJson,
		})
		// 解析模板
		t, _ = t.Parse(apiDoc)

		buffer := bytes.NewBuffer([]byte{})
		sort.Slice(globalHandlerInfos[:], func(i, j int) bool {
			return globalHandlerInfos[i].Method+globalHandlerInfos[i].URL < globalHandlerInfos[j].Method+globalHandlerInfos[j].URL
		})
		t.Execute(buffer, map[string]interface{}{
			"title": "API Doc",
			"apis":  globalHandlerInfos,
		})
		c.Writer.Header().Set("Content-Type", "text/html")
		c.String(200, buffer.String())
	})
}

func extractJson(curl string) string {
	start := strings.Index(curl, "{")
	if start == -1 {
		return ""
	}
	end := strings.LastIndex(curl, "}")
	if end == -1 {
		return ""
	}
	return curl[start : end+1]
}

// Execute starts listening and serving HTTP requests.
func (g *gWeb) Run(cmd Command) error {
	gin.SetMode(ReleaseMode) // disable gin's debug output
	if g.environ.Env == "live" {
		SetMode(ReleaseMode)
	}
	g.initBeforeRun(cmd)
	if cmd.PreRun != nil {
		if err := cmd.PreRun(g.router); err != nil {
			debugPrint("service PreRun() error: %v", err.Error())
			return err
		}
	}
	for _, m := range cmd.Modules {
		if err := g.RegisterModule(m); err != nil {
			debugPrint("register module failed: %v", err.Error())
			return err
		}
	}

	addrs := []string{}
	if g.opts.Address != "" {
		addrs = append(addrs, g.opts.Address)
	}
	address := util.ResolveAddress(addrs)
	server := &http.Server{
		Addr:    address,
		Handler: g.engine,
	}

	if serviceMode == DebugMode {
		g.registerAPIView()
		debugPrint("API docs address: %v/gwapi", address)
	}
	go func() {
		debugPrint("Listening on %v", address)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("client listen error: %s\n", err)
		}
		g.errChan <- err
	}()
	if cmd.PostRun != nil {
		if err := cmd.PostRun(); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		debugPrint("client shutdown error: %v", err.Error())
		retErr = err
	}
	if cmd.PostStop != nil {
		if err = cmd.PostStop(); err != nil {
			debugPrint("service AfterStop() error: %v", err.Error())
			return err
		}
	}
	return retErr
}

// Stop terminates the application.
func (g *gWeb) Stop() {
	go func() { g.stopChan <- true }()
}

func (g *gWeb) GetHdlSimpleNameByUrl(url string, method string) string {
	key := method + ":" + url
	v, ok := pathHandlerMap[key]
	if !ok {
		return ""
	}
	return v
}
