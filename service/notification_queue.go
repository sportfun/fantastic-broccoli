package service

import (
	"sync"

	"github.com/xunleii/fantastic-broccoli/notification"
)

type NotificationQueue struct {
	locker        sync.Mutex
	notifications map[string][]*notification.Notification
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{locker: sync.Mutex{}, notifications: map[string][]*notification.Notification{}}
}

func (queue *NotificationQueue) Notify(notification *notification.Notification) {
	queue.locker.Lock()
	defer queue.locker.Unlock()

	queue.notifications[notification.To()] = append(queue.notifications[notification.To()], notification)
}

func (queue *NotificationQueue) Notifications(name string) []*notification.Notification {
	queue.locker.Lock()
	defer queue.locker.Unlock()

	arr, ok := queue.notifications[name]

	if !ok {
		return []*notification.Notification{}
	}
	queue.notifications[name] = []*notification.Notification{}
	return arr
}
