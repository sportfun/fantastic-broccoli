package utils

import (
	. "github.com/onsi/gomega"
	"sync"
	"testing"
)

func TestVolatile(t *testing.T) {
	RegisterTestingT(t)
	volatileChecking(NewVolatile)
}

func TestOneTimeVolatile(t *testing.T) {
	RegisterTestingT(t)

	volatileChecking(NewOneTimeVolatile)

	volatile := NewOneTimeVolatile(nil)
	Expect(volatile.Set(nil)).Should(Succeed())

	Expect(volatile.Set("another data")).Should(MatchError("volatile already set"))
	Expect(volatile.Get()).ShouldNot(Equal("another data"))
	Expect(volatile.Get()).Should(BeNil())
}

func TestIncrementVolatile(t *testing.T) {
	RegisterTestingT(t)

	volatile := NewIncrementVolatile(0).(Incremental)
	Expect(volatile.Set(0)).Should(Succeed())

	Expect(volatile.Set("another type")).Should(MatchError("increment volatile can be only set with integer"))
	Expect(volatile.Get()).Should(Equal(0))

	volatile.Inc(1)
	Expect(volatile.Get()).Should(Equal(1))

	volatile.Inc(3)
	Expect(volatile.Get()).Should(Equal(4))
}

func TestVolatile_RaceCondition(t *testing.T) {
	volatileRaceCondition(NewVolatile(nil))
}

func TestOneTimeVolatile_RaceCondition(t *testing.T) {
	volatileRaceCondition(NewOneTimeVolatile(nil))
}

func TestIncrementVolatile_RaceCondition(t *testing.T) {
	volatile := NewIncrementVolatile(0).(Incremental)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()

		for i := 0; i < 0xff; i++ {
			volatile.Inc(1)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 0xff; i++ {
			volatile.Get()
		}
	}()

	wg.Wait()

}

func volatileChecking(newVolatile func(interface{}) Volatile) {
	for _, value := range []interface{}{
		"string",
		true,
		0xFF,
		struct{}{},
		&struct{}{},
		0 + 0i,
	} {
		volatile := newVolatile(value)
		Expect(volatile.Get()).Should(Equal(value))
	}

	volatile := newVolatile(0)
	Expect(volatile.Set(nil)).Should(Succeed())
	Expect(volatile.Get()).Should(BeNil())
}

func volatileRaceCondition(volatile Volatile) {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()

		for i := 0; i < 0xff; i++ {
			volatile.Set(i)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 0xff; i++ {
			volatile.Get()
		}
	}()

	wg.Wait()
}
