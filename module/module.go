package module

import (
	. "fantastic-broccoli"
	"go.uber.org/zap"
)

type Module interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *Properties) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() Name
	State() State
}
