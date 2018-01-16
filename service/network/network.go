package network

import (
	"fmt"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/kernel"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
)

type Network struct {
	state  byte
	linkId string

	logger        log.Logger
	client        *gosocketio.Client
	notifications *service.NotificationQueue
}

func init() {
	kernel.RegisterService(&Network{})
}

func (service *Network) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.state = env.StartedState

	service.notifications = notifications
	service.logger = logger

	return nil
}

func (service *Network) Configure(config *config.GAkisitorConfig) error {
	var err error

	if config == nil {
		service.state = env.PanicState
		return fmt.Errorf("configuration not defined")
	}

	service.linkId = config.System.LinkID
	service.client, err = gosocketio.Dial(
		gosocketio.GetUrl(string(config.System.ServerIP), int(config.System.ServerPort), config.System.ServerSSL),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		service.state = env.PanicState
		return err
	}

	initiated :=
		service.on(gosocketio.OnConnection, service.onConnectionHandler) &&
			service.on(gosocketio.OnDisconnection, service.onDisconnectionHandler) &&
			service.on(OnCommand, service.onCommandChanHandler) &&
			service.emit(OnCommand, object.NewCommandObject(env.StateCmd, "started"))

	if !initiated {
		service.state = env.PanicState
		return fmt.Errorf("impossible to initialise network")
	}

	service.state = env.IdleState
	return nil
}

func (service *Network) Process() error {
	service.state = env.WorkingState
	for _, n := range service.notifications.Notifications(env.NetworkServiceEntity) {
		if err := service.handle(n); err != nil {
			return err
		}
	}
	service.state = env.IdleState
	return nil
}

func (service *Network) Stop() error {
	service.state = env.StoppedState

	if service.client != nil {
		service.client.Close()
	}
	return nil
}

func (service *Network) Name() string {
	return env.NetworkServiceEntity
}

func (service *Network) State() byte {
	return service.state
}
