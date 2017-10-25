package object

type NetworkObject struct {
	command string
	args    []string
}

func NewNetworkObject(command string, args ...string) *NetworkObject {
	return &NetworkObject{command: command, args: args}
}

func (m *NetworkObject) AddArgument(args ...string) *NetworkObject {
	m.args = append(m.args, args...)
	return m
}

func (m *NetworkObject) Command() string {
	return m.command
}

func (m *NetworkObject) Arguments() []string {
	return m.args
}
