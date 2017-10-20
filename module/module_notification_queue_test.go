package module

import (
	"testing"
	"fantastic-broccoli/model"
	"fantastic-broccoli/utils"
)

// -- Test

func TestNotificationQueue_NotifyError(t *testing.T) {
	q := NotificationQueue{}

	for i := 0; i < 0xFF; i++ {
		q.NotifyError(model.FATAL, "Message")
	}

	utils.AssertEquals(t, 0xFF, len(q._errors))
}

// -- Benchmark

func BenchmarkNotificationQueue(b *testing.B) {
	q := NotificationQueue{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			q.NotifyMetric("ModuleName", "Value")
		} else {
			q.NotifyError(model.ERROR, "Message")
		}
	}
}
