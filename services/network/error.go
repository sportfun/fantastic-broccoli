package network

import "go.uber.org/zap"

type errorType int
type netError func(*Service, interface{}, error)

const (
	SocketOn   errorType = iota
	SocketEmit
)

func (service *Service) checkIf(x interface{}, err error, fnc netError) bool {
	if err == nil {
		return true
	}

	fnc(service, x, err)
	return false
}

func IsEmitted(service *Service, x interface{}, err error) {
	service.logger.Error("Failed to emit message", zap.String("reason", err.Error()))
}

func IsLitening(service *Service, x interface{}, err error) {
	service.logger.Error("Failed to create channel handler", zap.String("reason", err.Error()))
}