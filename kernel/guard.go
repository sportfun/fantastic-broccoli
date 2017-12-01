package kernel

import (
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/service"
)

type serviceErrorHandler func(*Core, service.Service, error)

type guard struct {
	component

	last error
}

var (
	internError = log.NewArgumentBinder("internal error from '%s' (%s): %s")
	baseError   = log.NewArgumentBinder("%s")
)

func (g *guard) hasFailed(e error) bool {
	if e == nil {
		return false
	}

	g.last = e
	switch e := e.(type) {
	case *internalError:
		g.core.logger.Error(internError.Bind(e.origin, e.level, e.Error()))
	default:
		g.core.logger.Error(baseError.Bind(e.Error()))
	}

	return true
}

func (g *guard) hasPanic(e error) bool {
	if !g.hasFailed(e) {
		return false
	}

	switch e := e.(type) {
	case *internalError:
		return e.level == env.FatalLevel
	default:
		return false
	}
}

func (g *guard) checkIf(s service.Service, h serviceErrorHandler, e error) bool {
	if e == nil {
		return true
	}

	g.last = e
	h(g.core, s, e)
	return false
}
