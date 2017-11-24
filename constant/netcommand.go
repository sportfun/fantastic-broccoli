package constant

import "github.com/xunleii/fantastic-broccoli/common/types"

var NetCommand = struct {
	State          types.CommandName
	Link           types.CommandName
	StartSession   types.CommandName
	EndSession     types.CommandName
	RestartService types.CommandName
}{
	State:          "state",
	Link:           "link",
	StartSession:   "start_session",
	EndSession:     "end_session",
	RestartService: "restart_service",
}
