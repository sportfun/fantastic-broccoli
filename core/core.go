package core

import (
	. "fantastic-broccoli"
	"go.uber.org/zap"
	"fantastic-broccoli/network"
	"fantastic-broccoli/module"
)

type Core struct {
	services      []Service
	notifications NotificationQueue
	logger        *zap.Logger
}

func (c *Core) Configure(p *Properties, l *zap.Logger) {
	services(c)
	c.notifications = NotificationQueue{}

	for _, v := range c.services {
		err := v.Start(&c.notifications, l)
		if err != nil {
			// TODO: Check error
		}

		err = v.Configure(p)
		if err != nil {
			// TODO: Check error
		}
	}
}

func (c *Core) Run() {

}

func services(c *Core) {
	c.services = append(c.services, new(network.Service))
	c.services = append(c.services, new(module.Service))
}
