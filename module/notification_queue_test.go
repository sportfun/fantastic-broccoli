package module

import (
	"fmt"
	. "github.com/onsi/gomega"
	. "github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/utils"
	"sync"
	"testing"
)

func TestErrorObject_ErrorLevel(t *testing.T) {
	RegisterTestingT(t)
	obj := ErrorObject{ErrorObject: *object.NewErrorObject("", fmt.Errorf("")), errorLevel: WarningLevel}

	Expect(obj.ErrorLevel()).Should(Equal(WarningLevel))
}

func TestNotificationQueue(t *testing.T) {
	RegisterTestingT(t)

	queue := NewNotificationQueue()
	NotifyError := func(l string, f string, a ...interface{}) { queue.NotifyError(nil, l, f, a...) }
	testCases := []struct {
		Notifier  func(string, string, ...interface{})
		FParam    string
		Format    string
		Arguments []interface{}
		Object    interface{}
	}{
		{Notifier: queue.NotifyData, FParam: "module name", Format: "data value: %v", Arguments: []interface{}{87 + 6i}, Object: object.NewDataObject("module name", "data value: (87+6i)")},
		{Notifier: queue.NotifyData, FParam: "", Format: "data value: %v", Arguments: []interface{}{87 + 6i}, Object: object.NewDataObject("", "data value: (87+6i)")},
		{Notifier: queue.NotifyData, FParam: "module name", Format: "value", Arguments: []interface{}{}, Object: object.NewDataObject("module name", "value")},

		{Notifier: NotifyError, FParam: ErrorLevel, Format: "error message: %s", Arguments: []interface{}{"failure"}, Object: &ErrorObject{*object.NewErrorObject("notification_queue_test.go:25", fmt.Errorf("error message: failure")), ErrorLevel, nil}},
		{Notifier: NotifyError, FParam: WarningLevel, Format: "warning message: %s", Arguments: []interface{}{"failure"}, Object: &ErrorObject{*object.NewErrorObject("notification_queue_test.go:25", fmt.Errorf("warning message: failure")), WarningLevel, nil}},
		{Notifier: NotifyError, FParam: CriticalLevel, Format: "critical error", Arguments: []interface{}{}, Object: &ErrorObject{*object.NewErrorObject("notification_queue_test.go:25", fmt.Errorf("critical error")), CriticalLevel, nil}},
	}

	for _, tc := range testCases {
		tc.Notifier(tc.FParam, tc.Format, tc.Arguments...)
		Expect(queue.Notifications()).Should(ConsistOf(notification.NewNotification("", "", tc.Object)))
		Expect(queue.notifications).Should(BeEmpty())
	}
}

func TestNotificationQueue_RaceCondition(t *testing.T) {
	queue := NewNotificationQueue()
	wg := sync.WaitGroup{}
	inc := utils.NewIncrementVolatile(0).(utils.Incremental)

	wg.Add(0xF0)
	for i := 0; i < 0xFF; i++ {
		go func() {
			queue.NotifyData("origin", "%v", false)
			if inc.Get().(int) < 0xF0 {
				inc.Inc(1)
				wg.Done()
			}
		}()
	}
	wg.Wait()
	queue.Notifications()
}
