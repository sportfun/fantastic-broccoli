package module

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/constant"
	"fmt"
	"go.uber.org/zap"
	"fantastic-broccoli/common/types/notification/object"
)

func (s *Service) notificationHandler(n *notification.Notification) {
	switch string(n.From()) {
	case constant.NetworkService:
		netNotificationHandler(s, n)
	default:
		s.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Service, n *notification.Notification) {
	m := n.Content().(object.NetworkObject)

	switch m.Command() {
	case "new_session":
		c.logger.Debug("new session")
		for _, m := range c.modules {
			m.StartSession()
		}
	case "end_session":
		c.logger.Debug("end session")
		for _, m := range c.modules {
			m.StopSession()
		}
	default:
		c.logger.Error("unknown network command",
			zap.String("where", string(constant.Core)),
			zap.String("command", m.Command()))
	}
}
