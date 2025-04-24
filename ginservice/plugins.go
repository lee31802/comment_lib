package ginservice

import (
	"github.com/codegangsta/inject"
)

// Plugin is a module without router, use for integrating outside components into application.
type Plugin interface {
	Install(inject.TypeMapper, func(Stopper))
}
