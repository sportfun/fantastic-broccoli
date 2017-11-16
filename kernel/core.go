package kernel

import (
	"go.uber.org/zap"

	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/errors"
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
	core.services = services
	core.logger = logger
	core.notifications = service.NotificationQueue{}

	core.internal = nil
	logger.Info("start services")
	for _, s := range services {
		if !core.checkIf(s, s.Start(&core.notifications, logger), IsStarted) ||
			!core.checkIf(s, s.Configure(props), IsConfigured) {
			return errors.NewInternalError(core.internal, constant.ErrorLevels.Fatal, errors.OriginList.Core, "")
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
