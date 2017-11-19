package network

import (
	"testing"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/utils"
	"sync"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"time"
	"fmt"
	"github.com/graarh/golang-socketio"
	"go.uber.org/zap"
	"github.com/mitchellh/mapstructure"
)

type watcher struct {
	LinkId string
	Data   interface{}
}

func onConnectionReceiver(client *gosocketio.Channel, args interface{}, logger *zap.Logger, waitGroup *sync.WaitGroup) {
	logger.Info(fmt.Sprintf("[Server] New client connected, client id is '%s' (%s)", client.Id(), client.Ip()))

	packets := []webPacket{
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.Link)},
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.StartSession)},
		{LinkId: "", Body: *object.NewCommandObject(constant.NetCommand.EndSession)},
	}

	for _, packet := range packets {
		time.Sleep(500 * time.Millisecond)
		logger.Info(fmt.Sprintf("[Server] Send packet to client '%s' (%v)", client.Id(), packet))
		client.Emit(constant.Channels.Command, packet)
	}

	time.Sleep(2 * time.Second)
	waitGroup.Done()
}

func onDisconnectionReceiver(client *gosocketio.Channel, args interface{}, logger *zap.Logger) {
	logger.Info(fmt.Sprintf("[Server] %s (%s) disconnected", client.Id(), client.Ip()))
}

func onGenericReceiver(
	args interface{},
	l *zap.Logger, watch *watcher,
	receiver func(*zap.Logger, webPacket, *watcher) error,
) {
	var web webPacket

	switch {
	case mapstructure.Decode(args, &web) == nil:
		watch.LinkId = web.LinkId

		if err := receiver(l, web, watch); err != nil {
			l.Warn("[Server] Message handled, but unknown web packet body type",
				zap.String("packet_body", fmt.Sprintf("%v", web.Body)),
				zap.String("reason", err.Error()),
			)
		}

	default:
		l.Warn("[Server] Message handled, but unknown packet type",
			zap.String("packet", fmt.Sprintf("%v", args)),
		)
	}
}

func onCommandReceiver(logger *zap.Logger, web webPacket, watch *watcher) error {
	var netObj object.CommandObject
	if err := mapstructure.Decode(web.Body, &netObj); err != nil {
		return err
	}
	watch.Data = netObj

	logger.Info("[Server] Command handled",
		zap.String("link_id", web.LinkId),
		zap.String("command", netObj.Command),
		zap.String("args", fmt.Sprint(netObj.Args)),
	)
	return nil
}

func onDataReceiver(logger *zap.Logger, web webPacket, watch *watcher) error {
	var dataObj object.DataObject
	if err := mapstructure.Decode(web.Body, &dataObj); err != nil {
		return err
	}
	watch.Data = dataObj

	logger.Info("[Server] Data handled",
		zap.String("link_id", web.LinkId),
		zap.String("module", dataObj.Module),
		zap.String("value", fmt.Sprint(dataObj.Value)),
	)
	return nil
}

func onErrorReceiver(logger *zap.Logger, web webPacket, watch *watcher) error {
	var errObj object.ErrorObject
	if err := mapstructure.Decode(web.Body, &errObj); err != nil {
		return err
	}
	watch.Data = errObj

	logger.Info("[Server] Error handled",
		zap.String("link_id", web.LinkId),
		zap.String("origin", errObj.Origin),
		zap.String("reason", fmt.Sprint(errObj.Reason)),
	)
	return nil
}

func TestService(t *testing.T) {
	l := utils.Default.Logger()
	pkt := watcher{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	r := map[string]func(*gosocketio.Channel, interface{}){
		gosocketio.OnConnection:    func(c *gosocketio.Channel, a interface{}) { onConnectionReceiver(c, a, l, &wg) },
		gosocketio.OnDisconnection: func(c *gosocketio.Channel, a interface{}) { onDisconnectionReceiver(c, a, l) },
		constant.Channels.Command:  func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onCommandReceiver) },
		constant.Channels.Data:     func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onDataReceiver) },
		constant.Channels.Error:    func(c *gosocketio.Channel, a interface{}) { onGenericReceiver(a, l, &pkt, onErrorReceiver) },
	}

	go utils.Default.SocketIOServer(r)

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
		l.Fatal(err.Error())
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
