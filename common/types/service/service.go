package service

import (
	"go.uber.org/zap"
	"fantastic-broccoli/model"
)

type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	Name() string
	State() int
}
