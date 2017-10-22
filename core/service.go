package core

import (
	. "fantastic-broccoli"
	"go.uber.org/zap"
)

type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *Properties) error
	Process() error
	Stop() error

	Name() Name
	State() State
}
