package network

import (
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"fmt"
	"github.com/graarh/golang-socketio"
	"go.uber.org/zap"
	"github.com/mitchellh/mapstructure"
)

var serviceCommandMap = map[string]string{
	constant.NetCommand.Link:         constant.EntityNames.Core,
	constant.NetCommand.StartSession: constant.EntityNames.Services.Module,
	constant.NetCommand.EndSession:   constant.EntityNames.Services.Module,
}

func (service *Service) onConnectionHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Info("successfully connected to the server", zap.String("id", client.Id()))
}

func (service *Service) onDisconnectionHandler(client *gosocketio.Channel) {
	service.logger.Debug("disconnection handled", zap.String("id", client.Id()))
	if service.state != constant.States.Stopped {
		service.notifications.Notify(notification.NewNotification(service.Name(), constant.EntityNames.Core, constant.NetCommand.RestartService))
	}
}

func (service *Service) onCommandChanHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Debug("command handled",
		zap.String("id", client.Id()),
		zap.String("packet", fmt.Sprint(args)),
	)

	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		webPacketHandler(service, web)
	default:
		service.logger.Warn("unknown packet type",
			zap.String("where", service.Name()),
			zap.String("packet", fmt.Sprintf("%v", args)),
		)
	}
}

func webPacketHandler(s *Service, packet webPacket) {
	var netObj object.CommandObject

	switch {
	case mapstructure.Decode(packet.Body, &netObj) == nil:
		if target, exist := serviceCommandMap[netObj.Command]; exist {
			s.logger.Debug("command handled", zap.String("target", target), zap.Any("object", netObj))
			s.notifications.Notify(notification.NewNotification(s.Name(), target, netObj))
		} else {
			s.logger.Warn("unknown command",
				zap.String("where", s.Name()),
				zap.String("command", fmt.Sprintf("%s", netObj.Command)),
			)
		}
	default:
		s.logger.Warn("unknown web packet body type",
			zap.String("where", s.Name()),
			zap.String("packet_body", fmt.Sprintf("%v", netObj)),
		)
	}
}
