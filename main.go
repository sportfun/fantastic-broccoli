package main

import (
	"flag"

	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/errors"
	"github.com/xunleii/fantastic-broccoli/kernel"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/services/module"
	"github.com/xunleii/fantastic-broccoli/services/network"
	"github.com/xunleii/fantastic-broccoli/utils"
)

var propertiesPath string
var nRetryMax int
var props *properties.Properties
var logger log.Logger

var (
	simpleError   = log.NewArgumentBinder("%s")
	internalError = log.NewArgumentBinder("internal error from '%s' (%s): %s")
)

func init() {
	flag.StringVar(&propertiesPath, "p", "/etc/sportsfun/config.json", "path where file is configured (shorthand)")
	flag.StringVar(&propertiesPath, "properties", "/etc/sportsfun/config.json", "path where file is configured")
	flag.IntVar(&nRetryMax, "m", 5, "number max of retry before failure (shorthand)")
	flag.IntVar(&nRetryMax, "max-retry", 5, "number max of retry before failure")
}

func main() {
	flag.Parse()
	core := kernel.Core{}
	services := []service.Service{&network.Service{}, &module.Service{}}
	props = properties.LoadFrom(propertiesPath)
	logger = log.NewLogger.Prod(props)

configuration:
	if hasFailed(core.Configure(services, props, logger)) {
		core.Stop()
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
	case *errors.InternalError:
		logger.Error(internalError.Bind(err.Origin, err.Level, err.Error()))
	default:
		logger.Error(simpleError.Bind(err.Error()))
	}

	return true
}

func hasPanic(err error) bool {
	if !hasFailed(err) {
		return false
	}

	switch err := err.(type) {
	case *errors.InternalError:
		return err.Level == constant.ErrorLevels.Fatal
	default:
		return false
	}
}
