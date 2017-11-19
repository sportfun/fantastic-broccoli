package main

import (
	"github.com/xunleii/fantastic-broccoli/utils"
	"github.com/graarh/golang-socketio"
	"github.com/xunleii/fantastic-broccoli/constant"
	"fmt"
	"time"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/mitchellh/mapstructure"
)

var online = map[string]*client{}

type client struct {
	isOnline bool
	data     [][]byte
}

type webPacket struct {
	LinkId string      `json:"link_id" mapstructure:"link_id"`
	Body   interface{} `json:"body" mapstructure:"body"`
}

func OnConnection(c *gosocketio.Channel, a interface{}) {
	fmt.Printf("[%s] New connection\n", c.Id())

	online[c.Id()] = &client{isOnline: true, data: [][]byte{}}
	fmt.Printf("[%s] Start testing loop\n", c.Id())

	for online[c.Id()].isOnline {
		time.Sleep(5 * time.Second)
		fmt.Printf("[%s] > send '%s'\n", c.Id(), constant.NetCommand.StartSession)
		c.Emit(constant.Channels.Command, cmdToWSObj(constant.NetCommand.StartSession))
		time.Sleep(10 * time.Second)
		fmt.Printf("[%s] > send '%s'\n", c.Id(), constant.NetCommand.EndSession)
		c.Emit(constant.Channels.Command, cmdToWSObj(constant.NetCommand.EndSession))
		// Statistiques
	}

	fmt.Printf("[%s] End testing loop\n", c.Id())
}

func OnDisconnection(c *gosocketio.Channel, a interface{}) {
	fmt.Printf("[%s] Disconnected\n", c.Id())
	online[c.Id()].isOnline = false
}

func OnCommand(c *gosocketio.Channel, a interface{}) {
	fmt.Printf("[%s] Command received\n", c.Id())

	var obj interface{}
	if obj = toWSObj(c, a); obj == nil {
		return
	}

	var object object.CommandObject
	if err := mapstructure.Decode(obj, &object); err != nil {
		fmt.Printf("[%s]\t Invalid command (%#v) (%s)\n", c.Id(), obj, err.Error())
		return
	}

	fmt.Printf("[%s]\t %s -> %v\n", c.Id(), object.Command, object.Args)
}

func OnData(c *gosocketio.Channel, a interface{}) {
	fmt.Printf("[%s] Data received\n", c.Id())

	var obj interface{}
	if obj = toWSObj(c, a); obj == nil {
		return
	}

	var object object.DataObject
	if err := mapstructure.Decode(obj, &object); err != nil {
		fmt.Printf("[%s]\t Invalid data (%#v) (%s)\n", c.Id(), obj, err.Error())
		return
	}

	fmt.Printf("[%s]\t %s -> %#v\n", c.Id(), object.Module, object.Value)
}

func OnError(c *gosocketio.Channel, a interface{}) {
	fmt.Printf("[%s] Error received\n", c.Id())

	var obj interface{}
	if obj = toWSObj(c, a); obj == nil {
		return
	}

	var object object.ErrorObject
	if err := mapstructure.Decode(obj, &object); err != nil {
		fmt.Printf("[%s]\t Invalid error (%#v) (%s)\n", c.Id(), obj, err.Error())
		return
	}

	fmt.Printf("[%s]\t %s -> %s", c.Id(), object.Origin, object.Reason)
}

func cmdToWSObj(command string, args ...string) webPacket {
	return webPacket{"XXXX-XXXX-XXXX-XXXX", object.CommandObject{Command: command, Args: args}}
}

func toWSObj(c *gosocketio.Channel, packet interface{}) interface{} {
	var ws webPacket

	if err := mapstructure.Decode(packet, &ws); err != nil {
		fmt.Printf("[%s]\t Invalid packet (%#v) (%s)\n", c.Id(), ws, err)
		return nil
	}

	return ws.Body
}

func main() {
	utils.Default.SocketIOServer(utils.WSReceivers{
		gosocketio.OnConnection:    OnConnection,
		gosocketio.OnDisconnection: OnDisconnection,
		constant.Channels.Command:  OnCommand,
		constant.Channels.Data:     OnData,
		constant.Channels.Error:    OnError,
	})
}
