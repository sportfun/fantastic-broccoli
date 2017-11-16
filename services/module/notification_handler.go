package module

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
	"fmt"
	"go.uber.org/zap"
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

func netNotificationHandler(s *Service, n *notification.Notification) {
	m := n.Content().(object.NetworkObject)

	switch m.Command {
	case constant.CommandStartSession:
		s.logger.Debug("start session")
		for _, m := range s.modules {
			m.StartSession()
		}
	case constant.CommandEndSession:
		s.logger.Debug("end session")
		for _, m := range s.modules {
			m.StopSession()
		}
	default:
		s.logger.Error("unknown network command",
			zap.String("where", string(constant.Core)),
			zap.String("command", m.Command))
	}
}
