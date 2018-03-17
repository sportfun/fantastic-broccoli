package bus

import (
	"errors"
	"time"
)

type ReplyHandler interface {
	consume(<-chan interface{}, <-chan error, error, time.Duration)
}

type syncReply struct{ ReplyConsumer }
type asyncReply struct{ ReplyConsumer }

type ReplyConsumer func(interface{}, error)

var ErrReplyTimeout = errors.New("reply has timeout")

func SyncReplyHandler(consumer ReplyConsumer) ReplyHandler  { return &syncReply{consumer} }
func AsyncReplyHandler(consumer ReplyConsumer) ReplyHandler { return &asyncReply{consumer} }

func (r *syncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	consume(r.ReplyConsumer, c, e, err, timeout)
}

func (r *asyncReply) consume(c <-chan interface{}, e <-chan error, err error, timeout time.Duration) {
	if err != nil {
		r.ReplyConsumer(nil, err)
	}

	go consume(r.ReplyConsumer, c, e, nil, timeout)
}

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
