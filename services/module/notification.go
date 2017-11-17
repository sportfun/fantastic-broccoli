package module

import (
	"fmt"
	"go.uber.org/zap"

	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
)

func (service *Service) handle(n *notification.Notification) {
	switch string(n.From()) {
	case constant.EntityNames.Services.Network:
		netNotificationHandler(service, n)
	default:
		service.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())))
	}
}

func netNotificationHandler(s *Service, n *notification.Notification) {
	m := n.Content().(object.NetworkObject)

	switch m.Command {
	case constant.NetCommand.StartSession:
		s.logger.Debug("start session")
		for _, m := range s.modules {
			m.StartSession()
		}
	case constant.NetCommand.EndSession:
		s.logger.Debug("end session")
		for _, m := range s.modules {
			m.StopSession()
		}
	default:
		s.logger.Error("unknown network command",
			zap.String("where", string(constant.EntityNames.Services.Module)),
			zap.String("command", m.Command))
	}
}
