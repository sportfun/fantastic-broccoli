package utils

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
)

type WSReceiver func(*gosocketio.Channel, interface{})
type WSReceivers map[string]func(*gosocketio.Channel, interface{})

func SocketIOServer(receivers WSReceivers) {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// create receiver
	for method, receiver := range receivers {
		server.On(method, receiver)
	}

	// start server
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)
	log.Fatal(http.ListenAndServe(":8080", serveMux))
}
