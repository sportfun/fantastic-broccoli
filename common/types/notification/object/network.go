package object

type NetworkObject struct {
	Command string   `json:"command" mapstructure:"command"`
	Args    []string `json:"args" mapstructure:"args"`
}

func NewNetworkObject(command string, args ...string) *NetworkObject {
	return &NetworkObject{Command: command, Args: args}
}

func (networkObject *NetworkObject) AddArgument(args ...string) *NetworkObject {
	networkObject.Args = append(networkObject.Args, args...)
	return networkObject
}
