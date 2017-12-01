package module_test

import (
	"testing"
	"time"

	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/module"
	"github.com/sportfun/gakisitor/log"
)

type definitionFactory func(interface{}) *config.ModuleDefinition
type preTest func(*testing.T, module.Module)
type postTest func(*testing.T, int, module.Module, *module.NotificationQueue)

type environment struct {
	definition definitionFactory
	tick       time.Duration

	test struct {
		pre  preTest
		post postTest
	}

	module module.Module
	queue  *module.NotificationQueue
	logger log.Logger
}

// Create new test environment for module
func NewEnvironment(factory definitionFactory, pre preTest, post postTest, tick time.Duration) *environment {
	return &environment{
		definition: factory,
		tick:       tick,
		test: struct {
			pre  preTest
			post postTest
		}{pre: pre, post: post},
	}
}
