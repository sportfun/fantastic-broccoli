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
	case _const.NETWORK_SERVICE:
		netNotificationHandler(c, n)
	default:
		c.logger.Warn("Unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Core, n *notification.Notification) {
	m := n.Content().(network.Message)

	switch m.Command() {
	case "link":
		c.notifications.Notify(notification.NewNotification(_const.CORE, _const.NETWORK_SERVICE, c.properties.System.LinkID))
	default:
		c.logger.Error("Unknown network command",
			zap.String("where", string(_const.CORE)),
			zap.String("command", m.Command()))
	}
}
