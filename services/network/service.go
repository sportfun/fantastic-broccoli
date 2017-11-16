package network

import (
	"fantastic-broccoli/common/types/notification/object"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/properties"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"go.uber.org/zap"
)

type Service struct {
	state  int
	linkId string

	logger        *zap.Logger
	client        *gosocketio.Client
	notifications *service.NotificationQueue
}

func (s *Service) Start(notifications *service.NotificationQueue, logger *zap.Logger) error {
	s.state = constant.Started

	s.notifications = notifications
	s.logger = logger

	return nil
}

func (s *Service) Configure(props *properties.Properties) error {
	var err error

	s.linkId = props.System.LinkID
	s.client, err = gosocketio.Dial(
		gosocketio.GetUrl(string(props.System.ServerIP), int(props.System.ServerPort), props.System.ServerSSL),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		return err
	}

	initiated :=
		s.on(gosocketio.OnConnection, s.onConnectionHandler) &&
			s.on(gosocketio.OnDisconnection, s.onDisconnectionHandler) &&
			s.on(constant.CommandChan, s.onCommandChanHandler) &&
			s.emit(constant.CommandChan, object.NewNetworkObject(constant.CommandState, "started"))

	if !initiated {
		return fmt.Errorf("impossible to initialise network")
	}

	s.state = constant.Idle
	return nil
}

func (s *Service) Process() error {
	s.state = constant.Working
	for _, n := range s.notifications.Notifications(constant.NetworkService) {
		if err := s.notificationHandler(n); err != nil {
			return err
		}
	}
	s.state = constant.Idle
	return nil
}

func (s *Service) Stop() error {
	s.state = constant.Stopped
	s.client.Close()
	return nil
}

func (s *Service) Name() string {
	return "Network"
}

func (s *Service) State() int {
	return s.state
}
