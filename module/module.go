package module

import (
	"fantastic-broccoli/model"
	"fantastic-broccoli/notification"
)

type Name string

type Module interface {
	Start(queue NotificationQueue) error
	Configure(properties model.Properties) error
	Process() error
	Stop() error

	Name() Name
	NotificationCaster() notification.Caster
	State() model.State
}
