package event

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Event struct {
	payload interface{}
	reply   chan interface{}
	error   chan error
}
type subscriber struct {
	id     string
	ch     chan<- *Event
	cancel func()
}

type Bus struct {
	subscribers map[string][]subscriber
	ids         map[string]interface{}
	sync        sync.Mutex
}
type EventConsumer func(event *Event, err error)

const publishTimeout = 25 * time.Millisecond
const replyTimeout = 25 * time.Millisecond

var ErrPublishTimeout = errors.New("publish timeout")
var ErrSubscriberClosed = errors.New("subscriber closed")
var ErrChannelNotFound = errors.New("channel not found")
var ErrChannelSubscriberNotFound = errors.New("channel subscriber not found")
var ErrChannelSubscriberAlreadyExists = errors.New("channel subscriber already exists")

func (event *Event) Message() interface{}      { return event.payload }
func (event *Event) Reply() chan<- interface{} { return event.reply }

func NewBus() *Bus { return &Bus{subscribers: map[string][]subscriber{}, ids: map[string]interface{}{}, sync: sync.Mutex{}} }

func (bus *Bus) Publish(channel string, data interface{}, handler ReplyHandler) {
	if _, exists := bus.subscribers[channel]; !exists {
		if handler != nil {
			handler.consume(nil, nil, ErrChannelNotFound, 0)
		}
		return
	}

	bus.sync.Lock()
	defer bus.sync.Unlock()

	for _, evChannel := range bus.subscribers[channel] {
		event := &Event{payload: data, reply: make(chan interface{}), error: make(chan error)}
		go func(ch chan<- *Event, e *Event) {
			select {
			case ch <- e:
			case <-time.After(replyTimeout):
				event.error <- ErrPublishTimeout
			}
		}(evChannel.ch, event)

		if handler != nil {
			handler.consume(event.reply, event.error, nil, replyTimeout)
		}
	}
}

func (bus *Bus) Subscribe(channel string, handler EventConsumer) error {
	ch := make(chan *Event)
	ctx, cnl := context.WithCancel(context.Background())

	go func(channel <-chan *Event, ctx context.Context) {
		defer handler(nil, ErrSubscriberClosed)

		for {
			select {
			case <-ctx.Done():
				return
			case event, o := <-channel:
				if !o {
					return
				}

				handler(event, nil)
			}
		}
	}(ch, ctx)

	id := id(channel, handler)

	bus.sync.Lock()
	defer bus.sync.Unlock()

	if _, exists := bus.ids[id]; exists {
		return ErrChannelSubscriberAlreadyExists
	}

	bus.subscribers[channel] = append(bus.subscribers[channel], subscriber{id: id, ch: ch, cancel: cnl})
	bus.ids[id] = nil
	return nil
}

func (bus *Bus) Unsubscribe(channel string, handler EventConsumer) error {
	if sub, exists := bus.subscribers[channel]; !exists {
		return ErrChannelNotFound
	} else {
		id := id(channel, handler)

		for i, sbcr := range sub {
			if sbcr.id == id {
				bus.sync.Lock()
				sbcr.cancel()
				bus.subscribers[channel] = append(sub[:i], sub[i+1:]...)
				delete(bus.ids, id)

				if len(sub) == 0 {
					delete(bus.subscribers, channel)
				}
				bus.sync.Unlock()
				return nil
			}
		}

		return ErrChannelSubscriberNotFound
	}
}

func id(c string, h interface{}) string { return c + ":" + runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() }
