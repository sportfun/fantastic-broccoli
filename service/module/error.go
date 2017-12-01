package module

import (
	"fmt"

	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/module"
	"github.com/xunleii/fantastic-broccoli/notification/object"
)

type pluginError byte
type moduleError func(*Service, module.Module, error, *object.ErrorObject)

const (
	NoModule = iota
	PluginLoading
	SymbolLoading
	ModuleLoading
)

var pluginFailureLog = log.NewArgumentBinder("%s")

func (service *Service) pluginFailure(t pluginError, err error, a ...interface{}) {
	obj := object.NewErrorObject(env.ModuleServiceEntity)

	switch t {
	case NoModule:
		obj.Why(err)
	case PluginLoading:
		obj.Why(fmt.Errorf("failure during plugin loading ('%s'): %s", a, err.Error()))
	case SymbolLoading:
		obj.Why(fmt.Errorf("failure during module loading ('%s'): %s", a, err.Error()))
	}
	service.logger.Error(pluginFailureLog.Bind(obj.Reason))
	service.notifications.Notify(netBuilder.With(obj).Build())
}

func (service *Service) checkIf(mod module.Module, err error, fnc moduleError) bool {
	if err == nil {
		return true
	}

	obj := object.NewErrorObject(env.ModuleServiceEntity)
	fnc(service, mod, err, obj)
	service.notifications.Notify(netBuilder.With(obj).Build())
	return false
}

func IsStarted(service *Service, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') starting: %s", mod.Name(), err.Error()))
}

func IsConfigured(service *Service, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') configuration: %s", mod.Name(), err.Error()))
}

func IsProcessed(service *Service, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') processing: %s", mod.Name(), err.Error()))
}

func IsStopped(service *Service, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') stopping: %s", mod.Name(), err.Error()))
}
