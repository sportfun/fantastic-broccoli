package module

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/const"
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/network"
)

type ErrorType int

var netBuilder *notification.Builder = notification.Builder{}.From(_const.MODULE_SERVICE).To(_const.NETWORK_SERVICE)

const (
	ModuleConfiguration ErrorType = iota
	ModuleProcess
	ModuleStarting
	NoModule
	PluginLoading
	Stop
	SymbolLoading
)

func (s *Service) errorHandler(t ErrorType, e error) {
	if e == nil {
		return
	}

	m := network.NewMessage("error").AddArgument(e.Error())
	switch t {
	case ModuleConfiguration:
	case ModuleProcess:
	case ModuleStarting:
	case NoModule:
		//TODO: Blink Failure LED
		s.notifications.Notify(netBuilder.With(m).Build())
	case Stop:
	default:
		s.notifications.Notify(netBuilder.With(m).Build())
	}
}
