package example

import (
	"go.uber.org/zap"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/model"
	"fantastic-broccoli/constant"
	"fmt"
	"time"
)

type ModuleExample struct {
	state      int
	queue      *module.notificationQueue
	logger     *zap.Logger
	data       chan string
	end_runner chan bool
}

func (m *ModuleExample) Start(q *module.notificationQueue, l *zap.Logger) error {
	m.logger = l
	m.queue = q
	m.state = constant.STARTED

	l.Info("Module 'Example' started")
	return nil
}
func (m *ModuleExample) Configure(properties *model.Properties) error {
	m.logger.Info("Module 'Example' configured")
	m.state = constant.IDLE
	return nil
}
func (m *ModuleExample) Process() error {
	if m.state == constant.IDLE {
		m.logger.Info("session not started")
		return nil
	}

	var st string
	var end = true
	for ; end; {
		select {
		case s := <-m.data:
			st += s
		default:
			end = false
		}
	}
	m.logger.Info("process ended", zap.Int("nb_value", len(st)))

	return nil
}
func (m *ModuleExample) Stop() error {
	m.state = constant.STOPPED
	return nil
}

func (m *ModuleExample) StartSession() error {
	if m.state == constant.WORKING {
		return fmt.Errorf("previous session has not been ended")
	}

	m.data = make(chan string, 0xFF)
	m.end_runner = make(chan bool, 1)
	go func() {
		defer m.logger.Info("end goroutine")
		defer close(m.data)

		for {
			select {
			case <-m.end_runner:
				return
			default:
				m.data <- "|"
			}

			time.Sleep(time.Nanosecond)
		}
	}()

	m.state = constant.WORKING
	return nil
}
func (m *ModuleExample) StopSession() error {
	m.end_runner <- true

	close(m.end_runner)
	m.state = constant.IDLE
	return nil
}

func (m *ModuleExample) Name() string {
	return "ModuleExample"
}
func (m *ModuleExample) State() int {
	return m.state
}
