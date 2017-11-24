package object

import "github.com/xunleii/fantastic-broccoli/common/types"

type CommandObject struct {
	Command types.CommandName `json:"command" mapstructure:"command"`
	Args    []string          `json:"args" mapstructure:"args"`
}

func NewCommandObject(command types.CommandName, args ...string) *CommandObject {
	return &CommandObject{Command: command, Args: args}
}

func (networkObject *CommandObject) AddArgument(args ...string) *CommandObject {
	networkObject.Args = append(networkObject.Args, args...)
	return networkObject
}
