package module

import (
	"github.com/xunleii/fantastic-broccoli/properties"
	"go.uber.org/zap"
)

type Module interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *properties.Properties) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() string
	State() byte
}
