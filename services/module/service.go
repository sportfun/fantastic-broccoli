package module

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/common/types/network"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/const"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"plugin"
	"errors"
)

var netBuilder *notification.Builder = notification.Builder{}.
	From(_const.ModuleService).
	To(_const.NetworkService)

type Service struct {
	modules   map[types.Name]module.Module
	sessionId int
	state     types.State

	messages      module.NotificationQueue
	notifications *service.NotificationQueue
	logger        *zap.Logger
}

func (s *Service) Start(q *service.NotificationQueue, l *zap.Logger) error {
	s.modules = map[types.Name]module.Module{}
	s.sessionId = -1
	s.state = _const.STARTED

	s.messages = module.NotificationQueue{}
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
		s.errorHandler(NoModule, errors.New("no module charged"))
	}
	return nil
}

func (s *Service) Process() error {
	for _, n := range s.notifications.Notifications(s.Name()) {
		s.notificationHandler(n)
	}

	for _, e := range s.modules {
		s.errorHandler(ModuleProcess, e.Process())
	}

	for _, e := range s.messages.NotificationsError() {
		m := network.NewMessage("error").
			AddArgument(string(e.From())).
			AddArgument(e.Content().(string))
		s.notifications.Notify(netBuilder.With(m).Build())

		if e.Content().(module.ErrorObject).ErrorLevel == _const.FATAL {
			s.errorHandler(ModuleStop, s.modules[e.From()].Stop(), e.From())
		}
	}

	for _, e := range s.messages.NotificationsData() {
		m := network.NewMessage("data").
			AddArgument(e.Content().(string))
		s.notifications.Notify(netBuilder.With(m).Build())
	}

	return nil
}

func (s *Service) Stop() error {
	for n, e := range s.modules {
		s.errorHandler(ModuleStop, e.Stop(), n)
	}
	return nil
}

func (s *Service) Name() types.Name {
	return _const.ModuleService
}

func (s *Service) State() types.State {
	return s.state
}
