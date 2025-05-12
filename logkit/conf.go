package logkit

import (
	"fmt"
	"github.com/lee31802/comment_lib/conf"
	env2 "github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/util"
	"os"
	"path"
)

func initConfiguration(configPath string) (bool, *conf.Configuration) {
	env := env2.Environment()
	config := conf.NewConfiguration()
	configName := fmt.Sprintf("config_%v.yml", env)
	if configPath == "" {
		appPath := util.GetWorkDir()
		configPath = path.Join(appPath, "conf")
	}
	logOpts.configPath = configPath
	configPath = path.Join(configPath, configName)
	has := false

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		DebugPrint("load app config error: %v", err.Error())
	} else {
		DebugPrint("load app config: %v", configPath)
		has = true
		config.Apply(configPath)
	}

	return has, config
}
