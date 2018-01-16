package module

import (
	"fmt"
	"plugin"

	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/kernel"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/module"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
)

var netBuilder = notification.NewBuilder().
	From(env.ModuleServiceEntity).
	To(env.NetworkServiceEntity)

type moduleContainer map[string]module.Module

type Manager struct {
	modules moduleContainer
	state   byte

	messages      *module.NotificationQueue
	notifications *service.NotificationQueue
	logger        log.Logger
}

func init() {
	kernel.RegisterService(&Manager{})
}

func (service *Manager) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.modules = map[string]module.Module{}
	service.state = env.StartedState

	service.messages = module.NewNotificationQueue()
	service.notifications = notifications
	service.logger = logger

	return nil
}

func loadModules(service *Manager, config *config.GAkisitorConfig) moduleContainer {
	modules := moduleContainer{}

	for _, moduleDefinition := range config.Modules {
		plug, err := plugin.Open(moduleDefinition.Path)
		if err != nil {
			service.pluginFailure(PluginLoading, err, moduleDefinition.Path)
			continue
		}

		exporter, err := plug.Lookup("ExportModule")
		if err != nil {
			service.pluginFailure(SymbolLoading, err, moduleDefinition.Name)
			continue
		}

		module := exporter.(func() module.Module)()
		if module == nil {
			service.pluginFailure(ModuleLoading, nil, moduleDefinition.Name)
			continue
		}

		modules[module.Name()] = module

		if !service.checkIf(module, module.Start(service.messages, service.logger), isStarted) {
			delete(modules, module.Name())
			continue
		}

		if !service.checkIf(module, module.Configure(&moduleDefinition), isConfigured) {
			service.checkIf(module, module.Stop(), isStopped)
			delete(modules, module.Name())
			continue
		}
	}

	return modules
}

func (service *Manager) Configure(config *config.GAkisitorConfig) error {
	service.modules = loadModules(service, config)

	if len(service.modules) == 0 {
		err := fmt.Errorf("no module charged")
		service.pluginFailure(NoModule, err, "")
		service.state = env.PanicState
		return err
	}

	service.state = env.IdleState
	return nil
}

func (service *Manager) Process() error {
	service.state = env.WorkingState
	for _, notif := range service.notifications.Notifications(service.Name()) {
		service.handle(notif)
	}

	for _, mod := range service.modules {
		service.checkIf(mod, mod.Process(), isProcessed)
	}

	for _, notif := range service.messages.Notifications() {
		switch obj := notif.Content().(type) {
		case *module.ErrorObject:
			netBuilder.With(&obj.ErrorObject)

			if obj.ErrorLevel() == env.FatalLevel {
				mod := obj.From()
				if mod != nil {
					service.checkIf(mod, mod.Stop(), isStopped)
				}
			}
		case *object.DataObject:
			netBuilder.With(obj)
		default:
			continue
		}
		service.notifications.Notify(netBuilder.Build())
	}

	service.state = env.IdleState
	return nil
}

func (service *Manager) Stop() error {
	for _, m := range service.modules {
		service.checkIf(m, m.Stop(), isStopped)
	}
	return nil
}

func (service *Manager) Name() string {
	return env.ModuleServiceEntity
}

func (service *Manager) State() byte {
	return service.state
}
