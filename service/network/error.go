package network

import (
	"github.com/sportfun/gakisitor/log"
)

type errorType int
type netError func(*Network, interface{}, error)

var (
	failedToEmit              = log.NewArgumentBinder("failed to emit message: %s")
	failedToCreateChanHandler = log.NewArgumentBinder("failed to create channel handler: %s")
)

func (service *Network) checkIf(fnc netError, x interface{}, err error) bool {
	if err == nil {
		return true
	}

	fnc(service, x, err)
	return false
}

func isEmitted(service *Network, x interface{}, err error) {
	service.logger.Error(failedToEmit.Bind(err.Error()))
}

func isListening(service *Network, x interface{}, err error) {
	service.logger.Error(failedToCreateChanHandler.Bind(err.Error()))
}
