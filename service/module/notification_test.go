package module

import (
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/notification"
	"github.com/sportfun/gakisitor/notification/object"
	"github.com/sportfun/gakisitor/service"
	"testing"
)

func TestManager_Notification(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	module := &tModule{name: "module"}
	manager := &Manager{logger: log.NewTest(&buffer), notifications: service.NewNotificationQueue(), modules: moduleContainer{"module": module}}
	builder := notification.NewBuilder().From(env.NetworkServiceEntity).To(env.ModuleServiceEntity)

	testCases := []struct {
		Notification *notification.Notification
		Predicate    func() bool
		Log          string
	}{
		{Notification: builder.With(object.NewCommandObject(env.StartSessionCmd)).Build(), Predicate: func() bool { return module.InSession }, Log: "DEBUG	start session"},
		{Notification: builder.With(object.NewCommandObject(env.EndSessionCmd)).Build(), Predicate: func() bool { return !module.InSession }, Log: "DEBUG	stop session"},

		{Notification: nil, Predicate: func() bool { return true }, Log: "WARN	unexpected nil notification"},
		{Notification: notification.NewNotification("...", "", nil), Predicate: func() bool { return true }, Log: "WARN	unhandled notification origin (...)"},
		{Notification: builder.With("none").Build(), Predicate: func() bool { return true }, Log: `ERROR	invalid network notification	{"content": "none"}`},
		{Notification: builder.With(object.NewCommandObject("potatoes")).Build(), Predicate: func() bool { return true }, Log: "ERROR	unknown network command (potatoes)"},
	}

	for _, tc := range testCases {
		buffer = ""
		manager.handle(tc.Notification)

		Expect(tc.Predicate()).Should(BeTrue())
		Expect(buffer).Should(Equal(tc.Log))
	}
}
