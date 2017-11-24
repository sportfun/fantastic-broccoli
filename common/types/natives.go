package types

type ErrorLevel string
func (o ErrorLevel) String() string { return string(o) }

type CommandName string
func (o CommandName) String() string { return string(o) }

type ChannelName string
func (o ChannelName) String() string { return string(o) }

type StateType byte
func (o StateType) Byte() byte { return byte(o) }