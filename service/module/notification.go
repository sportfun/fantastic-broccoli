package module

import (
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/notification"
	"github.com/xunleii/fantastic-broccoli/notification/object"
	"github.com/xunleii/fantastic-broccoli/env"
)

var (
	debugNotificationHandled      = log.NewArgumentBinder("notification handled")
	debugStartSessionNotification = log.NewArgumentBinder("start session")
	debugEndSessionNotification   = log.NewArgumentBinder("stop session")

	unhandledNotificationOrigin = log.NewArgumentBinder("unhandled notification origin (%s)")
	invalidNetworkNotification  = log.NewArgumentBinder("invalid network notification")
	unknownNetworkCommand       = log.NewArgumentBinder("unknown network command (%s)")
)

func (service *Service) handle(n *notification.Notification) {
	service.logger.Debug(debugNotificationHandled.More("notification", *n))

	switch string(n.From()) {
	case env.NetworkServiceEntity:
		netNotificationHandler(service, n)
	default:
		service.logger.Warn(unhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
}

func netNotificationHandler(service *Service, n *notification.Notification) {
	var commandObject *object.CommandObject

	switch obj := n.Content().(type) {
	case *object.CommandObject:
		commandObject = obj
	default:
		service.logger.Error(invalidNetworkNotification.More("content", n.Content()))
		return
	}

	switch commandObject.Command {
	case env.StartSessionCmd:
		service.logger.Debug(debugStartSessionNotification)
		for _, m := range service.modules {
			m.StartSession()
		}
	case env.EndSessionCmd:
		service.logger.Debug(debugEndSessionNotification)
		for _, m := range service.modules {
			m.StopSession()
		}
	default:
		service.logger.Error(unknownNetworkCommand.Bind(commandObject.Command))
	}
}
