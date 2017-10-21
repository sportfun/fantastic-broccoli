package module

import (
	"fantastic-broccoli/model"
	"fantastic-broccoli/core"
	"plugin"
)

type Manager struct {
	modules   map[Name]Module
	sessionId int
	state     model.State

	notifications core.NotificationQueue
	messages      NotificationQueue
}

func (m *Manager) Start(q core.NotificationQueue) error {
	m.modules = map[Name]Module{}
	m.sessionId = -1
	m.state = model.STARTED
	m.messages = NotificationQueue{}
	m.notifications = q

	return nil
}

func (m *Manager) Configure(props model.Properties) error {
	for _, e := range props.Modules {
		p, err := plugin.Open(string(e.Path))
		if err != nil {
			// TODO: Error management
		}

		ex, err := p.Lookup("Export")
		if err != nil {
			// TODO: Error management
		}

		m.modules[Name(e.Name)] = ex.(func() (Module))()

		if err = m.modules[Name(e.Name)].Start(m.messages); err != nil {
			// TODO: Error management
		}
		if err = m.modules[Name(e.Name)].Configure(props); err != nil {
			// TODO: Error management
		}
	}

	if len(m.modules) == 0 {
		// TODO: Error management
	}
	return nil
}

func (m *Manager) Process() error {
	n := m.notifications.Notifications(m.Name())
	for _, e := range n {
		// Interpret notification
		e.To()
	}

	for _, e := range m.modules {
		e.Process()
	}

	for _, e := range m.messages.NotificationsError() {
		//Write notification for Network
		e.To()
	}

	for _, e := range m.messages.NotificationsMetric() {
		//Write notification for Network
		e.To()
	}

	return nil
}

func (m *Manager) Stop() error {
	for _, e := range m.modules {
		err := e.Stop()
		err.Error()
		// TODO: Error management
	}
	return nil
}

func (m *Manager) Name() core.Name {
	return model.MODULE_MANAGER
}

func (m *Manager) State() model.State {
	return m.state
}
