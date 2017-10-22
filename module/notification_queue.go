package module

import (
	. "fantastic-broccoli"
	"fantastic-broccoli/notification"
)

type NotificationQueue struct {
	data   []notification.Notification
	errors []notification.Notification
}

type ErrorObject struct {
	ErrorType ErrorType
	Message   string
}

type DataObject struct {
	From  notification.Origin
	Value string
}

func (q *NotificationQueue) NotifyError(t ErrorType, m string) {
	q.errors = append(q.errors, *notification.NewNotification("", "", ErrorObject{t, m}))
}

func (q *NotificationQueue) NotificationsError() []notification.Notification {
	arr := q.errors
	q.errors = []notification.Notification{}
	return arr
}

func (q *NotificationQueue) NotifyData(f notification.Origin, v string) {
	q.data = append(q.data, *notification.NewNotification("", "", DataObject{f, v}))
}

func (q *NotificationQueue) NotificationsData() []notification.Notification {
	arr := q.data
	q.data = []notification.Notification{}
	return arr
}
