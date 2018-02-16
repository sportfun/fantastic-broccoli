package event

import (
	"errors"
	"time"
)

type ReplyHandler interface {
	consume(<-chan interface{}, <-chan error, error, time.Duration)
}

type syncReply struct{ ReplyConsumer }
type asyncReply struct{ *syncReply }

type ReplyConsumer func(interface{}, error)

var ErrReplyTimeout = errors.New("reply handler has timeout")

func SyncReplyHandler(consumer ReplyConsumer) ReplyHandler  { return &syncReply{consumer} }
func AsyncReplyHandler(consumer ReplyConsumer) ReplyHandler { return &asyncReply{&syncReply{consumer}} }

func (r *syncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	if err != nil {
		r.ReplyConsumer(nil, err)
	}

	select {
	case v, o := <-c:
		if o {
			r.ReplyConsumer(v, nil)
		}
	case e, o := <-e:
		if o {
			r.ReplyConsumer(nil, e)
		}
	case <-time.After(timeout):
		r.ReplyConsumer(nil, ErrReplyTimeout)
	}
}

func (r *asyncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	if err != nil {
		r.ReplyConsumer(nil, err)
	}

	go r.syncReply.consume(c, e, nil, timeout)
}
