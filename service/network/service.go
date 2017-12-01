package network

import (
	"fmt"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/notification/object"
	"github.com/xunleii/fantastic-broccoli/service"
	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/kernel"
)

type Service struct {
	state  byte
	linkId string

	logger        log.Logger
	client        *gosocketio.Client
	notifications *service.NotificationQueue
}

func init() {
	kernel.RegisterService(&Service{})
}

func (service *Service) Start(notifications *service.NotificationQueue, logger log.Logger) error {
	service.state = env.StartedState

	service.notifications = notifications
	service.logger = logger

	return nil
}

func (service *Service) Configure(config *config.GAkisitorConfig) error {
	var err error

	service.linkId = config.System.LinkID
	service.client, err = gosocketio.Dial(
		gosocketio.GetUrl(string(config.System.ServerIP), int(config.System.ServerPort), config.System.ServerSSL),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		return err
	}

	initiated :=
		service.on(gosocketio.OnConnection, service.onConnectionHandler) &&
			service.on(gosocketio.OnDisconnection, service.onDisconnectionHandler) &&
			service.on(OnCommand, service.onCommandChanHandler) &&
			service.emit(OnCommand, object.NewCommandObject(env.StateCmd, "started"))

	if !initiated {
		return fmt.Errorf("impossible to initialise network")
	}

	service.state = env.IdleState
	return nil
}

func (service *Service) Process() error {
	service.state = env.WorkingState
	for _, n := range service.notifications.Notifications(env.NetworkServiceEntity) {
		if err := service.handle(n); err != nil {
			return err
		}
	}
	service.state = env.IdleState
	return nil
}

func (service *Service) Stop() error {
	service.state = env.StoppedState

	if service.client != nil {
		service.client.Close()
	}
	return nil
}

func (service *Service) Name() string {
	return env.NetworkServiceEntity
}

func (service *Service) State() byte {
	return service.state
}
