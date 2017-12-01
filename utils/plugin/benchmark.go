package plugin

import (
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/module"
)

// Custom module benchmark tool
func Benchmark(t *testing.T, mod module.Module, env *environment) {

	t.Logf("--- Execute Benchmark ---")
	bench := testing.Benchmark(func(b *testing.B) {

		// Init environment
		queue := module.NewNotificationQueue()
		logger := log.NewDevelopment()

		// Start & Configure module
		if err := mod.Start(queue, logger); err != nil {
			t.Fatalf("! Failure during module starting - %s\n", err.Error())
		}

		if err := mod.Configure(env.definition(b)); err != nil {
			t.Fatalf("! Failure during module configuration - %s\n", err.Error())
		}

		// Start bench
		if err := mod.StartSession(); err != nil {
			t.Fatalf("! Failure during starting session - %s\n", err.Error())
		}

		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			if err := mod.Process(); err != nil {
				// TODO: Failure only if critical (v.1.X)
				t.Logf("! Failure during processing - %s\n", err.Error())
			}
			time.Sleep(env.tick)
		}
		b.StopTimer()

		mod.Stop()
		time.Sleep(env.tick)

	})

	t.Logf("--- Benchmark results ---")

	realTime := bench.T - env.tick*time.Duration(bench.N)
	opPerSecond := time.Second / (realTime / time.Duration(bench.N))
	t.Logf("time: %dns\t\t\t%d ns/op (%d ops)", realTime, realTime/time.Duration(bench.N), opPerSecond)
	t.Logf("memory: %dB\t\t\t%d allocs", bench.MemBytes, bench.MemAllocs)
}
