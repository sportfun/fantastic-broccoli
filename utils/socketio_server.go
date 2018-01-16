package utils

import (
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
)

type WSReceiver func(*gosocketio.Channel, interface{})
type WSReceivers map[string]WSReceiver

func SocketIOServer(receivers WSReceivers, port int32) {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// create receiver
	for method, receiver := range receivers {
		server.On(method, receiver)
	}

	// start server
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)
	log.Print(http.ListenAndServe(fmt.Sprintf(":%d", port), serveMux))
}
