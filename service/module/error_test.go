package module

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
	"testing"
)

func TestManager_PluginFailure(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	manager := Manager{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue()}
	builder := notification.NewBuilder().From(env.ModuleServiceEntity).To(env.NetworkServiceEntity)

	testCases := []struct {
		Type  pluginError
		Error error
		Name  string

		Reason string
	}{
		{Type: NoModule, Error: fmt.Errorf("error#1"), Name: "plg#1", Reason: "error#1"},
		{Type: PluginLoading, Error: fmt.Errorf("error#2"), Name: "plg#2", Reason: "failure during plugin loading ('plg#2'): error#2"},
		{Type: SymbolLoading, Error: fmt.Errorf("error#3"), Name: "plg#3", Reason: "failure during symbol loading ('plg#3'): error#3"},
		{Type: 0xff, Error: fmt.Errorf("error#4"), Name: "plg#4", Reason: "unknown error type from ('plg#4'): error#4"},
	}

	for _, tc := range testCases {
		buffer = ""
		manager.pluginFailure(tc.Type, tc.Error, tc.Name)

		ntfns := manager.notifications.Notifications(env.NetworkServiceEntity)
		Expect(ntfns).Should(And(
			HaveLen(1),
			ContainElement(builder.With(object.NewErrorObject(env.ModuleServiceEntity, fmt.Errorf(tc.Reason))).Build()),
		))
		Expect(buffer).Should(Equal("ERROR	" + tc.Reason))
	}
}

func TestManager_ModuleError(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	manager := &Manager{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue()}
	module := &tModule{name: "module name"}
	obj := object.NewErrorObject("")

	testCases := []struct {
		Fnc    moduleError
		Error  error
		Reason string
	}{
		{Fnc: isStarted, Error: fmt.Errorf("err#1"), Reason: "failure during module ('module name') starting: err#1"},
		{Fnc: isConfigured, Error: fmt.Errorf("err#2"), Reason: "failure during module ('module name') configuration: err#2"},
		{Fnc: isProcessed, Error: fmt.Errorf("err#3"), Reason: "failure during module ('module name') processing: err#3"},
		{Fnc: isStopped, Error: fmt.Errorf("err#4"), Reason: "failure during module ('module name') stopping: err#4"},
	}

	for _, tc := range testCases {
		tc.Fnc(manager, module, tc.Error, obj)
		Expect(obj.Reason).Should(Equal(tc.Reason))
	}
}

func TestManager_CheckIf(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	manager := &Manager{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue()}
	module := &tModule{name: "module name"}
	builder := notification.NewBuilder().From(env.ModuleServiceEntity).To(env.NetworkServiceEntity)

	testCases := []struct {
		Fnc    moduleError
		Error  error
		Reason string
	}{
		{Fnc: nil, Error: nil, Reason: ""},
		{Fnc: isStarted, Error: fmt.Errorf("err#1"), Reason: "failure during module ('module name') starting: err#1"},
		{Fnc: isConfigured, Error: fmt.Errorf("err#2"), Reason: "failure during module ('module name') configuration: err#2"},
		{Fnc: isProcessed, Error: fmt.Errorf("err#3"), Reason: "failure during module ('module name') processing: err#3"},
		{Fnc: isStopped, Error: fmt.Errorf("err#4"), Reason: "failure during module ('module name') stopping: err#4"},
	}

	for _, tc := range testCases {
		buffer = ""
		manager.checkIf(module, tc.Error, tc.Fnc)

		ntfns := manager.notifications.Notifications(env.NetworkServiceEntity)
		switch {
		case tc.Fnc == nil:
			Expect(ntfns).Should(BeEmpty())
		default:
			Expect(ntfns).Should(And(
				HaveLen(1),
				ContainElement(builder.With(object.NewErrorObject(env.ModuleServiceEntity, fmt.Errorf(tc.Reason))).Build()),
			))
		}
	}

}
