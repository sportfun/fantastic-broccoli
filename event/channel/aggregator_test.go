package channel

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestAggregate(t *testing.T) {
	type T struct{ int }
	tPtr := &T{}
	cases := []struct {
		value        interface{}
		valueMatcher types.GomegaMatcher
	}{
		{value: 0xEF, valueMatcher: Equal(0xEF)},
		{value: T{0xAB}, valueMatcher: Equal(T{0xAB})},
		{value: tPtr, valueMatcher: Equal(tPtr)},
		{value: make(chan int), valueMatcher: WithTransform(func(c chan int) reflect.Kind { return reflect.TypeOf(c).Kind() }, Equal(reflect.Chan))},
	}

	Aggregate(nil)

	RegisterTestingT(t)
	for _, test := range cases {
		out := make(chan interface{})
		ins := []chan interface{}{make(chan interface{}), make(chan interface{}), make(chan interface{})}
		inss := []<-chan interface{}{ins[0], ins[1], ins[2]}
		Expect(Aggregate(out, inss...))

		for _, ch := range ins {
			ch <- test.value
		}

		for _, ch := range ins {
			Expect(<-out).Should(test.valueMatcher)
			close(ch)
		}
	}
}
