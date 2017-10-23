package notification

import "fantastic-broccoli/common/types"

type Object interface{}

type Notification struct {
	from    types.Name
	to      types.Name
	content Object
}

func NewNotification(from types.Name, to types.Name, content Object) *Notification {
	return &Notification{from, to, content}
}

func (n *Notification) From() types.Name {
	return n.from
}

func (n *Notification) To() types.Name {
	return n.to
}

func (n *Notification) Content() Object {
	return n.content
}
