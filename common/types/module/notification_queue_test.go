package module

import (
	"fantastic-broccoli/constant"
	"fantastic-broccoli/utils"
	"runtime"
	"testing"
	"fantastic-broccoli/common/types/notification/object"
)

// -- Test

func TestNotificationQueueNotifyError(t *testing.T) {
	q := NewNotificationQueue()

	for i := 0; i < 0xFF; i++ {
		q.NotifyError(constant.FATAL, "Error message %s %s", "can be", "formatted")
	}

	utils.AssertEquals(t, 0xFF, len(q.errors))
	errors := q.NotificationsError()
	utils.AssertEquals(t, 0xFF, len(errors))
	utils.AssertEquals(t, 0, len(q.errors))

	o := errors[0].Content().(ErrorObject)
	_, origin, _, _ := runtime.Caller(0)
	utils.AssertEquals(t, origin, o.Origin())
	utils.AssertEquals(t, "Error message can be formatted", o.Reason().Error())
	utils.AssertEquals(t, constant.FATAL, o.ErrorLevel())
}

func TestNotificationQueueNotifyData(t *testing.T) {
	q := NewNotificationQueue()

	for i := 0; i < 0xFF; i++ {
		q.NotifyData("ModuleName", "%d RPM", 1000)
	}

	utils.AssertEquals(t, 0xFF, len(q.data))
	data := q.NotificationsData()
	utils.AssertEquals(t, 0xFF, len(data))
	utils.AssertEquals(t, 0, len(q.data))

	o := data[0].Content().(object.DataObject)
	utils.AssertEquals(t, "ModuleName", o.Module())
	utils.AssertEquals(t, "1000 RPM", o.Value())
}

// -- Benchmark

func BenchmarkNotificationQueue(b *testing.B) {
	q := notificationQueue{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			q.NotifyData("ModuleName", "Value")
		} else {
			q.NotifyError(constant.ERROR, "Message")
		}
	}
}
