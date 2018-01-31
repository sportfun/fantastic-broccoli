package engine

import (
	"reflect"
	"sync"
	"errors"
	"time"
	"log"
)

type (
	signal byte
	link map[string]*reflect.Value
	linkMap map[string]link
	workerFactory func(links linkMap, flow WorkerFlow) error
)

// a worker is a container representing a goroutine. This container is used by the scheduler to
// 'manage' theses goroutine and try to prevent some failure (like panic or anything else).
type worker struct {
	name    string        // worker name
	factory workerFactory // worker factory
	links   linkMap       // list of links registered for this worker
	alive   bool          // state of the worker
	sync    sync.Mutex    // sync mutex

	shutdownWorkerChannel chan interface{} // chan for worker shutdown
	pulseInWorkerChannel  chan interface{} // chan for heartbeat system
	pulseOutWorkerChannel chan interface{} // chan for heartbeat system
}

// a flow structure contains all channels used by the worker to communicate with the scheduler.
type WorkerFlow struct {
	shutdown <-chan interface{} // chan for worker shutdown
	pulseIn  <-chan interface{} // chan for heartbeat system
	pulseOut chan<- interface{} // chan for heartbeat system
}

var ErrWorkerAlreadySpawned = errors.New("worker already spawned")

// spawn a new instance (goroutine) of the worker.
func (worker *worker) Spawn() error {
	if worker.IsAlive() {
		return ErrWorkerAlreadySpawned
	}

	worker.shutdownWorkerChannel = make(chan interface{})
	worker.pulseInWorkerChannel = make(chan interface{})
	worker.pulseOutWorkerChannel = make(chan interface{})

	workerFlow := WorkerFlow{
		shutdown: worker.shutdownWorkerChannel,
		pulseIn:  worker.pulseInWorkerChannel,
		pulseOut: worker.pulseOutWorkerChannel,
	}

	if err := worker.factory(worker.links, workerFlow); err != nil {
		return err
	}
	worker.alive = true

	// goroutine used to ping the worker instance (goroutine) in order
	// to define if it was stopped or stuck.
	go func() {
		var pulseOut chan<- interface{} = worker.pulseInWorkerChannel
		var pulseIn <-chan interface{} = worker.pulseOutWorkerChannel

		//TODO: LOG (defer) :: DEBUG - Heartbeat controller for X stopped
		defer log.Printf("{%s[heartbeat]}[DEBUG]	Heartbeat controller stopped", worker.name)
		//TODO: LOG :: DEBUG - Heartbeat controller for X started
		log.Printf("{%s[heartbeat]}[DEBUG]	Heartbeat controller started", worker.name)
		for worker.IsAlive() {
			select {
			case pulseOut <- nil:
			case <-time.After(TTL): // the worker is stuck or crashed
				//TODO: LOG :: DEBUG - No signal from the worker (stuck or crashed)
				log.Printf("{%s[heartbeat]}[DEBUG]	No signal from the worker (stuck or crashed)", worker.name)
				worker.Kill()
				return
			}

			select {
			case <-pulseIn:
			case <-time.After(time.Microsecond): // the worker has crashed
				//TODO: LOG :: DEBUG - No response from the worker (crashed/not implemented)
				log.Printf("{%s[heartbeat]}[DEBUG]	No response from the worker (crashed/not implemented)", worker.name)
				worker.Kill()
				return
			}

			time.Sleep(TTW) // avoid CPU overloading if the worker do nothing
		}
	}()

	//TODO: LOG :: INFO - Worker X spawned
	log.Printf("{%s[worker]}[INFO]		Worker spawned", worker.name)
	return nil
}

// kill the worker. (it send a signal to stop it, but don't wait
// if it will be handled)
func (worker *worker) Kill() {
	defer func() { recover() }() // panic if shutdownWorkerChannel is closed

	if !worker.IsAlive() {
		return
	}

	worker.sync.Lock()
	worker.alive = false
	worker.sync.Unlock()

	select {case worker.shutdownWorkerChannel <- true:default:} // Can't wait here (main 'thread')

	close(worker.shutdownWorkerChannel)

	//TODO: LOG :: INFO - Worker X killed
	log.Printf("{%s[worker]}[INFO]		Worker killed", worker.name)
}

// return if the worker is alive.
func (worker *worker) IsAlive() bool {
	worker.sync.Lock()
	defer worker.sync.Unlock()
	return worker.alive
}
