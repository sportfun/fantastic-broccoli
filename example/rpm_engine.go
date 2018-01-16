package main

import "math/rand"

type rpmEngine struct {
	min       float64
	max       float64
	step      float64
	precision float64

	rand    *rand.Rand
	lastval int
}

func (e *rpmEngine) NewValue() float64 {
	rpm := int(e.rand.Float64() * (e.max - e.min) * e.precision)

	if e.lastval == 0 {
		e.lastval = rpm
		return e.min + float64(rpm/int(e.precision))
	}

	e.lastval = (int(e.lastval) - rpm) % int(e.step)

	return e.min + float64(e.lastval/int(e.precision))
}
