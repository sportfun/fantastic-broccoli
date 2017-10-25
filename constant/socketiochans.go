package constant

import (
	"github.com/graarh/golang-socketio"
)

const (
	CommandChan = "command"
	DataChan    = "data"
	ErrorChan   = gosocketio.OnError
)
