package gakisitor

import (
	"context"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/pkg/errors"
	"github.com/sportfun/gakisitor/event"
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
		gosocketio.GetUrl(string(Profile.Network.HostAddress), Profile.Network.Port, Profile.Network.EnableSsl),
		transport.GetDefaultWebsocketTransport(),
	); err != nil {
		return err
	}

	defer net.unsubscribe()

	for _, fnc := range []func() error{
		func() error { return net.client.On(onConnection, net.onConnectionHandler) },
		func() error { return net.client.On(onDisconnection, net.onDisconnectionHandler) },
		func() error { return net.client.On(onCommand, net.onCommandHandler) },

		func() error { return bus.Subscribe(":data", net.busDataHandler) },
		func() error { return bus.Subscribe(":status", net.busStatusHandler) },
		func() error { return bus.Subscribe(":error", net.busErrorHandler) },

		func() error { return net.client.Emit(onCommand, nil) },
	} {
		if err = fnc(); err != nil {
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

func (net *network) onConnectionHandler()    {}
func (net *network) onDisconnectionHandler() {}
func (net *network) onCommandHandler()       {}

func (net *network) busDataHandler(event event.Event, err error)   {}
func (net *network) busStatusHandler(event event.Event, err error) {}
func (net *network) busErrorHandler(event event.Event, err error)  {}
