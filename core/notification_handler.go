package core

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/constant"
	"fmt"
	"go.uber.org/zap"
	"fantastic-broccoli/common/types/notification/object"
)

func (c *Core) notificationHandler(n *notification.Notification) {
	// TODO: Check notif type
	switch string(n.From()) {
	case constant.NetworkService:
		netNotificationHandler(c, n)
	default:
		c.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Core, n *notification.Notification) {
	m := n.Content().(object.NetworkObject)

	switch m.Command() {
	case "link":
		c.notifications.Notify(notification.NewNotification(constant.Core, constant.NetworkService, c.properties.System.LinkID))
	default:
		c.logger.Error("unknown network command",
			zap.String("where", string(constant.Core)),
			zap.String("command", m.Command()))
	}
}
