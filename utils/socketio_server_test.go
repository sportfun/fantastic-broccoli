package utils

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"sync"
	"testing"
	"time"
	. "github.com/onsi/gomega"
)

func TestSocketIOServer(t *testing.T) {
	RegisterTestingT(t)

	v := NewVolatile(nil)
	x := sync.NewCond(&sync.Mutex{})
	var receivers = WSReceivers{
		"data": func(c *gosocketio.Channel, d interface{}) {
			x.Broadcast()
			v.Set(d)
		},
	}

	go SocketIOServer(receivers, 8080)
	client, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 8080, false),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// error: invalid channel/method (timeout because of deadlock)
	TimeoutPrecision = 250 * time.Microsecond
	ReleaseIfTimeout(t, 250*time.Millisecond, func(testing.TB) {
		x.L.Lock()
		client.Emit("data", "args")
		x.Wait()
		Expect(v.Get()).Should(Equal("args"))
	})
}
