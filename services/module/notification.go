package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
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
	case constant.EntityNames.Services.Network:
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
	}

	switch commandObject.Command {
	case constant.NetCommand.StartSession:
		service.logger.Debug(debugStartSessionNotification)
		for _, m := range service.modules {
			m.StartSession()
		}
	case constant.NetCommand.EndSession:
		service.logger.Debug(debugEndSessionNotification)
		for _, m := range service.modules {
			m.StopSession()
		}
	default:
		service.logger.Error(unknownNetworkCommand.Bind(commandObject.Command))
	}
}
