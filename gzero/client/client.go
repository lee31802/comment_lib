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

func initConfiguration(appPath, configPath string, env *env.Environ) *conf.Configuration {
	config := conf.NewConfiguration()
	configName := fmt.Sprintf("config_%v.yml", env.Env)
	if configPath == "" {
		configPath = path.Join(appPath, "conf")
	}
	configPath = path.Join(configPath, configName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		gzero.DebugPrint("load app config error: %v", err.Error())
	} else {
		gzero.DebugPrint("load app config: %v", configPath)
	}

	return config
}

func (c *client) initConfig() {
	if c.opts.serviceName == "" {
		panic("service name not set")
	}
	configPath := c.opts.configPath
	*c.config = *initConfiguration(c.appPath, configPath, c.environ)
	c.opts.updateFromConfig(c.config)
}
