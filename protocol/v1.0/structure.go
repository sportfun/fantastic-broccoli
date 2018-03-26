package v1_0

import "github.com/sportfun/gakisitor/plugin"

// CommandPacket implements the command packet of the protocol.
type CommandPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Command string        `json:"command"`
		Args    []interface{} `json:"args"`
	} `json:"body"`
}

// DataPacket implements the data packet of the protocol.
type DataPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Module string      `json:"module"`
		Value  interface{} `json:"value"`
	} `json:"body"`
}

// ErrorPacket implements the error packet of the protocol.
type ErrorPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Origin string `json:"origin"`
		Reason string `json:"reason"`
	} `json:"body"`
}

type channelID byte

// List of channels id
const (
	Command channelID = iota
	Data
	Error
)

// List of channel names
var Channels = map[channelID]string{
	Command: "command",
	Data:    "data",
	Error:   "error",
}


// List of protocol instructions
var Instructions = map[string]plugin.Instruction{
	"start_session": plugin.StartSessionInstruction,
	"end_session":   plugin.StopSessionInstruction,
}
