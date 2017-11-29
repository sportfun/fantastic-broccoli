package network

import (
	"testing"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	def "github.com/xunleii/fantastic-broccoli/utils/default"
	"sync"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"time"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/mitchellh/mapstructure"
	"github.com/xunleii/fantastic-broccoli/utils"
	"github.com/xunleii/fantastic-broccoli/log"
)

type watcher struct {
	LinkId string
	Data   interface{}
}

var (
	infoClientConnected    = log.NewArgumentBinder("[Server] New client connected, client id is '%s' (%s)")
	infoClientDisconnected = log.NewArgumentBinder("[Server] %s (%s) disconnected")
	infoSendPacket         = log.NewArgumentBinder("[Server] Send packet to client '%s' (%v)")

	infoCommandHandled = log.NewArgumentBinder("[Server] Command handled")
	infoDataHandled    = log.NewArgumentBinder("[Server] Data handled")
	infoErrorHandled   = log.NewArgumentBinder("[Server] Error handled")

	errorUnknownPacketBody = log.NewArgumentBinder("[Server] Message handled, but unknown web packet body type")
	errorUnknownPacketType = log.NewArgumentBinder("[Server] Message handled, but unknown packet type")
)

func onConnectionReceiver(client *gosocketio.Channel, args interface{}, logger log.Logger, waitGroup *sync.WaitGroup) {
	logger.Info(infoClientConnected.Bind(client.Id(), client.Ip()))

	packets := []webPacket{
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.Link)},
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.StartSession)},
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.EndSession)},
	}

	for _, packet := range packets {
		time.Sleep(500 * time.Millisecond)
		logger.Info(infoSendPacket.Bind(client.Id(), packet))
		client.Emit(constant.Channels.Command.String(), packet)
	}

	time.Sleep(2 * time.Second)
	waitGroup.Done()
}

func onDisconnectionReceiver(client *gosocketio.Channel, args interface{}, logger log.Logger) {
	logger.Info(infoClientDisconnected.Bind(client.Id(), client.Ip()))
}

func onGenericReceiver(
	args interface{},
	l log.Logger, watch *watcher,
	receiver func(log.Logger, webPacket, *watcher) error,
) {
	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		watch.LinkId = web.LinkId

		if err := receiver(l, web, watch); err != nil {
			l.Warn(errorUnknownPacketBody.More("packer_body", web.Body).More("reason", err.Error()))
		}

	default:
		l.Warn(errorUnknownPacketType.More("packet", args))
	}
}

func onCommandReceiver(logger log.Logger, web webPacket, watch *watcher) error {
	var netObj object.CommandObject
	if err := mapstructure.Decode(web.Body, &netObj); err != nil {
		return err
	}
	watch.Data = netObj

	logger.Info(infoCommandHandled.More("link_id", web.LinkId).More("command", netObj.Command.String()).More("args", netObj.Args))
	return nil
}

func onDataReceiver(logger log.Logger, web webPacket, watch *watcher) error {
	var dataObj object.DataObject
	if err := mapstructure.Decode(web.Body, &dataObj); err != nil {
		return err
	}
	watch.Data = dataObj

	logger.Info(infoDataHandled.More("link_id", web.LinkId).More("module", dataObj.Module).More("value", dataObj.Value))
	return nil
}

func onErrorReceiver(logger log.Logger, web webPacket, watch *watcher) error {
	var errObj object.ErrorObject
	if err := mapstructure.Decode(web.Body, &errObj); err != nil {
		return err
	}
	watch.Data = errObj

	logger.Info(infoErrorHandled.More("link_id", web.LinkId).More("origin", errObj.Origin).More("reason", errObj.Reason))
	return nil
}

// TODO: Remove logger
func TestService(t *testing.T) {
	l := log.NewLogger.Dev(nil)
	pkt := watcher{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	r := map[string]func(*gosocketio.Channel, interface{}){
		gosocketio.OnConnection:            func(c *gosocketio.Channel, a interface{}) { onConnectionReceiver(c, a, l, &wg) },
		gosocketio.OnDisconnection:         func(c *gosocketio.Channel, a interface{}) { onDisconnectionReceiver(c, a, l) },
		constant.Channels.Command.String(): func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onCommandReceiver) },
		constant.Channels.Data.String():    func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onDataReceiver) },
		constant.Channels.Error.String():   func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onErrorReceiver) },
	}

	go def.SocketIOServer(r)

	s := Service{}
	p := properties.Properties{
		System: properties.SystemDefinition{
			LinkID:     "70ed3820-d487-42b9-92a8-ae9cbf55918c",
			ServerIP:   "localhost",
			ServerPort: 8080,
			ServerSSL:  false,
		},
	}
	q := service.NewNotificationQueue()

	s.Start(q, l)
	if err := s.Configure(&p); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
	coreNotifications := q.Notifications(constant.EntityNames.Core)
	moduleNotifications := q.Notifications(constant.EntityNames.Services.Module)

	utils.AssertEquals(t, 1, len(coreNotifications))
	utils.AssertEquals(t, 2, len(moduleNotifications))

	objNet := coreNotifications[0].Content().(object.CommandObject)
	utils.AssertEquals(t, constant.NetCommand.Link, objNet.Command)
	objNet = moduleNotifications[0].Content().(object.CommandObject)
	utils.AssertEquals(t, constant.NetCommand.StartSession, objNet.Command)
	objNet = moduleNotifications[1].Content().(object.CommandObject)
	utils.AssertEquals(t, constant.NetCommand.EndSession, objNet.Command)

	q.Notify(notification.NewNotification(constant.EntityNames.Services.Module, constant.EntityNames.Services.Network, object.NewDataObject("Example", "256")))
	s.Process()
	time.Sleep(time.Second)
	do := pkt.Data.(object.DataObject)
	utils.AssertEquals(t, p.System.LinkID, pkt.LinkId)
	utils.AssertEquals(t, "Example", do.Module)
	utils.AssertEquals(t, "256", do.Value)

	q.Notify(notification.NewNotification(constant.EntityNames.Services.Module, constant.EntityNames.Services.Network, object.NewErrorObject("Origin", fmt.Errorf("error"))))
	s.Process()
	time.Sleep(time.Second)
	de := pkt.Data.(object.ErrorObject)
	utils.AssertEquals(t, p.System.LinkID, pkt.LinkId)
	utils.AssertEquals(t, "Origin", de.Origin)
	utils.AssertEquals(t, "error", de.Reason)

	s.Stop()
	time.Sleep(time.Second)
}
