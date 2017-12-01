package plugin

import (
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/module"
)

type InternalLogger func(format string, a ...interface{})
type definitionFactory func(interface{}) *config.ModuleDefinition
type preTest func(*testing.T, InternalLogger, module.Module)
type postTest func(*testing.T, InternalLogger, int, module.Module, *module.NotificationQueue)

type environment struct {
	definition definitionFactory
	tick       time.Duration

	test struct {
		pre  preTest
		post postTest
	}
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
