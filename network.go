package main

import (
	"context"
	"fmt"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

// Register the network as a worker
func init() {
	Gakisitor.RegisterWorker("network", networkTask)
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
		gosocketio.GetUrl(Gakisitor.Network.HostAddress, Gakisitor.Network.Port, Gakisitor.Network.EnableSsl),
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

// unsubscribe unsubscribes all bus handlers.
func (net *network) unsubscribe() {
	net.bus.Unsubscribe(":data", net.busDataHandler)
	net.bus.Unsubscribe(":error", net.busErrorHandler)

	net.client.Close()
}

// onConnectionHandler handles connection from socketIO and
// sent signal to the server.
func (net *network) onConnectionHandler(*gosocketio.Channel) {
	logrus.Infof("Successfully connected to %s:%d", Gakisitor.Network.HostAddress, Gakisitor.Network.Port) // LOG :: INFO - Successfully connected to {host}:{port}
	if err := net.client.Emit(
		Channels[Command],
		CommandPacket{
			Type:   "hardware",
			LinkID: Gakisitor.LinkID,
			Body: struct {
				Command string        `json:"command"`
				Args    []interface{} `json:"args"`
			}{Command: "link", Args: nil},
		},
	); err != nil {
		panic(fmt.Sprintf("Failed to send message to the server: %s", err)) // Panic - Failed to send message to the server: {error}
	}
}

// onDisconnectionHandler handles disconnection from socketIO.
func (net *network) onDisconnectionHandler(*gosocketio.Channel) {
	logrus.Infof("Disconnected from %s:%d", Gakisitor.Network.HostAddress, Gakisitor.Network.Port) // LOG :: INFO - Client disconnected
	close(net.disconnected)
}

// onCommandHandler handles command from socketIO.
func (net *network) onCommandHandler(_ *gosocketio.Channel, p CommandPacket) {
	net.bus.Publish(":instruction", p.Body.Command, bus.SyncReplyHandler(func(_ interface{}, e error) {
		if e != nil && e != bus.ErrReplyTimeout {
			logrus.Errorf("Failed to publish: %s", e) // LOG :: ERROR - Failed to publish: X
		}
	}))
}

// busDataHandler handles event from ':data' on the event bus.
func (net *network) busDataHandler(event *bus.Event, err error) {
	if err != nil {
		if err != bus.ErrSubscriberDeleted {
			logrus.Errorf("Bus handler for ':data' failed: %s", err) // LOG :: ERROR - Bus handler for ':data' failed: {error}
		}
		return
	}

	if data, valid := event.Message().(struct {
		name  string
		value interface{}
	}); !valid {
		logrus.Errorf("Invalid data type: %#v", event.Message()) // LOG :: ERROR - Invalid data type: {message}
	} else {
		if err := net.client.Emit(
			Channels[Data],
			DataPacket{
				Type:   "hardware",
				LinkId: Gakisitor.LinkID,
				Body: struct {
					Module string      `json:"module"`
					Value  interface{} `json:"value"`
				}{Module: data.name, Value: data.value},
			},
		); err != nil {
			logrus.Errorf("Failed to send message to the server: %s", err) // LOG :: ERROR - Failed to send message to the server: {error}
		}
	}
}

// busErrorHandler handles event from ':error' on the event bus.
func (net *network) busErrorHandler(event *bus.Event, err error) {
	if err != nil {
		if err != bus.ErrSubscriberDeleted {
			logrus.Errorf("Bus handler for ':error' failed: %s", err) // LOG :: ERROR - Bus handler for ':error' failed: {error}
		}
		return
	}

	if err, valid := event.Message().(struct {
		origin string
		error  error
	}); !valid {
		logrus.Errorf("Invalid error type: %v", event.Message()) // LOG :: ERROR - Invalid error type: {message}
	} else {
		if err := net.client.Emit(
			Channels[Error],
			ErrorPacket{
				Type:   "hardware",
				LinkID: Gakisitor.LinkID,
				Body: struct {
					Origin string `json:"origin"`
					Reason string `json:"reason"`
				}{Origin: err.origin, Reason: err.error.Error()},
			},
		); err != nil {
			logrus.Errorf("Failed to send message to the server: %s", err) // LOG :: ERROR - Failed to send message to the server: {error}
		}
	}
}
