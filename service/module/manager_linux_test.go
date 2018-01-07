package module

import (
	"github.com/sportfun/gakisitor/config"
	"fmt"
	"github.com/sportfun/gakisitor/env"
	"github.com/sportfun/gakisitor/log"
	"testing"
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/service"
)

var (
	validConfiguration = config.GAkisitorConfig{
		Modules: []config.ModuleDefinition{
			{
				Name: "RPM Generator",
				Path: "../../.resources/example.so",
				Config: map[string]interface{}{
					"rpm.min":       100.,
					"rpm.max":       250.,
					"rpm.step":      10.,
					"rpm.precision": 1000.,
				},
			},
		},
	}
	noFile = config.GAkisitorConfig{
		Modules: []config.ModuleDefinition{
			{
				Name: "RPM Generator",
				Path: "plugin.so",
				Config: map[string]interface{}{
					"rpm.min":       100.,
					"rpm.max":       250.,
					"rpm.step":      10.,
					"rpm.precision": 1000.,
				},
			},
		},
	}
	noExport = config.GAkisitorConfig{
		Modules: []config.ModuleDefinition{
			{
				Name: "RPM Generator",
				Path: "../../.resources/example_no_export.so",
				Config: map[string]interface{}{
					"rpm.min":       100.,
					"rpm.max":       250.,
					"rpm.step":      10.,
					"rpm.precision": 1000.,
				},
			},
		},
	}
	noConfig = config.GAkisitorConfig{
		Modules: []config.ModuleDefinition{
			{
				Name: "RPM Generator",
				Path: "../../.resources/example.so",
			},
		},
	}
)

func TestManager_Configure(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := log.NewTest(&buffer)

	testCases := []struct {
		Configuration config.GAkisitorConfig
		State         byte
		Error         error
		Log           string
	}{
		{Configuration: validConfiguration, State: env.IdleState, Error: nil, Log: "DEBUG	module 'RPM Generator' configured"},

		{Configuration: noFile, State: env.PanicState, Error: fmt.Errorf("no module charged"), Log: ""},
		{Configuration: noExport, State: env.PanicState, Error: fmt.Errorf("no module charged"), Log: ""},
		{Configuration: noConfig, State: env.PanicState, Error: fmt.Errorf("no module charged"), Log: ""},
	}

	for _, tc := range testCases {
		buffer = ""
		manager := Manager{}

		Expect(manager.Start(service.NewNotificationQueue(), logger)).Should(Succeed())

		switch {
		case tc.Error == nil:
			Expect(manager.Configure(&tc.Configuration)).Should(Succeed())
		default:
			Expect(manager.Configure(&tc.Configuration)).Should(MatchError(tc.Error))
		}

		Expect(buffer).Should(MatchRegexp(tc.Log))
		Expect(manager.State()).Should(Equal(tc.State))
	}
}
