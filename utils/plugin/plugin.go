package plugin

import (
	"fmt"
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/utils"
)

type InternalLogger func(format string, a ...interface{})
type definitionFactory func(interface{}) properties.ModuleDefinition
type preTest func(*testing.T, InternalLogger, module.Module)
type postTest func(*testing.T, InternalLogger, int, module.Module, *module.NotificationQueue)

type testEnvironment struct {
	definition definitionFactory
	tick       time.Duration

	test struct {
		pre  preTest
		post postTest
	}
}

func NewEnvironment(factory definitionFactory, pre preTest, post postTest, tick time.Duration) *testEnvironment {
	return &testEnvironment{
		definition: factory,
		tick:       tick,
		test: struct {
			pre  preTest
			post postTest
		}{pre: pre, post: post},
	}
}

func Benchmark(t *testing.T, mod module.Module, env *testEnvironment) {
	t.Logf("--- Execute Benchmark ---")
	bench := testing.Benchmark(func(b *testing.B) {
		t.Logf("-------------------------")
		// Init environment
		queue := module.NewNotificationQueue()
		logger := log.NewLogger.Dev(nil)

		// Start & Configure module
		t.Logf("- Start module '%s'\n", mod.Name())
		if err := mod.Start(queue, logger); err != nil {
			t.Fatalf("! Failure during module starting - %s\n", err.Error())
		}

		t.Logf("- Configure module '%s'\n", mod.Name())
		if err := mod.Configure(env.definition(b)); err != nil {
			t.Fatalf("! Failure during module configuration - %s\n", err.Error())
		}

		// Start bench
		t.Logf("- Start benchmark\n")

		t.Logf("\t- Start new session\n")
		if err := mod.StartSession(); err != nil {
			t.Fatalf("! Failure during starting session - %s\n", err.Error())
		}

		b.ResetTimer()
		t.Logf("\t- Processing loops [%d time(s)]\n", b.N)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			if err := mod.Process(); err != nil {
				// TODO: Failure only if critical (v.1.X)
				t.Logf("! Failure during processing - %s\n", err.Error())
			}
			time.Sleep(env.tick)
		}
		b.StopTimer()

		t.Logf("- Stop module '%s'\n", mod.Name())
		mod.Stop()

		time.Sleep(env.tick)
		t.Logf("-------------------------")
	})

	t.Logf("--- Benchmark results ---")

	realTime := bench.T - env.tick*time.Duration(bench.N)
	opPerSecond := time.Second / (realTime / time.Duration(bench.N))
	t.Logf("time: %dns\t\t\t%d ns/op (%d ops)", realTime, realTime/time.Duration(bench.N), opPerSecond)
	t.Logf("memory: %dB\t\t\t%d allocs", bench.MemBytes, bench.MemAllocs)
}

func Test(t *testing.T, mod module.Module, nprocesses int, env *testEnvironment) {
	// Init environment
	queue := module.NewNotificationQueue()
	logger := log.NewLogger.Dev(nil)
	inlog := func(format string, a ...interface{}) { t.Logf("\t> %s", fmt.Sprintf(format, a...)) }



	// Start & Configure module
	t.Logf("- Start module '%s'\n", mod.Name())
	if err := mod.Start(queue, logger); err != nil {
		t.Fatalf("! Failure during module starting - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Started, mod.State())

	// Start pre tests
	if env.test.pre != nil {
		t.Logf("- Start pre tests\n")
		env.test.pre(t, inlog, mod)
	}

	t.Logf("- Start module '%s'\n", mod.Name())
	if err := mod.Configure(env.definition(t)); err != nil {
		t.Fatalf("! Failure during module configuration - %s\n", err.Error())
	}



	// Start tests
	t.Logf("- Start tests\n")
	// Error verification (Process without session)
	err := mod.Process()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.ErrorLevels.Warning, err.ErrorLevel())

	t.Logf("\t- Start new session\n")
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	// Error verification (Session already exist)
	err = mod.StartSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.ErrorLevels.Warning, err.ErrorLevel())

	t.Logf("\t- Processing loops [%d time(s)]\n", nprocesses)
	for i := 0; i < nprocesses; i++ {
		time.Sleep(env.tick)
		if err := mod.Process(); err != nil {
			// TODO: Failure only if critical (v.1.X)
			t.Logf("! Failure during processing - %s\n", err.Error())
		}
	}

	t.Logf("\t- Stop session\n")
	if err := mod.StopSession(); err != nil {
		t.Fatalf("! Failure during ending session - %s\n", err.Error())
	}
	// Error verification (Process without session)
	err = mod.Process()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.ErrorLevels.Warning, err.ErrorLevel())
	// Error verification (Session already stopped)
	err = mod.StopSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.ErrorLevels.Warning, err.ErrorLevel())



	// Start post tests
	if env.test.post != nil {
		t.Logf("- Start post tests\n")
		env.test.post(t, inlog, nprocesses, mod, queue)
	}



	// Stop all
	t.Logf("- Stop module '%s'\n", mod.Name())
	if err := mod.Stop(); err != nil {
		t.Fatalf("! Failure during ending module - %s\n", err.Error())
	}
	time.Sleep(2 * env.tick)
}
