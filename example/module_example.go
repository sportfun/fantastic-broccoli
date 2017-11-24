package main

import (
	"time"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type ModuleExample struct {
	logger        log.Logger
	notifications *module.NotificationQueue

	buffer    string
	data      chan string
	endRunner chan bool
	state     types.StateType
}

var (
	LogModuleStarted    = log.NewArgumentBinder("module 'Example' started")
	LogModuleConfigured = log.NewArgumentBinder("module 'Example' started")
)

func (m *ModuleExample) Start(q *module.NotificationQueue, l log.Logger) error {
	m.logger = l
	m.notifications = q
	m.state = constant.States.Started

	l.Info(LogModuleStarted)
	return nil
}

func (m *ModuleExample) Configure(properties *properties.Properties) error {
	m.logger.Info(LogModuleConfigured)
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
		return nil
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

func (m *ModuleExample) State() types.StateType {
	return m.state
}

func ExportModule() module.Module {
	return &ModuleExample{}
}
