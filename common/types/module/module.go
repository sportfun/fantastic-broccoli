package module

import (
	"fantastic-broccoli/model"
	"go.uber.org/zap"
)

type Module interface {
	Start(queue *notificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() string
	State() int
}
