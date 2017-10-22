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
	ErrorType types.ErrorType
	Message   string
}

type DataObject struct {
	From  notification.Origin
	Value string
}

func (q *NotificationQueue) NotifyError(t types.ErrorType, f string, p ...interface{}) {
	q.errors = append(q.errors, *notification.NewNotification("", "", ErrorObject{t, fmt.Sprintf(f, p)}))
}

func (q *NotificationQueue) NotificationsError() []notification.Notification {
	arr := q.errors
	q.errors = []notification.Notification{}
	return arr
}

func (q *NotificationQueue) NotifyData(o notification.Origin, f string, p ...interface{}) {
	q.data = append(q.data, *notification.NewNotification("", "", DataObject{o, fmt.Sprintf(f, p)}))
}

func (q *NotificationQueue) NotificationsData() []notification.Notification {
	arr := q.data
	q.data = []notification.Notification{}
	return arr
}
