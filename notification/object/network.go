package object

type CommandObject struct {
	Command string   `json:"command" mapstructure:"command"`
	Args    []string `json:"args" mapstructure:"args"`
}

func NewCommandObject(command string, args ...string) *CommandObject {
	if len(args) == 0 {
		return &CommandObject{Command: command, Args: []string{}}
	}
	return &CommandObject{Command: command, Args: args}
}

func (networkObject *CommandObject) AddArgument(args ...string) *CommandObject {
	networkObject.Args = append(networkObject.Args, args...)
	return networkObject
}
