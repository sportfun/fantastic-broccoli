package module

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/const"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"plugin"
	"errors"
)

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
			s.errorHandler(PluginLoading, err)
			continue
		}

		ex, err := p.Lookup("ExportModule")
		if err != nil {
			s.errorHandler(SymbolLoading, err)
			continue
		}

		mod := ex.(func() (module.Module))()
		s.modules[e.Name] = mod

		s.errorHandler(ModuleStarting, mod.Start(&s.messages, s.logger))
		s.errorHandler(ModuleConfiguration, mod.Configure(props))
	}

	if len(s.modules) == 0 {
		s.errorHandler(NoModule, errors.New("no module charged"))
	}
	return nil
}

func (s *Service) Process() error {
	for _, e := range s.notifications.Notifications(s.Name()) {
		// TODO: Notification interpretation -> New session (Network) | End session (Network)
		e.To()
	}

	for n, e := range s.modules {
		err := e.Process()
		// TODO: Error management (Module processing)
		err.Error()
	}

	for _, e := range s.messages.NotificationsError() {
		// TODO: Write notification for Network (if error is FATAL, call system)
		e.To()
	}

	for _, e := range s.messages.NotificationsData() {
		// TODO: Write notification for Network
		e.To()
	}

	return nil
}

func (s *Service) Stop() error {
	for _, e := range s.modules {
		err := e.Stop()
		err.Error()
		// TODO: Error management
	}
	return nil
}

func (s *Service) Name() types.Name {
	return "Module"
}

func (s *Service) State() types.State {
	return s.state
}
