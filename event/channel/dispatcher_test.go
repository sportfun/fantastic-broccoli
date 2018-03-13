package channel

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestDispatch(t *testing.T) {
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

	Dispatch(nil)

	RegisterTestingT(t)
	for _, test := range cases {
		in := make(chan interface{})
		outs := []chan interface{}{make(chan interface{}), make(chan interface{}), make(chan interface{})}
		outss := []chan<- interface{}{outs[0], outs[1], outs[2]}
		Expect(Dispatch(in, outss...))

		in <- test.value

		for _, ch := range outs {
			Expect(<-ch).Should(test.valueMatcher)
		}
		close(in)
	}
}
