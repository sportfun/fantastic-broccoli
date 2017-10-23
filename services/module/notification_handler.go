package module

import (
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/common/types/network"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/const"
	"fmt"
	"go.uber.org/zap"
)

func (s *Service) notificationHandler(n *notification.Notification) {
	switch types.Name(n.From()) {
	case _const.NetworkService:
		netNotificationHandler(s, n)
	default:
		s.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(c *Service, n *notification.Notification) {
	m := n.Content().(network.Message)

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
			zap.String("where", string(_const.Core)),
			zap.String("command", m.Command()))
	}
}
