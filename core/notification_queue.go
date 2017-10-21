package core

import "fantastic-broccoli/notification"

type NotificationQueue struct {
	notifications map[string][]notification.Notification
}

func (q *NotificationQueue) Notify(n notification.Notification) {
	k := string(n.To())
	q.notifications[k] = append(q.notifications[k], n)
}

func (q *NotificationQueue) Notifications(n Name) []notification.Notification {
	k := string(n)
	arr, ok := q.notifications[k]

	if !ok {
		return []notification.Notification{}
	}
	q.notifications[k] = []notification.Notification{}
	return arr
}
