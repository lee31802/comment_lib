package logkit

import (
	"fmt"
	"github.com/lee31802/comment_lib/env"
	"os"
	"strings"
)

func DebugPrint(format string, values ...interface{}) {
	if env.Environment() == "test" {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[logkit] "+format, values...)
	}
}
