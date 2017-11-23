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
	debugDisconnectionHandled = log.NewArgumentBinder("disconnection handled")
	debugCommandHandled       = log.NewArgumentBinder("command packet handled")
	debugCommand              = log.NewArgumentBinder("valid command handled")

	successfullyConnected = log.NewArgumentBinder("successfully connected to the server")
	unknownPacketType     = log.NewArgumentBinder("unknown packet type")
	unknownCommand        = log.NewArgumentBinder("unknown command '%s'")
	unknownWebPacketBody  = log.NewArgumentBinder("unknown web packet body type")
)

func (service *Service) onConnectionHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Info(successfullyConnected.More("session_id", client.Id()))
}

func (service *Service) onDisconnectionHandler(client *gosocketio.Channel) {
	service.logger.Debug(debugDisconnectionHandled.More("session_id", client.Id()))

	if service.state != constant.States.Stopped {
		service.notifications.Notify(notification.NewNotification(service.Name(), constant.EntityNames.Core, constant.NetCommand.RestartService))
	}
}

func (service *Service) onCommandChanHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Debug(debugCommandHandled.More("session_id", client.Id()).More("packet", args))

	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		webPacketHandler(service, web)
	default:
		service.logger.Warn(unknownPacketType.More("session_id", client.Id()))
	}
}

func webPacketHandler(service *Service, packet webPacket) {
	var netObj object.CommandObject

	switch {
	case mapstructure.Decode(packet.Body, &netObj) == nil:
		if target, exist := serviceCommandMap[netObj.Command]; exist {
			service.logger.Debug(debugCommand.More("target", target).More("object", netObj))
			service.notifications.Notify(notification.NewNotification(service.Name(), target, netObj))
		} else {
			service.logger.Warn(unknownCommand.Bind(netObj.Command))
		}
	default:
		service.logger.Warn(unknownWebPacketBody.More("packet_body", netObj))
	}
}
