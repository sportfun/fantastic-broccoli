package network

import (
	"github.com/xunleii/fantastic-broccoli/log"
)

type errorType int
type netError func(*Service, interface{}, error)

const (
	SocketOn errorType = iota
	SocketEmit
)

var (
	failedToEmit              = log.NewArgumentBinder("failed to emit message: %s")
	failedToCreateChanHandler = log.NewArgumentBinder("failed to create channel handler: %s")
)

func (service *Service) checkIf(x interface{}, err error, fnc netError) bool {
	if err == nil {
		return true
	}

	fnc(service, x, err)
	return false
}

func IsEmitted(service *Service, x interface{}, err error) {
	service.logger.Error(failedToEmit.Bind(err.Error()))
}

func IsListening(service *Service, x interface{}, err error) {
	service.logger.Error(failedToCreateChanHandler.Bind(err.Error()))
}
