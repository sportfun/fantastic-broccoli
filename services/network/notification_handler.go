package network

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/constant"
	"go.uber.org/zap"
	"fmt"
	"fantastic-broccoli/common/types/notification/object"
)

func (s *Service) messageHandler(m *object.NetworkObject) {

}

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
		succeed = s.Emit(constant.CommandChan, o)
	case *object.DataObject:
		succeed = s.Emit(constant.DataChan, o)
	case *object.ErrorObject:
		succeed = s.Emit(constant.ErrorChan, o)
	default:
		s.logger.Warn("unknown content type",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%v", n.Content())),
		)
	}

	if !succeed {
		// TODO: Write error
		return fmt.Errorf("")
	}
	return nil
}
