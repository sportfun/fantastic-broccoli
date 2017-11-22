package kernel

import (
	"github.com/xunleii/fantastic-broccoli/errors"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

type serviceError func(*Core, service.Service, error)

var (
	FailureServiceStarting      = log.NewArgumentBinder("failure during '%s' starting")
	FailureServiceConfiguration = log.NewArgumentBinder("failure during '%s' configuration")
	FailureServiceProcessing    = log.NewArgumentBinder("failure during '%s' processing")
	FailureServiceEnding        = log.NewArgumentBinder("failure during '%s' ending")
	KernelPanic                 = log.NewArgumentBinder("panic during %s %s")
)

func (core *Core) checkIf(srv service.Service, err error, fnc serviceError) bool {
	core.internal = err
	if err == nil {
		return true
	}

	fnc(core, srv, err)
	return false
}

func IsStarted(core *Core, srv service.Service, err error) {
	core.logger.Error(FailureServiceStarting.Bind(srv.Name()).More("error", err))

	switch err := err.(type) {
	case *errors.InternalError:
		if err.Level == constant.ErrorLevels.Fatal {
			_panic(core, err, srv.Name(), "start")
		}
		core.internal = err
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Fatal, errors.OriginList.Service, srv.Name())
	}
}

func IsConfigured(core *Core, srv service.Service, err error) {
	core.logger.Error(FailureServiceConfiguration.Bind(srv.Name()).More("error", err))

	switch err := err.(type) {
	case *errors.InternalError:
		if err.Level == constant.ErrorLevels.Fatal {
			_panic(core, err, srv.Name(), "configuration")
		}
		core.internal = err
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Fatal, errors.OriginList.Service, srv.Name())
	}
}

func IsProcessed(core *Core, srv service.Service, err error) {
	core.logger.Error(FailureServiceProcessing.Bind(srv.Name()).More("error", err))

	switch err := err.(type) {
	case *errors.InternalError:
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Error, errors.OriginList.Service, srv.Name())
	}
}

func IsStopped(core *Core, srv service.Service, err error) {
	core.logger.Error(FailureServiceEnding.Bind(srv.Name()).More("error", err))
}

func _panic(core *Core, err error, name string, when string) {
	core.logger.Error(KernelPanic.Bind(name, when).More("error", err))

	core.state = constant.States.Panic
}
