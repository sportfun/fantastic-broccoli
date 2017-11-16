package network

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
	"fmt"
	"go.uber.org/zap"
)

func (s *Service) notificationHandler(n *notification.Notification) error {
	switch string(n.From()) {
	case constant.ModuleService:
		fallthrough
	case constant.Core:
		return serviceNotificationHandler(s, n)
	default:
		s.logger.Warn("unhandled notification",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())),
		)
	}
	return nil
}

func serviceNotificationHandler(s *Service, n *notification.Notification) error {
	var succeed = true

	switch o := n.Content().(type) {
	case *object.NetworkObject:
		succeed = s.emit(constant.CommandChan, o)
	case *object.DataObject:
		succeed = s.emit(constant.DataChan, o)
	case *object.ErrorObject:
		succeed = s.emit(constant.ErrorChan, o)
	default:
		s.logger.Warn("unknown content type",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())),
		)
	}

	if !succeed {
		return fmt.Errorf("failed to emit message")
	}
	return nil
}
