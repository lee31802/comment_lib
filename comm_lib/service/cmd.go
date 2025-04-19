package service

type RootCmd struct {
	BaseCmd
}

type BaseCmd interface {
	Name() string
	BeforeStart(Router) error
	AfterStart() error
	BeforeStop() error
	AfterStop() error
	Modules() []Module
}
