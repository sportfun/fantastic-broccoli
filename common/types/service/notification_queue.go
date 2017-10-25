package service

import "fantastic-broccoli/common/types/notification"

type NotificationQueue struct {
	notifications map[string][]*notification.Notification
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{map[string][]*notification.Notification{}}
}

func (queue *NotificationQueue) Notify(notification *notification.Notification) {
	queue.notifications[notification.To()] = append(queue.notifications[notification.To()], notification)
}

func (queue *NotificationQueue) Notifications(name string) []*notification.Notification {
	arr, ok := queue.notifications[name]

	if !ok {
		return []*notification.Notification{}
	}
	queue.notifications[name] = []*notification.Notification{}
	return arr
}
