// Package plugin_test provides tools used to tests plugins.
package plugin_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/plugin"
	"github.com/sportfun/gakisitor/profile"
)

// PluginTestDesc contains an example of a test plugin configuration and a
// matcher to test values, used by the PluginValidityChecker.
type PluginTestDesc struct {
	ConfigJSON   string
	ValueChecker OmegaMatcher
}

// PluginValidityChecker checks if the plugin works as expected
func PluginValidityChecker(t *testing.T, pluginInstance *plugin.Plugin, desc PluginTestDesc) {
	RegisterTestingT(t)

	ctx, cancel := context.WithCancel(context.Background())

	chData := make(chan interface{})
	chInst := make(chan plugin.Instruction)
	chStat := make(chan plugin.State)
	channels := plugin.Chan{
		Data:        chData,
		Instruction: chInst,
		Status:      chStat,
	}

	profile := profile.Plugin{
		Name:   pluginInstance.Name,
		Path:   "",
		Config: map[string]interface{}{},
	}
	Expect(json.Unmarshal([]byte(desc.ConfigJSON), &profile.Config)).To(Succeed())

	go func() {
		time.Sleep(15 * time.Second)
		panic("Plugin test has timeout (> 15s)")
	}()

	go func() {
		chInst <- plugin.StatusPluginInstruction
		goExpect(<-chStat, Equal(plugin.IdleState))

		chInst <- plugin.StartSessionInstruction
		chInst <- plugin.StatusPluginInstruction
		goExpect(<-chStat, Equal(plugin.InSessionState))

		// clean fist data
		<-chData
		for i := 0; i < 5; i++ {
			goExpect(<-chData, desc.ValueChecker)
		}

		chInst <- plugin.StopSessionInstruction
		chInst <- plugin.StatusPluginInstruction
		Expect(<-chStat).To(Equal(plugin.IdleState))

		cancel()
	}()

	Expect(pluginInstance.Instance(ctx, profile, channels)).Should(Succeed())
}

func goExpect(v interface{}, matcher OmegaMatcher) {
	if succeed, _ := matcher.Match(v); !succeed {
		panic(matcher.FailureMessage(v))
	}
}
