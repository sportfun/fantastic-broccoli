package main

import (
	"context"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/pkg/errors"
	"github.com/sportfun/main/event"
	. "github.com/sportfun/main/protocol/v1.0"
	"log"
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
	bus          *event.Bus
	disconnected chan struct{}
}

var errNetworkDisconnected = errors.New("client disconnected")

func networkTask(ctx context.Context, bus *event.Bus) error {
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
	//TODO: LOG :: INFO - Successfully connected to X:Y
	log.Printf("{network}[INFO]			Successfully connected to %s:%d", Profile.Network.HostAddress, Profile.Network.Port)
}

func (net *network) onDisconnectionHandler(*gosocketio.Channel) {
	//TODO: LOG :: INFO - Client disconnected
	log.Printf("{network}[INFO]			Disconnected from %s:%d", Profile.Network.HostAddress, Profile.Network.Port)
	close(net.disconnected)
}

func (net *network) onCommandHandler(_ *gosocketio.Channel, p CommandPacket) {
	if _, exists := Instructions[p.Body.Command]; !exists {
		//TODO: LOG :: ERROR - Unknown instruction X
		log.Printf("{network}[ERROR]				Unknown instruction '%s'", p.Body.Command)
	} else {
		net.bus.Publish(":instruction", p.Body, event.SyncReplyHandler(func(_ interface{}, e error) {
			if e != nil && e != event.ErrReplyTimeout {
				//TODO: LOG :: ERROR - Failed to publish: X
				log.Printf("{network}[ERROR]			Failed to publish: %s", e)
			}
		}))
	}
}

func (net *network) busDataHandler(event *event.Event, err error) {
	if data, valid := event.Message().(struct {
		name  string
		value interface{}
	}); !valid {
		//TODO: LOG :: ERROR - Invalid plugin data: X
		log.Printf("{network}[ERROR]			Invalid data type: %v", event.Message())
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
			//TODO: LOG :: ERROR - Failed to send message to the server: X
			log.Printf("{network}[ERROR]			Failed to send message to the server: %s", err)
		}
	}
}

func (net *network) busErrorHandler(event *event.Event, err error) {
	if error, valid := event.Message().(struct {
		origin string
		error  error
	}); !valid {
		//TODO: LOG :: ERROR - Invalid error type: X
		log.Printf("{network}[ERROR]			Invalid error type: %v", event.Message())
	} else {
		if err := net.client.Emit(
			Channels[Data],
			ErrorPacket{
				LinkId: Profile.LinkID,
				Body: struct {
					Origin string `json:"origin"`
					Reason string `json:"reason"`
				}{Origin: error.origin, Reason: error.error.Error()},
			},
		); err != nil {
			//TODO: LOG :: ERROR - Failed to send message to the server: X
			log.Printf("{network}[ERROR]			Failed to send message to the server: %s", err)
		}
	}
}
