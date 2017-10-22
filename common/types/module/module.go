package module

import (
	"fantastic-broccoli/model"
	"fantastic-broccoli/common/types"
	"go.uber.org/zap"
)

type Module interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() types.Name
	State() types.State
}
