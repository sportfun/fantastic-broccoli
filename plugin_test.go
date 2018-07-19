package main

import (
	"sync"
	"testing"
	"time"

	"github.com/sportfun/gakisitor/profile"

	"github.com/onsi/gomega"

	"github.com/sportfun/gakisitor/event/bus"
	. "github.com/sportfun/gakisitor/plugin"
	. "github.com/sportfun/gakisitor/protocol/v1.0"
)

func TestPlugin_Load(t *testing.T) {
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
