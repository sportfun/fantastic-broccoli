package plugin

import (
	"github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
	"context"
	"time"
)

// Plugin exports the example plugin
var Plugin = plugin.Plugin{
	Name: "Plugin Example",
	Instance: func(ctx context.Context, profile profile.Plugin, channels plugin.Chan) error {
		var inSession bool
		var state = plugin.IdleState
		dataMarshallable := struct {
			A int       `json:"a"`
			B time.Time `json:"b"`
		}{0, time.Now()}

		// configuration value shared into the profile (Config > ManyItems > ThisItem)
		_, e := profile.AccessTo("ManyItems", "ThisItem")
		if e != nil {
			return e
		}

		dataTicker := time.Tick(200 * time.Millisecond)
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
				case plugin.StatusPluginInstruction:
					channels.Status <- state
				case plugin.StartSessionInstruction:
					inSession = true
					state = plugin.InSessionState
				case plugin.StopSessionInstruction:
					inSession = false
					state = plugin.IdleState
				}

				// example of data sending
			case <-dataTicker:
				if inSession {
					dataMarshallable.B = time.Now()
					channels.Data <- dataMarshallable
				}
			}
		}
	},
}
