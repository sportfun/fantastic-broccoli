package plugin

import (
	"context"
	"time"

	"github.com/sportfun/gakisitor/profile"
)

func ExamplePlugin_basic() {
	_ = Plugin{
		Name: "ExamplePlugin",
		Instance: func(ctx context.Context, profile profile.Plugin, channels Chan) error {
			var inSession bool
			var state = IdleState
			dataMarshallable := struct {
				A int     `json:"a"`
				B float64 `json:"b"`
			}{0, 0.0}

			// configuration value shared into the profile (Config > ManyItems > ThisItem)
			_, e := profile.AccessTo("ManyItems", "ThisItem")
			if e != nil {
				return e
			}

			// plugin main loop
			for {
				select {
				// closing context here ... DO NOT FORGET IT
				case <-ctx.Done():
					return nil

					// interpret instruction here
				case instruction, valid := <-channels.Instruction:
					// if channel is closed, you must stop the plugin
					if !valid {
						return nil
					}

					switch instruction {
					case StatusPluginInstruction:
						channels.Status <- state
					case StartSessionInstruction:
						inSession = true
						state = InSessionState
					case StopSessionInstruction:
						inSession = false
						state = IdleState
					}

					// example of data sending
				case <-time.Tick(time.Millisecond):
					if inSession {
						channels.Data <- dataMarshallable
					}
				}
			}
		},
	}
}
