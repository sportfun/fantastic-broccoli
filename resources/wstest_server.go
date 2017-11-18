package main

import (
	"fantastic-broccoli/utils"
	"github.com/graarh/golang-socketio"
	"fantastic-broccoli/constant"
	"fmt"
)

func OnConnection(c *gosocketio.Channel, a interface{}) {

}

func OnCommand(c *gosocketio.Channel, a interface{}) {

}

func OnData(c *gosocketio.Channel, a interface{}) {

}

func OnError(c *gosocketio.Channel, a interface{}) {

}

func OnDefault(c *gosocketio.Channel, a interface{}) {
	fmt.Print(*c)
}

func main() {
	utils.Default.SocketIOServer(utils.WSReceivers{
		gosocketio.OnConnection:   OnDefault,
		constant.Channels.Command: OnDefault,
		constant.Channels.Data:    OnDefault,
		constant.Channels.Error:   OnDefault,
	})
}
