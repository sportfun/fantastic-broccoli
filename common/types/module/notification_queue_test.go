package module

import (
	"testing"
	"fantastic-broccoli/const"
	"fantastic-broccoli/utils"
)

// -- Test

func TestNotificationQueue_NotifyError(t *testing.T) {
	q := NotificationQueue{}

	for i := 0; i < 0xFF; i++ {
		q.NotifyError(_const.FATAL, "Message")
	}

	utils.AssertEquals(t, 0xFF, len(q.errors))
}

// -- Benchmark

func BenchmarkNotificationQueue(b *testing.B) {
	q := NotificationQueue{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			q.NotifyData("ModuleName", "Value")
		} else {
			q.NotifyError(_const.ERROR, "Message")
		}
	}
}
