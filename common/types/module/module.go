package module

import (
	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type Module interface {
	Start(queue *NotificationQueue, logger log.Logger) Error
	Configure(properties properties.ModuleDefinition) Error
	Process() Error
	Stop() Error

	StartSession() Error
	StopSession() Error

	Name() string
	State() types.StateType
}
