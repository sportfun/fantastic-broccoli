package kernel

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
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
	var commandObject object.CommandObject

	switch obj := n.Content().(type) {
	case object.CommandObject:
		commandObject = obj
	default:
		c.logger.Error("invalid network command",
			zap.String("where", string(constant.EntityNames.Core)),
			zap.Any("content", n.Content()))

	}

	switch commandObject.Command {
	case constant.NetCommand.Link:
		c.notifications.Notify(notification.NewNotification(
			constant.EntityNames.Core,
			constant.EntityNames.Services.Network,
			c.properties.System.LinkID,
		))
	default:
		c.logger.Error("unknown network command",
			zap.String("where", string(constant.EntityNames.Core)),
			zap.String("command", commandObject.Command))
	}
}
