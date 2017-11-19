package main

import (
	"testing"
	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"go.uber.org/zap"
	"github.com/xunleii/fantastic-broccoli/properties"
	"time"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/utils"
	"github.com/xunleii/fantastic-broccoli/constant"
)

func ImplSpecTest(t *testing.T, nProcess int, logger *zap.Logger, queue *module.NotificationQueue) {
	notifications := queue.Notifications()

	utils.AssertEquals(t, nProcess, len(notifications))
	for _, notification := range notifications {
		o := notification.Content().(object.DataObject)

		utils.AssertEquals(t, 5, len(o.Value.(string)))
		utils.AssertEquals(t, "|||||", o.Value.(string))
		logger.Info("data notified", zap.String("from", o.Module), zap.String("value", o.Value.(string)))
	}
}

func ImplCustomProperty() properties.Properties {
	return properties.Properties{}
}

func StartModuleTest(t *testing.T, m module.Module) (module.Module, *module.NotificationQueue, *zap.Logger) {
	q := module.NewNotificationQueue()
	l, _ := zap.NewProduction()

	l.Info("Start module test", zap.String("moduleName", m.Name()))
	if err := m.Start(q, l); err != nil {
		l.Fatal(err.Error())
	}
	utils.AssertEquals(t, constant.States.Started, m.State())
	if err := m.Configure(&properties.Properties{}); err != nil {
		l.Fatal(err.Error())
	}

	return m, q, l
}

func TestModuleExample(t *testing.T) {
	m, q, l := StartModuleTest(t, &ModuleExample{})

	// Error verification (Process without session)
	// TODO: Explicit error knowable by the manager
	m.Process()

	m.StartSession()
	// Error verification (Session already exist)
	// TODO: Explicit error knowable by the manager
	m.StartSession()

	for i := 0; i < 0xF; i++ {
		time.Sleep(time.Millisecond * 500)
		m.Process()
	}

	m.StopSession()
	// Error verification (Process without session)
	// TODO: Explicit error knowable by the manager
	m.Process()

	ImplSpecTest(t, 0xF, l, q)
	time.Sleep(time.Second)
}