package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
)

func TestScheduler_RegisterWorker(t *testing.T) {
	RegisterTestingT(t)

	sch := &scheduler{workers: map[string]*worker{}, bus: bus.New(), ctx: context.Background(), deadSig: make(chan string)}

	cases := []struct {
		name           string
		factory        workerTask
		workersMatcher OmegaMatcher
		panicMatcher   OmegaMatcher
	}{
		{"", nil, BeEmpty(), Equal("worker name '' is invalid")},
		{"    ", nil, BeEmpty(), Equal("worker name '' is invalid")},
		{"::", nil, BeEmpty(), Equal("worker name '::' is invalid")},
		{"worker", nil, BeEmpty(), Equal("worker task can't be nil")},
		{"group/worker/1", nil, BeEmpty(), Equal("worker task can't be nil")},

		{"group/worker/1", IFTTT(whenTick(time.Microsecond), doNothing()), HaveLen(1), BeNil()},
		{"group/worker/1", IFTTT(whenTick(time.Microsecond), doNothing()), HaveLen(1), Equal("worker 'group/worker/1' already registered")},
		{"    group/worker/1   ", IFTTT(whenTick(time.Microsecond), doNothing()), HaveLen(1), Equal("worker 'group/worker/1' already registered")},

		{"group/worker/2", IFTTT(whenTick(time.Microsecond), doNothing()), HaveLen(2), BeNil()},
	}

	for _, test := range cases {
		func() {
			defer func() { Expect(recover()).Should(test.panicMatcher) }()
			sch.RegisterWorker(test.name, test.factory)
			Expect(sch.workers).Should(test.workersMatcher)
		}()
	}
}
func TestScheduler_Run_DirectCall(t *testing.T) {
	RegisterTestingT(t)

	type workerDefinition struct {
		name string
		task workerTask
	}
	var unrealisticError = errors.New("unrealistic error")

	cases := []struct {
		workers       []workerDefinition
		nWorker       int32
		returnMatcher OmegaMatcher
	}{
		{
			[]workerDefinition{
				{"t/worker", errorTask(unrealisticError)},
			},
			0,
			Equal(errors.New("worker 't/worker' has been restarted too many times")),
		},
		{
			[]workerDefinition{
				{"t/worker", panicTask(unrealisticError)},
			},
			0,
			Equal(errors.New("worker 't/worker' has been restarted too many times")),
		},
		{
			[]workerDefinition{
				{"t/worker/1", IFTTT(whenTick(time.Millisecond), sendTo(":worker:2"))},
				{"t/worker/2", IFTTT(ifReceiveFrom(":worker:2"), sendTo(":worker:3"))},
				{"t/worker/3", IFTTT(ifReceiveFrom(":worker:3"), sendTo(":worker:4"))},
				{"t/worker/4", IFTTT(ifReceiveFrom(":worker:4"), sleep(10*time.Millisecond, sendTo(":worker:2")))},
			},
			4,
			BeNil(),
		},
	}

	for _, test := range cases {
		ctx, cancel := context.WithCancel(context.Background())
		sch := &scheduler{workers: map[string]*worker{}, bus: bus.New(), ctx: ctx, deadSig: make(chan string), workerRetryMax: 5, workerRetryInterval: 200 * time.Millisecond}

		for _, worker := range test.workers {
			sch.RegisterWorker(worker.name, worker.task)
		}

		go func(nWorker int32) {
			time.Sleep(100 * time.Millisecond)
			Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(nWorker))
			cancel()
		}(test.nWorker)

		Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(int32(0)))
		Expect(sch.runUntilClosed()).Should(test.returnMatcher)
		time.Sleep(100 * time.Millisecond)
		Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(int32(0)))
	}
}
func TestScheduler_Run(t *testing.T) {
	RegisterTestingT(t)

	ctx, cancel := context.WithCancel(context.Background())
	sch := &scheduler{workers: map[string]*worker{}, bus: bus.New(), ctx: ctx, deadSig: make(chan string), workerRetryMax: 5, workerRetryInterval: 200 * time.Millisecond}

	sch.RegisterWorker("t/worker", panicTask(errors.New("unrealistic error")))
	Expect(<-sch.Run()).Should(BeFalse())
	cancel()
}

// Tasks generation
var onlineWorker = new(int32)

func errorTask(err error) workerTask {
	id := rand.Uint32()
	logger := logrus.WithField("id", id)

	return func(ctx context.Context, bus *bus.Bus) error {
		logger.Printf("Start worker")
		defer func() {
			logger.Printf("Stop worker")
		}()
		return err
	}
}
func panicTask(err error) workerTask {
	id := rand.Uint32()
	logger := logrus.WithField("id", id)

	return func(ctx context.Context, bus *bus.Bus) error {
		logger.Printf("Start worker")
		defer func() {
			logger.Printf("Stop worker")
		}()
		panic(err)
	}
}
func IFTTT(ifThis func(context.Context, *bus.Bus) <-chan interface{}, thenThat func(interface{}, context.Context, *bus.Bus)) workerTask {
	return func(ctx context.Context, bus *bus.Bus) error {
		logger := logrus.WithField("id", ctx.Value(workerContextKey("name")))

		logger.Println("Start worker")
		atomic.AddInt32(onlineWorker, 1)
		defer func() {
			logger.Println("Stop worker")
			atomic.AddInt32(onlineWorker, -1)
		}()

		trg := ifThis(ctx, bus)
		for {
			select {
			case <-ctx.Done():
				return nil
			case v, o := <-trg:
				if !o {
					logger.Errorln("Trigger close")
					return nil
				}
				thenThat(v, ctx, bus)
			}
		}
	}
}

func whenTick(d time.Duration) func(context.Context, *bus.Bus) <-chan interface{} {
	return func(context.Context, *bus.Bus) <-chan interface{} {
		c := make(chan interface{})
		go func() { c <- <-time.Tick(d) }()
		return c
	}
}
func ifReceiveFrom(from string) func(context.Context, *bus.Bus) <-chan interface{} {
	return func(_ context.Context, b *bus.Bus) <-chan interface{} {
		c := make(chan interface{})
		b.Subscribe(from, func(event *bus.Event, err error) {
			Expect(err).Should(BeNil())
			c <- event
		})
		return c
	}
}

func doNothing() func(interface{}, context.Context, *bus.Bus) {
	return func(v interface{}, ctx context.Context, _ *bus.Bus) {
		logrus.WithField("id", ctx.Value(workerContextKey("name"))).Println("Do nothing")
	}
}
func sendTo(to string) func(interface{}, context.Context, *bus.Bus) {
	return func(v interface{}, ctx context.Context, b *bus.Bus) {
		logger := logrus.WithField("id", ctx.Value(workerContextKey("name")))

		if e, ok := v.(*bus.Event); ok {
			e.Reply() <- fmt.Sprintf("received by %s", ctx.Value(workerContextKey("name")))
			v = e.Message()
		}
		logger.Printf("Received %#v", v)
		b.Publish(to, v, bus.SyncReplyHandler(func(data interface{}, err error) {
			if err != nil {
				logger.Errorln(err.Error())
			} else {
				logger.Printf("Response: %v", data)
			}
		}))
	}
}
func sleep(d time.Duration, f func(interface{}, context.Context, *bus.Bus)) func(interface{}, context.Context, *bus.Bus) {
	return func(v interface{}, ctx context.Context, b *bus.Bus) {
		time.Sleep(d)
		f(v, ctx, b)
	}
}
