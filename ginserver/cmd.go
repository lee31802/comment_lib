package ginserver

// Structure to manage groups for commands
type Group struct {
	ID    string
	Title string
}

type Command struct {
	Name string
	// PreRun: children of this command will not inherit.
	PreRun func(router Router) error
	// PostRun: run after the Run command.
	PostRun  func() error
	PreStop  func() error
	PostStop func() error
	Modules  []Module
}

func (cmd Command) Execute() error {
	return s.Run(cmd)
}
