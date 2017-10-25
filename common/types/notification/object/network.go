package object

type NetworkObject struct {
	command string
	args    []string
}

func NewNetworkObject(command string, args ...string) *NetworkObject {
	return &NetworkObject{command: command, args: args}
}

func (networkObject *NetworkObject) AddArgument(args ...string) *NetworkObject {
	networkObject.args = append(networkObject.args, args...)
	return networkObject
}

func (networkObject *NetworkObject) Command() string {
	return networkObject.command
}

func (networkObject *NetworkObject) Arguments() []string {
	return networkObject.args
}
