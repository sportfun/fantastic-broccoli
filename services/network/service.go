package network

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/const"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
)

type Service struct {
	state types.State

	notifications *service.NotificationQueue
	logger        *zap.Logger
}

func (m *Service) Start(q *service.NotificationQueue, l *zap.Logger) error {
	m.state = _const.STARTED

	m.notifications = q
	m.logger = l

	return nil
}

func (m *Service) Configure(props *model.Properties) error {
	return nil
}

func (m *Service) Process() error {
	return nil
}

func (m *Service) Stop() error {
	return nil
}

func (m *Service) Name() types.Name {
	return "Network"
}

func (m *Service) State() types.State {
	return m.state
}
