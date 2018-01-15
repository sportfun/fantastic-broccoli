package network

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/service"
	"github.com/sportfun/gakisitor/log"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/env"
	"encoding/json"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/utils"
	"sync"
	"time"
)

const (
	IPAddress = "localhost"
	Port      = 8888
)

func TestManager_Notification(t *testing.T) {
	RegisterTestingT(t)

	bufferSocket := ""
	wg := sync.WaitGroup{}
	receiver := func(channel *gosocketio.Channel, i interface{}) { bufferSocket = marshall(i); wg.Done() }
	receivers := utils.WSReceivers{
		OnData:    receiver,
		OnCommand: receiver,
		OnError:   receiver,
	}
	go utils.SocketIOServer(receivers, Port)

	bufferLog := ""
	socket, err := gosocketio.Dial(
		gosocketio.GetUrl(IPAddress, Port, false),
		transport.GetDefaultWebsocketTransport(),
	)
	Expect(err).Should(Succeed())

	network := &Network{logger: log.NewTest(&bufferLog), notifications: service.NewNotificationQueue(), client: socket}

	testCases := []struct {
		Failed  bool
		Origin  string
		Content interface{}

		Channel string
		Data    string
		Log     string
	}{
		{Failed: true, Origin: env.NetworkServiceEntity, Content: nil, Channel: "", Data: ``, Log: `WARN	unhandled notification origin (network_manager)`},
		{Failed: true, Origin: env.CoreEntity, Content: "", Channel: "", Data: ``, Log: `WARN	unknown content type	{"packet": ""}`},

		{Origin: env.CoreEntity, Content: object.NewCommandObject("command", "a", "b"), Channel: OnCommand, Data: `{"body":{"args":["a","b"],"command":"command"},"link_id":""}`, Log: `DEBUG	notification handled	{"notification": {}}`},
		{Origin: env.ModuleServiceEntity, Content: object.NewCommandObject("command", "a"), Channel: OnCommand, Data: `{"body":{"args":["a"],"command":"command"},"link_id":""}`, Log: `DEBUG	notification handled	{"notification": {}}`},
		{Origin: env.ModuleServiceEntity, Content: object.NewCommandObject("command"), Channel: OnCommand, Data: `{"body":{"args":[],"command":"command"},"link_id":""}`, Log: `DEBUG	notification handled	{"notification": {}}`},
		{Origin: env.ModuleServiceEntity, Content: object.NewDataObject("name", 0), Channel: OnData, Data: `{"body":{"module":"name","value":0},"link_id":""}`, Log: `DEBUG	notification handled	{"notification": {}}`},
		{Origin: env.ModuleServiceEntity, Content: object.NewErrorObject("origin", gosocketio.ErrorSendTimeout), Channel: OnError, Data: `{"body":{"origin":"origin","reason":"Timeout"},"link_id":""}`, Log: `DEBUG	notification handled	{"notification": {}}`},
	}

	for _, tc := range testCases {
		utils.ReleaseIfTimeout(t, 150*time.Millisecond, func(testing.TB) {
			bufferSocket = ""
			bufferLog = ""
			if !tc.Failed {
				wg.Add(1)
				Expect(network.handle(notification.NewNotification(tc.Origin, env.NetworkServiceEntity, tc.Content))).Should(Succeed())
				wg.Wait()

				Expect(receivers).Should(HaveKey(tc.Channel))
				Expect(bufferSocket).Should(Equal(tc.Data))
			} else {
				Expect(network.handle(notification.NewNotification(tc.Origin, env.NetworkServiceEntity, tc.Content))).Should(Succeed())
			}
			Expect(bufferLog).Should(Equal(tc.Log))
		})
	}
}

func marshall(i interface{}) string { jobj, _ := json.Marshal(i); return string(jobj) }
