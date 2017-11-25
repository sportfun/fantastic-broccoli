package main

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/xunleii/fantastic-broccoli/common/types/module"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/properties"
	"github.com/xunleii/fantastic-broccoli/utils"
	"github.com/xunleii/fantastic-broccoli/utils/plugin"
)

const nProcesses = 5
var environment = plugin.NewEnvironment(definitionFactoryImpl, preTestImpl, postTestImpl, 300*time.Millisecond)

func definitionFactoryImpl(_type interface{}) properties.ModuleDefinition {
	var v interface{}

	switch _type.(type) {
	// For testing engine
	case *testing.T:
		json.Unmarshal([]byte("{\"rpm.min\":120,\"rpm.max\":500,\"rpm.step\":25,\"rpm.precision\":10}"), &v)
	case *testing.B:
		json.Unmarshal([]byte("{\"rpm.min\":0,\"rpm.max\":1200.0,\"rpm.step\":250,\"rpm.precision\":1000}"), &v)

	// For pre testing
	case string:
		json.Unmarshal([]byte(_type.(string)), &v)

	default:
		log.Fatalf("unknown %#v, impossible to generate module definition", v)
	}

	return properties.ModuleDefinition{
		Name: "RPM Generator",
		Conf: v,
	}
}


func preTestImpl(t *testing.T, log plugin.InternalLogger, module module.Module) {
	failure_l58 := definitionFactoryImpl("{\"rpm.max\":1200}")
	failure_l63 := definitionFactoryImpl("{\"rpm.min\":\"0\",\"rpm.max\":\"1200\",\"rpm.step\":\"250\",\"rpm.precision\":\"1000\"}")
	failure_l71 := properties.ModuleDefinition{Conf: nil}
	failure_l74 := definitionFactoryImpl("{}")

	// failure at l.58
	utils.AssertNotEquals(t, nil, module.Configure(failure_l58))
	// failure at l.63
	utils.AssertNotEquals(t, nil, module.Configure(failure_l63))
	// failure at l.71
	utils.AssertNotEquals(t, nil, module.Configure(failure_l71))
	// failure at l.74
	utils.AssertNotEquals(t, nil, module.Configure(failure_l74))
}

func postTestImpl(t *testing.T, log plugin.InternalLogger, nprocesses int, module module.Module, queue *module.NotificationQueue) {
	notifications := queue.Notifications()

	utils.AssertEquals(t, nprocesses, len(notifications))
	for _, notification := range notifications {
		o := notification.Content().(*object.DataObject)

		utils.AssertEquals(t, 5, len(o.Value.(string)), func(a interface{}, b interface{}) bool { return a.(int) <= b.(int) })
		log("data notified : {%#v} from '%s'", o.Value, o.Module)
	}
}


func TestModule(t *testing.T) {
	plugin.Test(t, ExportModule(), nProcesses, environment)
}

func TestBenchmarkModule(t *testing.T) {
	plugin.Benchmark(t, ExportModule(), environment)
}
