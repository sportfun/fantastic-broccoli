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

type rpmEngine struct {
	min       float64
	max       float64
	step      float64
	precision float64

	rand    *rand.Rand
	lastval int
}

type rpmGenerator struct {
	logger        log.Logger
	notifications *module.NotificationQueue
	state         types.StateType

	engine rpmEngine
	data   chan float64
	done   chan bool
}

const tick = 50 * time.Millisecond

var (
	debugModuleStarted    = log.NewArgumentBinder("module '%s' started")
	debugModuleConfigured = log.NewArgumentBinder("module '%s' configured")
	debugRpmCalculated    = log.NewArgumentBinder("new rpm calculated")
)

// - Start

func (m *rpmGenerator) Start(q *module.NotificationQueue, l log.Logger) error {
	m.logger = l
	m.notifications = q
	m.state = constant.States.Started
	m.engine.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	l.Debug(debugModuleStarted.Bind(m.Name()))
	return nil
}

// - Configure

func loadConfItem(items map[string]interface{}, itemName string) (float64, error) {
	_, ok := items[itemName]
	if !ok {
		return 0, fmt.Errorf("invalid value of '%s' in configuration", itemName)
	}

	v, ok := items[itemName].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid value of '%s' in configuration", itemName)
	}

	return v, nil
}

func (m *rpmGenerator) Configure(properties properties.ModuleDefinition) error {
	if m.state == constant.States.Stopped {
		return nil
	}

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
		if *v, err = loadConfItem(items, k); err != nil {
			m.state = constant.States.Panic
			return err
		}
	}

	m.logger.Debug(debugModuleConfigured.Bind(m.Name()))
	m.state = constant.States.Idle
	return nil
}

// - Process

func (m *rpmGenerator) Process() error {
	if m.state == constant.States.Stopped {
		return nil
	}

	if m.state != constant.States.Working {
		m.state = constant.States.Failed
		return fmt.Errorf("session not started")
	}

	var rpm float64
	var nvalue int

aggregator:
	for {
		select {
		case val := <-m.data:
			rpm += val
			nvalue++
		default:
			break aggregator
		}
	}

	rpm = rpm / float64(nvalue)
	m.logger.Debug(debugRpmCalculated.More("nb_value", nvalue).More("value", rpm))

	m.notifications.NotifyData(m.Name(), "%f", rpm)
	return nil
}

// - Stop

func (m *rpmGenerator) Stop() error {
	if m.state == constant.States.Working {
		m.StopSession()
	}

	m.state = constant.States.Stopped
	return nil
}

// - Session

func (m *rpmGenerator) StartSession() error {
	if m.state == constant.States.Stopped {
		return nil
	}

	if m.state == constant.States.Working {
		m.StopSession()
		return fmt.Errorf("session already exist")
	}

	m.data = make(chan float64, 0xff)
	m.done = make(chan bool, 1)

	// TODO : Don't create multi goroutine !!! FIX IT
	go func() {
		defer log.NewLogger.Dev(nil).Debug(log.NewArgumentBinder("exit goroutine"))
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
	if m.state == constant.States.Stopped {
		return nil
	}

	if m.state != constant.States.Working {
		m.state = constant.States.Idle
		return fmt.Errorf("session not started")
	}

	m.done <- true
	close(m.done)

	m.state = constant.States.Idle
	return nil
}

// - Properties

func (m *rpmGenerator) Name() string {
	return "RPM Generator"
}

func (m *rpmGenerator) State() types.StateType {
	return m.state
}

// - Engine

func (e *rpmEngine) NewValue() float64 {
	rpm := int(e.rand.Float64() * (e.max - e.min) * e.precision)

	if e.lastval == 0 {
		e.lastval = rpm
		return e.min + float64(rpm/int(e.precision))
	}

	e.lastval = (int(e.lastval) - rpm) % int(e.step)

	return e.min + float64(e.lastval/int(e.precision))
}

// - Exporter

func ExportModule() module.Module {
	return &rpmGenerator{}
}
