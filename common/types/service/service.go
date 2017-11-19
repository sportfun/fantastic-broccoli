package service

import (
	"github.com/xunleii/fantastic-broccoli/properties"
	"go.uber.org/zap"
)

type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *properties.Properties) error
	Process() error
	Stop() error

	Name() string
	State() byte
}
