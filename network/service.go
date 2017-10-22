package network

import (
	. "fantastic-broccoli"
	"fantastic-broccoli/core"
	"go.uber.org/zap"
)

type Service struct {
	state State

	notifications *core.NotificationQueue
	logger        *zap.Logger
}

func (m *Service) Start(q *core.NotificationQueue, l *zap.Logger) error {
	m.state = STARTED

	m.notifications = q
	m.logger = l

	return nil
}

func (m *Service) Configure(props *Properties) error {
	return nil
}

func (m *Service) Process() error {
	return nil
}

func (m *Service) Stop() error {
	return nil
}

func (m *Service) Name() Name {
	return MODULE_SERVICE
}

func (m *Service) State() State {
	return m.state
}
