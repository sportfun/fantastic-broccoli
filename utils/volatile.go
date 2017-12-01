package utils

import (
	"sync"
	"fmt"
)

type Volatile interface {
	Get() interface{}
	Set(interface{}) error
}

type volatile struct {
	m sync.Mutex
	v interface{}
}

type oneTimeVolatile struct {
	volatile
	alreadySet bool
}

func NewVolatile(value interface{}) Volatile {
	return &volatile{v: value}
}

func NewOneTimeVolatile(value interface{}) Volatile {
	return &oneTimeVolatile{volatile: volatile{v: value}}
}

func (s *volatile) Get() interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v
}

func (s *volatile) Set(value interface{}) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.v = value
	return nil
}

func (s *oneTimeVolatile) Set(value interface{}) error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.alreadySet {
		return fmt.Errorf("volatile already set")
	}

	s.alreadySet = true
	s.v = value
	return nil
}
