package network

import (
	"fmt"

	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"time"
)

var (
	debugNotificationHandled    = log.NewArgumentBinder("notification handled")
	unhandledNotificationOrigin = log.NewArgumentBinder("unhandled notification origin (%s)")
	unknownContentType          = log.NewArgumentBinder("unknown content type")
)

func (service *Network) handle(n *notification.Notification) error {
	//service.logger.Debug(debugNotificationHandled.More("notification", fmt.Sprintf("%#v", n)))
	time.Sleep(5 * time.Millisecond)

	switch string(n.From()) {
	case env.ModuleServiceEntity:
		fallthrough
	case env.CoreEntity:
		return defaultNotificationHandler(service, n)
	default:
		service.logger.Warn(unhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
	return nil
}

func defaultNotificationHandler(service *Network, n *notification.Notification) error {
	var succeed = true

	switch o := n.Content().(type) {
	case object.CommandObject:
		succeed = service.emit(OnCommand, o)
	case object.DataObject:
		succeed = service.emit(OnData, o)
	case object.ErrorObject:
		succeed = service.emit(OnError, o)
	default:
		service.logger.Warn(unknownContentType.More("packet", n.Content()))
	}

	if !succeed {
		return fmt.Errorf("failed to emit message")
	}
	return nil
}
