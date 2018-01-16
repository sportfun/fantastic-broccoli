package service

import (
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/notification"
	"testing"
)

func TestNotificationQueue(t *testing.T) {
	RegisterTestingT(t)

	queue := NewNotificationQueue()
	Expect(queue.Notifications("")).Should(BeEmpty())

	testCases := []struct {
		Notifications []*notification.Notification
		Result        map[string]struct {
			Number  int
			Content []interface{}
		}
	}{
		// Simple tests
		{
			Notifications: []*notification.Notification{
				notification.NewNotification("", "AAA", 0),
				notification.NewNotification("", "BBB", 1),
				notification.NewNotification("", "CCC", 1),
				notification.NewNotification("", "BBB", 3+2i),
				notification.NewNotification("", "AAA", "data"),
				notification.NewNotification("", "AAA", ""),
				notification.NewNotification("", "BBB", struct{ A int }{A: 29}),
			},
			Result: map[string]struct {
				Number  int
				Content []interface{}
			}{
				"AAA": {Number: 3, Content: []interface{}{0, "data", ""}},
				"BBB": {Number: 3, Content: []interface{}{1, 3 + 2i, struct{ A int }{A: 29}}},
				"CCC": {Number: 1, Content: []interface{}{1}},
			},
		},
	}

	for _, tc := range testCases {
		for _, notification := range tc.Notifications {
			queue.Notify(notification)
		}

		for dst, rslt := range tc.Result {
			ntfns := queue.Notifications(dst)

			Expect(ntfns).Should(And(
				HaveLen(rslt.Number),
				WithTransform(func(n []*notification.Notification) []interface{} {
					content := make([]interface{}, len(n))
					for idx, ntfn := range n {
						content[idx] = ntfn.Content()
					}
					return content
				}, ConsistOf(rslt.Content...)),
			))
		}
	}
}
