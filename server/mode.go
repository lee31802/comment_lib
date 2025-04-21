package server

// Ginweb mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

var (
	serviceMode = DebugMode
)

// SetMode changes ginweb's mode.
func SetMode(value string) {
	serviceMode = value
}
