package module

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fmt"
	"runtime"
)

var builder = notification.NewBuilder().
	From("").
	To("")

type notificationQueue struct {
	data   []*notification.Notification
	errors []*notification.Notification
}

type ErrorObject struct {
	object.ErrorObject
	errorLevel int
}

func NewNotificationQueue() *notificationQueue {
	return &notificationQueue{}
}

func (queue *notificationQueue) NotifyError(level int, format string, a ...interface{}) {
	_, origin, _, _ := runtime.Caller(1)
	errorObject := ErrorObject{*object.NewErrorObject(origin, fmt.Errorf(format, a...)), level}
	queue.errors = append(queue.errors, builder.With(errorObject).Build())
}

func (queue *notificationQueue) NotificationsError() []*notification.Notification {
	arr := queue.errors
	queue.errors = []*notification.Notification{}
	return arr
}

func (queue *notificationQueue) NotifyData(origin string, format string, a ...interface{}) {
	dataObject := object.NewDataObject(origin, fmt.Sprintf(format, a...))
	queue.data = append(queue.data, builder.With(dataObject).Build())
}

func (queue *notificationQueue) NotificationsData() []*notification.Notification {
	arr := queue.data
	queue.data = []*notification.Notification{}
	return arr
}

func (obj *ErrorObject) ErrorLevel() int {
	return obj.errorLevel
}
