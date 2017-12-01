package module

import (
	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/module"
	"github.com/xunleii/fantastic-broccoli/notification"
	"github.com/xunleii/fantastic-broccoli/notification/object"
	"github.com/xunleii/fantastic-broccoli/service"
	"github.com/xunleii/fantastic-broccoli/utils"
	"strconv"
	"testing"
)

type ModuleImpl struct {
	n *module.NotificationQueue
	s byte
}

func (m *ModuleImpl) Start(queue *module.NotificationQueue, logger log.Logger) error {
	m.n = queue
	m.s = env.IdleState
	return nil
}

func (*ModuleImpl) Configure(*config.ModuleDefinition) error {
	return nil
}

func (m *ModuleImpl) Process() error {
	if m.s == env.WorkingState {
		m.n.NotifyData(m.Name(), "#####")
	}
	return nil
}

func (m *ModuleImpl) Stop() error {
	m.StopSession()
	return nil
}

func (m *ModuleImpl) StartSession() error {
	m.s = env.WorkingState
	return nil
}

func (m *ModuleImpl) StopSession() error {
	m.s = env.IdleState
	return nil
}

func (*ModuleImpl) Name() string {
	return ""
}

func (m *ModuleImpl) State() byte {
	return m.s
}

func TestService(t *testing.T) {
	s := Service{}
	ms := []module.Module{&ModuleImpl{}, &ModuleImpl{}}
	p := &config.ModuleDefinition{}
	q := service.NewNotificationQueue()
	b := notification.NewBuilder().From(env.NetworkServiceEntity).To(env.ModuleServiceEntity)
	l := log.NewDevelopment()

	s.Start(q, l)
	// Manually configuration
	for i, m := range ms {
		m.Start(s.messages, l)
		m.Configure(p)
		s.modules[m.Name()+strconv.Itoa(i)] = m
	}
	s.state = env.IdleState

	// Invalid: Session not started
	s.Process()

	q.Notify(b.With(object.NewCommandObject(env.StartSessionCmd)).Build())

	s.Process()
	s.Process()

	// 2 process call = (2 * nb_module) notifications for the network
	utils.AssertEquals(t, 2*len(ms), len(q.Notifications(env.NetworkServiceEntity)))
	s.Process()

	d := s.notifications.Notifications(env.NetworkServiceEntity)
	utils.AssertEquals(t, len(ms), len(d))
	o := d[0].Content().(*object.DataObject)
	utils.AssertEquals(t, ms[0].Name(), o.Module)
	utils.AssertEquals(t, "#####", o.Value.(string))

	q.Notify(b.With(*object.NewCommandObject(env.EndSessionCmd)).Build())

	// Invalid: Session not started
	s.Process()

	s.Stop()
}
