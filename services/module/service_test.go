package module

import (
	"fantastic-broccoli/constant"
	"fantastic-broccoli/common/types/module"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/example"
	"fantastic-broccoli/properties"
	"fantastic-broccoli/utils"
	"strconv"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	s := Service{}
	ms := []module.Module{&example.ModuleExample{}, &example.ModuleExample{}}
	p := properties.Properties{}
	q := service.NewNotificationQueue()
	b := notification.NewBuilder().From(constant.EntityNames.Services.Network).To(constant.EntityNames.Services.Module)
	l := utils.Default.Logger()

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
	time.Sleep(250 * time.Millisecond)

	q.Notify(b.With(*object.NewNetworkObject(constant.NetCommand.StartSession)).Build())
	// Invalid: Session already started ... Only in logs
	q.Notify(b.With(*object.NewNetworkObject(constant.NetCommand.StartSession)).Build())

	s.Process()
	s.Process()
	// 2 process call = (2 * nb_module) notifications for the network
	utils.AssertEquals(t, 2*len(ms), len(q.Notifications(constant.EntityNames.Services.Network)))

	time.Sleep(3 * time.Second)
	s.Process()
	d := s.notifications.Notifications(constant.EntityNames.Services.Network)
	utils.AssertEquals(t, len(ms), len(d))
	o := d[0].Content().(object.DataObject)
	utils.AssertEquals(t, ms[0].Name(), o.Module)
	utils.AssertEquals(t, 10, len(o.Value.(string)))

	time.Sleep(250 * time.Millisecond)
	q.Notify(b.With(*object.NewNetworkObject(constant.NetCommand.EndSession)).Build())
	// Invalid: Session not started
	s.Process()

	s.Stop()

	// Clean all threads
	time.Sleep(500 * time.Millisecond)
}
