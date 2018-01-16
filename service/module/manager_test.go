package module

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/config"
	. "github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/module"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
	"github.com/sportfun/gakisitor/utils/module_test"
	"testing"
)

type tModule struct {
	name   string
	logger log.Logger
	queue  *module.NotificationQueue
	state  byte

	InSession    bool
	ForceFailure string
	NbProcess    int
}

var none = fmt.Errorf("")

func TestManager_TestModuleValidity(t *testing.T) {
	environment := module_test.NewEnvironment(func(interface{}) *config.ModuleDefinition { return nil }, nil, nil, 0)
	module_test.Test(t, &tModule{}, environment)
}

func TestManager_Start(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	manager := Manager{}

	Expect(manager.Start(queue, logger)).Should(Succeed())
	Expect(manager.State()).Should(Equal(StartedState))
}

func TestManager_Process(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	manager := Manager{}
	mod := tModule{name: "test"}

	Expect(manager.Start(queue, logger)).Should(Succeed())
	manager.modules = moduleContainer{"test": &mod}
	mod.Start(manager.messages, manager.logger)
	manager.state = IdleState

	Expect(mod.InSession).Should(BeFalse())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(BeEmpty())

	manager.notifications.Notify(notification.NewNotification(NetworkServiceEntity, ModuleServiceEntity, object.NewCommandObject(StartSessionCmd)))
	Expect(manager.Process()).Should(Succeed())
	Expect(mod.InSession).Should(BeTrue())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(ConsistOf(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewDataObject("test", "1"))))

	Expect(manager.Process()).Should(Succeed())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(ConsistOf(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewDataObject("test", "2"))))

	mod.ForceFailure = WarningLevel
	Expect(manager.Process()).Should(Succeed())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(ConsistOf(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewErrorObject("manager_test.go:146", fmt.Errorf("forced failure")))))

	mod.ForceFailure = FatalLevel
	Expect(manager.Process()).Should(Succeed())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(ConsistOf(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewErrorObject("manager_test.go:146", fmt.Errorf("forced failure")))))
	Expect(mod.State()).Should(Equal(StoppedState))

	manager.notifications.Notify(notification.NewNotification(NetworkServiceEntity, ModuleServiceEntity, object.NewCommandObject(StartSessionCmd)))
	Expect(manager.Process()).Should(Succeed())
	Expect(mod.InSession).Should(BeTrue())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(ConsistOf(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewDataObject("test", "3"))))

	manager.notifications.Notify(notification.NewNotification(NetworkServiceEntity, ModuleServiceEntity, object.NewCommandObject(EndSessionCmd)))
	Expect(manager.Process()).Should(Succeed())
	Expect(mod.InSession).Should(BeFalse())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(BeEmpty())
}

func TestManager_Stop(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	manager := Manager{}
	mod := tModule{name: "test"}

	Expect(manager.Start(queue, logger)).Should(Succeed())
	manager.modules = moduleContainer{"test": &mod}
	mod.Start(manager.messages, manager.logger)
	manager.state = IdleState

	Expect(mod.InSession).Should(BeFalse())
	Expect(manager.notifications.Notifications(NetworkServiceEntity)).Should(BeEmpty())

	manager.notifications.Notify(notification.NewNotification(NetworkServiceEntity, ModuleServiceEntity, object.NewCommandObject(StartSessionCmd)))
	Expect(manager.Process()).Should(Succeed())
	Expect(mod.InSession).Should(BeTrue())

	Expect(manager.Stop()).Should(Succeed())
	Expect(mod.InSession).Should(BeFalse())
	Expect(mod.State()).Should(Equal(StoppedState))
}

func (m *tModule) Start(q *module.NotificationQueue, l log.Logger) error {
	if q == nil || l == nil {
		m.state = PanicState
		return none
	}

	m.queue = q
	m.logger = l
	m.state = StartedState

	return nil
}
func (m *tModule) Configure(*config.ModuleDefinition) error { m.state = IdleState; return nil }
func (m *tModule) StartSession() error {
	if m.state == WorkingState {
		m.StopSession()
		return fmt.Errorf("")
	}
	m.state = WorkingState
	m.InSession = true
	return nil
}
func (m *tModule) Process() error {
	if m.state != WorkingState {
		return nil
	}

	if m.ForceFailure != "" {
		m.queue.NotifyError(m, m.ForceFailure, "forced failure")
		m.ForceFailure = ""
		return nil
	}

	m.NbProcess++
	m.queue.NotifyData(m.Name(), "%d", m.NbProcess)
	return nil
}
func (m *tModule) StopSession() error {
	defer func() { m.state = IdleState }()
	if m.state != WorkingState {
		return none
	}

	m.InSession = false
	return nil
}
func (m *tModule) Stop() error { m.StopSession(); m.state = StoppedState; return nil }

func (m *tModule) Name() string { return m.name }
func (m *tModule) State() byte  { return m.state }
