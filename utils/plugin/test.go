package plugin

import (
	"encoding/json"
	"testing"
	"time"
	"fmt"

	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/utils"
)

var NProcesses = 5

func testStart(t *testing.T, mod module.Module, queue *module.NotificationQueue, logger log.Logger) {
	t.Logf("- Start module '%s'\n", mod.Name())

	// Start the module
	if err := mod.Start(queue, logger); err != nil {
		t.Fatalf("! Failure during module starting - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Started, mod.State())
}

func newDefinition(definition string) properties.ModuleDefinition {
	var conf interface{}

	if definition == "" {
		conf = nil
	} else {
		json.Unmarshal([]byte(definition), &conf)
	}

	return properties.ModuleDefinition{Conf: definition}
}

func testConfigure(t *testing.T, mod module.Module, env *environment) {
	t.Logf("- Configure module '%s'\n", mod.Name())

	nil_definition := newDefinition("")
	empty_definition := newDefinition("{}")
	invalid_definition := newDefinition("{\"no_key_def\":true}")

	// Configuration failure : NIL definition
	err := mod.Configure(nil_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, constant.States.Panic, mod.State())

	// Configuration failure : empty definition
	err = mod.Configure(empty_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, constant.States.Panic, mod.State())

	// Configuration failure : invalid definition
	err = mod.Configure(invalid_definition)
	utils.AssertNotEquals(t, nil, err)
	utils.AssertEquals(t, constant.States.Panic, mod.State())

	// Configuration succeed
	if err := mod.Configure(env.definition(t)); err != nil {
		t.Fatalf("! Failure during module configuration - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Idle, mod.State())
}

func testStartSession(t *testing.T, mod module.Module) {

	// Process failure : no session started
	err := mod.Process()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.States.Idle, mod.State())

	// Start session successfully
	t.Logf("\t- Start new session\n")
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Working, mod.State())

	// Starting session failure : session already started
	err = mod.StartSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.States.Idle, mod.State())

	// Start session successfully
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Working, mod.State())
}

func testProcess(t *testing.T, mod module.Module, env *environment) {
	t.Logf("\t- Processing loops [%d time(s)]\n", NProcesses)

	for i := 0; i < NProcesses; i++ {
		time.Sleep(env.tick)
		if err := mod.Process(); err != nil {
			switch mod.State() {
			case constant.States.Panic:
				t.Fatalf("! Panic during processing - %s\n", err.Error())
			case constant.States.Working:
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
	utils.AssertEquals(t, constant.States.Idle, mod.State())

	// Stopping session failure : session already stopped
	err = mod.StopSession()
	utils.AssertNotEquals(t, err, nil)
	utils.AssertEquals(t, constant.States.Idle, mod.State())
}

func testStop(t *testing.T, mod module.Module, env *environment) {
	// Start session successfully
	if err := mod.StartSession(); err != nil {
		t.Fatalf("! Failure during starting session - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Working, mod.State())

	// Stop session & module
	t.Logf("- Stop module '%s'\n", mod.Name())
	if err := mod.Stop(); err != nil {
		t.Fatalf("! Failure during ending module - %s\n", err.Error())
	}
	utils.AssertEquals(t, constant.States.Stopped, mod.State())

	// Wait all goroutine
	time.Sleep(2 * env.tick)
}

// Custom module test tool
func Test(t *testing.T, mod module.Module, env *environment) {
	// Init environment
	queue := module.NewNotificationQueue()
	logger := log.NewLogger.Dev(nil)
	inlog := func(format string, a ...interface{}) { t.Logf("\t> %s", fmt.Sprintf(format, a...)) }

	// Start & Configure module
	testStart(t, mod, queue, logger)

	// Start pre tests
	if env.test.pre != nil {
		t.Logf("- Start pre tests\n")
		env.test.pre(t, inlog, mod)
	}

	// Execute tests
	t.Logf("- Start tests\n")
	testConfigure(t, mod, env)
	testStartSession(t, mod)
	testProcess(t, mod, env)
	testStopSession(t, mod)

	// Start post tests
	if env.test.post != nil {
		t.Logf("- Start post tests\n")
		env.test.post(t, inlog, NProcesses, mod, queue)
	}

	// Stop all
	testStop(t, mod, env)
}
