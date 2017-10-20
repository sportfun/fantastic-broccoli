package module

import (
	"fantastic-broccoli/notification"
	"fantastic-broccoli/model"
)

type NotificationQueue struct {
	_metrics []notification.Notification
	_errors  []notification.Notification
}

type ErrorObject struct {
	ErrorType model.ErrorType
	Message   string
}

type MetricObject struct {
	From  notification.Origin
	Value string
}

func (nq *NotificationQueue) NotifyError(t model.ErrorType, m string) {
	nq._errors = append(nq._errors, *notification.NewNotification("", "", ErrorObject{t, m}))
}

func (nq *NotificationQueue) NotificationsError() []notification.Notification {
	arr := nq._errors
	nq._errors = []notification.Notification{}
	return arr
}

func (nq *NotificationQueue) NotifyMetric(f notification.Origin, v string) {
	nq._metrics = append(nq._metrics, *notification.NewNotification("", "", MetricObject{f, v}))
}

func (nq *NotificationQueue) NotificationsMetric() []notification.Notification {
	arr := nq._metrics
	nq._metrics = []notification.Notification{}
	return arr
}
