package example

import (
	"testing"
	"fantastic-broccoli/common/types/module"
	"go.uber.org/zap"
	"fantastic-broccoli/model"
	"time"
)

func TestModuleExample(t *testing.T) {
	m := ModuleExample{}
	q := module.notificationQueue{}
	l, _ := zap.NewProduction()
	p := model.Properties{}

	m.Start(&q, l)
	m.Configure(&p)

	m.Process()

	m.StartSession()
	m.Process()
	time.Sleep(1 * time.Second)
	m.Process()
	m.Process()
	m.StartSession()
	m.StopSession()
	m.Process()

	m.Stop()
	time.Sleep(250 * time.Millisecond)
}