package gweb

import (
	"os"
	"strings"
)

var (
	env   = getEnvLower("env")
	host  = getEnvLower("HOST")
	port  = getEnvLower("PORT")
	appID = getEnvLower("MARATHON_APP_ID")
)

// Environ holds basic environment variables.
type Environ struct {
	Env   string `json:"env"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	AppID string `json:"appid"`
}

func defaultEnviron() *Environ {
	return &Environ{
		AppID: AppID(),
		Env:   Environment(),
		Host:  Host(),
		Port:  Port(),
	}
}

func getEnvLower(key string) string {
	return strings.TrimSpace(strings.ToLower(os.Getenv(key)))
}

// Environment gets the environment from the environmental variable 'env'.
// It always returns lower case.
func Environment() string {
	return env
}

// Port gets the port from the environmental variable 'PORT'.
func Port() string {
	return port
}

// Host gets the host from the environmental variable 'HOST'.
// It always returns lower case.
func Host() string {
	return host
}

// AppID gets the host from the environmental variable 'MARATHON_APP_ID'.
// It always returns lower case.
func AppID() string {
	return appID
}
