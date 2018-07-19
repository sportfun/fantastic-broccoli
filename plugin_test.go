package main

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sportfun/gakisitor/profile"

	"github.com/onsi/gomega"

	"github.com/sportfun/gakisitor/event/bus"
	. "github.com/sportfun/gakisitor/plugin"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestPlugin_Task(t *testing.T) {
	skipCI(t)

	gomega.RegisterTestingT(t)
	Gakisitor.Plugins = append(Gakisitor.Plugins, profile.Plugin{
		Name: "...",
		Path: "plugin.so",
	})

	ctx, cnl := context.WithTimeout(context.Background(), time.Second)
	go func() {
		time.Sleep(15 * time.Second)
		cnl()
		panic("watchdog: stucked task")
	}()
	gomega.Expect(pluginTask(ctx, bus.New())).Should(gomega.Succeed())
}

func TestPlugin_Load(t *testing.T) {
	skipCI(t)

	gomega.RegisterTestingT(t)

	plg := &plugin{
		plugins: map[string]*pluginDefinition{
			"already_exists": &pluginDefinition{},
		},
	}

	cases := []struct {
		profile profile.Plugin
		err     string
	}{
		{
			profile: profile.Plugin{Path: ""},
			err:     "realpath failed",
		},
		{
			profile: profile.Plugin{Path: ".........."},
			err:     "realpath failed",
		},
		{
			profile: profile.Plugin{Path: ".."},
			err:     "cannot read file data: Is a directory",
		},
		{
			profile: profile.Plugin{Path: "main.go"},
			err:     "invalid ELF header",
		},
		{
			profile: profile.Plugin{Path: "no_symbol.so"},
			err:     "plugin: symbol Plugin not found",
		},
		{
			profile: profile.Plugin{Path: "invalid_type.so"},
			err:     "invalid symbol type",
		},
		{
			profile: profile.Plugin{Path: "already_exists.so"},
			err:     "plugin 'already_exists' already loaded",
		},
		{
			profile: profile.Plugin{Path: "plugin.so"},
		},
	}

	for _, test := range cases {
		result := gomega.Expect(plg.load(test.profile))
		if test.err == "" {
			result.Should(gomega.Succeed())
		} else {
			result.Should(gomega.MatchError(gomega.MatchRegexp(test.err)))
		}
	}
}

func TestPlugin_Run(t *testing.T) {
	gomega.RegisterTestingT(t)

	plg := plugin{
		bus:     bus.New(),
		plugins: map[string]*pluginDefinition{},
		data:    make(chan interface{}),
		status:  make(chan interface{}),
		sync:    sync.Mutex{},
		active:  sync.WaitGroup{},
		deadSig: make(chan string),
	}
	ctx, cnl := context.WithCancel(context.Background())
	started := sync.WaitGroup{}

	started.Add(1)
	plg.run(ctx, &pluginDefinition{
		instance: &Plugin{
			Instance: func(ctx context.Context, p profile.Plugin, c Chan) error {
				started.Done()
				select {
				case <-ctx.Done():
					return nil
				}
			},
		},
		profile: profile.Plugin{
			Name: "...",
		},
	})

	started.Wait()
	gomega.Expect(plg.instruction).Should(gomega.HaveLen(1))
	cnl()
	plg.active.Wait()
	time.Sleep(time.Second)

	var x string
	gomega.Expect(plg.instruction).Should(gomega.BeEmpty())
	gomega.Expect(plg.deadSig).Should(gomega.Receive(&x))
	gomega.Expect(x).Should(gomega.Equal("..."))
}

func TestPlugin_BusHandler(t *testing.T) {
	gomega.RegisterTestingT(t)

	inst := make(chan Instruction, 0)
	plg := &plugin{
		bus:         bus.New(),
		instruction: []chan<- Instruction{inst},
		sync:        sync.Mutex{},
	}

	plg.bus.Subscribe(":instruction", plg.busInstructionHandler)
	channel := gomega.Eventually(func() chan Instruction { return inst }, 250*time.Millisecond)

	cases := []struct {
		message       interface{}
		shouldReceive bool
		instruction   Instruction
	}{
		{
			message:       0,
			shouldReceive: false,
		},
		{
			message:       "unknown",
			shouldReceive: false,
		},
		{
			message:       "start_game",
			shouldReceive: true,
			instruction:   Instructions["start_game"],
		},
	}

	for _, test := range cases {
		plg.bus.Publish(":instruction", test.message, nil)
		if test.shouldReceive {
			var x interface{}

			channel.Should(gomega.Receive(&x))
			gomega.Expect(x).Should(gomega.Equal(test.instruction))
		} else {
			channel.ShouldNot(gomega.Receive())
		}
	}
}
