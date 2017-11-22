package network

import (
	"github.com/graarh/golang-socketio"
	"github.com/mitchellh/mapstructure"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

var serviceCommandMap = map[types.CommandName]string{
	constant.NetCommand.Link:         constant.EntityNames.Core,
	constant.NetCommand.StartSession: constant.EntityNames.Services.Module,
	constant.NetCommand.EndSession:   constant.EntityNames.Services.Module,
}

var (
	DebugDisconnectionHandled = log.NewArgumentBinder("disconnection handled")
	DebugCommandHandled       = log.NewArgumentBinder("command packet handled")
	DebugCommand              = log.NewArgumentBinder("valid command handled")

	SuccessfullyConnected = log.NewArgumentBinder("successfully connected to the server")
	UnknownPacketType     = log.NewArgumentBinder("unknown packet type")
	UnknownCommand        = log.NewArgumentBinder("unknown command '%s'")
	UnknownWebPacketBody  = log.NewArgumentBinder("unknown web packet body type")
)

func (service *Service) onConnectionHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Info(SuccessfullyConnected.More("session_id", client.Id()))
}

func (service *Service) onDisconnectionHandler(client *gosocketio.Channel) {
	service.logger.Debug(DebugDisconnectionHandled.More("session_id", client.Id()))

	if service.state != constant.States.Stopped {
		service.notifications.Notify(notification.NewNotification(service.Name(), constant.EntityNames.Core, constant.NetCommand.RestartService))
	}
}

func (service *Service) onCommandChanHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Debug(DebugCommandHandled.More("session_id", client.Id()).More("packet", args))

	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		webPacketHandler(service, web)
	default:
		service.logger.Warn(UnknownPacketType.More("session_id", client.Id()))
	}
}

func webPacketHandler(service *Service, packet webPacket) {
	var netObj object.CommandObject

	switch {
	case mapstructure.Decode(packet.Body, &netObj) == nil:
		if target, exist := serviceCommandMap[netObj.Command]; exist {
			service.logger.Debug(DebugCommand.More("target", target).More("object", netObj))
			service.notifications.Notify(notification.NewNotification(service.Name(), target, netObj))
		} else {
			service.logger.Warn(UnknownCommand.Bind(netObj.Command))
		}
	default:
		service.logger.Warn(UnknownWebPacketBody.More("packet_body", netObj))
	}
}
