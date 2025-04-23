package server

// Ginweb mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

var (
	serverMode = DebugMode
)

// SetMode changes ginweb's mode.
func SetMode(value string) {
	serverMode = value
}
