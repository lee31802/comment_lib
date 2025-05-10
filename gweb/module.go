package gweb

import (
	"github.com/lee31802/comment_lib/env"
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
	Environ env.Environ
}
