package module

import (
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
	"fmt"
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

func (s *Service) errorHandler(t errorType, err error, a ...interface{}) {
	obj := object.NewErrorObject(constant.ModuleService)

	switch t {
	case ModuleConfiguration:
		obj.Why(fmt.Errorf("failure during module ('%s') configuration: %s", a, err.Error()))
	case ModuleProcess:
		obj.Why(fmt.Errorf("failure during module ('%s') processing: %s", a, err.Error()))
	case ModuleStarting:
		obj.Why(fmt.Errorf("failure during module ('%s') starting: %s", a, err.Error()))
	case ModuleStop:
		obj.Why(fmt.Errorf("failure during module ('%s') stopping: %s", a, err.Error()))
	case NoModule:
		obj.Why(err)
	case PluginLoading:
		obj.Why(fmt.Errorf("failure during plugin ('%s') loading: %s", a, err.Error()))
	case SymbolLoading:
		obj.Why(fmt.Errorf("failure during module ('%s') loading: %s", a, err.Error()))
	}
	s.notifications.Notify(netBuilder.With(obj).Build())
}
