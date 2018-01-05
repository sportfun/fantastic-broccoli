package utils

import (
	"testing"
	"time"
	"runtime"
	"reflect"
)

var TimeoutPrecision = 50 * time.Millisecond

func ReleaseIfTimeout(t testing.TB, duration time.Duration, fnc func(t testing.TB)) {
	volatile := NewOneTimeVolatile(nil)

	go func() {
		time.Sleep(duration)
		volatile.Set(false)
	}()

	go func() {
		fnc(t)
		volatile.Set(true)
	}()

	for volatile.Get() == nil {
		time.Sleep(TimeoutPrecision)
	}

	if !volatile.Get().(bool) {
		t.Fatalf("function '%s' has timeout", runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name())
	}
}
