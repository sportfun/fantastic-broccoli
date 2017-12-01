package module

import (
	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/log"
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
