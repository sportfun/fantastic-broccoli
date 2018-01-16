package network

import (
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/config"
	. "github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
	"github.com/sportfun/gakisitor/utils"
	"testing"
)

const (
	IPAddress = "localhost"
	Port      = 8888
)

func init() {
	go utils.SocketIOServer(utils.WSReceivers{}, Port)
}

func TestNetwork_Start(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	network := Network{}

	Expect(network.Start(queue, logger)).Should(Succeed())
	Expect(network.State()).Should(Equal(StartedState))
}

func TestNetwork_Configure(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	network := Network{}

	Expect(network.Start(queue, logger)).Should(Succeed())

	Expect(network.Configure(nil)).Should(MatchError("configuration not defined"))
	Expect(network.State()).Should(Equal(PanicState))

	network.state = StartedState
	Expect(network.Configure(&config.GAkisitorConfig{System: config.SystemDefinition{}})).ShouldNot(Succeed())
	Expect(network.State()).Should(Equal(PanicState))

	network.state = StartedState
	Expect(network.Configure(&config.GAkisitorConfig{
		System: config.SystemDefinition{
			LinkID:     "0000",
			ServerIP:   IPAddress,
			ServerPort: Port,
		},
	})).Should(Succeed())
	Expect(network.State()).Should(Equal(IdleState))
}

func TestNetwork_Process(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	network := Network{}

	Expect(network.Start(queue, logger)).Should(Succeed())

	queue.Notify(notification.NewNotification(ModuleServiceEntity, NetworkServiceEntity, object.NewDataObject("", "")))
	Expect(network.Process()).Should(MatchError("failed to emit message"))

	Expect(network.Configure(&config.GAkisitorConfig{
		System: config.SystemDefinition{
			LinkID:     "0000",
			ServerIP:   IPAddress,
			ServerPort: Port,
		},
	})).Should(Succeed())

	queue.Notify(notification.NewNotification(ModuleServiceEntity, network.Name(), object.NewDataObject("", "")))
	Expect(network.Process()).Should(Succeed())
}

func TestNetwork_Stop(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)
	queue := service.NewNotificationQueue()
	network := Network{}

	Expect(network.Start(queue, logger)).Should(Succeed())
	Expect(network.Configure(&config.GAkisitorConfig{
		System: config.SystemDefinition{
			LinkID:     "0000",
			ServerIP:   IPAddress,
			ServerPort: Port,
		},
	})).Should(Succeed())
	Expect(network.Stop()).Should(Succeed())
}
