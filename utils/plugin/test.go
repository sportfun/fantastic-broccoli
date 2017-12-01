package plugin

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/module"
	"github.com/xunleii/fantastic-broccoli/config"
	"github.com/xunleii/fantastic-broccoli/utils"
	"github.com/xunleii/fantastic-broccoli/env"
)

var NProcesses = 5

func testStart(t *testing.T, mod module.Module, queue *module.NotificationQueue, logger log.Logger) {
	t.Logf("- Start module '%s'\n", mod.Name())

	// Start the module
	if err := mod.Start(queue, logger); err != nil {
		t.Fatalf("! Failure during module starting - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.StartedState, mod.State())
}

func newDefinition(definition string) config.ModuleDefinition {
	var conf interface{}

	if definition == "" {
		conf = nil
	} else {
		json.Unmarshal([]byte(definition), &conf)
	}

	return config.ModuleDefinition{Config: definition}
}

func testConfigure(t *testing.T, mod module.Module, environment *environment) {
	t.Logf("- Configure module '%s'\n", mod.Name())

	nil_definition := newDefinition("")
	empty_definition := newDefinition("{}")
	invalid_definition := newDefinition("{\"no_key_def\":true}")

	// Configuration failure : NIL definition
	err := mod.Configure(&nil_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, env.PanicState, mod.State())

	// Configuration failure : empty definition
	err = mod.Configure(&empty_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, env.PanicState, mod.State())

	// Configuration failure : invalid definition
	err = mod.Configure(&invalid_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, env.PanicState, mod.State())

	// Configuration succeed
	if err := mod.Configure(environment.definition(t)); err != nil {
		t.Fatalf("! Failure during module configuration - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.IdleState, mod.State())
}

func testStartSession(t *testing.T, mod module.Module) {

	// Process failure : no session started
	err := mod.Process()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, env.IdleState, mod.State())

	// Start session successfully
	t.Logf("\t- Start new session\n")
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.WorkingState, mod.State())

	// Starting session failure : session already started
	err = mod.StartSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, env.IdleState, mod.State())

	// Start session successfully
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.WorkingState, mod.State())
}

func testProcess(t *testing.T, mod module.Module, environment *environment) {
	t.Logf("\t- Processing loops [%d time(s)]\n", NProcesses)

	for i := 0; i < NProcesses; i++ {
		time.Sleep(environment.tick)
		if err := mod.Process(); err != nil {
			switch mod.State() {
			case env.PanicState:
				t.Fatalf("! Panic during processing - %s\n", err.Error())
			case env.WorkingState:
				t.Logf("! Failure during processing - %s\n", err.Error())
			default:
				t.Fatalf("! Failure during processing - %s (invalid module state (0x%X))\n", err.Error(), mod.State())
			}
		}
	}
}

func testStopSession(t *testing.T, mod module.Module) {
	t.Logf("\t- Stop session\n")

	// Stop sessions successfully
	if err := mod.StopSession(); err != nil {
		t.Fatalf("! Failure during ending session - %s\n", err.Error())
	}

	// Process failure : no session started
	err := mod.Process()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, env.IdleState, mod.State())

	// Stopping session failure : session already stopped
	err = mod.StopSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, env.IdleState, mod.State())
}

func testStop(t *testing.T, mod module.Module, environment *environment) {
	// Start session successfully
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.WorkingState, mod.State())

	// Stop session & module
	t.Logf("- Stop module '%s'\n", mod.Name())
	if err := mod.Stop(); err != nil {
		t.Fatalf("! Failure during ending module - %s\n", err.Error())
	}
	utils.AssertEquals(t, env.StoppedState, mod.State())

	// Wait all goroutine
	time.Sleep(2 * environment.tick)
}

// Custom module test tool
func Test(t *testing.T, mod module.Module, environment *environment) {
	// Init environment
	queue := module.NewNotificationQueue()
	logger := log.NewDevelopment()
	inlog := func(format string, a ...interface{}) { t.Logf("\t> %s", fmt.Sprintf(format, a...)) }

	// Start & Configure module
	testStart(t, mod, queue, logger)

	// Start pre tests
	if environment.test.pre != nil {
		t.Logf("- Start pre tests\n")
		environment.test.pre(t, inlog, mod)
	}

	// Execute tests
	t.Logf("- Start tests\n")
	testConfigure(t, mod, environment)
	testStartSession(t, mod)
	testProcess(t, mod, environment)
	testStopSession(t, mod)

	// Start post tests
	if environment.test.post != nil {
		t.Logf("- Start post tests\n")
		environment.test.post(t, inlog, NProcesses, mod, queue)
	}

	// Stop all
	testStop(t, mod, environment)
}
