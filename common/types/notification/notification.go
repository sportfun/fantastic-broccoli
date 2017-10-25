package notification

type Notification struct {
	from    string
	to      string
	content interface{}
}

func NewNotification(from string, to string, content interface{}) *Notification {
	return &Notification{from, to, content}
}

func (n *Notification) From() string {
	return n.from
}

func (n *Notification) To() string {
	return n.to
}

func (n *Notification) Content() interface{} {
	return n.content
}
