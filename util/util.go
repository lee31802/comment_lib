package util

import (
	"fmt"
	"os"
)

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if servicePort := os.Getenv("PORT"); servicePort != "" {
			return ":" + servicePort
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

func GetWorkDir() string {
	wd := os.Getenv("APP_PATH")
	if len(wd) > 0 {
		return wd
	}
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

func ColorFormat(color string, s interface{}) string {
	if ss, ok := s.(string); ok && len(ss) == 0 {
		return ""
	}
	if color == "green" {
		color = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	} else if color == "white" {
		color = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	} else if color == "yellow" {
		color = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	} else if color == "blue" {
		color = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	} else if color == "cyan" {
		color = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	} else {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("%s %v %s", color, s, string([]byte{27, 91, 48, 109}))
}
