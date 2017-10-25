package network

import (
	"github.com/graarh/golang-socketio"
	"fantastic-broccoli/common/types/notification"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/common/types/notification/object"
	"go.uber.org/zap"
	"fmt"
)

var serviceCommandMap = map[string]string{
	constant.CommandLink:         constant.Core,
	constant.CommandStartSession: constant.ModuleService,
	constant.CommandEndSession:   constant.ModuleService,
}

func (s *Service) onConnectionHandler(c *gosocketio.Channel, args interface{}) {
	s.logger.Info(fmt.Sprintf("successfully connected to the server (%s)", c.Id()))
}

func (s *Service) onDisconnectionHandler(c *gosocketio.Channel) {
	s.notifications.Notify(notification.NewNotification(s.Name(), constant.Core, constant.CommandRestartService))
}

func (s *Service) onCommandChanHandler(c *gosocketio.Channel, args interface{}) {
	switch c := args.(type) {
	case webPacket:
		webPacketHandler(s, c)
	default:
		s.logger.Warn("unknown packet type",
			zap.String("where", s.Name()),
			zap.String("packet", fmt.Sprintf("%v", c)),
		)
	}
}

func webPacketHandler(s *Service, packet webPacket) {
	switch b := packet.body.(type) {
	case object.NetworkObject:
		if target, exist := serviceCommandMap[b.Command()]; exist {
			s.notifications.Notify(notification.NewNotification(s.Name(), target, b))
		} else {
			s.logger.Warn("unknown command",
				zap.String("where", s.Name()),
				zap.String("command", fmt.Sprintf("%s", b.Command())),
			)
		}
	default:
		s.logger.Warn("unknown web packet body type",
			zap.String("where", s.Name()),
			zap.String("packet_body", fmt.Sprintf("%v", b)),
		)
	}
}
