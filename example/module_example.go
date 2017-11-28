package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type rpmGenerator struct {
	logger        log.Logger
	notifications *module.NotificationQueue
	state         types.StateType

	engine rpmEngine
	data   chan float64
	done   chan struct{}
}

const tick = 50 * time.Millisecond

var (
	debugModuleStarted    = log.NewArgumentBinder("module '%s' started")
	debugModuleConfigured = log.NewArgumentBinder("module '%s' configured")
	debugRpmCalculated    = log.NewArgumentBinder("new rpm calculated")
	debugSessionStarted   = log.NewArgumentBinder("session started")
	debugSessionStopped   = log.NewArgumentBinder("session stopped")
	debugModuleStopped    = log.NewArgumentBinder("module '%s' stopped")
)

func (m *rpmGenerator) isSet(a interface{}, name string) (error, bool) {
	if a != nil {
		return nil, true
	}

	m.state = constant.States.Panic
	return fmt.Errorf("%s is not set", name), false
}

func (m *rpmGenerator) Start(q *module.NotificationQueue, l log.Logger) error {
	if err, isSet := m.isSet(q, "notification queue"); !isSet {
		return err
	}
	if err, isSet := m.isSet(l, "logger"); !isSet {
		return err
	}

	m.logger = l
	m.notifications = q
	m.state = constant.States.Started
	m.engine.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	l.Debug(debugModuleStarted.Bind(m.Name()))
	return nil
}

func loadConfigurationItem(items map[string]interface{}, name string) (float64, error) {
	_, ok := items[name]
	if !ok {
		return 0, fmt.Errorf("invalid value of '%s' in configuration", name)
	}

	v, ok := items[name].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid value of '%s' in configuration", name)
	}

	return v, nil
}

func (m *rpmGenerator) Configure(properties properties.ModuleDefinition) error {
	if properties.Conf == nil {
		m.state = constant.States.Panic
		return fmt.Errorf("configuration needed for this module. RTFM")
	}

	items, ok := properties.Conf.(map[string]interface{})
	if !ok {
		m.state = constant.States.Panic
		return fmt.Errorf("valid configuration needed for this module. RTFM")
	}

	var err error
	for k, v := range map[string]*float64{
		"rpm.min":       &m.engine.min,
		"rpm.max":       &m.engine.max,
		"rpm.step":      &m.engine.step,
		"rpm.precision": &m.engine.precision,
	} {
		if *v, err = loadConfigurationItem(items, k); err != nil {
			m.state = constant.States.Panic
			return err
		}
	}

	m.logger.Debug(debugModuleConfigured.Bind(m.Name()))
	m.state = constant.States.Idle
	return nil
}

func (m *rpmGenerator) calcRpm() (float64, int) {
	rpm := 0.
	nvalue := 0

	for {
		select {
		case val := <-m.data:
			rpm += val
			nvalue++
		default:
			return rpm / float64(nvalue), nvalue
		}
	}
}

func (m *rpmGenerator) Process() error {
	if m.state != constant.States.Working || m.data == nil {
		return fmt.Errorf("session not started")
	}

	rpm, nvalue := m.calcRpm()
	m.logger.Debug(debugRpmCalculated.More("nb_value", nvalue).More("value", rpm))

	m.notifications.NotifyData(m.Name(), "%f", rpm)
	return nil
}

func (m *rpmGenerator) Stop() error {
	if m.state == constant.States.Working {
		m.StopSession()
	}

	m.logger.Debug(debugModuleStopped.Bind(m.Name()))
	m.state = constant.States.Stopped
	return nil
}

func (m *rpmGenerator) StartSession() error {
	if m.state == constant.States.Working || m.data != nil {
		m.StopSession()
		return fmt.Errorf("session already exist")
	}

	m.logger.Debug(debugSessionStarted)
	m.data, m.done = make(chan float64, 0xff), make(chan struct{}, 1)
	go func() {
		defer close(m.data)

		for {
			select {
			case <-m.done:
				return
			default:
				m.data <- m.engine.NewValue()
			}

			time.Sleep(tick)
		}
	}()

	m.state = constant.States.Working
	return nil
}

func (m *rpmGenerator) StopSession() error {
	if m.state != constant.States.Working || m.done == nil {
		m.state = constant.States.Idle
		return fmt.Errorf("session not started")
	}

	close(m.done)
	m.done = nil
	m.data = nil

	m.logger.Debug(debugSessionStopped)
	m.state = constant.States.Idle
	return nil
}

func (m *rpmGenerator) Name() string {
	return "RPM Generator"
}

func (m *rpmGenerator) State() types.StateType {
	return m.state
}
