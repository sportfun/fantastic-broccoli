package network

import (
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/constant"
	"fmt"
	"github.com/graarh/golang-socketio"
	"go.uber.org/zap"
	"github.com/mitchellh/mapstructure"
)

var serviceCommandMap = map[string]string{
	constant.CommandLink:         constant.Core,
	constant.CommandStartSession: constant.ModuleService,
	constant.CommandEndSession:   constant.ModuleService,
}

func (s *Service) onConnectionHandler(client *gosocketio.Channel, args interface{}) {
	s.logger.Info("successfully connected to the server", zap.String("id", client.Id()))
}

func (s *Service) onDisconnectionHandler(client *gosocketio.Channel) {
	s.logger.Debug("disconnection handled", zap.String("id", client.Id()))
	if s.state != constant.Stopped {
		s.notifications.Notify(notification.NewNotification(s.Name(), constant.Core, constant.CommandRestartService))
	}
}

func (s *Service) onCommandChanHandler(client *gosocketio.Channel, args interface{}) {
	s.logger.Debug("command handled",
		zap.String("id", client.Id()),
		zap.String("packet", fmt.Sprint(args)),
	)

	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		webPacketHandler(s, web)
	default:
		s.logger.Warn("unknown packet type",
			zap.String("where", s.Name()),
			zap.String("packet", fmt.Sprintf("%v", args)),
		)
	}
}

//TODO: Risk of concurrency (notification.Notify)
func webPacketHandler(s *Service, packet webPacket) {
	var netObj object.NetworkObject

	switch {
	case mapstructure.Decode(packet.Body, &netObj) == nil:
		if target, exist := serviceCommandMap[netObj.Command]; exist {
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
