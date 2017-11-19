package main

import (
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/properties"
	"go.uber.org/zap"
	"time"
)

type ModuleExample struct {
	state         byte
	notifications *module.NotificationQueue
	logger        *zap.Logger
	data          chan string
	endRunner     chan bool
	buffer        string
}

func (m *ModuleExample) Start(q *module.NotificationQueue, l *zap.Logger) error {
	m.logger = l
	m.notifications = q
	m.state = constant.States.Started

	l.Info("module 'Example' started")
	return nil
}

func (m *ModuleExample) Configure(properties *properties.Properties) error {
	m.logger.Info("module 'Example' configured")
	m.state = constant.States.Idle
	return nil
}

func (m *ModuleExample) Process() error {
	if m.state == constant.States.Idle {
		// Session not started
		return nil
	}

aggregator:
	for {
		select {
		case val := <-m.data:
			m.buffer += val
		default:
			break aggregator
		}
	}

	if len(m.buffer) > 5 {
		m.notifications.NotifyData(m.Name(), m.buffer)
		m.buffer = ""
	}

	return nil
}

func (m *ModuleExample) Stop() error {
	if m.state == constant.States.Working {
		m.StopSession()
	}

	m.state = constant.States.Stopped
	return nil
}

func (m *ModuleExample) StartSession() error {
	if m.state == constant.States.Working {
		// Previous session has not been ended
		m.StopSession()
	}

	// Chan where we buffer 0x9 char
	m.data = make(chan string, 0x9)
	m.endRunner = make(chan bool, 1)
	go func() {
		defer close(m.data)

		for {
			select {
			case <-m.endRunner:
				return
			default:
				m.data <- "|"
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	m.state = constant.States.Working
	return nil
}

func (m *ModuleExample) StopSession() error {
	if m.state == constant.States.Idle {
		// Session already stopped
		return nil
	}

	m.endRunner <- true

	close(m.endRunner)
	m.state = constant.States.Idle
	return nil
}

func (m *ModuleExample) Name() string {
	return "ModuleExample"
}

func (m *ModuleExample) State() byte {
	return m.state
}

func ExportModule() module.Module {
	return &ModuleExample{}
}
