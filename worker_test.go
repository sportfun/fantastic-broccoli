package gakisitor

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var (
	errUnrealistic = errors.New("unrealistic error... but it just a test")
)

var (
	naiveWorkerFactoryBuilder = func(t *testing.T) workerFactory {
		id := rand.Uint32()

		return func(links linkMap, flow workerFlow) error {
			//TODO: LOG :: Use *testing.T when the logger will be implemented
			log.Printf("> [%d]				start simple worker", id)

			go func() {
				defer func() { log.Printf("> [%d]				stop simple worker (shutdown)", id) }()
				for {
					select {
					case <-flow.shutdown:
						return
					case <-flow.pulseIn:
						log.Printf("> [%d]				heartbeat handled", id)
						flow.pulseOut <- nil
						log.Printf("> [%d]				heartbeat sent", id)
					}
				}
			}()
			return nil
		}
	}
	completeWorkerFactoryBuilder = func(t *testing.T, targetWorker, inName, outName string) workerFactory {
		id := rand.Uint32()

		return func(links linkMap, flow workerFlow) error {
			//TODO: LOG :: Use *testing.T when the logger will be implemented
			log.Printf("> [%d]				start worker", id)

			go func() {
				defer func() { log.Printf("> [%d]				stop worker (shutdown)", id) }()
				var in <-chan string = links[targetWorker][inName].Interface().(chan string)
				var out chan<- string = links[targetWorker][outName].Interface().(chan string)

				for {
					select {
					case <-flow.shutdown:
						return
					case <-flow.pulseIn:
						log.Printf("> [%d]				heartbeat handled", id)
						flow.pulseOut <- nil
						log.Printf("> [%d]				heartbeat sent", id)
					case v := <-in:
						log.Printf("> [%d]				message from '%s' handled", id, targetWorker)
						out <- v
						log.Printf("> [%d]				message for '%s' sent", id, targetWorker)
					}
				}
			}()
			return nil
		}
	}
)

func TestWorker_SpawnKill(t *testing.T) {
	RegisterTestingT(t)
	TTL = 250 * time.Microsecond
	numGoroutine := runtime.NumGoroutine()

	for id, tcase := range []struct {
		spawnMustFail bool
		startAlive    bool
		factory       workerFactory
		err           error
	}{
		{true, true, naiveWorkerFactoryBuilder(t), ErrWorkerAlreadySpawned},
		{true, false, func(_ linkMap, _ workerFlow) error { return errUnrealistic }, errUnrealistic},
		{false, false, naiveWorkerFactoryBuilder(t), nil},
	} {
		worker := &worker{
			factory: tcase.factory,
			alive:   tcase.startAlive,
			sync:    sync.Mutex{},

			shutdownWorkerChannel: make(chan interface{}),
			pulseInWorkerChannel:  make(chan interface{}),
			pulseOutWorkerChannel: make(chan interface{}),
		}

		if tcase.spawnMustFail {
			Expect(worker.Spawn()).Should(MatchError(tcase.err))
			Expect(worker.IsAlive()).Should(Equal(tcase.startAlive))
			worker.Kill() // to be sure
			Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond).Should(
				Equal(numGoroutine),
				fmt.Sprintf("goroutine not ended (testcase#%d)", id+1),
			)
			continue
		}

		Expect(worker.Spawn()).Should(Succeed())
		Expect(worker.IsAlive()).Should(BeTrue())

		worker.Kill()
		Expect(worker.IsAlive()).Should(BeFalse())
		Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond).Should(
			Equal(numGoroutine),
			fmt.Sprintf("goroutine not ended (testcase#%d)", id+1),
		)
	}
}

func TestWorker_KillAccuracy(t *testing.T) {
	RegisterTestingT(t)
	TTL = 25 * time.Microsecond
	numGoroutine := runtime.NumGoroutine()

	worker := &worker{
		factory: naiveWorkerFactoryBuilder(t),
		sync:    sync.Mutex{},

		shutdownWorkerChannel: make(chan interface{}),
		pulseInWorkerChannel:  make(chan interface{}),
		pulseOutWorkerChannel: make(chan interface{}),
	}

	worker.Spawn()
	Eventually(func() int { return runtime.NumGoroutine() }, time.Millisecond, TTL).Should(Equal(numGoroutine + 2))

	worker.Kill()
	Eventually(func() int { return runtime.NumGoroutine() }, 250*time.Millisecond, TTL).Should(Equal(numGoroutine))
}

func TestWorker_PingerAccurancy(t *testing.T) {
	RegisterTestingT(t)
	TTL = 25 * time.Microsecond
	numGoroutine := runtime.NumGoroutine()

	locker := make(chan string)
	unlocker := make(chan string)
	lockerValue := reflect.ValueOf(locker)
	unlockerValue := reflect.ValueOf(unlocker)

	worker := &worker{
		factory: completeWorkerFactoryBuilder(t, "...", "<", ">"),
		sync:    sync.Mutex{},
		links:   linkMap{"...": link{"<": &lockerValue, ">": &unlockerValue}},

		shutdownWorkerChannel: make(chan interface{}),
		pulseInWorkerChannel:  make(chan interface{}),
		pulseOutWorkerChannel: make(chan interface{}),
	}

	worker.Spawn()
	time.Sleep(2 * TTL)
	Expect(worker.IsAlive()).Should(Equal(true))

	locker <- ""
	Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond, TTL).Should(Equal(numGoroutine + 1))

	Expect(unlocker).Should(Receive()) // unlock worker
	Eventually(func() int { return runtime.NumGoroutine() }, 250*time.Millisecond, TTL).Should(Equal(numGoroutine))
}
