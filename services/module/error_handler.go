package module

import (
	"fantastic-broccoli/common/types/network"
	"fantastic-broccoli/const"
	"fmt"
)

type ErrorType int

const (
	ModuleConfiguration ErrorType = iota
	ModuleProcess
	ModuleStarting
	ModuleStop
	NoModule
	PluginLoading
	SymbolLoading
)

func (s *Service) errorHandler(t ErrorType, e error, p ...interface{}) {
	if e == nil {
		return
	}

	m := network.NewMessage("error").
		AddArgument(string(_const.ModuleService))

	switch t {
	case ModuleConfiguration:
		m.AddArgument(fmt.Sprintf("failure during module ('%s') configuration: %s", p, e.Error()))
	case ModuleProcess:
		m.AddArgument(fmt.Sprintf("failure during module ('%s') processing: %s", p, e.Error()))
	case ModuleStarting:
		m.AddArgument(fmt.Sprintf("failure during module ('%s') starting: %s", p, e.Error()))
	case ModuleStop:
		m.AddArgument(fmt.Sprintf("failure during module ('%s') stopping: %s", p, e.Error()))
	case NoModule:
		//TODO: Blink Failure LED
		m.AddArgument(e.Error())
	case PluginLoading:
		m.AddArgument(fmt.Sprintf("failure during plugin ('%s') loading: %s", p, e.Error()))
	case SymbolLoading:
		m.AddArgument(fmt.Sprintf("failure during module ('%s') loading: %s", p, e.Error()))
	}
	s.notifications.Notify(netBuilder.With(m).Build())
}
