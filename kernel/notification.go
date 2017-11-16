package kernel

import (
	"fmt"
	"go.uber.org/zap"

	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
)

func (core *Core) handle(n *notification.Notification) {

	switch string(n.From()) {
	case constant.EntityNames.Services.Network:
		netNotificationHandler(core, n)
	default:
		core.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Core, n *notification.Notification) {
	m := n.Content().(object.NetworkObject)

	switch m.Command {
	case constant.NetCommand.Link:
		c.notifications.Notify(notification.NewNotification(
			constant.EntityNames.Core,
			constant.EntityNames.Services.Network,
			c.properties.System.LinkID,
		))
	default:
		c.logger.Error("unknown network command",
			zap.String("where", string(constant.EntityNames.Core)),
			zap.String("command", m.Command))
	}
}
