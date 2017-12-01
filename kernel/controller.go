package kernel

import (
	"fmt"

	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/service"
)

type controller struct {
	component
	guard *guard
}

var (
	infoStartServices   = log.NewArgumentBinder("start services")
	infoServicesStarted = log.NewArgumentBinder("services successfully started (%d services)")
	infoStopServices    = log.NewArgumentBinder("stop services")
	infoServicesStopped = log.NewArgumentBinder("services successfully stopped (%d services)")
)

func (c *controller) tryConfigure() error {
	if !c.core.config.IsLoaded() {
		return fmt.Errorf("properties not loaded")
	}
	c.core.notifications = service.NewNotificationQueue()

	c.core.logger.Info(infoStartServices)
	for _, service := range services {
		if !c.guard.checkIf(service, isStarted, service.Start(c.core.notifications, c.core.logger)) ||
			!c.guard.checkIf(service, isConfigured, service.Configure(&c.core.config)) {
			return c.guard.last
		}
	}
	c.core.logger.Info(infoServicesStarted.Bind(len(services)))

	c.core.state = env.IdleState
	return nil
}

func (c *controller) configure() {
	c.core.config.Load()
	c.core.logger = log.NewProduction(c.core.config.Log...)

configuration:
	if c.guard.hasFailed(c.tryConfigure()) {
		c.core.Stop()
		c.core.config.WaitReconfiguration()
		goto configuration
	}
}

func (c *controller) process() error {
	c.core.state = env.WorkingState

	if len(services) == 0 {
		return NewInternalError(fmt.Errorf("no service found"), env.FatalLevel, "kernel", env.CoreEntity)
	}
	for _, service := range services {
		if !c.guard.checkIf(service, isProcessed, service.Process()) {
			return c.guard.last
		}

		for _, notification := range c.core.notifications.Notifications(env.CoreEntity) {
			c.core.handle(notification)
		}
	}
	c.core.state = env.IdleState
	return nil
}

func (c *controller) run() {
	c.core.retry.current = 0

processing:
	for c.core.isRunning() {
		if c.guard.hasPanic(c.process()) {
			c.core.Restart()
			c.core.retry.current++

			if c.core.retry.current < c.core.retry.max {
				goto processing
			}

			//TODO: Maintenance mode(blink led)
			panic("maintenance mode not implemented")
		}
		c.core.retry.current = 0
	}
	c.core.Stop()
}

func (c *controller) stop() {
	c.core.logger.Info(infoStopServices)

	for _, service := range services {
		if !c.guard.checkIf(service, isProcessed, service.Stop()) {
			c.core.panic(c.guard.last, service.Name(), "stop")
		}
	}
	c.core.logger.Info(infoServicesStopped.Bind(len(services)))
	c.core.state = env.StoppedState
}
