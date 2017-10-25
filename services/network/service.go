package network

import (
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/model"
	"go.uber.org/zap"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"fantastic-broccoli/common/types/notification/object"
	"fmt"
)

type Service struct {
	state int

	notifications *service.notificationQueue
	logger        *zap.Logger
	client        *gosocketio.Client
	linkId        string
	messages      []*object.NetworkObject
}

func (s *Service) Start(q *service.notificationQueue, l *zap.Logger) error {
	s.state = constant.STARTED

	s.notifications = q
	s.logger = l

	return nil
}

func (s *Service) Configure(props *model.Properties) error {
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
		s.On(gosocketio.OnConnection, s.onConnectionHandler) &&
			s.On(gosocketio.OnDisconnection, s.onDisconnectionHandler) &&
			s.On(constant.CommandChan, s.onCommandChanHandler) &&
			s.Emit(constant.CommandChan, object.NewNetworkObject(constant.CommandState, "started"))

	if !initiated {
		//TODO: Write error
		return fmt.Errorf("")
	}

	s.state = constant.IDLE
	return nil
}

func (s *Service) Process() error {
	s.state = constant.WORKING
	for _, m := range s.messages {
		s.messageHandler(m)
	}

	for _, n := range s.notifications.Notifications(constant.NetworkService) {
		if err := s.notificationHandler(n); err != nil {
			return err
		}
	}
	s.state = constant.IDLE
	return nil
}

func (s *Service) Stop() error {
	s.client.Close()
	s.state = constant.STOPPED
	return nil
}

func (s *Service) Name() string {
	return "Network"
}

func (s *Service) State() int {
	return s.state
}
