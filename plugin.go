package main

import (
	"context"
	"errors"

	"github.com/sportfun/main/event"
	. "github.com/sportfun/main/plugin"
	"github.com/sportfun/main/profile"
	"log"
	sysplugin "plugin"
	"fmt"
	. "github.com/sportfun/main/protocol/v1.0"
	"sync"
)

func init() {
	Scheduler.RegisterWorker("plugin", pluginTask)
}

type pluginDefinition struct {
	instance *Plugin
	profile  profile.Plugin
	cancel   func()
}
type plugin struct {
	bus     *event.Bus
	plugins map[string]*pluginDefinition

	data        chan interface{}
	status      chan State
	instruction []chan<- Instruction
	sync        sync.Mutex

	deadSig chan string
}

var errNoPluginLoaded = errors.New("no plugin loaded")

func pluginTask(ctx context.Context, bus *event.Bus) error {
	var err error
	var plg plugin

	plg.plugins = map[string]*pluginDefinition{}
	plg.sync = sync.Mutex{}
	plg.bus = bus
	plg.deadSig = make(chan string)
	plg.status = make(chan State)
	plg.data = make(chan interface{})

	for _, plugin := range Profile.Plugins {
		plg.load(plugin)
	}

	if len(plg.plugins) == 0 {
		return errNoPluginLoaded
	}

	defer plg.unsubscribe()
	if bus.Subscribe(":instruction", plg.busInstructionHandler); err != nil {
		return err
	}

	for _, plugin := range plg.plugins {
		plg.run(plugin, ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-plg.data:
			plg.bus.Publish(":data", data, nil)
		case name := <-plg.deadSig:
			//TODO: Check if it broken (reset too quickly) ... If yes, disable it and check if plugin are available
			plg.run(plg.plugins[name], ctx)
			//TODO: Add signal state management
		}
	}
}

func (plg *plugin) unsubscribe() {
	plg.bus.Unsubscribe(":instruction", plg.busInstructionHandler)
}

func (plg *plugin) load(profile profile.Plugin) {
	var p *sysplugin.Plugin
	var s sysplugin.Symbol
	var v *Plugin

	for _, step := range []func() error{
		func() error { var err error; p, err = sysplugin.Open(profile.Path); return err },
		func() error { var err error; s, err = p.Lookup("Plugin"); return err },
		func() error {
			var valid bool
			if v, valid = s.(*Plugin); !valid {
				return fmt.Errorf("invalid symbol type (need %T, but get %T)", s, v)
			}
			return nil
		},
		func() error {
			if _, exists := plg.plugins[v.Name]; exists {
				return fmt.Errorf("plugin '%s' already loaded", v.Name)
			}
			return nil
		},
	} {
		if err := step(); err != nil {
			//TODO: LOG :: ERROR - Failed to load plugin X: Y
			log.Printf("{plugin}[ERROR]			Failed to load plugin '%s': %s", profile.Name, err)
			return
		}
	}

	//TODO: LOG :: INFO - Plugin X successfully loaded
	log.Printf("{plugin}[INFO]			Plugin '%s' successfully loaded", profile.Name)

	plg.plugins[v.Name] = &pluginDefinition{instance: v, profile: profile, cancel: nil}
}

func (plg *plugin) run(pluginX *pluginDefinition, parent context.Context) {
	ctx, cnl := context.WithCancel(parent)

	pluginX.cancel = cnl
	go func(pluginY *Plugin, profile profile.Plugin, ctx context.Context) {
		data := make(chan interface{})
		go func(in <-chan interface{}, out chan<- interface{}) {
			for v := range in {
				out <- struct {
					name  string
					value interface{}
				}{name: pluginY.Name, value: v}
			}
		}(data, plg.data)
		defer func(c chan interface{}) { close(c) }(data)

		status := make(chan State)
		go func(in <-chan State, out chan<- State) {
			for v := range in {
				out <- v
			}
		}(status, plg.status)
		defer func(c chan State) { close(c) }(status)

		inst := make(chan Instruction)
		plg.sync.Lock()
		plg.instruction = append(plg.instruction, inst)
		plg.sync.Unlock()
		defer func(c chan Instruction) { close(inst) }(inst)
		defer func(p *plugin, c chan Instruction) {
			p.sync.Lock()
			defer p.sync.Unlock()
			for i, x := range p.instruction {
				if x == c {
					p.instruction = append(p.instruction[:i-1], p.instruction[i:]...)
				}
			}
		}(plg, inst)

		if err := pluginY.Instance(ctx, profile, Chan{Data: data, Status: status, Instruction: inst}); err != nil {
			//TODO: LOG :: ERROR - Plugin has crashed: X
			log.Printf("{plugin#%s}[ERROR]		Plugin has crashed: %s", pluginY.Name, err.Error())
			plg.bus.Publish(":error", struct {
				origin string
				error  error
			}{origin: pluginY.Name, error: err}, nil)
		}

		plg.deadSig <- pluginY.Name
	}(pluginX.instance, pluginX.profile, ctx)
}

func (plg *plugin) busInstructionHandler(event *event.Event, err error) {
	if inst, exists := Instructions[event.Message().(string)]; !exists {
		//TODO: LOG :: ERROR - Unknown instruction X
		log.Printf("{network}[ERROR]				Unknown instruction '%s'", event.Message().(string))
	} else {
		plg.sync.Lock()
		defer plg.sync.Unlock()

		for _, x := range plg.instruction {
			go func(c chan<- Instruction) { c <- inst }(x)
		}
	}
}
