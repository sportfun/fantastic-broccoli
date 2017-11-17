package constant

import (
	"github.com/graarh/golang-socketio"
)

var Channels = struct {
	Command string
	Data    string
	Error   string
}{
	Command: "command",
	Data:    "data",
	Error:   gosocketio.OnError,
}
