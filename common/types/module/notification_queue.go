package module

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fmt"
	"path"
	"runtime"
)

var builder = notification.NewBuilder().
	From("").
	To("")

type NotificationQueue struct {
	data   []*notification.Notification
	errors []*notification.Notification
}

type ErrorObject struct {
	object.ErrorObject
	errorLevel int
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (queue *NotificationQueue) NotifyError(level int, format string, a ...interface{}) {
	_, caller, line, _ := runtime.Caller(1)
	origin := fmt.Sprintf("%s:%d", path.Base(caller), line)
	errorObject := ErrorObject{*object.NewErrorObject(origin, fmt.Errorf(format, a...)), level}
	queue.errors = append(queue.errors, builder.With(errorObject).Build())
}

func (queue *NotificationQueue) NotificationsError() []*notification.Notification {
	arr := queue.errors
	queue.errors = []*notification.Notification{}
	return arr
}

func (queue *NotificationQueue) NotifyData(origin string, format string, a ...interface{}) {
	dataObject := *object.NewDataObject(origin, fmt.Sprintf(format, a...))
	queue.data = append(queue.data, builder.With(dataObject).Build())
}

func (queue *NotificationQueue) NotificationsData() []*notification.Notification {
	arr := queue.data
	queue.data = []*notification.Notification{}
	return arr
}

func (obj *ErrorObject) ErrorLevel() int {
	return obj.errorLevel
}
