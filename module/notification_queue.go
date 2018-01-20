package module

import (
	"fmt"
	"path"
	"runtime"
	"sync"

	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
)

var builder = notification.NewBuilder().
	From("").
	To("")

type ErrorObject struct {
	object.ErrorObject
	errorLevel string
	origin     Module
}

func (obj *ErrorObject) ErrorLevel() string {
	return obj.errorLevel
}

func (obj *ErrorObject) From() Module {
	return obj.origin
}

type NotificationQueue struct {
	notifications []*notification.Notification
	locker        sync.Mutex
}

func NewNotificationQueue() *NotificationQueue {
	return &NotificationQueue{}
}

func (queue *NotificationQueue) NotifyError(module Module, level string, format string, a ...interface{}) {
	queue.locker.Lock()
	defer queue.locker.Unlock()

	_, caller, line, _ := runtime.Caller(1)
	origin := fmt.Sprintf("%s:%d", path.Base(caller), line)

	errorObject := ErrorObject{*object.NewErrorObject(origin, fmt.Errorf(format, a...)), level, module}
	queue.notifications = append(queue.notifications, builder.With(errorObject).Build())
}

func (queue *NotificationQueue) NotifyData(origin string, format string, a ...interface{}) {
	queue.locker.Lock()
	defer queue.locker.Unlock()

	dataObject := *object.NewDataObject(origin, fmt.Sprintf(format, a...))
	queue.notifications = append(queue.notifications, builder.With(dataObject).Build())
}

func (queue *NotificationQueue) Notifications() []*notification.Notification {
	queue.locker.Lock()
	defer queue.locker.Unlock()

	arr := queue.notifications
	queue.notifications = []*notification.Notification{}
	return arr
}
