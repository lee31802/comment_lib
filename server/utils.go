package server

import (
	"fmt"
	"os"
	"strings"
)

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

func getWorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

func debugPrint(format string, values ...interface{}) {
	if serverMode == DebugMode {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[Ginweb] "+format, values...)
	}
}
