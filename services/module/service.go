package module

import (
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/model"
	"fmt"
	"go.uber.org/zap"
	"plugin"
)

var netBuilder = notification.NewBuilder().
	From(constant.ModuleService).
	To(constant.NetworkService)

type Service struct {
	modules map[string]module.Module
	state   int

	messages      *module.NotificationQueue
	notifications *service.NotificationQueue
	logger        *zap.Logger
}

func (s *Service) Start(notification *service.NotificationQueue, logger *zap.Logger) error {
	s.modules = map[string]module.Module{}
	s.state = constant.Started

	s.messages = module.NewNotificationQueue()
	s.notifications = notification
	s.logger = logger

	return nil
}

func (s *Service) Configure(props *model.Properties) error {
	for _, e := range props.Modules {
		p, err := plugin.Open(string(e.Path))
		if err != nil {
			s.errorHandler(PluginLoading, err, e.Path)
			continue
		}

		ex, err := p.Lookup("ExportModule")
		if err != nil {
			s.errorHandler(SymbolLoading, err, e.Name)
			continue
		}

		mod := ex.(func() (module.Module))()
		s.modules[mod.Name()] = mod

		if err = mod.Start(s.messages, s.logger); err != nil {
			s.errorHandler(ModuleStarting, err, e.Name)
			delete(s.modules, mod.Name())
			continue
		}

		if err = mod.Configure(props); err != nil {
			s.errorHandler(ModuleConfiguration, err, e.Name)
			if err = mod.Stop(); err != nil {
				s.errorHandler(ModuleStop, err, e.Name)
			}
			delete(s.modules, mod.Name())
		}
	}

	if len(s.modules) == 0 {
		err := fmt.Errorf("no module charged")
		s.errorHandler(NoModule, err)
		s.state = constant.Stopped
		return err
	}

	s.state = constant.Idle
	return nil
}

func (s *Service) Process() error {
	s.state = constant.Working
	for _, n := range s.notifications.Notifications(s.Name()) {
		s.notificationHandler(n)
	}

	for _, m := range s.modules {
		if err := m.Process(); err != nil {
			s.errorHandler(ModuleProcess, err)
		}
	}

	for _, n := range s.messages.Notifications() {
		switch obj := n.Content().(type) {
		case module.ErrorObject:
			netBuilder.With(obj.ErrorObject)

			if obj.ErrorLevel() == constant.Fatal {
				if err := s.modules[n.From()].Stop(); err != nil {
					s.errorHandler(ModuleStop, err, n.From())
				}
			}
		case object.DataObject:
			netBuilder.With(obj)
		default:
			continue
		}
		s.notifications.Notify(netBuilder.Build())
	}

	s.state = constant.Idle
	return nil
}

func (s *Service) Stop() error {
	for n, m := range s.modules {
		if err := m.Stop(); err != nil {
			s.errorHandler(ModuleStop, err, n)
		}
	}
	return nil
}

func (s *Service) Name() string {
	return constant.ModuleService
}

func (s *Service) State() int {
	return s.state
}
