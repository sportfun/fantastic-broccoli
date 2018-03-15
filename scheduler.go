package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sportfun/main/event"
)

type workerContextKey string
type workerTask func(ctx context.Context, bus *event.Bus) error

type worker struct {
	run workerTask

	lastRetry time.Time
	numRetry  *int32
}

type scheduler struct {
	workers map[string]*worker
	bus     *event.Bus

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
	//TODO: LOG :: INFO - Worker X registered
	log.Printf("{scheduler}[INFO]			Worker '%s' registered", name)
}

// Prepare and start the worker scheduler.
func (scheduler *scheduler) Run() (err error) {
	//TODO: LOG :: INFO - Start scheduler
	log.Printf("{scheduler}[INFO]			Start scheduler")
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
		//TODO: LOG (defer) :: INFO - scheduler stopped
		log.Printf("{scheduler}[INFO]			Stop scheduler")
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

func (scheduler *scheduler) spawnWorker(name string) {
	worker, exists := scheduler.workers[name]
	ctx := context.WithValue(scheduler.ctx, workerContextKey("name"), name)

	if !exists {
		panic("worker '" + name + "' doesn't exists")
	}

	if atomic.LoadInt32(worker.numRetry) > 5 {
		panic("worker '" + name + "' has been restarted too many times")
	}

	go func(name string) {
		defer func(name string) {
			if r := recover(); r != nil {
				//TODO: LOG :: ERROR - Worker '{name}' has panicked: {reason}
				log.Printf("{scheduler}[ERROR]			Worker '%s' has panicked: '%s'", name, r)
				if time.Since(worker.lastRetry) < time.Second {
					atomic.AddInt32(worker.numRetry, 1)
				} else {
					atomic.AddInt32(worker.numRetry, -1)
				}
				worker.lastRetry = time.Now()
				scheduler.deadSig <- name
			}
		}(name)

		if err := worker.run(ctx, scheduler.bus); err != nil {
			//TODO: LOG :: ERROR - Worker '{name}' has failed: {reason}
			log.Printf("{scheduler}[ERROR]			Worker '%s' has failed: '%s'", name, err)
			if time.Since(worker.lastRetry) < time.Second {
				atomic.AddInt32(worker.numRetry, 1)
			} else {
				atomic.AddInt32(worker.numRetry, -1)
			}
			worker.lastRetry = time.Now()
			scheduler.deadSig <- name
			return
		}
		//TODO: LOG :: INFO - Worker '{name}' successfully stopped
		log.Printf("{scheduler}[INFO]			Worker '%s' successfully stopped", name)
	}(name)

	//TODO: LOG :: INFO - Worker '{name}' has been launched
	log.Printf("{scheduler}[ERROR]			Worker '%s' has been launched", name)
}
func (scheduler *scheduler) workerValidity(names ...string) {
	for _, name := range names {
		if !workerNameFilter.MatchString(name) {
			panic("worker name '" + name + "' is invalid")
		}
	}
}
