package utils

import (
	"go.uber.org/zap"
	"github.com/graarh/golang-socketio"
	"log"
	"github.com/graarh/golang-socketio/transport"
	"net/http"
	"go.uber.org/zap/zapcore"
	"os"
)

type _default struct{}
type webPacket struct {
	LinkId string      `json:"link_id" mapstructure:"link_id"`
	Body   interface{} `json:"body" mapstructure:"body"`
}

var Default = _default{}

func (d *_default) Logger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logger := zap.New(core)
	defer logger.Sync()
	return logger
}

func (d *_default) SocketIOServer(receivers map[string]func(*gosocketio.Channel, interface{})) {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	// create receiver
	for method, receiver := range receivers {
		server.On(method, receiver)
	}

	// start server
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server)
	log.Fatal(http.ListenAndServe(":8080", serveMux))
}
