package event

type Event struct{}
type Bus struct{}

type EventConsumer func(event Event, err error)

func (event *Event) Message() interface{} { panic("not implemented") }
func (event *Event) Reply(interface{})    { panic("not implemented") }

func (bus *Bus) Publish(channel string, data interface{}, handler ...ReplyHandler)                { panic("not implemented") }
func (bus *Bus) Subscribe(channel string, handler EventConsumer, handlers ...EventConsumer) error { panic("not implemented") }
func (bus *Bus) Unsubscribe(channel string, handler EventConsumer) error                          { panic("not implemented") }
