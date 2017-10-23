package module

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/notification"
	"fmt"
)

type NotificationQueue struct {
	data   []notification.Notification
	errors []notification.Notification
}

type ErrorObject struct {
	ErrorLevel types.ErrorLevel
	Message    string
}

type DataObject struct {
	From  types.Name
	Value string
}

func (q *NotificationQueue) NotifyError(t types.ErrorLevel, f string, p ...interface{}) {
	q.errors = append(q.errors, *notification.NewNotification("", "", ErrorObject{t, fmt.Sprintf(f, p)}))
}

func (q *NotificationQueue) NotificationsError() []notification.Notification {
	arr := q.errors
	q.errors = []notification.Notification{}
	return arr
}

func (q *NotificationQueue) NotifyData(o types.Name, f string, p ...interface{}) {
	q.data = append(q.data, *notification.NewNotification("", "", DataObject{o, fmt.Sprintf(f, p)}))
}

func (q *NotificationQueue) NotificationsData() []notification.Notification {
	arr := q.data
	q.data = []notification.Notification{}
	return arr
}
