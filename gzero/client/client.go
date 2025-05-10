package client

import (
	"fmt"
	"github.com/lee31802/comment_lib/conf"
	"github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/gzero"
	"github.com/lee31802/comment_lib/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"os"
	"path"
)

type client struct {
	appPath     string
	serviceName string
	environ     *env.Environ
	config      *conf.Configuration
	opts        *Options
}

func init() {
	logx.Disable()
}

func newGoZeroClient(options ...Option) *client {
	// Default config
	opts := newOptions()
	for _, setter := range options {
		setter(opts)
	}
	appPath := util.GetWorkDir()
	defaultEnv := env.DefaultEnviron()
	config := conf.NewConfiguration()
	goClient := &client{
		appPath: appPath,
		opts:    opts,
		config:  config,
		environ: defaultEnv,
	}
	return goClient
}

// var DefaultClient *zrpc.Client
func (c *client) initClient() zrpc.Client {
	c.initConfig()
	return zrpc.MustNewClient(c.opts.Client)
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

func (c *client) initConfig() {
	if c.opts.serviceName == "" {
		panic("service name not set")
	}
	appPath := c.opts.appPath
	if appPath == "" {
		appPath = util.GetWorkDir()
	}
	c.appPath = appPath
	*c.config = *initConfiguration(appPath, c.environ)
	c.opts.updateFromConfig(c.config)
}
