package service

import (
	"github.com/xunleii/fantastic-broccoli/notification"
	"github.com/xunleii/fantastic-broccoli/utils"
	"testing"
)

func TestNotificationQueueNotify(t *testing.T) {
	q := NewNotificationQueue()
	n := notification.NewNotification("from", "to", nil)

	ns := q.Notifications("to")
	utils.AssertEquals(t, 0, len(ns))

	q.Notify(n)
	ns = q.Notifications("to")
	utils.AssertEquals(t, 1, len(ns))
	utils.AssertEquals(t, n, ns[0])
}
