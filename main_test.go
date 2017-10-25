package fantastic_broccoli

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"fantastic-broccoli/model"
	"fantastic-broccoli/core"
	"log"
	"net/http"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"fantastic-broccoli/constant"
	"time"
	"sync"
	"fantastic-broccoli/common/types/service"
	"fantastic-broccoli/services/network"
)

func NewLogger() *zap.Logger {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.Lock(os.Stdout)
	topicErrors := zapcore.Lock(os.Stderr)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	defer logger.Sync()
	return logger
}

func NewServer() {
	//create server instance, you can setup transport parameters or get the default one
	//look at websocket.go for parameters description
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// --- caller is default handlers

	//on connection handler, occurs once for each connected client
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, args interface{}) {
		//client id is unique
		log.Printf("[Server] New client connected, client id is '%s'", c.Id())

		c.Join(constant.CommandChan)
	})
	//on disconnection handler, if client hangs connection unexpectedly, it will still occurs
	//you can omit function args if you do not need them
	//you can return string value for ack, or return nothing for emit
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		//caller is not necessary, client will be removed from rooms
		//automatically on disconnect
		//but you can remove client from room whenever you need to
		c.Leave(constant.CommandChan)

		log.Printf("[Server] %s (%s) disconnected", c.Id(), c.Ip())
	})

	//error catching handler
	server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("Error occurs")
	})

	// --- caller is custom handler

	server.On(constant.CommandChan, func(c *gosocketio.Channel, args interface{}) {
		log.Printf("[Server] Something successfully handled (%v)", args)
		c.Emit(constant.CommandChan, "ok")
	})

	//setup http server like caller for handling connections
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)
	log.Fatal(http.ListenAndServe(":80", serveMux))
}

func BenchmarkServer(b *testing.B) {
	var wg sync.WaitGroup
	go NewServer()

	wg.Add(0xF)
	for i := 0; i < 0xF; i++ {
		go func(i int) {
			c, err := gosocketio.Dial(
				gosocketio.GetUrl("localhost", 80, false),
				transport.GetDefaultWebsocketTransport(),
			)

			if err != nil {
				log.Fatal(err)
			}

			c.On(constant.CommandChan, func(x *gosocketio.Channel, args interface{}) {
				log.Printf("[%s] Received something '%v' from %s", c.Id(), args, x.Id())
			})
			ms := &struct {
				LinkId  string
				Command string
				Args    []string
			}{"0x", "command", []string{"ok"}}
			c.Emit(constant.CommandChan, ms)
			time.Sleep(10 * time.Second)
			c.Close()
			time.Sleep(time.Second)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkCore_Configure(b *testing.B) {
	c := core.Core{}
	p := model.Properties{
		System: model.SystemDefinition{
			LinkID:     "c0f629ed-4b1a-4fea-a88e-5d9070807112",
			ServerIP:   "localhost/socket.io/",
			ServerPort: 80,
			ServerSSL:  false,
		},
	}
	l := NewLogger()
	go NewServer()

	s := []service.Service{new(network.Service)}

	c.Configure(s, &p, l)
	c.Run()
}

func TestLol(t *testing.T) {
	NewLogger().Info("constructed a logger")
}
