package constant

import "github.com/xunleii/fantastic-broccoli/common/types"

var ErrorLevels = struct {
	Warning  types.ErrorLevel
	Error    types.ErrorLevel
	Critical types.ErrorLevel
	Fatal    types.ErrorLevel
}{
	Warning:  "Warning",
	Error:    "Error",
	Critical: "Critical",
	Fatal:    "Fatal",
}
