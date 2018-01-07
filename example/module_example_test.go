package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/sportfun/gakisitor/config"
	"github.com/sportfun/gakisitor/module"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/utils/module_test"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/notification"
)

var environment = module_test.NewEnvironment(definitionFactoryImpl, preTestImpl, postTestImpl, tick*5)

func TestModule(t *testing.T) {
	module_test.Test(t, ExportModule(), environment)
}

func definitionFactoryImpl(obj interface{}) *config.ModuleDefinition {
	var v interface{}

	switch obj.(type) {
	case *testing.T:
		json.Unmarshal([]byte("{\"rpm.min\":120,\"rpm.max\":500,\"rpm.step\":25,\"rpm.precision\":10}"), &v)
	case *testing.B:
		json.Unmarshal([]byte("{\"rpm.min\":0,\"rpm.max\":1200.0,\"rpm.step\":250,\"rpm.precision\":1000}"), &v)

	case string:
		json.Unmarshal([]byte(obj.(string)), &v)
	case nil:
		v = nil

	default:
		log.Fatalf("unknown %#v, impossible to generate module definition", v)
	}

	return &config.ModuleDefinition{
		Name:   "RPM Generator",
		Config: v,
	}
}

func preTestImpl(t *testing.T, module module.Module) {
	unmarshall := func(s string) interface{} { var v interface{}; json.Unmarshal([]byte(s), &v); return v }

	nilDefinition := &config.ModuleDefinition{Config: nil}
	emptyDefinition := &config.ModuleDefinition{Config: unmarshall("{}")}
	invalidDefinition := &config.ModuleDefinition{Config: unmarshall("{\"no_key_def\":true}")}

	failure_l57 := definitionFactoryImpl("{\"rpm.max\":1200}")
	failure_l62 := definitionFactoryImpl("{\"rpm.min\":\"0\",\"rpm.max\":\"1200\",\"rpm.step\":\"250\",\"rpm.precision\":\"1000\"}")
	failure_l70 := definitionFactoryImpl(nil)

	Expect(module.Configure(nilDefinition)).Should(module_test.ExpectFor(module).Panic())     // failed: NIL definition
	Expect(module.Configure(emptyDefinition)).Should(module_test.ExpectFor(module).Panic())   // failed: empty definition
	Expect(module.Configure(invalidDefinition)).Should(module_test.ExpectFor(module).Panic()) // failed: invalid definition

	Expect(module.Configure(failure_l57)).Should(module_test.ExpectFor(module).Panic()) // failure at l.57
	Expect(module.Configure(failure_l62)).Should(module_test.ExpectFor(module).Panic()) // failure at l.63
	Expect(module.Configure(failure_l70)).Should(module_test.ExpectFor(module).Panic()) // failure at l.70

}

func postTestImpl(t *testing.T, nprocesses int, module module.Module, queue *module.NotificationQueue) {
	notifications := queue.Notifications()

	Expect(notifications).Should(HaveLen(nprocesses))

	for _, n := range notifications {
		Expect(n).Should(WithTransform(
			func(n *notification.Notification) string { return n.Content().(*object.DataObject).Value.(string) },
			HaveLen(10),
		))
	}
}
