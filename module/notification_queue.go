package module

import (
	"fantastic-broccoli/notification"
	"fantastic-broccoli/model"
)

type NotificationQueue struct {
	metrics []notification.Notification
	errors  []notification.Notification
}

type ErrorObject struct {
	ErrorType model.ErrorType
	Message   string
}

type MetricObject struct {
	From  notification.Origin
	Value string
}

func (q *NotificationQueue) NotifyError(t model.ErrorType, m string) {
	q.errors = append(q.errors, *notification.NewNotification("", "", ErrorObject{t, m}))
}

func (q *NotificationQueue) NotificationsError() []notification.Notification {
	arr := q.errors
	q.errors = []notification.Notification{}
	return arr
}

func (q *NotificationQueue) NotifyMetric(f notification.Origin, v string) {
	q.metrics = append(q.metrics, *notification.NewNotification("", "", MetricObject{f, v}))
}

func (q *NotificationQueue) NotificationsMetric() []notification.Notification {
	arr := q.metrics
	q.metrics = []notification.Notification{}
	return arr
}
