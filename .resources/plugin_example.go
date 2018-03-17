package main

import (
	"github.com/sportfun/main/plugin"
	"github.com/sportfun/main/profile"
	"context"
	"time"
)

var Plugin = plugin.Plugin{
	Name: "Plugin Example",
	Instance: func(ctx context.Context, profile profile.Plugin, channels plugin.Chan) error {
		var inSession bool
		var state = plugin.RunningState
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
				case plugin.StatusPluginInstruction:
					channels.Status <- state
				case plugin.StartSessionInstruction:
					inSession = true
					state = plugin.InSessionState
				case plugin.StopSessionInstruction:
					inSession = false
					state = plugin.StoppedState
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
