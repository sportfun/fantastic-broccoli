package main

import (
	"context"
	"errors"
	"fmt"
	sysplugin "plugin"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
	. "github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

func init() {
	Gakisitor.RegisterWorker("plugin", pluginTask)
}

type pluginDefinition struct {
	instance *Plugin
	profile  profile.Plugin
	cancel   func()
}
type plugin struct {
	bus     *bus.Bus
	plugins map[string]*pluginDefinition

	data        chan interface{}
	status      chan interface{}
	instruction []chan<- Instruction
	sync        sync.Mutex

	active  sync.WaitGroup
	deadSig chan string
}

var errNoPluginLoaded = errors.New("no plugin loaded")

func pluginTask(ctx context.Context, bus *bus.Bus) error {
	var err error
	var plg = plugin{
		bus:     bus,
		plugins: map[string]*pluginDefinition{},
		data:    make(chan interface{}),
		status:  make(chan interface{}),
		sync:    sync.Mutex{},
		active:  sync.WaitGroup{},
		deadSig: make(chan string),
	}

	for _, plugin := range Gakisitor.Plugins {
		if err = plg.load(plugin); err != nil {
			logrus.Errorf("Failed to load plugin '%s': %s", plugin.Name, err) // LOG :: ERROR - Failed to load plugin {name}: {err}
		}
	}

	if len(plg.plugins) == 0 {
		return errNoPluginLoaded
	}

	defer plg.unsubscribe()
	if bus.Subscribe(":instruction", plg.busInstructionHandler); err != nil {
		return err
	}

	for _, plugin := range plg.plugins {
		plg.run(ctx, plugin)
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Debug("Closed by context, wait all plugins")
			plg.active.Wait()
			logrus.Debug("All plugins stopped")
			plg.active.Wait()
			return nil
		case data := <-plg.data:
			plg.bus.Publish(":data", data, nil)
		case status := <-plg.status:
			logrus.Infof("Status received: %v", status)
		case name := <-plg.deadSig:
			plg.run(ctx, plg.plugins[name])
		}
	}
}

// unsubscribe unsubscribes all bus handler.
func (plg *plugin) unsubscribe() {
	plg.bus.Unsubscribe(":instruction", plg.busInstructionHandler)
}

// load loads one plugin.
func (plg *plugin) load(profile profile.Plugin) error {
	var p *sysplugin.Plugin
	var s sysplugin.Symbol
	var v *Plugin

	for _, step := range []func() error{
		func() error { var err error; p, err = sysplugin.Open(profile.Path); return err },
		func() error { var err error; s, err = p.Lookup("Plugin"); return err },
		func() error {
			var valid bool
			if v, valid = s.(*Plugin); !valid {
				return fmt.Errorf("invalid symbol type (need %T, but get %T)", v, s)
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
			return err
		}
	}

	logrus.Infof("Plugin %s successfully loaded", profile.Name) // LOG :: INFO - Plugin {name} successfully loaded
	plg.plugins[v.Name] = &pluginDefinition{instance: v, profile: profile, cancel: nil}
	return nil
}

// run manages one plugin.
func (plg *plugin) run(parentCtx context.Context, def *pluginDefinition) {
	if def == nil {
		return
	}

	ctx, cnl := context.WithCancel(parentCtx)

	def.cancel = cnl
	plg.active.Add(1)
	go func(p *Plugin, profile profile.Plugin, ctx context.Context) {
		defer func(p *plugin) {
			if err := recover(); err != nil {
				logrus.Errorf("Panic recovered into plugin service: %s", err) // LOG :: ERROR - Panic recovered into plugin service: {reason}
			}
		}(plg)
		defer func() { plg.deadSig <- profile.Name }()
		defer func() { plg.active.Done() }()

		data := make(chan interface{})
		go func(in <-chan interface{}, out chan<- interface{}) {
			for v := range in {
				out <- struct {
					name  string
					value interface{}
				}{name: p.Name, value: v}
			}
		}(data, plg.data)
		defer func(c chan interface{}) { close(c) }(data)

		status := make(chan State)
		go func(in <-chan State, out chan<- interface{}) {
			for v := range in {
				out <- struct {
					name  string
					state State
				}{name: p.Name, state: v}
			}
		}(status, plg.status)
		defer func(c chan State) { close(c) }(status)

		inst := make(chan Instruction)
		plg.sync.Lock()
		plg.instruction = append(plg.instruction, inst)
		plg.sync.Unlock()

		defer func(p *plugin, c chan Instruction) {
			p.sync.Lock()
			defer p.sync.Unlock()
			if len(p.instruction) == 1 {
				p.instruction = nil
			}

			for i := len(p.instruction) - 1; i >= 0; i-- {
				if p.instruction[i] == c {
					p.instruction = append(p.instruction[:i-1], p.instruction[i:]...)
					close(c)
				}
			}
		}(plg, inst)

		if err := p.Instance(ctx, profile, Chan{Data: data, Status: status, Instruction: inst}); err != nil {
			logrus.Errorf("Plugin '%s' has crashed: %s", p.Name, err) // LOG :: ERROR - Plugin {name} has crashed: {err}
			plg.bus.Publish(":error", struct {
				origin string
				error  error
			}{origin: p.Name, error: err}, nil)
		}
	}(def.instance, def.profile, ctx)
}

// busInstructionHandler handles event from ':instruction' on the event bus.
func (plg *plugin) busInstructionHandler(event *bus.Event, err error) {
	if err != nil {
		if err != bus.ErrSubscriberDeleted {
			logrus.Errorf("Bus handler for ':instruction' failed: %s", err) // LOG :: ERROR - Bus handler for ':instruction' failed: {error}
		}
		return
	}

	if name, valid := event.Message().(string); !valid {
		logrus.Errorf("Invalid instruction type: %#v", event.Message()) // LOG :: ERROR - Invalid instruction type: {message}
		return
	} else if inst, exists := Instructions[name]; !exists {
		logrus.Errorf("Unknown instruction '%s'", event.Message().(string)) // LOG :: ERROR - Unknown instruction {message}
	} else {
		plg.sync.Lock()
		defer plg.sync.Unlock()

		for _, x := range plg.instruction {
			go func(c chan<- Instruction) { c <- inst }(x)
		}
	}
}
