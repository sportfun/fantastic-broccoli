package utils

import "sync"

type ConditionalWaitGroup struct {
	sync.WaitGroup
}

func (wg *ConditionalWaitGroup) AddIf(delta int, condition bool) {
	if condition {
		wg.Add(delta)
	}
}

func (wg *ConditionalWaitGroup) WaitIf(condition bool) {
	if condition {
		wg.Wait()
	}
}

func (wg *ConditionalWaitGroup) DoneIf(condition bool) {
	if condition {
		wg.Done()
	}
}