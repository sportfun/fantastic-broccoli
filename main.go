package main

import (
	"context"
	"log"

	"github.com/sportfun/main/event"
	"github.com/sportfun/main/profile"
)

var Scheduler = &scheduler{bus: event.NewBus(), workers: map[string]*worker{}, ctx: context.Background(), deadSig: make(chan string)}
var Profile = &profile.Profile{
	Network: struct {
		HostAddress string `json:"host_address"` // host address (IPv4 / IPv6)
		Port        int    `json:"port"`         // host port
		EnableSsl   bool   `json:"enable_ssl"`   // enable SSL (if required)
	}{HostAddress: "localhost", Port: 8080, EnableSsl: false},
}

func main() {
	log.Print(Scheduler.Run())
}
