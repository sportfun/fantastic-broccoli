package module

import (
	"fantastic-broccoli/common/types/module"

	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"plugin"
	"fmt"
	"fantastic-broccoli/common/types/notification/object"
)

var netBuilder = new(notification.Builder).
	From(constant.ModuleService).
	To(constant.NetworkService)

type Service struct {
	modules   map[string]module.Module
	sessionId int
	state     int

	messages      module.notificationQueue
	notifications *service.NotificationQueue
	logger        *zap.Logger
}

func (s *Service) Start(q *service.NotificationQueue, l *zap.Logger) error {
	s.modules = map[string]module.Module{}
	s.sessionId = -1
	s.state = constant.STARTED

	s.messages = module.notificationQueue{}
	s.notifications = q
	s.logger = l

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

		s.errorHandler(ModuleStarting, mod.Start(&s.messages, s.logger), e.Name)
		s.errorHandler(ModuleConfiguration, mod.Configure(props), e.Name)
	}

	if len(s.modules) == 0 {
		s.errorHandler(NoModule, fmt.Errorf("no module charged"))
	}

	s.state = constant.IDLE
	return nil
}

func (s *Service) Process() error {
	s.state = constant.WORKING
	for _, n := range s.notifications.Notifications(s.Name()) {
		s.notificationHandler(n)
	}

	for _, e := range s.modules {
		s.errorHandler(ModuleProcess, e.Process())
	}

	for _, e := range s.messages.NotificationsError() {
		m := object.NewErrorObject(e.From(), fmt.Errorf(e.Content().(string)))
		s.notifications.Notify(netBuilder.With(m).Build())

		if e.Content().(module.ErrorObject).ErrorLevel == constant.FATAL {
			// TODO: Stop module
			s.errorHandler(ModuleStop, s.modules[e.From()].Stop(), e.From())
		}
	}

	for _, e := range s.messages.NotificationsData() {
		m := object.NewDataObject(e.From(), e.Content())
		s.notifications.Notify(netBuilder.With(m).Build())
	}

	s.state = constant.IDLE
	return nil
}

func (s *Service) Stop() error {
	for n, e := range s.modules {
		s.errorHandler(ModuleStop, e.Stop(), n)
	}
	return nil
}

func (s *Service) Name() string {
	return constant.ModuleService
}

func (s *Service) State() int {
	return s.state
}
