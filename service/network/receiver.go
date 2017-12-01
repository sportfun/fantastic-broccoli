package network

import (
	"github.com/graarh/golang-socketio"
	"github.com/mitchellh/mapstructure"

	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/notification"
	"github.com/xunleii/fantastic-broccoli/notification/object"
	"github.com/xunleii/fantastic-broccoli/env"
)

var serviceCommandMapper = map[string]string{
	env.LinkCmd:         env.CoreEntity,
	env.StartSessionCmd: env.ModuleServiceEntity,
	env.EndSessionCmd:   env.ModuleServiceEntity,
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

	if service.state != env.StoppedState {
		service.notifications.Notify(notification.NewNotification(service.Name(), env.CoreEntity, env.RestartServiceCmd))
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
		if target, exist := serviceCommandMapper[netObj.Command]; exist {
			service.logger.Debug(debugCommand.More("target", target).More("object", netObj))
			service.notifications.Notify(notification.NewNotification(service.Name(), target, netObj))
		} else {
			service.logger.Warn(unknownCommand.Bind(netObj.Command))
		}
	default:
		service.logger.Warn(unknownWebPacketBody.More("packet_body", netObj))
	}
}
