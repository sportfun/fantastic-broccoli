package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestScheduler_RegisterWorker(t *testing.T) {
	RegisterTestingT(t)
	scheduler := &scheduler{linksCache: linkMap{}, workers: map[string]*worker{}}

	for _, tcase := range []struct {
		name           string
		factory        workerFactory
		workersMatcher types.GomegaMatcher
		panicMatcher   types.GomegaMatcher
	}{
		{"", nil, BeEmpty(), Equal("worker name can't be empty")},
		{"    ", nil, BeEmpty(), Equal("worker name can't be empty")},
		{"::", nil, BeEmpty(), Equal(`worker name can't contain ':/\|"'?!@#$%^&*()+='`)},
		{"worker1", nil, BeEmpty(), Equal("worker factory can't be nil")},
		{"worker1", naiveWorkerFactoryBuilder(t), HaveLen(1), BeNil()},
		{"worker1", naiveWorkerFactoryBuilder(t), HaveLen(1), Equal("worker 'worker1' already registered")},
		{"    worker1   ", naiveWorkerFactoryBuilder(t), HaveLen(1), Equal("worker 'worker1' already registered")},
		{"worker2", naiveWorkerFactoryBuilder(t), HaveLen(2), BeNil()},
	} {
		func() {
			defer func() { Expect(recover()).Should(tcase.panicMatcher) }()
			scheduler.RegisterWorker(tcase.name, tcase.factory)
			Expect(scheduler.workers).Should(tcase.workersMatcher)
		}()
	}
}

func TestScheduler_RegisterLink(t *testing.T) {
	RegisterTestingT(t)
	scheduler := &scheduler{linksCache: linkMap{}, workers: map[string]*worker{}}
	nilValue := reflect.ValueOf(nil)

	for _, tcase := range []struct {
		name         string
		origin       string
		destination  string
		link         reflect.Value
		linksMatcher types.GomegaMatcher
		panicMatcher types.GomegaMatcher
	}{
		{"", "", "", nilValue, BeEmpty(), Equal("worker name can't be empty")},
		{"", "::", "", nilValue, BeEmpty(), Equal(`worker name can't contain ':/\|"'?!@#$%^&*()+='`)},
		{"", "worker1", "::", nilValue, BeEmpty(), Equal(`worker name can't contain ':/\|"'?!@#$%^&*()+='`)},
		{"", "worker1", "worker1", nilValue, BeEmpty(), Equal("worker can't be linked with himself")},
		{"", "worker1", "worker2", nilValue, BeEmpty(), Equal("worker link name can't be empty")},
		{"link", "worker1", "worker2", reflect.ValueOf(0), BeEmpty(), Equal("worker link must be a channel")},
		{"link", "worker1", "worker2", reflect.ValueOf(make(chan interface{})), And(HaveLen(1), HaveKey("worker2:worker1")), BeNil()},
		{"link", "worker1", "worker2", reflect.ValueOf(make(chan interface{})), And(HaveLen(1), HaveKey("worker2:worker1")), BeNil()},
		{"link", "worker2", "worker1", reflect.ValueOf(make(chan interface{})), And(HaveLen(1), HaveKey("worker2:worker1")), BeNil()},
		{"link", "worker1", "worker2", reflect.ValueOf(make(chan int)), And(HaveLen(1), HaveKey("worker2:worker1")), Equal("worker link 'link' already exists between 'worker1' and 'worker2', but type is different")},
		{"link", "worker3", "worker2", reflect.ValueOf(make(chan chan string)), And(HaveLen(2), HaveKey("worker3:worker2")), BeNil()},
	} {
		func() {
			defer func() { Expect(recover()).Should(tcase.panicMatcher) }()
			scheduler.RegisterLink(tcase.origin, tcase.destination, tcase.name, tcase.link)
			Expect(scheduler.linksCache).Should(tcase.linksMatcher)
		}()
	}
}

