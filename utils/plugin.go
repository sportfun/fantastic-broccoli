package utils

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type pluginTools struct{}
type PropertyFactory func() *properties.Properties
type SpecializedTest func(*testing.T, int, *module.NotificationQueue)

var Plugin = pluginTools{}

func (pgt *pluginTools) Benchmark(
	t *testing.T,
	m module.Module,
	pfactory PropertyFactory,
	tick time.Duration,
) {
	t.Logf("--- Execute Benchmark ---")
	bench := testing.Benchmark(func(b *testing.B) {
		t.Logf("-------------------------")
		// Init environment
		queue := module.NewNotificationQueue()
		log, _ := zap.NewProduction()

		// Start & Configure module
		t.Logf("- Start module '%s'\n", m.Name())
		if err := m.Start(queue, log); err != nil {
			t.Fatalf("! Failure during module starting - %s\n", err.Error())
		}

		t.Logf("- Configure module '%s'\n", m.Name())
		if err := m.Configure(pfactory()); err != nil {
			t.Fatalf("! Failure during module configuration - %s\n", err.Error())
		}

		// Start bench
		t.Logf("- Start benchmark\n")

		t.Logf("\t- Start new session\n")
		if err := m.StartSession(); err != nil {
			t.Fatalf("! Failure during starting session - %s\n", err.Error())
		}

		b.ResetTimer()
		t.Logf("\t- Processing loops [%d time(s)]\n", b.N)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			if err := m.Process(); err != nil {
				// TODO: Failure only if critical (v.1.X)
				t.Logf("! Failure during processing - %s\n", err.Error())
			}
			time.Sleep(tick)
		}
		b.StopTimer()

		t.Logf("- Stop module '%s'\n", m.Name())
		m.Stop()

		time.Sleep(tick)
		t.Logf("-------------------------")
	})

	t.Logf("--- Benchmark results ---")

	realTime := bench.T - tick*time.Duration(bench.N)
	opPerSecond := time.Second / (realTime/time.Duration(bench.N))
	t.Logf("time: %dns\t\t\t%d ns/op (%d ops)", realTime, realTime/time.Duration(bench.N), opPerSecond)
	t.Logf("memory: %dB\t\t\t%d allocs", bench.MemBytes, bench.MemAllocs)
}

func (pgt *pluginTools) Test(
	t *testing.T,
	m module.Module,
	pfactory PropertyFactory,
	sptest SpecializedTest,
	nprocesses int,
	tick time.Duration,
) {
	// Init environment
	queue := module.NewNotificationQueue()
	log, _ := zap.NewProduction()

	// Start & Configure module
	t.Logf("- Start module '%s'\n", m.Name())
	if err := m.Start(queue, log); err != nil {
		t.Fatalf("! Failure during module starting - %s\n", err.Error())
	}

	AssertEquals(t, constant.States.Started, m.State())
	if err := m.Configure(pfactory()); err != nil {
		t.Fatalf("! Failure during module configuration - %s\n", err.Error())
	}

	// Start tests

	t.Logf("- Start tests\n")
	// Error verification (Process without session)
	// TODO: Explicit error knowable by the manager
	// m.Process()

	t.Logf("\t- Start new session\n")
	if err := m.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	// Error verification (Session already exist)
	// TODO: Explicit error knowable by the manager
	// m.StartSession()

	t.Logf("\t- Processing loops [%d time(s)]\n", nprocesses)
	for i := 0; i < nprocesses; i++ {
		time.Sleep(tick)
		if err := m.Process(); err != nil {
			// TODO: Failure only if critical (v.1.X)
			t.Logf("! Failure during processing - %s\n", err.Error())
		}
	}

	t.Logf("\t- Stop session\n")
	if err := m.StopSession(); err != nil {
		t.Fatalf("! Failure during ending session - %s\n", err.Error())
	}
	// Error verification (Process without session)
	// TODO: Explicit error knowable by the manager
	// m.Process()

	if sptest != nil {
		t.Logf("\t- Launch specialized tests\n")
		sptest(t, nprocesses, queue)
	}

	t.Logf("- Stop module '%s'\n", m.Name())
	if err := m.Stop(); err != nil {
		t.Fatalf("! Failure during ending module - %s\n", err.Error())
	}
	time.Sleep(2 * tick)
}
