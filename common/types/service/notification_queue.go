package service

import "fantastic-broccoli/common/types/notification"

type notificationQueue struct {
	notifications map[string][]*notification.Notification
}

func NewNotificationQueue() *notificationQueue {
	return &notificationQueue{map[string][]*notification.Notification{}}
}

func (queue *notificationQueue) Notify(notification *notification.Notification) {
	queue.notifications[notification.To()] = append(queue.notifications[notification.To()], notification)
}

func (queue *notificationQueue) Notifications(name string) []*notification.Notification {
	arr, ok := queue.notifications[name]

	if !ok {
		return []*notification.Notification{}
	}
	queue.notifications[name] = []*notification.Notification{}
	return arr
}
