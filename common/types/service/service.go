package service

import (
	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type Service interface {
	Start(queue *NotificationQueue, logger log.Logger) error
	Configure(properties *properties.Properties) error
	Process() error
	Stop() error

	Name() string
	State() types.StateType
}
