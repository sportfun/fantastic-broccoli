package bus

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestEvent(t *testing.T) {
	RegisterTestingT(t)

	rep := make(chan interface{})
	evnt := Event{payload: "payload", reply: rep}
	Expect(evnt.Message()).Should(Equal("payload"))
	go func() { evnt.Reply() <- nil }()
	Expect(<-rep).Should(BeNil())
}

func TestBus_Subscribe(t *testing.T) {
	RegisterTestingT(t)

	var x int32
	handlerA := func(*Event, error) { atomic.AddInt32(&x, 1) }
	handlerB := func(*Event, error) { handlerA(nil, nil) }

	bus := New()
	bus.Subscribe(":subscribe", handlerA)
	bus.Subscribe(":subscribe", handlerB)

	bus.Publish(":subscribe", nil, nil)
	Eventually(func() int32 { return atomic.LoadInt32(&x) }, 150*time.Millisecond).Should(Equal(int32(2)))
}

func TestBus_Subscribe_error(t *testing.T) {
	RegisterTestingT(t)

	handler := func(*Event, error) {}

	bus := New()
	bus.Subscribe(":subscribe", handler)

	Expect(bus.Subscribe(":subscribe", handler)).Should(MatchError(ErrChannelSubscriberAlreadyExists))
}

func TestBus_Unsubscribe(t *testing.T) {
	RegisterTestingT(t)

	var x int32
	handlerA := func(_ *Event, err error) {
		if err == ErrSubscriberDeleted {
			atomic.AddInt32(&x, -1)
		} else {
			atomic.AddInt32(&x, 1)
		}
	}
	handlerB := func(*Event, error) { t.Log("ok"); atomic.AddInt32(&x, 1) }

	bus := New()
	bus.Subscribe(":unsubscribe", handlerA)
	bus.Subscribe(":unsubscribe", handlerB)
	bus.Unsubscribe(":unsubscribe", handlerA)
	bus.Publish(":unsubscribe", nil, nil)

	Eventually(func() int32 { return atomic.LoadInt32(&x) }, 150*time.Millisecond).Should(Equal(int32(0)))
}

func TestBus_Unsubscribe_closeChannel(t *testing.T) {
	RegisterTestingT(t)

	handler := func(*Event, error) {}

	bus := New()
	bus.Subscribe(":unsubscribe", handler)
	bus.Unsubscribe(":unsubscribe", handler)

	bus.Publish(":unsubscribe", nil, SyncReplyHandler(func(_ interface{}, err error) {
		Expect(err).Should(MatchError(ErrChannelNotFound))
	}))
}

func TestBus_Unsubscribe_error(t *testing.T) {
	RegisterTestingT(t)

	handlerA := func(*Event, error) {}
	handlerB := func(*Event, error) {}

	bus := New()
	bus.Subscribe(":unsubscribe", handlerA)

	Expect(bus.Unsubscribe("::", nil)).Should(MatchError(ErrChannelNotFound))
	Expect(bus.Unsubscribe(":unsubscribe", nil)).Should(MatchError(ErrChannelSubscriberNotFound))
	Expect(bus.Unsubscribe(":unsubscribe", handlerB)).Should(MatchError(ErrChannelSubscriberNotFound))
}

func TestBus_SyncPublish(t *testing.T) {
	RegisterTestingT(t)

	var x int32
	handler := func(event *Event, err error) { event.Reply() <- int32(1) }

	bus := New()
	bus.Subscribe(":sync", handler)
	bus.Publish(":sync", nil, SyncReplyHandler(func(data interface{}, err error) { atomic.AddInt32(&x, data.(int32)) }))

	Expect(atomic.LoadInt32(&x)).Should(Equal(int32(1)))
}

func TestBus_AsyncPublish(t *testing.T) {
	RegisterTestingT(t)

	var x int32
	handler := func(event *Event, err error) { event.Reply() <- int32(1) }

	bus := New()
	bus.Subscribe(":async", handler)
	bus.Publish(":async", nil, AsyncReplyHandler(func(data interface{}, err error) { time.Sleep(time.Second); atomic.AddInt32(&x, data.(int32)) }))

	Expect(atomic.LoadInt32(&x)).Should(Equal(int32(0)))
}

func TestBus_Publish_error(t *testing.T) {
	RegisterTestingT(t)

	bus := New()
	bus.Publish(":error", nil, SyncReplyHandler(func(data interface{}, err error) {
		Expect(err).Should(MatchError(ErrChannelNotFound))
	}))
}
