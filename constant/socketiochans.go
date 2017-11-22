package constant

import (
	"github.com/graarh/golang-socketio"
	"github.com/xunleii/fantastic-broccoli/common/types"
)

var Channels = struct {
	Command types.ChannelName
	Data    types.ChannelName
	Error   types.ChannelName
}{
	Command: "command",
	Data:    "data",
	Error:   gosocketio.OnError,
}
