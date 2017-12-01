package module

import (
	"fmt"
	"plugin"

	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/module"
	"github.com/xunleii/fantastic-broccoli/notification"
	"github.com/xunleii/fantastic-broccoli/notification/object"
	"github.com/xunleii/fantastic-broccoli/service"
	"github.com/xunleii/fantastic-broccoli/kernel"
)

var netBuilder = notification.NewBuilder().
	From(env.ModuleServiceEntity).
	To(env.NetworkServiceEntity)

type moduleContainer map[string]module.Module

type Service struct {
	modules moduleContainer
	state   byte

	messages      *module.NotificationQueue
	notifications *service.NotificationQueue
	logger        log.Logger
}

func init() {
	kernel.RegisterService(&Service{})
}

func (service *Service) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.modules = map[string]module.Module{}
	service.state = env.StartedState

	service.messages = module.NewNotificationQueue()
	service.notifications = notifications
	service.logger = logger

	return nil
}

func loadModules(service *Service, config *config.GAkisitorConfig) moduleContainer {
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

		if !service.checkIf(module, module.Start(service.messages, service.logger), IsStarted) {
			delete(modules, module.Name())
			continue
		}

		if !service.checkIf(module, module.Configure(&moduleDefinition), IsConfigured) {
			service.checkIf(module, module.Stop(), IsStopped)
			delete(modules, module.Name())
			continue
		}
	}

	return modules
}

func (service *Service) Configure(config *config.GAkisitorConfig) error {
	service.modules = loadModules(service, config)

	if len(service.modules) == 0 {
		err := fmt.Errorf("no module charged")
		service.pluginFailure(NoModule, err)
		service.state = env.PanicState
		return err
	}

	service.state = env.IdleState
	return nil
}

func (service *Service) Process() error {
	service.state = env.WorkingState
	for _, notif := range service.notifications.Notifications(service.Name()) {
		service.handle(notif)
	}

	for _, mod := range service.modules {
		service.checkIf(mod, mod.Process(), IsProcessed)
	}

	for _, notif := range service.messages.Notifications() {
		switch obj := notif.Content().(type) {
		case *module.ErrorObject:
			netBuilder.With(&obj.ErrorObject)

			if obj.ErrorLevel() == env.FatalLevel {
				mod := service.modules[notif.From()]
				service.checkIf(mod, mod.Stop(), IsStopped)
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

func (service *Service) Stop() error {
	for _, m := range service.modules {
		service.checkIf(m, m.Stop(), IsStopped)
	}
	return nil
}

func (service *Service) Name() string {
	return env.ModuleServiceEntity
}

func (service *Service) State() byte {
	return service.state
}
