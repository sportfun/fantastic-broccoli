package module_test

import (
	"testing"
	"time"

	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/module"
)

// Custom module benchmark tool
func Benchmark(root *testing.B, mod module.Module, env *environment) {

	root.Logf("--- Execute Benchmark ---")
	bench := testing.Benchmark(func(b *testing.B) {

		// Init environment
		queue := module.NewNotificationQueue()
		logger := log.NewDevelopment()

		// Start & Configure module
		if err := mod.Start(queue, logger); err != nil {
			root.Fatalf("! Failure during module starting - %s\n", err.Error())
		}

		if err := mod.Configure(env.definition(b)); err != nil {
			root.Fatalf("! Failure during module configuration - %s\n", err.Error())
		}

		// Start bench
		if err := mod.StartSession(); err != nil {
			root.Fatalf("! Failure during starting session - %s\n", err.Error())
		}

		b.ResetTimer()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			if err := mod.Process(); err != nil {
				// TODO: Failure only if critical (v.1.X)
				root.Logf("! Failure during processing - %s\n", err.Error())
			}
			time.Sleep(env.tick)
		}
		b.StopTimer()

		mod.Stop()
		time.Sleep(env.tick)

	})

	root.Logf("--- Benchmark results ---")

	realTime := bench.T - env.tick*time.Duration(bench.N)
	opPerSecond := time.Second / (realTime / time.Duration(bench.N))
	root.Logf("time: %dns\t\t\t%d ns/op (%d op/s)", realTime, realTime/time.Duration(bench.N), opPerSecond)
	root.Logf("memory: %dB\t\t\t%d allocs", bench.MemBytes, bench.MemAllocs)
}
