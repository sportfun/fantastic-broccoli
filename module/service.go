package module

import (
	. "fantastic-broccoli"
	"fantastic-broccoli/core"
	"plugin"
	"go.uber.org/zap"
)

type Service struct {
	modules   map[Name]Module
	sessionId int
	state     State

	messages      NotificationQueue
	notifications *core.NotificationQueue
	logger        *zap.Logger
}

func (m *Service) Start(q *core.NotificationQueue, l *zap.Logger) error {
	m.modules = map[Name]Module{}
	m.sessionId = -1
	m.state = STARTED

	m.messages = NotificationQueue{}
	m.notifications = q
	m.logger = l

	return nil
}

func (m *Service) Configure(props *Properties) error {
	for _, e := range props.Modules {
		p, err := plugin.Open(string(e.Path))
		if err != nil {
			// TODO: Error management (Plugin loading)
		}

		ex, err := p.Lookup("ExportModule")
		if err != nil {
			// TODO: Error management (Symbol loading)
		}

		mod := ex.(func() (Module))()
		m.modules[Name(e.Name)] = mod

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

func (m *Service) Name() Name {
	return MODULE_SERVICE
}

func (m *Service) State() State {
	return m.state
}
