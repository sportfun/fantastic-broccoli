package main

import (
	"log"
	"reflect"
	"strings"
	"time"
)

// The scheduler is an instance used to manage other worker. If a worker die, it will be
// allowed to recreate an instance of the dead worker.
// IT IS A NAIVE IMPLEMENTATION OF WHAT IS A SCHEDULER (and can't not really be defined by
// the term 'scheduler', but I had no idea how I choose this name)
type scheduler struct {
	linksCache               linkMap
	workers                  map[string]*worker
	shutdownSchedulerChannel chan interface{}
}

const invalidCharsInWorkerName = `:/\|"'?!@#$%^&*()+=` // List of invalid chars in the worker name
var (
	TTL        = 250 * time.Millisecond // (Time To Live) Time without ping before a worker was declared dead
	TTW        = 25 * time.Microsecond  // (Time To Wait) Time waiting between two ping
	TTR        = 250 * time.Millisecond // (Time To Refresh) Between to scheduler refresh
	GScheduler = &scheduler{
		linksCache:               linkMap{},
		workers:                  map[string]*worker{},
		shutdownSchedulerChannel: make(chan interface{}),
	}                                   // Singleton instance of the scheduler
)
// Register a new worker with its factory. It will be called and
// instanced in a goroutine only when the scheduler was ready.
func (scheduler *scheduler) RegisterWorker(name string, factory workerFactory) {
	name = strings.TrimSpace(name)

	workerValidity(name)
	if factory == nil {
		panic("worker factory can't be nil")
	}

	if _, exists := scheduler.workers[name]; exists {
		panic("worker '" + name + "' already registered")
	}

	scheduler.workers[name] = &worker{
		name:    name,
		factory: factory,
		links:   linkMap{},
	}
	//TODO: LOG :: INFO - Worker X registered
	log.Printf("{scheduler}[INFO]			Worker '%s' registered", name)
}

// Register a new communication link between two workers. This link is wrapped into a
// reflect.Value to force the targeted workers to cast it into the right type.
// (If not, panic is our best friend to avoid invalid communication)
func (scheduler *scheduler) RegisterLink(origin, destination, name string, link reflect.Value) {
	origin = strings.TrimSpace(origin)
	destination = strings.TrimSpace(destination)
	name = strings.TrimSpace(name)

	workerValidity(origin)
	workerValidity(destination)
	if origin == destination {
		panic("worker can't be linked with himself")
	}

	if name == "" {
		panic("worker link name can't be empty")
	}

	if link.Type().Kind() != reflect.Chan {
		panic("worker link must be a channel")
	}

	var keycode string

	if origin > destination {
		keycode = origin + ":" + destination
	} else {
		keycode = destination + ":" + origin
	}

	if _links, exists := scheduler.linksCache[keycode]; exists {
		if _link, exists := _links[name]; exists {
			if _link.Type() != link.Type() {
				panic("worker link '" + name + "' already exists between '" + origin + "' and '" + destination + "', but type is different")
			}
			return // Link already exists, our work is done
		}
	} else {
		scheduler.linksCache[keycode] = map[string]*reflect.Value{}
	}

	scheduler.linksCache[keycode][name] = &link
	//TODO: LOG :: INFO - Link between X and Y registered
	log.Printf("{scheduler}[INFO]			Link '%s' between '%s' and '%s' registered", name, origin, destination)
}

// Prepare and start the worker scheduler.
func (scheduler *scheduler) Run() {
	scheduler.mapLinks()
	scheduler.spawnAllWorkers()
	defer scheduler.killAllWorkers()

	//TODO: LOG (defer) :: INFO - Scheduler stopped
	log.Printf("{scheduler}[INFO]			Scheduler stopped")
	//TODO: LOG :: INFO - Start scheduler
	log.Printf("{scheduler}[INFO]			Start scheduler")
	for {
		select {
		case <-time.Tick(TTR): // Every X time, the scheduler checks if worker are alive
			for _, worker := range scheduler.workers {
				if !worker.IsAlive() {
					//TODO: LOG :: WARNING - Dead worker (X)
					log.Printf("{scheduler}[WARN]			Dead worker (%s)", worker.name)
					if err := worker.Spawn(); err != nil {
						//TODO: Better error management here
						panic("failed to spawn '" + worker.name + "': " + err.Error())
					}
				}
			}
		case <-scheduler.shutdownSchedulerChannel:
			//TODO: LOG :: INFO - Shutdown scheduler
			log.Printf("{scheduler}[INFO]			Shutdown scheduler")
			return
		}
	}
}

// Generate channels mapping between workers (from link cache)
func (scheduler *scheduler) mapLinks() {
	var exists bool
	var wrkA, wrkB *worker

	for keycode, links := range scheduler.linksCache {
		workerNames := strings.Split(keycode, ":")

		wrkA, exists = scheduler.workers[workerNames[0]]
		if !exists {
			panic("worker '" + workerNames[0] + "' not registered")
		}

		wrkB, exists = scheduler.workers[workerNames[1]]
		if !exists {
			panic("worker '" + workerNames[1] + "' not registered")
		}

		wrkA.links[wrkB.name] = links
		wrkB.links[wrkA.name] = links
	}
	scheduler.linksCache = linkMap{}
	//TODO: LOG :: INFO - Links map generated
	log.Printf("{scheduler}[INFO]			Link maps generated")
}

// Simple way to spawn all worker instances
func (scheduler *scheduler) spawnAllWorkers() {
	for name, worker := range scheduler.workers {
		if err := worker.Spawn(); err != nil {
			panic("failed to spawn '" + name + "': " + err.Error())
		}
	}
	//TODO: LOG :: INFO - All worker spawned
	log.Printf("{scheduler}[INFO]			All worker spawned")
}

func (scheduler *scheduler) killAllWorkers() {
	for _, worker := range scheduler.workers {
		worker.Kill()
	}
	//TODO: LOG :: INFO - All worker killed
	log.Printf("{scheduler}[INFO]			All worker killed")
}

// Check if a worker name is valid
func workerValidity(name string) {
	if name == "" {
		panic("worker name can't be empty")
	}

	if strings.ContainsAny(name, invalidCharsInWorkerName) {
		panic("worker name can't contain '" + invalidCharsInWorkerName + "'")
	}
}
