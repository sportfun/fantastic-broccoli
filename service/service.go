package service

import (
	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/log"
)

type Service interface {
	Start(*NotificationQueue, log.Logger) error
	Configure(*config.GAkisitorConfig) error
	Process() error
	Stop() error

	Name() string
	State() byte
}
