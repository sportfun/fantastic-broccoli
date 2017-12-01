package network

import (
	"github.com/graarh/golang-socketio"
	"github.com/mitchellh/mapstructure"

	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
)

var serviceCommandMapper = map[string]string{
	env.LinkCmd:         env.CoreEntity,
	env.StartSessionCmd: env.ModuleServiceEntity,
	env.EndSessionCmd:   env.ModuleServiceEntity,
}

const (
	debugDisconnectionHandled = "disconnection handled"
	debugCommandHandled       = "command packet handled"
	debugCommand              = "valid command handled"

	successfullyConnected = "successfully connected to the server"
	unknownPacketType     = "unknown packet type"
	unknownCommand        = "unknown command '%s'"
	unknownWebPacketBody  = "unknown web packet body type"
)

func (service *Network) onConnectionHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Info(log.NewArgumentBinder(successfullyConnected).More("session_id", client.Id()))
}

func (service *Network) onDisconnectionHandler(client *gosocketio.Channel) {
	service.logger.Debugf(debugDisconnectionHandled)

	if service.state != env.StoppedState {
		service.notifications.Notify(notification.NewNotification(service.Name(), env.CoreEntity, env.RestartServiceCmd))
	}
}

func (service *Network) onCommandChanHandler(client *gosocketio.Channel, args interface{}) {
	service.logger.Debug(log.NewArgumentBinder(debugCommandHandled).More("session_id", client.Id()).More("packet", args))

	var web websocket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		if web.LinkId == "" {
			break
		}
		webPacketHandler(service, web)
		return
	}
	service.logger.Warn(log.NewArgumentBinder(unknownPacketType).More("session_id", client.Id()))
}

func webPacketHandler(service *Network, packet websocket) {
	var netObj object.CommandObject

	switch {
	case mapstructure.Decode(packet.Body, &netObj) == nil:
		if netObj.Command == "" {
			break
		}

		if target, exist := serviceCommandMapper[netObj.Command]; exist {
			service.logger.Debug(log.NewArgumentBinder(debugCommand).More("target", target).More("object", netObj))
			service.notifications.Notify(notification.NewNotification(service.Name(), target, netObj))
		} else {
			service.logger.Warnf(unknownCommand, netObj.Command)
		}
		return
	}
	service.logger.Warn(log.NewArgumentBinder(unknownWebPacketBody).More("packet_body", netObj))
}
