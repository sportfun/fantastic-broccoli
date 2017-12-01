package kernel

import (
	"fmt"

	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/env"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/service"
)

type retryCounter struct {
	current int
	max     int
}

type component struct {
	core *Core
}

type Core struct {
	config config.GAkisitorConfig
	retry  retryCounter

	controller *controller
	guard      *guard

	logger        log.Logger
	notifications *service.NotificationQueue
	state         byte
}

var services []service.Service

func (core *Core) Parameter(name string) interface{} {
	switch name {
	case "config":
		return core.config.FilePtr()
	case "retry_max":
		return &core.retry.max
	default:
		panic(fmt.Sprintf("unkown parameter '%s' at init ... Shutdown service", name))
	}
}

func (core *Core) Init() {
	core.guard = &guard{component: component{core: core}}
	core.controller = &controller{component: component{core: core}, guard: core.guard}
}

func (core *Core) Run() {
	core.controller.configure()
	core.controller.run()
}

func (core *Core) Restart() {
	core.controller.stop()
	core.controller.configure()
}

func (core *Core) Stop() {
	core.controller.stop()
}

func (core *Core) isRunning() bool {
	return core.state != env.StoppedState && core.state != env.PanicState
}

func RegisterService(service service.Service) {
	services = append(services, service)
}
