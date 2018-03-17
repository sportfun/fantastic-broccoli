package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
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
	log.Infof("Worker '%s' registered", name) //LOG :: INFO - Worker {name} registered
}

// Prepare and start the worker scheduler.
func (scheduler *scheduler) Run() (err error) {
	log.Infof("Start scheduler") //LOG :: INFO - Start scheduler
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
		log.Infof("Stop scheduler") //LOG (defer) :: INFO - scheduler stopped
	}()

	for name := range scheduler.workers {
		scheduler.spawnWorker(name)
	}

	for {
		select {
		case <-scheduler.ctx.Done():
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

	if atomic.LoadInt32(worker.numRetry) > int32(Profile.Scheduler.Worker.Retry) {
		panic("worker '" + name + "' has been restarted too many times")
	}

	go func(name string) {
		defer func(name string) {
			if r := recover(); r != nil {
				log.Fatalf("Worker '%s' has panicked: %s", name, r) //LOG :: ERROR - Worker '{name}' has panicked: {reason}
				if time.Since(worker.lastRetry) < time.Millisecond*time.Duration(Profile.Scheduler.Worker.Interval) {
					atomic.AddInt32(worker.numRetry, 1)
				} else {
					atomic.AddInt32(worker.numRetry, -1)
				}
				worker.lastRetry = time.Now()
				scheduler.deadSig <- name
			}
		}(name)

		if err := worker.run(ctx, scheduler.bus); err != nil {
			log.Errorf("Worker '%s' has failed: %s", name, err) //LOG :: ERROR - Worker '{name}' has failed: {reason}
			if time.Since(worker.lastRetry) < 2*time.Second {
				atomic.AddInt32(worker.numRetry, 1)
			} else {
				atomic.AddInt32(worker.numRetry, -1)
			}
			worker.lastRetry = time.Now()
			scheduler.deadSig <- name
			return
		}
		log.Infof("Worker '%s' successfully stopped", name) //LOG :: INFO - Worker '{name}' successfully stopped
	}(name)

	log.Infof("Worker '%s' has been launched", name) //LOG :: INFO - Worker '{name}' has been launched
}

// workerValidity check if a worker name is valid
func (scheduler *scheduler) workerValidity(names ...string) {
	for _, name := range names {
		if !workerNameFilter.MatchString(name) {
			panic("worker name '" + name + "' is invalid")
		}
	}
}
