package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
)

type workerContextKey string
type workerTask func(ctx context.Context, bus *bus.Bus) error

type worker struct {
	run workerTask

	lastRetry time.Time
	numRetry  *int32
}

type scheduler struct {
	workers map[string]*worker
	bus     *bus.Bus

	ctx     context.Context
	deadSig chan string

	workerRetryMax      int32
	workerRetryInterval time.Duration
	workerOnline        sync.WaitGroup
}

// List of valid chars in the worker name
var workerNameFilter = regexp.MustCompile(`^[a-zA-Z0-9-/.]+$`)

// Register a new worker with its task. It will be called and
// instanced in a goroutine only when the scheduler was ready.
func (scheduler *scheduler) RegisterWorker(name string, task workerTask) {
	name = strings.TrimSpace(name)

	scheduler.workerValidity(name)
	if _, exists := scheduler.workers[name]; exists {
		panic("worker '" + name + "' already registered")
	}
	if task == nil {
		panic("worker task can't be nil")
	}

	scheduler.workers[name] = &worker{
		run:      task,
		numRetry: new(int32),
	}
	logrus.Debugf("Worker '%s' registered", name) // LOG :: Debug - Worker {name} registered
}

// Start the worker scheduler.
func (scheduler *scheduler) Run() (<-chan bool) {
	var restart = make(chan bool)

	go func(stopped chan<- bool) {
		err := scheduler.runUntilClosed()
		if err != nil {
			logrus.Error(err)
			restart <- false
		} else {
			restart <- true
		}
	}(restart)
	return restart
}

func (scheduler *scheduler) runUntilClosed() (err error) {
	logrus.Infof("Start scheduler") // LOG :: INFO - Start scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
		logrus.Infof("Stop scheduler") // LOG (defer) :: INFO - scheduler stopped
	}()

	scheduler.workerOnline = sync.WaitGroup{}
	for name := range scheduler.workers {
		scheduler.spawnWorker(name)
	}

	for {
		select {
		case <-scheduler.ctx.Done():
			logrus.Debug("Closed by context, wait all workers")
			scheduler.workerOnline.Wait()
			logrus.Debug("All workers stopped, stop scheduler")
			return

		case name, open := <-scheduler.deadSig:
			if !open {
				panic("unexpected closed channel (scheduler.deadSig)")
			}
			scheduler.spawnWorker(name)
		}
	}
}

// spawnWorker launch worker with the own manager
func (scheduler *scheduler) spawnWorker(name string) {
	worker, exists := scheduler.workers[name]
	ctx := context.WithValue(scheduler.ctx, workerContextKey("name"), name)

	if !exists {
		panic("worker '" + name + "' doesn't exists")
	}

	if atomic.LoadInt32(worker.numRetry) > scheduler.workerRetryMax {
		panic("worker '" + name + "' has been restarted too many times")
	}

	scheduler.workerOnline.Add(1)
	go func(name string) {
		defer func(name string) {
			if r := recover(); r != nil {
				logrus.WithField("stacktrace", string(debug.Stack())).Errorf("Worker '%s' has failed: %s", name, r) // LOG :: ERROR - Worker '{name}' has failed: {reason}
				if time.Since(worker.lastRetry) < scheduler.workerRetryInterval {
					atomic.AddInt32(worker.numRetry, 1)
				} else {
					atomic.AddInt32(worker.numRetry, -1)
				}
				worker.lastRetry = time.Now()
				scheduler.deadSig <- name
			}
		}(name)
		defer func() { scheduler.workerOnline.Done() }()

		if err := worker.run(ctx, scheduler.bus); err != nil {
			panic(err)
		}
		logrus.Infof("Worker '%s' successfully stopped", name) // LOG :: INFO - Worker '{name}' successfully stopped
	}(name)

	logrus.Infof("Worker '%s' has been launched", name) // LOG :: INFO - Worker '{name}' has been launched
}

// workerValidity check if a worker name is valid
func (scheduler *scheduler) workerValidity(names ...string) {
	for _, name := range names {
		if !workerNameFilter.MatchString(name) {
			panic("worker name '" + name + "' is invalid")
		}
	}
}
