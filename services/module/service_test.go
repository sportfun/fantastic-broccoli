package module

import (
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/properties"
	def "github.com/xunleii/fantastic-broccoli/utils/default"
	"strconv"
	"testing"
	"github.com/xunleii/fantastic-broccoli/utils"
	"go.uber.org/zap"
)

type ModuleImpl struct {
	n *module.NotificationQueue
	s byte
}

func (m *ModuleImpl) Start(queue *module.NotificationQueue, logger *zap.Logger) error {
	m.n = queue
	m.s = constant.States.Idle
	return nil
}

func (*ModuleImpl) Configure(properties *properties.Properties) error {
	return nil
}

func (m *ModuleImpl) Process() error {
	if m.s == constant.States.Working {
		m.n.NotifyData(m.Name(), "#####")
	}
	return nil
}

func (m *ModuleImpl) Stop() error {
	m.StopSession()
	return nil
}

func (m *ModuleImpl) StartSession() error {
	m.s = constant.States.Working
	return nil
}

func (m *ModuleImpl) StopSession() error {
	m.s = constant.States.Idle
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
	p := properties.Properties{}
	q := service.NewNotificationQueue()
	b := notification.NewBuilder().From(constant.EntityNames.Services.Network).To(constant.EntityNames.Services.Module)
	l := def.Logger()

	s.Start(q, l)
	// Manually configuration
	for i, m := range ms {
		m.Start(s.messages, l)
		m.Configure(&p)
		s.modules[m.Name()+strconv.Itoa(i)] = m
	}
	s.state = constant.States.Idle

	// Invalid: Session not started
	s.Process()

	q.Notify(b.With(*object.NewCommandObject(constant.NetCommand.StartSession)).Build())

	s.Process()
	s.Process()

	// 2 process call = (2 * nb_module) notifications for the network
	utils.AssertEquals(t, 2*len(ms), len(q.Notifications(constant.EntityNames.Services.Network)))
	s.Process()

	d := s.notifications.Notifications(constant.EntityNames.Services.Network)
	utils.AssertEquals(t, len(ms), len(d))
	o := d[0].Content().(*object.DataObject)
	utils.AssertEquals(t, ms[0].Name(), o.Module)
	utils.AssertEquals(t, "#####", o.Value.(string))

	q.Notify(b.With(*object.NewCommandObject(constant.NetCommand.EndSession)).Build())

	// Invalid: Session not started
	s.Process()

	s.Stop()
}
