package main

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
	"github.com/sportfun/gakisitor/profile"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
}

var Scheduler = &scheduler{bus: bus.New(), workers: map[string]*worker{}, ctx: context.Background(), deadSig: make(chan string)}
var Profile = &profile.Profile{
	LinkID: "0000-00000000-0000",
	Network: struct {
		HostAddress string `json:"host_address"` // host address (IPv4 / IPv6)
		Port        int    `json:"port"`         // host port
		EnableSsl   bool   `json:"enable_ssl"`   // enable SSL (if required)
	}{HostAddress: "localhost", Port: 8080, EnableSsl: false},
	Scheduler: struct {
		Worker struct {
			Retry    int `json:"retry"`
			Interval int `json:"interval"`
		} `json:"worker"`
	}{Worker: struct {
		Retry    int `json:"retry"`
		Interval int `json:"interval"`
	}{Retry: 5, Interval: 2000}},
	Plugins: []profile.Plugin{
		{
			Name: "Example plugin",
			Path: "./.resources/plugin_example.so",
			Config: map[string]interface{}{
				"ManyItems": map[string]interface{}{
					"ThisItem": 0,
				},
			},
		},
	},
}

func main() {
	log.Print(Scheduler.Run())
}
