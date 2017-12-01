package module

import (
	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/log"
)

type Module interface {
	Start(*NotificationQueue, log.Logger) error
	Configure(*config.ModuleDefinition) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() string
	State() byte
}
