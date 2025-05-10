package logkit

import (
	"fmt"
	"github.com/lee31802/comment_lib/conf"
	env2 "github.com/lee31802/comment_lib/env"
	"github.com/lee31802/comment_lib/util"
	"os"
	"path"
)

func initConfiguration(appPath string) (bool, *conf.Configuration) {
	if appPath == "" {
		appPath = util.GetWorkDir()
	}
	env := env2.Environment()
	config := conf.NewConfiguration()
	configName := fmt.Sprintf("config_%v.yml", env)
	configPath := path.Join(appPath, "conf", configName)
	defaultConfigPath := path.Join(appPath, "conf", "config.yml")
	has := false
	for _, path := range []string{configPath, defaultConfigPath} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			DebugPrint("load app config error: %v", err.Error())
		} else {
			DebugPrint("load app config: %v", path)
			has = true
			config.Apply(path)
			break
		}
	}
	return has, config
}
