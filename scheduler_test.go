package gakisitor

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestScheduler_RegisterWorker(t *testing.T) {
	RegisterTestingT(t)

	scheduler := &scheduler{linkCache: map[string]links{}, workers: map[string]*worker{}, ctx: context.Background(), deadSig: make(chan string)}

	for _, tcase := range []struct {
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

		{"group/worker/1", taskBuilder(triggerTicker(time.Microsecond), doNothing()), HaveLen(1), BeNil()},
		{"group/worker/1", taskBuilder(triggerTicker(time.Microsecond), doNothing()), HaveLen(1), Equal("worker 'group/worker/1' already registered")},
		{"    group/worker/1   ", taskBuilder(triggerTicker(time.Microsecond), doNothing()), HaveLen(1), Equal("worker 'group/worker/1' already registered")},

		{"group/worker/2", taskBuilder(triggerTicker(time.Microsecond), doNothing()), HaveLen(2), BeNil()},
	} {
		func() {
			defer func() { Expect(recover()).Should(tcase.panicMatcher) }()
			scheduler.RegisterWorker(tcase.name, tcase.factory)
			Expect(scheduler.workers).Should(tcase.workersMatcher)
		}()
	}
}
func TestScheduler_Run(t *testing.T) {
	RegisterTestingT(t)

	type workerDefinition struct {
		name string
		task workerTask
	}
	type linkDefinition struct {
		origin      string
		destination string
		name        string
	}
	var unrealisticError = errors.New("unrealistic error")

	for _, tcase := range []struct {
		workers       []workerDefinition
		links         []linkDefinition
		nWorker       int32
		returnMatcher OmegaMatcher
	}{
		{
			[]workerDefinition{
				{"worker", errorTask(unrealisticError)},
			},
			[]linkDefinition{},
			0,
			Equal(errors.New("worker 'worker' has been restarted too many times")),
		},
		{
			[]workerDefinition{
				{"worker", panicTask(unrealisticError)},
			},
			[]linkDefinition{},
			0,
			Equal(errors.New("worker 'worker' has been restarted too many times")),
		},

		{
			[]workerDefinition{
				{"test/worker/1", taskBuilder(triggerLinks("test/worker/2", "<"), doPingLinks("test/worker/2", ">"))},
				{"test/worker/2", taskBuilder(triggerLinks("test/worker/1", ">"), doPingLinks("test/worker/1", "<"))},
				{"test/worker/3", taskBuilder(triggerTicker(time.Second), doNothing())},
			},
			[]linkDefinition{
				{"test/worker/1", "test/worker/2", ">"},
				{"test/worker/1", "test/worker/2", "<"},
			},
			3,
			BeNil(),
		},
	} {
		ctx, cancel := context.WithCancel(context.Background())
		scheduler := &scheduler{workers: map[string]*worker{}, ctx: ctx, deadSig: make(chan string)}

		for _, worker := range tcase.workers {
			scheduler.RegisterWorker(worker.name, worker.task)
		}

		func() {
			go func() {
				time.Sleep(5 * time.Millisecond)
				Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(tcase.nWorker))
				cancel()
			}()

			Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(int32(0)))
			Expect(scheduler.Run()).Should(tcase.returnMatcher)
			time.Sleep(time.Millisecond)
			Expect(atomic.LoadInt32(onlineWorker)).Should(Equal(int32(0)))
		}()
	}
}

// Tasks generation
var onlineWorker = new(int32)

func taskBuilder(trigger func() <-chan interface{}, do func(interface{}, map[string]links, uint32)) workerTask {
	id := rand.Uint32()

	return func(ctx context.Context, links map[string]links) error {
		//TODO: LOG :: Use *testing.T when the logger will be implemented
		log.Printf("{worker#%d}			Start worker", id)
		atomic.AddInt32(onlineWorker, 1)
		defer func() {
			//TODO: LOG :: Use *testing.T when the logger will be implemented
			log.Printf("{worker#%d}			Stop worker", id)
			atomic.AddInt32(onlineWorker, -1)
		}()

		for {
			select {
			case <-ctx.Done():
				return nil

			case v, o := <-trigger(links):
				if !o {
					log.Printf("{worker#%d}			trigger closed", id)
					return nil
				}
				do(v, links, id)
			}
		}
		return nil
	}
}
func errorTask(err error) workerTask {
	id := rand.Uint32()

	return func(ctx context.Context, links map[string]links) error {
		//TODO: LOG :: Use *testing.T when the logger will be implemented
		log.Printf("{worker#%d}			Start worker", id)
		defer func() {
			//TODO: LOG :: Use *testing.T when the logger will be implemented
			log.Printf("{worker#%d}			Stop worker", id)
		}()

		return err
	}
}
func panicTask(err error) workerTask {
	id := rand.Uint32()

	return func(ctx context.Context, links map[string]links) error {
		//TODO: LOG :: Use *testing.T when the logger will be implemented
		log.Printf("{worker#%d}			Start worker", id)
		defer func() {
			//TODO: LOG :: Use *testing.T when the logger will be implemented
			log.Printf("{worker#%d}			Stop worker", id)
		}()

		panic(err)
	}
}

func triggerTicker(d time.Duration) func(map[string]links) <-chan interface{} {
	return func(links map[string]links) <-chan interface{} {
		trigger := make(chan interface{})
		go func() { time.Sleep(d); trigger <- time.Now() }()
		return trigger
	}
}
func triggerLinks(targetWorker, in string) func(map[string]links) <-chan interface{} {
	return func(links map[string]links) <-chan interface{} {
		return links[targetWorker][in]
	}
}

func doNothing() func(interface{}, map[string]links, uint32) {
	return func(v interface{}, links map[string]links, id uint32) {
		log.Printf("{worker#%d}			do nothing", id)
	}
}
func doPingLinks(targetWorker, out string) func(interface{}, map[string]links, uint32) {
	return func(v interface{}, links map[string]links, id uint32) {
		log.Printf("{worker#%d}			message from '%s' handled", id, targetWorker)
		links[targetWorker][out] <- v
		log.Printf("{worker#%d}			message sent to '%s'", id, targetWorker)
	}
}
