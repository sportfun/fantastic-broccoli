package network

type Message struct {
	command string
	args    []string
}

func NewMessage(command string) *Message  {
	m := new(Message)
	m.command = command
	return m
}

func (m *Message) AddArgument(arg string) *Message {
	m.args = append(m.args, arg)
	return m
}

func (m *Message) Command() string  {
	return m.command
}

func (m *Message) Arguments() []string  {
	return m.args
}