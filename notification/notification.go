package notification

type Notification struct {
	from    string
	to      string
	content interface{}
}

func NewNotification(from string, to string, content interface{}) *Notification {
	return &Notification{from, to, content}
}

func (notification *Notification) From() string {
	return notification.from
}

func (notification *Notification) To() string {
	return notification.to
}

func (notification *Notification) Content() interface{} {
	return notification.content
}
