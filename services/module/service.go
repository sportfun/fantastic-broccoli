package module

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/const"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"plugin"
)

type Service struct {
	modules   map[types.Name]module.Module
	sessionId int
	state     types.State

	messages      module.NotificationQueue
	notifications *service.NotificationQueue
	logger        *zap.Logger
}

func (m *Service) Start(q *service.NotificationQueue, l *zap.Logger) error {
	m.modules = map[types.Name]module.Module{}
	m.sessionId = -1
	m.state = _const.STARTED

	m.messages = module.NotificationQueue{}
	m.notifications = q
	m.logger = l

	return nil
}

func (m *Service) Configure(props *model.Properties) error {
	for _, e := range props.Modules {
		p, err := plugin.Open(string(e.Path))
		if err != nil {
			// TODO: Error management (Plugin loading)
		}

		ex, err := p.Lookup("ExportModule")
		if err != nil {
			// TODO: Error management (Symbol loading)
		}

		mod := ex.(func() (module.Module))()
		m.modules[e.Name] = mod

		if err = mod.Start(&m.messages, m.logger); err != nil {
			// TODO: Error management (Module starting)
		}
		if err = mod.Configure(props); err != nil {
			// TODO: Error management (Module configuration)
		}
	}

	if len(m.modules) == 0 {
		// TODO: Error management (No module charged)
	}
	return nil
}

func (m *Service) Process() error {
	for _, e := range m.notifications.Notifications(m.Name()) {
		// TODO: Notification interpretation -> New session (Network) | End session (Network)
		e.To()
	}

	for _, e := range m.modules {
		err := e.Process()
		// TODO: Error management (Module processing)
		err.Error()
	}

	for _, e := range m.messages.NotificationsError() {
		// TODO: Write notification for Network (if error is FATAL, call system)
		e.To()
	}

	for _, e := range m.messages.NotificationsData() {
		// TODO: Write notification for Network
		e.To()
	}

	return nil
}

func (m *Service) Stop() error {
	for _, e := range m.modules {
		err := e.Stop()
		err.Error()
		// TODO: Error management
	}
	return nil
}

func (m *Service) Name() types.Name {
	return "Module"
}

func (m *Service) State() types.State {
	return m.state
}
