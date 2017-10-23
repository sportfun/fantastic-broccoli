package core

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/network"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/const"
	"fmt"
	"go.uber.org/zap"
)

func (c *Core) notificationHandler(n *notification.Notification) {
	switch types.Name(n.From()) {
	case _const.NetworkService:
		netNotificationHandler(c, n)
	default:
		c.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Core, n *notification.Notification) {
	m := n.Content().(network.Message)

	switch m.Command() {
	case "link":
		c.notifications.Notify(notification.NewNotification(_const.Core, _const.NetworkService, c.properties.System.LinkID))
	default:
		c.logger.Error("unknown network command",
			zap.String("where", string(_const.Core)),
			zap.String("command", m.Command()))
	}
}
