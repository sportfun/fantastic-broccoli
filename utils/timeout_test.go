package utils

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

type exT struct {
	testing.T
	hasFailed bool
}

func (e *exT) Fatalf(f string, a ...interface{}) { e.hasFailed = true }

func TestReleaseIfTimeout(t *testing.T) {
	RegisterTestingT(t)

	e := exT{}
	ReleaseIfTimeout(&e, TimeoutPrecision, func(testing.TB) {})
	Expect(e.hasFailed).Should(BeFalse())
	ReleaseIfTimeout(&e, TimeoutPrecision, func(testing.TB) {
		time.Sleep(time.Second)
	})
	Expect(e.hasFailed).Should(BeTrue())
}
