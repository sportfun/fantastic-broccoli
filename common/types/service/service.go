package service

import (
	"fantastic-broccoli/properties"
	"go.uber.org/zap"
)

// TODO: [v1.x] Retourner une erreur avec type precis (pour la gestion d'erreur)
type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *properties.Properties) error
	Process() error
	Stop() error

	Name() string
	State() int
}
