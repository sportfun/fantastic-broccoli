package module

import (
	"fantastic-broccoli/model"
	"fantastic-broccoli/notification"
)

type Module interface {
	Start(q NotificationQueue) error
	Configure(p model.Properties) error
	Process() error
	Stop() error

	Name() string
	NotificationCaster() notification.Caster
	State() model.State
}
