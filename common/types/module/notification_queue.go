package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"fmt"
	"path"
	"runtime"
)

var builder = notification.NewBuilder().
	From("").
	To("")

type NotificationQueue struct {
	notifications []*notification.Notification
}

type ErrorObject struct {
	object.ErrorObject
	errorLevel string
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (queue *NotificationQueue) NotifyError(level string, format string, a ...interface{}) {
	_, caller, line, _ := runtime.Caller(1)
	origin := fmt.Sprintf("%s:%d", path.Base(caller), line)
	errorObject := &ErrorObject{*object.NewErrorObject(origin, fmt.Errorf(format, a...)), level}
	queue.notifications = append(queue.notifications, builder.With(errorObject).Build())
}

func (queue *NotificationQueue) NotifyData(origin string, format string, a ...interface{}) {
	dataObject := object.NewDataObject(origin, fmt.Sprintf(format, a...))
	queue.notifications = append(queue.notifications, builder.With(dataObject).Build())
}

func (queue *NotificationQueue) Notifications() []*notification.Notification {
	arr := queue.notifications
	queue.notifications = []*notification.Notification{}
	return arr
}

func (obj *ErrorObject) ErrorLevel() string {
	return obj.errorLevel
}
