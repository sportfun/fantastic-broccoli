package module

import (
	"fmt"

	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/module"
	"github.com/sportfun/gakisitor/notification/object"
)

type pluginError byte
type moduleError func(*Manager, module.Module, error, *object.ErrorObject)

const (
	NoModule = iota
	PluginLoading
	SymbolLoading
	ModuleLoading
)

var pluginFailureLog = log.NewArgumentBinder("%s")

func (service *Manager) pluginFailure(t pluginError, err error, pluginName string) {
	obj := object.NewErrorObject(env.ModuleServiceEntity)

	switch t {
	case NoModule:
		obj.Why(err)
	case PluginLoading:
		obj.Why(fmt.Errorf("failure during plugin loading ('%s'): %s", pluginName, err.Error()))
	case SymbolLoading:
		obj.Why(fmt.Errorf("failure during symbol loading ('%s'): %s", pluginName, err.Error()))
	default:
		obj.Why(fmt.Errorf("unknown error type from ('%s'): %s", pluginName, err.Error()))
	}
	service.logger.Error(pluginFailureLog.Bind(obj.Reason))
	service.notifications.Notify(netBuilder.With(obj).Build())
}

func (service *Manager) checkIf(mod module.Module, err error, fnc moduleError) bool {
	if err == nil {
		return true
	}

	obj := object.NewErrorObject(env.ModuleServiceEntity)
	fnc(service, mod, err, obj)
	service.notifications.Notify(netBuilder.With(obj).Build())
	return false
}

func isStarted(service *Manager, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') starting: %s", mod.Name(), err.Error()))
}

func isConfigured(service *Manager, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') configuration: %s", mod.Name(), err.Error()))
}

func isProcessed(service *Manager, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') processing: %s", mod.Name(), err.Error()))
}

func isStopped(service *Manager, mod module.Module, err error, obj *object.ErrorObject) {
	obj.Why(fmt.Errorf("failure during module ('%s') stopping: %s", mod.Name(), err.Error()))
}
