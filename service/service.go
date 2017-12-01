package service

import (
	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/log"
)

type Service interface {
	Start(*NotificationQueue, log.Logger) error
	Configure(*config.GAkisitorConfig) error
	Process() error
	Stop() error

	Name() string
	State() byte
}
