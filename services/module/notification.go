package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

var (
	DebugNotificationHandled      = log.NewArgumentBinder("notification handled")
	DebugStartSessionNotification = log.NewArgumentBinder("start session")
	DebugEndSessionNotification   = log.NewArgumentBinder("stop session")
	UnhandledNotificationOrigin   = log.NewArgumentBinder("unhandled notification origin (%s)")
	InvalidNetworkNotification    = log.NewArgumentBinder("invalid network notification")
	UnknownNetworkCommand         = log.NewArgumentBinder("unknown network command (%s)")
)

func (service *Service) handle(n *notification.Notification) {
	service.logger.Debug(DebugNotificationHandled.More("notification", *n))

	switch string(n.From()) {
	case constant.EntityNames.Services.Network:
		netNotificationHandler(service, n)
	default:
		service.logger.Warn(UnhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
}

func netNotificationHandler(service *Service, n *notification.Notification) {
	var commandObject *object.CommandObject

	switch obj := n.Content().(type) {
	case *object.CommandObject:
		commandObject = obj
	default:
		service.logger.Error(InvalidNetworkNotification.More("content", n.Content()))
	}

	switch commandObject.Command {
	case constant.NetCommand.StartSession:
		service.logger.Debug(DebugStartSessionNotification)
		for _, m := range service.modules {
			m.StartSession()
		}
	case constant.NetCommand.EndSession:
		service.logger.Debug(DebugEndSessionNotification)
		for _, m := range service.modules {
			m.StopSession()
		}
	default:
		service.logger.Error(UnknownNetworkCommand.Bind(commandObject.Command))
	}
}
