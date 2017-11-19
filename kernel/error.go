package kernel

import (
	"go.uber.org/zap"
	"fmt"

	"github.com/xunleii/fantastic-broccoli/errors"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
)

type serviceError func(*Core, service.Service, error)

func (core *Core) checkIf(srv service.Service, err error, fnc serviceError) bool {
	core.internal = err
	if err == nil {
		return true
	}

	fnc(core, srv, err)
	return false
}

func IsStarted(core *Core, srv service.Service, err error) {
	core.logger.Error(
		"failure during service start",
		zap.String("service", srv.Name()),
		zap.NamedError("error", err),
	)

	switch err := err.(type) {
	case *errors.InternalError:
		if err.Level == constant.ErrorLevels.Fatal {
			panic(core, err, srv.Name(), "start")
		}
		err.Level = constant.ErrorLevels.Fatal
		core.internal = err
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Fatal, errors.OriginList.Service, srv.Name())
	}
}

func IsConfigured(core *Core, srv service.Service, err error) {
	core.logger.Error(
		"failure during service configuration",
		zap.String("service", srv.Name()),
		zap.NamedError("error", err),
	)

	switch err := err.(type) {
	case *errors.InternalError:
		if err.Level == constant.ErrorLevels.Fatal {
			panic(core, err, srv.Name(), "configuration")
		}
		err.Level = constant.ErrorLevels.Fatal
		core.internal = err
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Fatal, errors.OriginList.Service, srv.Name())
	}
}

func IsProcessed(core *Core, srv service.Service, err error) {
	core.logger.Error(
		"failure during service processing",
		zap.String("service", srv.Name()),
		zap.NamedError("error", err),
	)

	switch err := err.(type) {
	case *errors.InternalError:
	default:
		core.internal = errors.NewInternalError(err, constant.ErrorLevels.Error, errors.OriginList.Service, srv.Name())
	}
}

func IsStopped(core *Core, srv service.Service, err error) {
	core.logger.Error(
		"failure during service ending",
		zap.String("service", srv.Name()),
		zap.NamedError("error", err),
	)
}

func panic(core *Core, err error, name string, when string) {
	core.logger.Error(
		fmt.Sprintf("panic during %s %s", name, when),
		zap.NamedError("error", err),
		zap.String("from", name),
		zap.String("when", when),
	)
	core.state = constant.States.Panic
}
