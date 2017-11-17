package kernel

import (
	"fmt"
	"go.uber.org/zap"

	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/properties"
)

type Core struct {
	services   []service.Service
	logger     *zap.Logger
	properties *properties.Properties

	notifications service.NotificationQueue
	internal      error
	state         byte
}

func (core *Core) Configure(services []service.Service, props *properties.Properties, logger *zap.Logger) error {
	// Property file can be not loaded (props.IsLoaded = false) if file not found or invalid
	if !props.IsLoaded() {
		return fmt.Errorf("properties not loaded")
	}

	core.services = services
	core.logger = logger
	core.notifications = service.NotificationQueue{}

	core.internal = nil
	logger.Info("start services")
	for _, s := range services {
		if !core.checkIf(s, s.Start(&core.notifications, logger), IsStarted) ||
			!core.checkIf(s, s.Configure(props), IsConfigured) {
			return core.internal
		}
	}
	logger.Info("services successfully started")

	core.state = constant.States.Idle
	return nil
}

func (core *Core) Run() error {
	core.state = constant.States.Working
	for _, s := range core.services {
		if core.checkIf(s, s.Process(), IsProcessed) {
			return core.internal
		}

		for _, n := range core.notifications.Notifications(constant.EntityNames.Core) {
			core.handle(n)
		}
	}
	core.state = constant.States.Idle
	return nil
}

func (core *Core) Stop() error {
	for _, s := range core.services {
		if core.checkIf(s, s.Stop(), IsStopped) {
			return core.internal
		}
	}
	core.state = constant.States.Stopped
	return nil
}

func (core *Core) State() byte {
	return core.state
}
