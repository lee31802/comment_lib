package gweb

import (
	"fmt"
	"github.com/lee31802/comment_lib/util"
	"os"
	"strings"
)

func debugPrint(format string, values ...interface{}) {
	if serviceMode == DebugMode {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[GinWeb] "+format, values...)
	}
}

func (infos handlerInfoList) prettyPrint(verbose bool) {
	for _, info := range infos {
		debugPrint(fmt.Sprintf("%s|%s|%s", util.ColorFormat("blue", info.HandlerName), util.ColorFormat("green", info.Method), util.ColorFormat("cyan", info.URL)))
		if info.Request != nil && verbose {
			for _, field := range info.Request.FieldInfos {
				debugPrint(fmt.Sprintf("   |- %-25s %-10v %-10v", util.ColorFormat("cyan", field.Name), util.ColorFormat("white", field.Typ), util.ColorFormat("white", field.Tag)))
			}
		}
	}
}
