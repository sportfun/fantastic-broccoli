package core

import "fantastic-broccoli/model"

type Name string

type Service interface {
	Start(queue NotificationQueue) error
	Configure(properties model.Properties) error
	Process() error
	Stop() error

	Name() Name
	State() model.State
}
