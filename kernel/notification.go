package kernel

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

var (
	debugNotificationHandled    = log.NewArgumentBinder("notification handled")
	unhandledNotificationOrigin = log.NewArgumentBinder("unhandled notification origin (%s)")
	invalidNetworkNotification  = log.NewArgumentBinder("invalid network notification")
	unknownNetworkCommand       = log.NewArgumentBinder("unknown network command (%s)")
)

func (core *Core) handle(n *notification.Notification) {
	core.logger.Debug(debugNotificationHandled.More("notification", *n))

	switch string(n.From()) {
	case constant.EntityNames.Services.Network:
		netNotificationHandler(core, n)
	default:
		core.logger.Warn(unhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
}

func netNotificationHandler(core *Core, n *notification.Notification) {
	var commandObject *object.CommandObject

	switch obj := n.Content().(type) {
	case *object.CommandObject:
		commandObject = obj
	default:
		core.logger.Error(invalidNetworkNotification.More("content", n.Content()))
	}

	switch commandObject.Command {
	case constant.NetCommand.Link:
		core.notifications.Notify(notification.NewNotification(
			constant.EntityNames.Core,
			constant.EntityNames.Services.Network,
			core.properties.System.LinkID,
		))
	default:
		core.logger.Error(unknownNetworkCommand.Bind(commandObject.Command))
	}
}
