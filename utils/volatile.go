package utils

import (
	"fmt"
	"sync"
)

type Volatile interface {
	Get() interface{}
	Set(interface{}) error
}

type Incremental interface {
	Volatile
	Inc(int)
}

type volatile struct {
	m sync.Mutex
	v interface{}
}

type oneTimeVolatile struct {
	volatile
	alreadySet bool
}

type incrementVolatile struct {
	volatile
}

func NewVolatile(value interface{}) Volatile {
	return &volatile{v: value}
}

func NewOneTimeVolatile(value interface{}) Volatile {
	return &oneTimeVolatile{volatile: volatile{v: value}}
}

func NewIncrementVolatile(value int) Volatile {
	return &incrementVolatile{volatile: volatile{v: value}}
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

func (s *incrementVolatile) Set(value interface{}) error {
	if _, isInteger := value.(int); !isInteger {
		return fmt.Errorf("increment volatile can be only set with integer")
	}
	return s.volatile.Set(value)
}

func (s *incrementVolatile) Inc(value int) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, isInteger := s.v.(int); !isInteger {
		s.v = 0
	}
	s.v = s.v.(int) + value
}
