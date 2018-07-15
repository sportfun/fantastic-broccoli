package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sportfun/gakisitor/event/bus"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/onsi/gomega"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

type wsServerMock struct {
	*testing.T
	*gosocketio.Server
}

func (s *wsServerMock) run(ctx context.Context) *sync.WaitGroup {
	Gakisitor.Network.HostAddress = "127.0.0.1"
	Gakisitor.Network.Port = 7482
	wg := &sync.WaitGroup{}
	wg.Add(1)

	srv := &http.Server{Addr: ":7482"}
	srv.Handler = s.Server

	go func() {
		srv.ListenAndServe()
	}()

	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(nil)
			wg.Done()
		}
	}()

	return wg
}

func (s *wsServerMock) bind(method string, bind *interface{}) *sync.WaitGroup {
	wg := &sync.WaitGroup{}

	if len(method) != 0 {
		wg.Add(1)
		s.On(method, func(c *gosocketio.Channel, actual map[string]interface{}) { *bind = actual; wg.Done() })
	}
	return wg
}

func TestNetwork_on_Dis_ConnectionHandler(t *testing.T) {
	gomega.RegisterTestingT(t)
	Gakisitor.LinkID = "TestNetwork_on_Dis_ConnectionHandler"

	cases := []struct {
		method  string
		handler func(net *network)
		expect  func(net *network, actual interface{})
	}{
		{
			method:  Channels[Command],
			handler: func(net *network) { net.onConnectionHandler(nil) },
			expect: func(_ *network, actual interface{}) {
				expected := map[string]interface{}{
					"type":    "hardware",
					"link_id": "TestNetwork_on_Dis_ConnectionHandler",
					"body": map[string]interface{}{
						"command": "link",
						"args":    nil,
					},
				}

				gomega.Expect(actual).Should(gomega.Equal(expected))
			},
		},
		{
			handler: func(net *network) { net.onDisconnectionHandler(nil) },
			expect: func(net *network, _ interface{}) {
				gomega.Expect(net.disconnected).Should(gomega.BeClosed())
			},
		},
	}

	for _, test := range cases {
		var err error
		var bind interface{}

		ctx, shutdown := context.WithCancel(context.Background())
		s := &wsServerMock{Server: gosocketio.NewServer(transport.GetDefaultWebsocketTransport())}
		wb := s.bind(test.method, &bind)
		wx := s.run(ctx)
		time.Sleep(500 * time.Millisecond)

		net := &network{
			disconnected: make(chan struct{}),
			bus:          nil,
		}
		net.client, err = gosocketio.Dial(
			gosocketio.GetUrl(Gakisitor.Network.HostAddress, Gakisitor.Network.Port, Gakisitor.Network.EnableSsl),
			transport.GetDefaultWebsocketTransport(),
		)
		gomega.Expect(err).Should(gomega.Succeed())

		test.handler(net)
		wb.Wait()
		test.expect(net, bind)

		shutdown()
		wx.Wait()
	}
}

