package gakisitor

import (
	"context"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/sportfun/gakisitor/event"
)

const (
	onConnection    = gosocketio.OnConnection
	onDisconnection = gosocketio.OnDisconnection
	onCommand       = "command"
)

type network struct {
	client *gosocketio.Client
}

func init() {
	Scheduler.RegisterWorker("network", networkTask)
}

func networkTask(ctx context.Context, bus *event.Bus) error {
	var err error
	var net network

	if net.client, err = gosocketio.Dial(
		gosocketio.GetUrl(string(Profile.Network.HostAddress), Profile.Network.Port, Profile.Network.EnableSsl),
		transport.GetDefaultWebsocketTransport(),
	); err != nil {
		return err
	}

	var onConnectionHandler interface{} = nil;
	var onDisconnectionHandler interface{} = nil;
	var onCommandHandler interface{} = nil;
	if err = net.client.On(onConnection, onConnectionHandler); err != nil {
		return err
	}
	if err = net.client.On(onDisconnection, onDisconnectionHandler); err != nil {
		return err
	}
	if err = net.client.On(onCommand, onCommandHandler); err != nil {
		return err
	}

	if err := net.client.Emit(onCommand, nil); err != nil {
		return err
	}

	var busDataHandler event.EventConsumer = nil
	var busErrorHandler event.EventConsumer = nil

	if err = bus.Subscribe(":data", busDataHandler); err != nil {
		return err
	}
	if err = bus.Subscribe(":error", busErrorHandler); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	}
}
