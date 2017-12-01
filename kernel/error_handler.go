package kernel

import (
	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/service"
)

var (
	serviceStartingFailed      = log.NewArgumentBinder("failure during '%s' starting")
	serviceConfigurationFailed = log.NewArgumentBinder("failure during '%s' configuration")
	serviceProcessingFailed    = log.NewArgumentBinder("failure during '%s' processing")
	serviceEndingFailed        = log.NewArgumentBinder("failure during '%s' ending")
	kernelPanic                = log.NewArgumentBinder("panic during %s %s")
)

func (core *Core) panic(err error, name string, when string) {
	core.logger.Error(kernelPanic.Bind(name, when).More("error", err))
	core.state = env.PanicState
}

func defaultErrorHandler(core *Core, name, when string, err error) {
	switch err := err.(type) {
	case *internalError:
		if err.level == env.FatalLevel {
			core.panic(err, name, when)
		}
	default:
		core.panic(err, name, when)
	}
}

func isStarted(core *Core, srv service.Service, err error) {
	core.logger.Error(serviceStartingFailed.Bind(srv.Name()).More("error", err))
	defaultErrorHandler(core, srv.Name(), "start", err)
}

func isConfigured(core *Core, srv service.Service, err error) {
	core.logger.Error(serviceConfigurationFailed.Bind(srv.Name()).More("error", err))
	defaultErrorHandler(core, srv.Name(), "configuration", err)
}

func isProcessed(core *Core, srv service.Service, err error) {
	core.logger.Error(serviceProcessingFailed.Bind(srv.Name()).More("error", err))
	defaultErrorHandler(core, srv.Name(), "process", err)
}

func isStopped(core *Core, srv service.Service, err error) {
	core.logger.Error(serviceEndingFailed.Bind(srv.Name()).More("error", err))
}
