package gzero

import (
	"fmt"
	"os"
	"strings"
)

func debugPrint(format string, values ...interface{}) {

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, "[GoZero] "+format, values...)

}
