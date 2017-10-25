package example

import (
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"time"
)

type ModuleExample struct {
	state         int
	notifications *module.NotificationQueue
	logger        *zap.Logger
	data          chan string
	endRunner     chan bool
}

func (m *ModuleExample) Start(q *module.NotificationQueue, l *zap.Logger) error {
	m.logger = l
	m.notifications = q
	m.state = constant.Started

	l.Info("module 'Example' started")
	return nil
}
func (m *ModuleExample) Configure(properties *model.Properties) error {
	m.logger.Info("module 'Example' configured")
	m.state = constant.Idle
	return nil
}
func (m *ModuleExample) Process() error {
	if m.state == constant.Idle {
		m.logger.Error("session not started")
		return nil
	}

	value := ""
aggregator:
	for {
		select {
		case s := <-m.data:
			value += s
		default:
			break aggregator
		}
	}
	m.logger.Info("process ended", zap.Int("nb_value", len(value)))
	m.notifications.NotifyData(m.Name(), value)

	return nil
}
func (m *ModuleExample) Stop() error {
	if m.state == constant.Working {
		m.StopSession()
	}

	m.state = constant.Stopped
	return nil
}

func (m *ModuleExample) StartSession() error {
	if m.state == constant.Working {
		m.logger.Error("previous session has not been ended")
		return nil
	}
	m.logger.Info("start new session")

	m.data = make(chan string, 0x9)
	m.endRunner = make(chan bool, 1)
	go func() {
		defer m.logger.Info("end goroutine")
		defer close(m.data)

		for {
			select {
			case <-m.endRunner:
				return
			default:
				m.data <- "|"
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	m.state = constant.Working
	return nil
}
func (m *ModuleExample) StopSession() error {
	m.logger.Info("end session")
	m.endRunner <- true

	close(m.endRunner)
	m.state = constant.Idle
	return nil
}

func (m *ModuleExample) Name() string {
	return "ModuleExample"
}
func (m *ModuleExample) State() int {
	return m.state
}