func TestNetwork_onCommandHandler(t *testing.T) {
	gomega.RegisterTestingT(t)
	Gakisitor.LinkID = "TestNetwork_onCommandHandler"
	var message string

	b := bus.New()
	net := &network{
		disconnected: make(chan struct{}),
		bus:          b,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	b.Subscribe(":instruction", func(event *bus.Event, err error) {
		message = event.Message().(string)
		wg.Done()
	})

	net.onCommandHandler(nil, CommandPacket{
		LinkID: Gakisitor.LinkID,
		Type:   "hardware",
		Body: struct {
			Command string        `json:"command"`
			Args    []interface{} `json:"args"`
		}{"INSTRUCTION", nil},
	})

	wg.Wait()
	gomega.Expect(message).Should(gomega.Equal("INSTRUCTION"))
}

func TestNetwork_busHandler(t *testing.T) {
	gomega.RegisterTestingT(t)
	Gakisitor.LinkID = "TestNetwork_busHandler"
	waitOrTimeout := func(wg *sync.WaitGroup) {
		c := make(chan struct{})
		go func() {
			defer close(c)
			wg.Wait()
		}()
		select {
		case <-c:
			return
		case <-time.After(250 * time.Millisecond):
			return
		}
	}

	cases := []struct {
		method  string
		prepare func(net *network)
		exec    func(net *network)
		expect  func(net *network, actual interface{})
	}{
		{
			prepare: func(net *network) { net.bus.Subscribe(":data", net.busDataHandler) },
			exec:    func(net *network) { net.busDataHandler(nil, fmt.Errorf("nothing")) },
			expect:  func(net *network, actual interface{}) { gomega.Expect(actual).Should(gomega.BeNil()) },
		},
		{
			prepare: func(net *network) { net.bus.Subscribe(":data", net.busDataHandler) },
			exec:    func(net *network) { net.bus.Publish(":data", struct{}{}, nil) },
			expect:  func(net *network, actual interface{}) { gomega.Expect(actual).Should(gomega.BeNil()) },
		},
		{
			method:  Channels[Data],
			prepare: func(net *network) { net.bus.Subscribe(":data", net.busDataHandler) },
			exec: func(net *network) {
				net.bus.Publish(":data", struct {
					name  string
					value interface{}
				}{"...", 5.8}, nil)
			},
			expect: func(_ *network, actual interface{}) {
				expected := map[string]interface{}{
					"type":    "hardware",
					"link_id": "TestNetwork_busHandler",
					"body": map[string]interface{}{
						"module": "...",
						"value":  5.8,
					},
				}

				gomega.Expect(actual).Should(gomega.Equal(expected))
			},
		},

		{
			prepare: func(net *network) { net.bus.Subscribe(":error", net.busErrorHandler) },
			exec:    func(net *network) { net.busErrorHandler(nil, fmt.Errorf("nothing")) },
			expect:  func(net *network, actual interface{}) { gomega.Expect(actual).Should(gomega.BeNil()) },
		},
		{
			prepare: func(net *network) { net.bus.Subscribe(":error", net.busErrorHandler) },
			exec:    func(net *network) { net.bus.Publish(":error", struct{}{}, nil) },
			expect:  func(net *network, actual interface{}) { gomega.Expect(actual).Should(gomega.BeNil()) },
		},
		{
			method:  Channels[Error],
			prepare: func(net *network) { net.bus.Subscribe(":error", net.busErrorHandler) },
			exec: func(net *network) {
				net.bus.Publish(":error", struct {
					origin string
					error  error
				}{"...", fmt.Errorf("TestNetwork_busHandler::Error::Reason")}, nil)
			},
			expect: func(_ *network, actual interface{}) {
				expected := map[string]interface{}{
					"type":    "hardware",
					"link_id": "TestNetwork_busHandler",
					"body": map[string]interface{}{
						"origin": "...",
						"reason": "TestNetwork_busHandler::Error::Reason",
					},
				}

				gomega.Expect(actual).Should(gomega.Equal(expected))
			},
		},
	}

	for _, test := range cases {
		var err error
		var bind interface{}

		ctx, shutdown := context.WithCancel(context.Background())
		s := &wsServerMock{Server: gosocketio.NewServer(transport.GetDefaultWebsocketTransport())}
		wb := s.bind(test.method, &bind)
		wx := s.run(ctx)
		time.Sleep(500 * time.Millisecond)

		net := &network{
			disconnected: make(chan struct{}),
			bus:          bus.New(),
		}
		net.client, err = gosocketio.Dial(
			gosocketio.GetUrl(Gakisitor.Network.HostAddress, Gakisitor.Network.Port, Gakisitor.Network.EnableSsl),
			transport.GetDefaultWebsocketTransport(),
		)
		gomega.Expect(err).Should(gomega.Succeed())

		test.prepare(net)
		test.exec(net)
		waitOrTimeout(wb)
		test.expect(net, bind)

		shutdown()
		wx.Wait()
	}
}
