package network

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
)

func (service *Service) handle(notif *notification.Notification) error {
	switch string(notif.From()) {
	case constant.EntityNames.Services.Module:
		fallthrough
	case constant.EntityNames.Core:
		return defaultNotificationHandler(service, notif)
	default:
		service.logger.Warn("unhandled notification",
			zap.String("where", string(notif.To())),
			zap.String("from", string(notif.From())),
			zap.String("message", fmt.Sprintf("%#v", notif.Content())),
		)
	}
	return nil
}

func defaultNotificationHandler(service *Service, n *notification.Notification) error {
	var succeed = true

	switch o := n.Content().(type) {
	case *object.CommandObject:
		succeed = service.emit(constant.Channels.Command, o)
	case *object.DataObject:
		succeed = service.emit(constant.Channels.Data, o)
	case *object.ErrorObject:
		succeed = service.emit(constant.Channels.Error, o)
	default:
		service.logger.Warn("unknown content type",
			zap.String("where", string(n.To())),
			zap.String("from", string(n.From())),
			zap.String("message", fmt.Sprintf("%#v", n.Content())),
		)
	}

	if !succeed {
		return fmt.Errorf("failed to emit message")
	}
	return nil
}
