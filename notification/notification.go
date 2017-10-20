package notification

type Origin string
type Destination string
type Object interface{}

type Notification struct {
	_from Origin
	_to   Destination
	_content Object
}

func (n *Notification) From() Origin  {
	return n._from
}

func (n *Notification) To() Destination  {
	return n._to
}

func (n *Notification) Content() Object  {
	return n._content
}

func (n *Notification) cast(c Caster) (Object, error) {
	return c.cast(n._from, n._content)
}