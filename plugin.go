package gakisitor

import (
	"context"
	"errors"

	"github.com/sportfun/gakisitor/event"
	. "github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
)

func init() {
	Scheduler.RegisterWorker("plugin", pluginTask)
}

type plugin struct {
	bus     *event.Bus
	plugins []Plugin
	stop    map[string]func()
}

var errNoPluginLoaded = errors.New("no plugin loaded")

func pluginTask(ctx context.Context, bus *event.Bus) error {
	var err error
	var plg plugin

	for _, plugin := range Profile.Plugins {
		plg.loadPlugin(plugin)
	}

	if len(plg.plugins) == 0 {
		return errNoPluginLoaded
	}

	defer plg.unsubscribe()

	if bus.Subscribe(":instruction", plg.busInstructionHandler); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (plg *plugin) unsubscribe() {
	plg.bus.Unsubscribe(":instruction", plg.busInstructionHandler)
}

func (plg *plugin) loadPlugin(profile profile.Plugin) {}

func (plg *plugin) busInstructionHandler(event event.Event, err error) {}
