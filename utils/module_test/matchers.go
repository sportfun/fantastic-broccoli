package module_test

import (
	"github.com/sportfun/gakisitor/module"
	. "github.com/sportfun/gakisitor/env"
	"github.com/onsi/gomega/types"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega"
)

var statesName = map[byte]string{
	UndefinedState: "UndefinedState",
	StartedState:   "StartedState",
	IdleState:      "IdleState",
	WorkingState:   "WorkingState",
	StoppedState:   "StoppedState",
	PanicState:     "PanicState",
}

type ExternalStateMatcher struct {
	// input
	State  byte
	Module module.Module
}

func (matcher *ExternalStateMatcher) Match(a interface{}) (success bool, err error) {
	if matcher.Module == nil {
		return a.(module.Module).State() == matcher.State, nil
	}
	return matcher.Module.State() == matcher.State, nil
}

func (matcher *ExternalStateMatcher) FailureMessage(a interface{}) (message string) {
	if matcher.Module == nil {
		return format.Message(statesName[a.(module.Module).State()], "to equal", statesName[matcher.State])
	}
	return format.Message(statesName[matcher.Module.State()], "to equal", statesName[matcher.State])
}

func (matcher *ExternalStateMatcher) NegatedFailureMessage(a interface{}) (message string) {
	if matcher.Module == nil {
		return format.Message(statesName[a.(module.Module).State()], "to equal", statesName[matcher.State])
	}
	return format.Message(statesName[matcher.Module.State()], "to equal", statesName[matcher.State])
}

func HaveState(state byte, module ...module.Module) types.GomegaMatcher {
	if len(module) > 0 {
		return &ExternalStateMatcher{
			Module: module[0],
			State:  state,
		}
	}

	return &ExternalStateMatcher{
		Module: nil,
		State:  state,
	}
}

type ModuleExpectation struct {
	module module.Module
}

func ExpectFor(m module.Module) ModuleExpectation { return ModuleExpectation{module: m} }

func (m ModuleExpectation) SucceedWith(state byte) types.GomegaMatcher {
	return And(
		Succeed(),
		HaveState(state, m.module),
	)
}

func (m ModuleExpectation) FailedWith(state byte) types.GomegaMatcher {
	return And(
		HaveOccurred(),
		HaveState(state, m.module),
	)
}

func (m ModuleExpectation) Panic() types.GomegaMatcher {
	return m.FailedWith(PanicState)
}
