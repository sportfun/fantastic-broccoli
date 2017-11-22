package kernel

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

var (
	DebugNotificationHandled    = log.NewArgumentBinder("notification handled")
	UnhandledNotificationOrigin = log.NewArgumentBinder("unhandled notification origin (%s)")
	InvalidNetworkNotification  = log.NewArgumentBinder("invalid network notification")
	UnknownNetworkCommand       = log.NewArgumentBinder("unknown network command (%s)")
)

func (core *Core) handle(n *notification.Notification) {
	core.logger.Debug(DebugNotificationHandled.More("notification", *n))

	switch string(n.From()) {
	case constant.EntityNames.Services.Network:
		netNotificationHandler(core, n)
	default:
		core.logger.Warn(UnhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
}

func netNotificationHandler(core *Core, n *notification.Notification) {
	var commandObject *object.CommandObject

	switch obj := n.Content().(type) {
	case *object.CommandObject:
		commandObject = obj
	default:
		core.logger.Error(InvalidNetworkNotification.More("content", n.Content()))
	}

	switch commandObject.Command {
	case constant.NetCommand.Link:
		core.notifications.Notify(notification.NewNotification(
			constant.EntityNames.Core,
			constant.EntityNames.Services.Network,
			core.properties.System.LinkID,
		))
	default:
		core.logger.Error(UnknownNetworkCommand.Bind(commandObject.Command))
	}
}
