package core

import (
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/model"
	"fantastic-broccoli/services/module"
	"fantastic-broccoli/services/network"
	"go.uber.org/zap"
)

type Core struct {
	services      []service.Service
	state         int
	notifications service.NotificationQueue
	logger        *zap.Logger
	properties    *model.Properties
}

func (c *Core) Configure(s []service.Service, p *model.Properties, l *zap.Logger) {
	c.services = s
	c.notifications = service.NotificationQueue{}
	c.logger = l

	l.Info("start services")
	for _, s := range c.services {
		c.serviceErrorHandler(ModuleStart, s.Start(&c.notifications, l))
		c.serviceErrorHandler(ModuleConfigure, s.Configure(p))
	}
	l.Info("services successfully started")
}

func (c *Core) Run() {
	for _, s := range c.services {
		c.serviceErrorHandler(ModuleProcess, s.Process())
		for _, n := range c.notifications.Notifications(constant.Core) {
			c.notificationHandler(n)
		}
	}
}

func (c *Core) Stop() {
	for _, s := range c.services {
		c.serviceErrorHandler(ModuleStop, s.Stop())
	}
	c.state = constant.STOPPED
}

func (c *Core) State() int {
	return c.state
}

func services(c *Core) {
	c.services = []service.Service{}
	c.services = append(c.services, new(network.Service))
	c.services = append(c.services, new(module.Service))
}
