package network

import (
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
	"github.com/sportfun/gakisitor/utils"
	"testing"
	"time"
)

func TestNetwork_SocketIO_On(t *testing.T) {
	RegisterTestingT(t)

	// Configure server
	receivers := utils.WSReceivers{
		"resend": func(channel *gosocketio.Channel, i interface{}) {
			o := i.(map[string]interface{})
			channel.Emit(o["method"].(string), o["data"])
		},
	}
	port := int32(Port + 0x1F)
	go utils.SocketIOServer(receivers, port)

	// Configure socket
	socket, err := gosocketio.Dial(
		gosocketio.GetUrl(IPAddress, int(port), false),
		transport.GetDefaultWebsocketTransport(),
	)
	Expect(err).Should(Succeed())

	// Configure manager
	logs := ""
	queue := service.NewNotificationQueue()
	network := &Network{logger: log.NewTest(&logs), notifications: queue, linkId: "0000"}

	// Test cases
	wg := utils.ConditionalWaitGroup{}
	expectationHandler := func(expectation interface{}) func(*gosocketio.Channel, interface{}) {
		return func(channel *gosocketio.Channel, i interface{}) { Expect(i).Should(Equal(expectation)) }
	}
	testCases := []struct {
		Client  *gosocketio.Client
		Handler utils.WSReceiver

		Data interface{}

		Failed       bool
		NeedSync     bool
		Log          string
		Notification *notification.Notification
	}{
		// Failure cases
		{Failed: true, Log: "ERROR	socket.io client not initialised"},
		{Failed: true, Client: socket, Log: "ERROR	failed to create channel handler: f is not function"},

		// Simple cases
		{NeedSync: true, Client: socket, Handler: expectationHandler(0.), Data: 0.},
		{NeedSync: true, Client: socket, Handler: expectationHandler("#2"), Data: "#2"},

		// Receiver cases
		{NeedSync: true, Client: socket, Handler: network.onConnectionHandler, Data: "", Log: "INFO	successfully connected to the server"},

		{NeedSync: true, Client: socket, Handler: network.onDisconnectionHandler, Data: "", Log: "DEBUG	disconnection handled", Notification: notification.NewNotification(env.NetworkServiceEntity, env.CoreEntity, env.RestartServiceCmd)},

		{NeedSync: true, Client: socket, Handler: network.onCommandChanHandler, Data: "invalid packet", Log: "WARN	unknown packet type"},
		{NeedSync: true, Client: socket, Handler: network.onCommandChanHandler, Data: websocket{}, Log: "WARN	unknown packet type"},
		{NeedSync: true, Client: socket, Handler: network.onCommandChanHandler, Data: websocket{LinkId: network.linkId}, Log: "WARN	unknown web packet body type"},
		{NeedSync: true, Client: socket, Handler: network.onCommandChanHandler, Data: websocket{LinkId: network.linkId, Body: object.NewCommandObject("NONE")}, Log: "WARN	unknown command 'NONE'"},
		{NeedSync: true, Client: socket, Handler: network.onCommandChanHandler, Data: websocket{LinkId: network.linkId, Body: object.NewCommandObject(env.LinkCmd)}, Log: "DEBUG	valid command handled", Notification: notification.NewNotification(env.NetworkServiceEntity, env.CoreEntity, *object.NewCommandObject(env.LinkCmd))},
	}

	for id, tc := range testCases {
		logs = ""
		method := fmt.Sprintf("#%d", id)

		network.client = tc.Client
		utils.ReleaseIfTimeout(t, 250*time.Millisecond, func(testing.TB) {
			// Create handler
			if tc.Handler == nil {
				Expect(network.on(method, nil)).ShouldNot(Equal(tc.Failed))
			} else {
				Expect(network.on(method, func(c *gosocketio.Channel, i interface{}) { tc.Handler(c, i); wg.DoneIf(tc.NeedSync) })).Should(Equal(!tc.Failed))

				wg.AddIf(1, tc.NeedSync)
				socket.Emit("resend", map[string]interface{}{
					"method": method,
					"data":   tc.Data,
				})
				wg.WaitIf(tc.NeedSync)
			}

			// Compare notification (if needed)
			if tc.Notification != nil {
				Expect(queue.Notifications(tc.Notification.To())).Should(ConsistOf(tc.Notification))
			}

			// Compare logs
			Expect(logs).Should(MatchRegexp(tc.Log))
		})
	}
}

func TestNetwork_SocketIO_Emit(t *testing.T) {
	RegisterTestingT(t)

	// Configure server
	safeData := utils.NewVolatile("")
	wg := utils.ConditionalWaitGroup{}
	receivers := utils.WSReceivers{
		"method": func(channel *gosocketio.Channel, i interface{}) {
			o := i.(map[string]interface{})
			safeData.Set(fmt.Sprintf("%v", o["body"]))
			wg.Done()
		},
	}
	port := int32(Port + 0x2F)
	go utils.SocketIOServer(receivers, port)

	// Configure socket
	socket, err := gosocketio.Dial(
		gosocketio.GetUrl(IPAddress, int(port), false),
		transport.GetDefaultWebsocketTransport(),
	)
	Expect(err).Should(Succeed())

	// Configure network
	logs := ""
	network := &Network{logger: log.NewTest(&logs), notifications: service.NewNotificationQueue()}
	testCases := []struct {
		Client *gosocketio.Client

		Data interface{}

		Failed   bool
		NeedSync bool
		Log      string
	}{
		{Failed: true, Log: "ERROR	socket.io client not initialised"},
		{NeedSync: true, Client: socket, Data: "succeed"},
	}

	for _, tc := range testCases {
		logs = ""

		network.client = tc.Client
		utils.ReleaseIfTimeout(t, 250*time.Millisecond, func(testing.TB) {
			// Emit data
			wg.AddIf(1, tc.NeedSync)
			Expect(network.emit("method", tc.Data)).ShouldNot(Equal(tc.Failed))
			wg.WaitIf(tc.NeedSync)

			// Compare handled data
			if tc.Data != nil {
				Expect(safeData.Get()).Should(Equal(fmt.Sprintf("%v", tc.Data)))
			}

			// Compare logs
			Expect(logs).Should(Equal(tc.Log))
		})
	}
}
