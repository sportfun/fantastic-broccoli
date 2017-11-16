package fantastic_broccoli

import (
	"flag"
	"go.uber.org/zap"

	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/constant"
	"fantastic-broccoli/errors"
	"fantastic-broccoli/kernel"
	"fantastic-broccoli/log"
	"fantastic-broccoli/properties"
	"fantastic-broccoli/services/module"
	"fantastic-broccoli/services/network"
	"fantastic-broccoli/utils"
)

var propertiesPath string
var nRetryMax int
var props *properties.Properties
var logger *zap.Logger

func init() {
	flag.StringVar(&propertiesPath, "properties", "/etc/sportsfun/acquisitor.json", "path where file is configured")
	flag.IntVar(&nRetryMax, "maxretry", 5, "number max of retry before failure")
	flag.Parse()
}

func main() {
	core := kernel.Core{}
	services := []service.Service{&network.Service{}, &module.Service{}}
	props = properties.LoadFrom(propertiesPath)
	logger = log.Configure(props)

configuration:
	if hasFailed(core.Configure(services, props, logger)) {
		properties.WaitReconfiguration(props) // Wait until properties file has been changed
		goto configuration
	}

	nRetry := 0
processing:
	for core.State() != constant.States.Stopped {
		if hasPanic(core.Run()) {
			hasFailed(core.Stop()) // Just used to display error if needed
			nRetry++

			// Retry n times before maintenance mode (system locked + LED blinked)
			if nRetry < nRetryMax {
				goto processing
			}
			utils.MaintenanceMode()
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
		logger.Error(err.Error(), zap.String("level", err.Level), zap.NamedError("error", &err))
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
