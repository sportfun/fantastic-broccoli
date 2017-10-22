package service

import (
	"go.uber.org/zap"
	"fantastic-broccoli/common/types"
	"fantastic-broccoli/model"
)

type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	Name() types.Name
	State() types.State
}
