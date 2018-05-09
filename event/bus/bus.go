package bus

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Event represents a message sent through the Bus. It provides methods to get
// the payload (ak. message) and to reply to the sender.
type Event struct {
	payload interface{}
	reply   chan interface{}
	error   chan error
}

// subscriber is an internal manager for the goroutine in charge of handling
// events received through the Bus.
type subscriber struct {
	id     string
	ch     chan<- *Event
	cancel func()
}

// Bus is an implementation of the Sub/Pub design pattern. It provides a simple
// way to send data to several handlers, at same times. It also manages the
// handlers to prevent some errors like handler's crash.
type Bus struct {
	subscribers map[string][]subscriber
	ids         map[string]interface{}
	sync        sync.Mutex
}

// EventConsumer is the handler in charge of receiving the message. The message
// is contained in the Event and, if an error occurs, the consumer was notified.
// Warning: an EventConsumer is only called when an event was sent, in a
// goroutine. DO NOT LOOP INFINITELY AND TAKE CARE OF CONCURRENCY.
type EventConsumer func(event *Event, err error)

// internal definitions of pub/sun timeout.
const (
	publishTimeout = 25 * time.Millisecond
	replyTimeout   = 25 * time.Millisecond
)

// ErrPublishTimeout occurs when the event publishing timeout.
var ErrPublishTimeout = errors.New("publish timeout")
// ErrSubscriberDeleted occurs when the subscriber was deleted (ak. unsubscribe).
var ErrSubscriberDeleted = errors.New("subscriber closed")
// ErrChannelNotFound occurs when the requested channel doesn't exist.
var ErrChannelNotFound = errors.New("channel not found")
// ErrChannelSubscriberNotFound occurs when a channel exists but no subscribers in.
var ErrChannelSubscriberNotFound = errors.New("channel subscriber not found")
// ErrChannelSubscriberAlreadyExists occurs when a subscriber already
// subscribed to the channel.
var ErrChannelSubscriberAlreadyExists = errors.New("channel subscriber already exists")

// Message return the payload (ak. message) of the event.
func (event *Event) Message() interface{} { return event.payload }

// Reply provide a channel to reply directly to the publisher.
func (event *Event) Reply() chan<- interface{} { return event.reply }

// New create a new instance of Bus.
func New() *Bus { return &Bus{subscribers: map[string][]subscriber{}, ids: map[string]interface{}{}, sync: sync.Mutex{}} }

// Publish publish a message to a channel. A reply handler can be provided in
// order to receive reply and/or catch errors.
func (bus *Bus) Publish(channel string, data interface{}, handler ReplyHandler) {
	bus.sync.Lock()
	defer bus.sync.Unlock()

	if _, exists := bus.subscribers[channel]; !exists {
		if handler != nil {
			handler.consume(nil, nil, ErrChannelNotFound, 0)
		}
		return
	}

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

// Subscribe links an handler to a channel. See EventConsumer for more
// information about the handler.
func (bus *Bus) Subscribe(channel string, handler EventConsumer) error {
	id := id(channel, handler)

	bus.sync.Lock()
	defer bus.sync.Unlock()

	if _, exists := bus.ids[id]; exists {
		return ErrChannelSubscriberAlreadyExists
	}

	ch := make(chan *Event)
	ctx, cnl := context.WithCancel(context.Background())

	go func(channel <-chan *Event, ctx context.Context) {
		defer handler(nil, ErrSubscriberDeleted)

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


	bus.subscribers[channel] = append(bus.subscribers[channel], subscriber{id: id, ch: ch, cancel: cnl})
	bus.ids[id] = nil
	return nil
}

// Unsubscribe removes a subscriber linked with a channel. If all subscribers
// linked with a channel are removed, the channel will be removed (and can
// create ErrChannelNotFound errors during publishing on the channel).
func (bus *Bus) Unsubscribe(channel string, handler EventConsumer) error {
	bus.sync.Lock()
	defer bus.sync.Unlock()

	if sub, exists := bus.subscribers[channel]; !exists {
		return ErrChannelNotFound
	} else {
		id := id(channel, handler)

		for i, sbcr := range sub {
			if sbcr.id == id {
				sbcr.cancel()
				if len(sub) == 1 {
					delete(bus.subscribers, channel)
				} else {
					bus.subscribers[channel] = append(sub[:i], sub[i+1:]...)
				}
				delete(bus.ids, id)

				return nil
			}
		}

		return ErrChannelSubscriberNotFound
	}
}

// id is an internal method to create unique id from the channel and the handler.
func id(c string, h interface{}) string { return c + ":" + runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name() }
