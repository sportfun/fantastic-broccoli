package v1_0

import "github.com/sportfun/gakisitor/plugin"

type CommandPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Command string        `json:"command"`
		Args    []interface{} `json:"args"`
	} `json:"body"`
}

type DataPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Module string      `json:"module"`
		Value  interface{} `json:"value"`
	} `json:"body"`
}

type ErrorPacket struct {
	LinkId string `json:"link_id"`
	Body struct {
		Origin string `json:"origin"`
		Reason string `json:"reason"`
	} `json:"body"`
}

type channelID byte

const (
	Command channelID = iota
	Data    channelID = iota
	Error   channelID = iota
)

var Channels = map[channelID]string{
	Command: "command",
	Data:    "data",
	Error:   "error",
}

var Instructions = map[string]plugin.Instruction{
	"start_session": plugin.StartSessionInstruction,
	"end_session":   plugin.StopSessionInstruction,
}
