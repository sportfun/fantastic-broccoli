package example

import (
	"go.uber.org/zap"
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/model"
	"fantastic-broccoli/const"
)

type ModuleExample struct {
	state  types.State
	queue  *module.NotificationQueue
	logger *zap.Logger
}

func (m *ModuleExample) Start(q *module.NotificationQueue, l *zap.Logger) error {
	m.logger = l
	m.queue = q
	m.state = _const.STARTED

	l.Info("Module 'Example' started")
	return nil
}
func (m *ModuleExample) Configure(properties *model.Properties) error {
	m.logger.Info("Module 'Example' configured")
	return nil
}
func (m *ModuleExample) Process() error {
	return nil
}
func (m *ModuleExample) Stop() error {
	m.state = _const.STOPPED
	return nil
}

func (m *ModuleExample) StartSession() error {
	return nil
}
func (m *ModuleExample) StopSession() error {
	return nil
}

func (m *ModuleExample) Name() types.Name {
	return "ModuleExample"
}
func (m *ModuleExample) State() types.State {
	return m.state
}
