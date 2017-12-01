package network

import "github.com/graarh/golang-socketio"

const (
	OnConnection    = gosocketio.OnConnection
	OnDisconnection = gosocketio.OnDisconnection
	OnCommand       = "command"
	OnData          = "data"
	OnError         = gosocketio.OnError
)
