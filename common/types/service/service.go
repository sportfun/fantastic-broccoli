package service

import (
	"fantastic-broccoli/model"
	"go.uber.org/zap"
)

// TODO: [v1.x] Retourner une erreur avec type precis (pour la gestion d'erreur)
type Service interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *model.Properties) error
	Process() error
	Stop() error

	Name() string
	State() int
}