func TestScheduler_MapLinks(t *testing.T) {
	RegisterTestingT(t)
	scheduler := &scheduler{linksCache: linkMap{}, workers: map[string]*worker{}}

	scheduler.RegisterWorker("worker1", naiveWorkerFactoryBuilder(t))
	scheduler.RegisterLink("worker1", "worker2", "_", reflect.ValueOf(make(chan interface{})))

	func() {
		defer func() { Expect(recover()).Should(Equal("worker 'worker2' not registered")) }()
		scheduler.mapLinks()
	}()
	scheduler.linksCache = linkMap{}

	scheduler.RegisterWorker("worker2", naiveWorkerFactoryBuilder(t))
	scheduler.RegisterWorker("worker3", naiveWorkerFactoryBuilder(t))
	scheduler.RegisterLink("worker1", "worker2", ".", reflect.ValueOf(make(chan interface{})))
	scheduler.RegisterLink("worker1", "worker2", "..", reflect.ValueOf(make(chan interface{})))

	scheduler.RegisterLink("worker1", "worker3", "...", reflect.ValueOf(make(chan interface{})))
	scheduler.RegisterLink("worker3", "worker2", "....", reflect.ValueOf(make(chan interface{})))

	func() {
		defer func() { Expect(recover()).Should(BeNil()) }()
		scheduler.mapLinks()
	}()

	Expect(scheduler.workers).Should(WithTransform(
		func(workers map[string]*worker) []string {
			var simplifiedLinks []string

			for _, worker := range workers {
				for linkedWorker, links := range worker.links {
					for linkName, _ := range links {
						simplifiedLinks = append(simplifiedLinks, fmt.Sprintf("%s>%s>%s", worker.name, linkedWorker, linkName))
					}
				}
			}
			return simplifiedLinks
		},
		ConsistOf(
			"worker1>worker2>.",
			"worker1>worker2>..",
			"worker1>worker3>...",

			"worker2>worker1>.",
			"worker2>worker1>..",
			"worker2>worker3>....",

			"worker3>worker1>...",
			"worker3>worker2>....",
		),
	))
}

func TestScheduler_Run(t *testing.T) {
	RegisterTestingT(t)
	TTL = 25 * time.Microsecond
	numGoroutine := runtime.NumGoroutine()

	shutdown := make(chan interface{})
	scheduler := &scheduler{linksCache: linkMap{}, workers: map[string]*worker{}, shutdownSchedulerChannel: shutdown}

	scheduler.RegisterWorker("workerFailure", func(links linkMap, flow WorkerFlow) error { return errUnrealistic })

	func() {
		defer func() { Expect(recover()).Should(Equal("failed to spawn 'workerFailure': " + errUnrealistic.Error())) }()
		scheduler.Run()
	}()
	scheduler.workers = map[string]*worker{}

	scheduler.RegisterWorker("worker1", completeWorkerFactoryBuilder(t, "worker2", "<", ">"))
	scheduler.RegisterWorker("worker2", completeWorkerFactoryBuilder(t, "worker1", ">", "<"))
	scheduler.RegisterLink("worker1", "worker2", ">", reflect.ValueOf(make(chan string)))
	scheduler.RegisterLink("worker1", "worker2", "<", reflect.ValueOf(make(chan string)))

	go func() {
		time.Sleep(25 * time.Millisecond)
		shutdown <- nil
	}()
	scheduler.Run()

	Eventually(func() int { return runtime.NumGoroutine() }, time.Second).Should(Equal(numGoroutine), "go routine not ended") // clean all goroutine
}

func TestScheduler_RunWithStuck(t *testing.T) {
	RegisterTestingT(t)
	TTL = time.Millisecond
	TTR = 500 * time.Microsecond
	numGoroutine := runtime.NumGoroutine()

	shutdown := make(chan interface{})
	scheduler := &scheduler{linksCache: linkMap{}, workers: map[string]*worker{}, shutdownSchedulerChannel: shutdown}

	in := make(chan string)
	out := make(chan string)
	scheduler.RegisterWorker("worker1", completeWorkerFactoryBuilder(t, "worker2", "<", ">"))
	scheduler.RegisterWorker("worker2", naiveWorkerFactoryBuilder(t))
	scheduler.RegisterLink("worker1", "worker2", "<", reflect.ValueOf(out))
	scheduler.RegisterLink("worker1", "worker2", ">", reflect.ValueOf(in))

	go func() {
		Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond).Should(Equal(numGoroutine + 6)) // 2 + 2 * (worker + heartbeat)

		out <- "..."                                                                                                  // stuck the worker 1
		Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond).Should(Equal(numGoroutine + 7)) // 2 + 2 * (worker + heartbeat) + stuck_worker

		<-in                                                                                                          // unstuck the worker
		Eventually(func() int { return runtime.NumGoroutine() }, 25*time.Millisecond).Should(Equal(numGoroutine + 6)) // 2 + 2 * (worker + heartbeat)
	}()
	go func() {
		time.Sleep(250 * time.Millisecond)
		shutdown <- nil
	}()
	scheduler.Run()

	Eventually(func() int { return runtime.NumGoroutine() }, time.Second).Should(Equal(numGoroutine), "go routine not ended") // clean all goroutine
}
