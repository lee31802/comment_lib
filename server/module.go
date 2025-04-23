package server

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	configDir = "conf"
	i18nDir   = "i18n"
)

// Module represents a sub-system that was attached to the application.
type Module interface {
	// Init will be invoked when an application tries to register the module intance to itself,
	// the first argument is the application instance. Module initialization should be done here.
	Init(Router)
}

// ModuleInfo contains informations about module.
type ModuleInfo struct {
	Module  Module
	AppPath string
	PkgPath string
	Environ Environ
}

// ConfigPath returns config location of the module.
func (m *ModuleInfo) ConfigPath() (string, bool) {
	configName := fmt.Sprintf("config_%v.yml", m.Environ.Env)
	if path, exists := m.configPath(configDir, configName); exists {
		return path, exists
	}
	return m.configPath(configDir, "config.yml") // fallback
}

func (m *ModuleInfo) configPath(configDir string, configName string) (string, bool) {
	var (
		modulePath string
		moduleName string
	)
	pkgs := strings.Split(m.PkgPath, string(os.PathSeparator))
	if len(pkgs) > 1 {
		moduleName = pkgs[len(pkgs)-1]
	}

	// 1. Search from Application-Path.
	appPath := m.AppPath
	if appPath != "" {
		absPath, err := filepath.Abs(appPath)
		if err == nil {
			modulePath = path.Join(absPath, moduleName)
			if cfgPath, err := m.searchConfig(modulePath, configDir, configName); err == nil {
				return cfgPath, true
			}
		}
	}
	// 2. Search from Work-Path and GOPATH.
	if filepath.IsAbs(m.PkgPath) {
		modulePath = m.PkgPath
	} else {
		workPath, _ := os.Getwd()
		modulePath = path.Join(workPath, moduleName)
		if _, err := os.Stat(modulePath); os.IsNotExist(err) {
			modulePath = path.Join(os.Getenv("GOPATH"), "src", m.PkgPath)
		}
	}
	if cfgPath, err := m.searchConfig(modulePath, configDir, configName); err == nil {
		return cfgPath, true
	}
	return "", false
}

func (m *ModuleInfo) searchConfig(modulePath string, configDir string, configName string) (string, error) {
	configPath := path.Join(modulePath, configDir, configName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", err
	}
	return configPath, nil
}
