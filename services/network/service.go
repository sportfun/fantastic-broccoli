package network

import (
	"fmt"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type Service struct {
	state  types.StateType
	linkId string

	logger        log.Logger
	client        *gosocketio.Client
	notifications *service.NotificationQueue
}

func (service *Service) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.state = constant.States.Started

	service.notifications = notifications
	service.logger = logger

	return nil
}

func (service *Service) Configure(props *properties.Properties) error {
	var err error

	service.linkId = props.System.LinkID
	service.client, err = gosocketio.Dial(
		gosocketio.GetUrl(string(props.System.ServerIP), int(props.System.ServerPort), props.System.ServerSSL),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		return err
	}

	initiated :=
		service.on(gosocketio.OnConnection, service.onConnectionHandler) &&
			service.on(gosocketio.OnDisconnection, service.onDisconnectionHandler) &&
			service.on(constant.Channels.Command.String(), service.onCommandChanHandler) &&
			service.emit(constant.Channels.Command.String(), object.NewCommandObject(constant.NetCommand.State, "started"))

	if !initiated {
		return fmt.Errorf("impossible to initialise network")
	}

	service.state = constant.States.Idle
	return nil
}

func (service *Service) Process() error {
	service.state = constant.States.Working
	for _, n := range service.notifications.Notifications(constant.EntityNames.Services.Network) {
		if err := service.handle(n); err != nil {
			return err
		}
	}
	service.state = constant.States.Idle
	return nil
}

func (service *Service) Stop() error {
	service.state = constant.States.Stopped
	service.client.Close()
	return nil
}

func (service *Service) Name() string {
	return constant.EntityNames.Services.Network
}

func (service *Service) State() types.StateType {
	return service.state
}
