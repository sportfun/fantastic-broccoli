package kernel

import (
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
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
	case env.NetworkServiceEntity:
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
	case env.LinkCmd:
		core.notifications.Notify(notification.NewNotification(
			env.CoreEntity,
			env.NetworkServiceEntity,
			core.config.System.LinkID,
		))
	default:
		core.logger.Error(unknownNetworkCommand.Bind(commandObject.Command))
	}
}
