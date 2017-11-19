package main

import (
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/utils"
)

func SpecializedTestImpl(t *testing.T, nprocesses int, queue *module.NotificationQueue) {
	notifications := queue.Notifications()

	utils.AssertEquals(t, nprocesses, len(notifications))
	for _, notification := range notifications {
		o := notification.Content().(*object.DataObject)

		utils.AssertEquals(t, 5, len(o.Value.(string)), func(a interface{}, b interface{}) bool { return a.(int) <= b.(int) })
		t.Logf("> Data notified : {%#v} from '%s'", o.Value, o.Module)
	}
}

func PropertyFactoryImpl() *properties.Properties {
	return &properties.Properties{}
}

func TestModule(t *testing.T) {
	utils.Plugin.Test(t, ExportModule(), PropertyFactoryImpl, SpecializedTestImpl, 5, 300*time.Millisecond)
}

func TestBenchmarkModule(t *testing.T) {
	utils.Plugin.Benchmark(t, ExportModule(), PropertyFactoryImpl, 300*time.Millisecond)
}
