// Package plugin_test provides tools used to tests plugins.
package plugin_test

import (
	"context"
	"encoding/json"
	"sync/atomic"
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

	var isFinished int32
	go func() {
		Expect(pluginInstance.Instance(ctx, profile, channels)).To(Succeed())
		atomic.AddInt32(&isFinished, 1)
	}()

	chInst <- plugin.StatusPluginInstruction
	Expect(<-chStat).To(Equal(plugin.IdleState))

	chInst <- plugin.StartSessionInstruction
	chInst <- plugin.StatusPluginInstruction
	Expect(<-chStat).To(Equal(plugin.InSessionState))

	Expect(<-chData).To(desc.ValueChecker)
	Expect(<-chData).To(desc.ValueChecker)

	chInst <- plugin.StopSessionInstruction
	chInst <- plugin.StatusPluginInstruction
	Expect(<-chStat).To(Equal(plugin.IdleState))

	cancel()

	Eventually(func() bool { return atomic.LoadInt32(&isFinished) == 1 }, 5*time.Second).Should(BeTrue(), "Plugin Validity Checker as timeout (> 5s)")
}
