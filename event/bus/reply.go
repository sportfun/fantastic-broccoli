package bus

import (
	"errors"
	"time"
)

// ReplyHandler provides an interface for handling replies during publishing
// (for Bus.Publish). It can't be overloaded, but two way are provides to
// generate ReplyHandler.
type ReplyHandler interface {
	consume(<-chan interface{}, <-chan error, error, time.Duration)
}

// internal type for synchronous reply handler.
type syncReply struct{ ReplyConsumer }
// internal type for asynchronous reply handler.
type asyncReply struct{ ReplyConsumer }

// ReplyConsumer is the handler in charge of receiving the reply. It call only
// once during the publishing process (after receiving reply or if an error
// occurs). DO NOT LOOP INFINITELY.
type ReplyConsumer func(data interface{}, err error)

// ErrReplyTimeout occurs when the subscriber (ak. the Bus.Subscribe handler)
// not sent reply or if it takes to mush time before sent it.
var ErrReplyTimeout = errors.New("reply has timeout")

// SyncReplyHandler generate a synchronous ReplyHandler. It means that if you
// use this ReplyHandler with Bus.Publish, your Bus.Publish call will never
// finished before the handler ending.
func SyncReplyHandler(consumer ReplyConsumer) ReplyHandler  { return &syncReply{consumer} }
// AsyncReplyHandler generate an asynchronous ReplyHandler. It means that if
// you use this ReplyHandler with Bus.Publish, your Bus.Publish call can be
// finished before the handler ending. TAKE CARE OF CONCURRENCY.
func AsyncReplyHandler(consumer ReplyConsumer) ReplyHandler { return &asyncReply{consumer} }

// internal implementation of the sync ReplyHandler.
func (r *syncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	consume(r.ReplyConsumer, c, e, err, timeout)
}

// internal implementation of the async ReplyHandler.
func (r *asyncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	if err != nil {
		r.ReplyConsumer(nil, err)
		return
	}

	go consume(r.ReplyConsumer, c, e, nil, timeout)
}

// internal generic implementation of the ReplyHandler.consume.
func consume(r ReplyConsumer, c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	if err != nil {
		r(nil, err)
		return
	}

	select {
	case v, o := <-c:
		if o {
			r(v, nil)
		}
	case e, o := <-e:
		if o {
			r(nil, e)
		}
	case <-time.After(timeout):
		r(nil, ErrReplyTimeout)
	}
}
