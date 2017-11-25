package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type Module interface {
	Start(queue *NotificationQueue, logger log.Logger) error
	Configure(properties properties.ModuleDefinition) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() string
	State() types.StateType
}
