package example

import (
	"testing"
	"fantastic-broccoli/common/types/module"
	"go.uber.org/zap"
	"fantastic-broccoli/model"
	"time"
	"fantastic-broccoli/common/types/notification/object"
)

func TestModuleExample(t *testing.T) {
	m := ModuleExample{}
	q := module.NewNotificationQueue()
	l, _ := zap.NewProduction()
	p := model.Properties{}

	m.Start(q, l)
	m.Configure(&p)

	m.Process()

	m.StartSession()
	m.Process()
	time.Sleep(time.Second)
	m.Process()
	time.Sleep(time.Second)
	m.Process()
	m.Process()

	m.StartSession()
	m.StopSession()
	m.Process()

	m.Stop()
	time.Sleep(250 * time.Millisecond)

	for _, d := range q.Notifications() {
		o := d.Content().(object.DataObject)
		l.Info("data notified", zap.String("from", o.Module()), zap.String("value", o.Value().(string)))
	}
}