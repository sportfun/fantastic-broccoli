package main

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/pkg/errors"
	"github.com/sportfun/gakisitor/event/bus"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

func init() {
	Scheduler.RegisterWorker("network", networkTask)
}

const (
	onConnection    = gosocketio.OnConnection
	onDisconnection = gosocketio.OnDisconnection
	onCommand       = "command"
)

type network struct {
	client       *gosocketio.Client
	bus          *bus.Bus
	disconnected chan struct{}
}

var errNetworkDisconnected = errors.New("client disconnected")

func networkTask(ctx context.Context, bus *bus.Bus) error {
	var err error
	var net network

	net.disconnected = make(chan struct{})
	net.bus = bus
	if net.client, err = gosocketio.Dial(
		gosocketio.GetUrl(Profile.Network.HostAddress, Profile.Network.Port, Profile.Network.EnableSsl),
		transport.GetDefaultWebsocketTransport(),
	); err != nil {
		return err
	}

	defer net.unsubscribe()

	for _, step := range []func() error{
		func() error { return net.client.On(onConnection, net.onConnectionHandler) },
		func() error { return net.client.On(onDisconnection, net.onDisconnectionHandler) },
		func() error { return net.client.On(onCommand, net.onCommandHandler) },

		func() error { return bus.Subscribe(":data", net.busDataHandler) },
		func() error { return bus.Subscribe(":error", net.busErrorHandler) },

		func() error { return net.client.Emit(onCommand, nil) },
	} {
		if err = step(); err != nil {
			return err
		}
	}

	select {
	case <-ctx.Done():
		return nil
	case <-net.disconnected:
		return errNetworkDisconnected
	}
}

func (net *network) unsubscribe() {
	net.bus.Unsubscribe(":data", net.busDataHandler)
	net.bus.Unsubscribe(":error", net.busErrorHandler)

	net.client.Close()
}

func (net *network) onConnectionHandler(*gosocketio.Channel) {
	log.Infof("Successfully connected to %s:%d", Profile.Network.HostAddress, Profile.Network.Port) //LOG :: INFO - Successfully connected to {host}:{port}
}

func (net *network) onDisconnectionHandler(*gosocketio.Channel) {
	log.Infof("Disconnected from %s:%d", Profile.Network.HostAddress, Profile.Network.Port) //LOG :: INFO - Client disconnected
	close(net.disconnected)
}

func (net *network) onCommandHandler(_ *gosocketio.Channel, p CommandPacket) {
	net.bus.Publish(":instruction", p.Body.Command, bus.SyncReplyHandler(func(_ interface{}, e error) {
		if e != nil && e != bus.ErrReplyTimeout {
			log.Errorf("Failed to publish: %s", e) //LOG :: ERROR - Failed to publish: X
		}
	}))
}

func (net *network) busDataHandler(event *bus.Event, err error) {
	if err != nil {
		if err != bus.ErrSubscriberClosed {
			log.Errorf("Bus handler for ':data' failed: %s", err) //LOG :: ERROR - Bus handler for ':data' failed: {error}
		}
		return
	}

	if data, valid := event.Message().(struct {
		name  string
		value interface{}
	}); !valid {
		log.Errorf("Invalid data type: %#v", event.Message()) //LOG :: ERROR - Invalid data type: {message}
	} else {
		if err := net.client.Emit(
			Channels[Data],
			DataPacket{
				LinkId: Profile.LinkID,
				Body: struct {
					Module string      `json:"module"`
					Value  interface{} `json:"value"`
				}{Module: data.name, Value: data.value},
			},
		); err != nil {
			log.Errorf("Failed to send message to the server: %s", err) //LOG :: ERROR - Failed to send message to the server: {error}
		}
	}
}

func (net *network) busErrorHandler(event *bus.Event, err error) {
	if err != nil {
		if err != bus.ErrSubscriberClosed {
			log.Errorf("Bus handler for ':error' failed: %s", err) //LOG :: ERROR - Bus handler for ':data' failed: {error}
		}
		return
	}

	if error, valid := event.Message().(struct {
		origin string
		error  error
	}); !valid {
		log.Errorf("Invalid error type: %v", event.Message()) //LOG :: ERROR - Invalid error type: {message}
	} else {
		if err := net.client.Emit(
			Channels[Error],
			ErrorPacket{
				LinkId: Profile.LinkID,
				Body: struct {
					Origin string `json:"origin"`
					Reason string `json:"reason"`
				}{Origin: error.origin, Reason: error.error.Error()},
			},
		); err != nil {
			log.Errorf("Failed to send message to the server: %s", err) //LOG :: ERROR - Failed to send message to the server: {error}
		}
	}
}
