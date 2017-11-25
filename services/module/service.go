package module

import (
	"fmt"
	"plugin"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

var netBuilder = notification.NewBuilder().
	From(constant.EntityNames.Services.Module).
	To(constant.EntityNames.Services.Network)

type moduleContainer map[string]module.Module

type Service struct {
	modules moduleContainer
	state   types.StateType

	messages      *module.NotificationQueue
	notifications *service.NotificationQueue
	logger        log.Logger
}

func (service *Service) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.modules = map[string]module.Module{}
	service.state = constant.States.Started

	service.messages = module.NewNotificationQueue()
	service.notifications = notifications
	service.logger = logger

	return nil
}

func loadModules(service *Service, props *properties.Properties) moduleContainer {
	modules := moduleContainer{}

	for _, moduleDefinition := range props.Modules {
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

		module := exporter.(func() (module.Module))()
		if module == nil {
			service.pluginFailure(ModuleLoading, nil, moduleDefinition.Name)
			continue
		}

		modules[module.Name()] = module

		if !service.checkIf(module, module.Start(service.messages, service.logger), IsStarted) {
			delete(modules, module.Name())
			continue
		}

		if !service.checkIf(module, module.Configure(moduleDefinition), IsConfigured) {
			service.checkIf(module, module.Stop(), IsStopped)
			delete(modules, module.Name())
			continue
		}
	}

	return modules
}

func (service *Service) Configure(props *properties.Properties) error {
	service.modules = loadModules(service, props)

	if len(service.modules) == 0 {
		err := fmt.Errorf("no module charged")
		service.pluginFailure(NoModule, err)
		service.state = constant.States.Panic
		return err
	}

	service.state = constant.States.Idle
	return nil
}

func (service *Service) Process() error {
	service.state = constant.States.Working
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

			if obj.ErrorLevel() == constant.ErrorLevels.Fatal {
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

	service.state = constant.States.Idle
	return nil
}

func (service *Service) Stop() error {
	for _, m := range service.modules {
		service.checkIf(m, m.Stop(), IsStopped)
	}
	return nil
}

func (service *Service) Name() string {
	return constant.EntityNames.Services.Module
}

func (service *Service) State() types.StateType {
	return service.state
}
