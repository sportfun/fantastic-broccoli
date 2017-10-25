package module

import (
	"fantastic-broccoli/constant"
	"fmt"
	"fantastic-broccoli/common/types/notification/object"
)

type errorType int

const (
	ModuleConfiguration errorType = iota
	ModuleProcess
	ModuleStarting
	ModuleStop
	NoModule
	PluginLoading
	SymbolLoading
)

func (s *Service) errorHandler(t errorType, e error, p ...interface{}) {
	if e == nil {
		return
	}

	m := object.NewErrorObject(constant.ModuleService)

	switch t {
	case ModuleConfiguration:
		m.Why(fmt.Errorf("failure during module ('%s') configuration: %s", p, e.Error()))
	case ModuleProcess:
		m.Why(fmt.Errorf("failure during module ('%s') processing: %s", p, e.Error()))
	case ModuleStarting:
		m.Why(fmt.Errorf("failure during module ('%s') starting: %s", p, e.Error()))
	case ModuleStop:
		m.Why(fmt.Errorf("failure during module ('%s') stopping: %s", p, e.Error()))
	case NoModule:
		//TODO: Blink Failure LED
		m.Why(e)
	case PluginLoading:
		m.Why(fmt.Errorf("failure during plugin ('%s') loading: %s", p, e.Error()))
	case SymbolLoading:
		m.Why(fmt.Errorf("failure during module ('%s') loading: %s", p, e.Error()))
	}
	s.notifications.Notify(netBuilder.With(m).Build())
}
