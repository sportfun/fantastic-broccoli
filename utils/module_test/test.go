package module_test

import (
	"testing"
	"time"

	. "github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/module"
	. "github.com/onsi/gomega"
)

var NProcesses = 5

// custom module test tool
func Test(t *testing.T, m module.Module, e *environment) {
	e.module = m
	e.queue = module.NewNotificationQueue()
	e.logger = log.NewDevelopment()

	t.Run("Start", func(t *testing.T) { eval(start, t, e) })

	if e.test.pre != nil {
		t.Run("PreTests", func(t *testing.T) {
			RegisterTestingT(t)
			e.test.pre(t, m)
		})
	}

	t.Run("Configure", func(t *testing.T) { eval(configure, t, e) })
	t.Run("StartSession", func(t *testing.T) { eval(startSessions, t, e) })
	t.Run("Process", func(t *testing.T) { eval(process, t, e) })
	t.Run("StopSession", func(t *testing.T) { eval(stopSession, t, e) })

	if e.test.post != nil {
		t.Run("PostTest", func(t *testing.T) {
			RegisterTestingT(t)
			e.test.post(t, NProcesses, e.module, e.queue)
		})
	}

	t.Run("Stop", func(t *testing.T) { eval(stop, t, e) })
}

func eval(fnc func(*testing.T, *environment), t *testing.T, environment *environment) {
	RegisterTestingT(t)
	fnc(t, environment)
}

func start(t *testing.T, e *environment) {
	Expect(e.module).Should(HaveState(UndefinedState))

	Expect(e.module.Start(nil, e.logger)).Should(HaveOccurred())             // failed: queue is nil
	Expect(e.module.Start(e.queue, nil)).Should(ExpectFor(e.module).Panic()) // failed: logger is nil

	Expect(e.module.Start(e.queue, e.logger)).Should(ExpectFor(e.module).SucceedWith(StartedState)) // succeed: start the module
}

func configure(t *testing.T, e *environment) {
	Expect(e.module.Configure(e.definition(t))).Should(ExpectFor(e.module).SucceedWith(IdleState)) // succeed: configure the module
}

func startSessions(t *testing.T, e *environment) {
	Expect(e.module).Should(HaveState(IdleState))

	Expect(e.module.Process()).Should(ExpectFor(e.module).SucceedWith(IdleState))         // failed: no session started, but no need to return error
	Expect(e.module.StartSession()).Should(ExpectFor(e.module).SucceedWith(WorkingState)) // succeed: start session successfully
	Expect(e.module.StartSession()).Should(ExpectFor(e.module).FailedWith(IdleState))     // failed: session already started
	Expect(e.module.StartSession()).Should(ExpectFor(e.module).SucceedWith(WorkingState)) // succeed: start session successfully
}

func process(t *testing.T, e *environment) {
	Expect(e.module).Should(HaveState(WorkingState))

	for i := 0; i < NProcesses; i++ {
		time.Sleep(e.tick)
		Expect(e.module.Process()).Should(Succeed())
	}
}

func stopSession(t *testing.T, e *environment) {
	Expect(e.module).Should(HaveState(WorkingState))

	Expect(e.module.StopSession()).Should(ExpectFor(e.module).SucceedWith(IdleState)) // succeed: stop session successfully
	Expect(e.module.Process()).Should(ExpectFor(e.module).SucceedWith(IdleState))     // failed: no session started, but no need to return error
	Expect(e.module.StopSession()).Should(ExpectFor(e.module).FailedWith(IdleState))  // failed: session already stopped
}

func stop(t *testing.T, e *environment) {
	Expect(e.module).Should(HaveState(IdleState))

	Expect(e.module.StartSession()).Should(ExpectFor(e.module).SucceedWith(WorkingState)) // succeed: start session successfully
	Expect(e.module.Stop()).Should(ExpectFor(e.module).SucceedWith(StoppedState))         // succeed: stop session & module

	time.Sleep(2 * e.tick) // wait all goroutine
}
