package fantastic_broccoli

import (
	"fantastic-broccoli/core"
	"fantastic-broccoli/services/network"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/services/module"
	"fantastic-broccoli/constant"
	"flag"
	"go.uber.org/zap"
	"errors"
	"fantastic-broccoli/utils"
)

var propertiesPath string
var nRetryMax int
var properties *properties.Properties
var logger *zap.Logger

func init() {
	flag.StringVar(&propertiesPath, "properties", "/etc/sportsfun/acquisitor.json", "path where file is configured")
	flag.IntVar(&nRetryMax, "maxretry", 5, "number max of retry before failure")
}

func main() {
	kernel := core.Core{}
	services := []service.Service{&network.Service{}, &module.Service{}}
	properties = properties.LoadFrom(propertiesPath)
	logger = logger.Configure(properties)

configuration:
	if hasFailed(kernel.Configure(services, properties, logger)) {
		properties.WaitReconfiguration(properties)
		goto configuration
	}

	nRetry := 0
infinit:
	for kernel.State() != constant.Stopped {
		if hasPanic(kernel.Run()) {
			hasFailed(kernel.Stop())
			if nRetry > nRetryMax {
				utils.MaintenanceMode()
			}
			nRetry++
			goto infinit
		}
		nRetry = 0
	}
}

func hasFailed(err error) bool {
	if err == nil {
		return false
	}

	switch err := err.(type) {
	case errors.InternalError:
		logger.Error(err.Error(), zap.String("level", err.Level), zap.NamedError("error", err))
	default:
		logger.Error(err.Error(), zap.NamedError("error", err))
	}

	return true
}

func hasPanic(err error) bool {
	if !hasFailed(err) {
		return false
	}

	switch err := err.(type) {
	case errors.InternalError:
		return err.Level == constant.ErrorLevels.Fatal
	default:
		return false
	}
}
