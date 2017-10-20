package notification

type Origin string
type Destination string
type Object interface{}

type Notification struct {
	from    Origin
	to      Destination
	content Object
}

func NewNotification(from Origin, to Destination, content Object) *Notification {
	return &Notification{from, to, content}
}

func (n *Notification) From() Origin {
	return n.from
}

func (n *Notification) To() Destination {
	return n.to
}

func (n *Notification) Content() Object {
	return n.content
}

func (n *Notification) Cast(c Caster) (Object, error) {
	return c.cast(n.from, n.content)
}
