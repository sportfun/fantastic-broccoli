package service

import (
	"fantastic-broccoli/model"
	"go.uber.org/zap"
)

type Service interface {
	Start(queue *notificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	Name() string
	State() int
}
