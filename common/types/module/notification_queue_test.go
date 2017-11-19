package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/utils"
	"fmt"
	"path"
	"runtime"
	"testing"
)

// -- Test

func TestNotificationQueueNotifyError(t *testing.T) {
	q := NewNotificationQueue()

	for i := 0; i < 0xFF; i++ {
		q.NotifyError(constant.Fatal, "Error message %s %s", "can be", "formatted")
	}
	_, caller, line, _ := runtime.Caller(0)
	origin := fmt.Sprintf("%s:%d", path.Base(caller), line-2)

	utils.AssertEquals(t, 0xFF, len(q.notifications))
	errors := q.Notifications()
	utils.AssertEquals(t, 0xFF, len(errors))
	utils.AssertEquals(t, 0, len(q.notifications))

	o := errors[0].Content().(ErrorObject)
	utils.AssertEquals(t, origin, o.Origin)
	utils.AssertEquals(t, "Error message can be formatted", o.Reason)
	utils.AssertEquals(t, constant.Fatal, o.ErrorLevel())
}

func TestNotificationQueueNotifyData(t *testing.T) {
	q := NewNotificationQueue()

	for i := 0; i < 0xFF; i++ {
		q.NotifyData("ModuleName", "%d RPM", 1000)
	}

	utils.AssertEquals(t, 0xFF, len(q.notifications))
	data := q.Notifications()
	utils.AssertEquals(t, 0xFF, len(data))
	utils.AssertEquals(t, 0, len(q.notifications))

	o := data[0].Content().(object.DataObject)
	utils.AssertEquals(t, "ModuleName", o.Module)
	utils.AssertEquals(t, "1000 RPM", o.Value)
}

// -- Benchmark

func BenchmarkNotificationQueueNotifyData(b *testing.B) {
	q := NewNotificationQueue()
	b.ResetTimer()

	origin := "ModuleName"
	data := "value"

	for i := 0; i < b.N; i++ {
		q.NotifyData(origin, data)
	}
}

func BenchmarkNotificationQueueNotifyError(b *testing.B) {
	q := NewNotificationQueue()
	b.ResetTimer()

	err := "error"

	for i := 0; i < b.N; i++ {
		q.NotifyError(constant.Critical, err)
	}

}
