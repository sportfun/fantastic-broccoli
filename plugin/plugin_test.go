package plugin

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/log"
	"github.com/sportfun/gakisitor/profile"
)

func TestState_Properties(t *testing.T) {
	RegisterTestingT(t)

	custom := NewState(0x30, "custom state")
	rewrote := NewState(RunningState.code, "rewrote state", "with raw")
	edited := RunningState.AddRaw("new raw")

	// check simple creation
	Expect(custom.Code()).Should(Equal(byte(0x30)))
	Expect(custom.Desc()).Should(Equal("custom state"))
	Expect(custom.Raw()).Should(BeNil())

	// check state rewriting (unsafe, see bellow)
	Expect(rewrote.Code()).Should(Equal(RunningState.Code()))
	Expect(rewrote.Desc()).ShouldNot(Equal(RunningState.Desc()))
	Expect(rewrote.Raw()).Should(Equal("with raw"))

	// check state edition
	Expect(edited.Code()).Should(Equal(RunningState.Code()))
	Expect(edited.Desc()).Should(Equal(RunningState.Desc()))
	Expect(edited.Raw()).Should(Equal("new raw"))

	// check state equality
	Expect(RunningState.Equal("...")).Should(BeFalse())
	Expect(RunningState.Equal(custom)).Should(BeFalse())
	Expect(RunningState.Equal(rewrote)).Should(BeTrue()) // equality is only based on the state code. DO NOT USE AN EXISTING STATE CODE
	Expect(RunningState.Equal(edited)).Should(BeTrue())

	// check state immutability
	edited.code = 0x00
	Expect(edited.Code()).Should(Equal(byte(0x00)))
	Expect(edited.Code()).ShouldNot(Equal(RunningState.Code()))
}

func ExamplePlugin_basic() {
	_ = Plugin{
		Name: "ExamplePlugin",
		Instance: func(profile profile.Plugin, log log.Log, channels Chan) error {
			var inSession bool
			var state = RunningState
			dataMarshallable := struct {
				A int     `json:"a"`
				B float64 `json:"b"`
			}{0, 0.0}

			// configuration value shared into the profile (Config > ManyItems > ThisItem)
			_, e := profile.AccessTo("ManyItems", "ThisItem")
			if e != nil {
				return e
			}

			// plugin main loop
			for {
				select {
				// interpret instruction here
				case instruction, valid := <-channels.Instruction:
					// if channel is closed, you must stop the plugin
					if !valid {
						return nil
					}

					switch instruction {
					case StatusPluginInstruction:
						channels.Status <- state
					case StartSessionInstruction:
						inSession = true
						state = InSessionState
					case StopSessionInstruction:
						inSession = false
						state = RunningState
					case StopPluginInstruction:
						return nil
					}

					// example of data sending
				case <-time.Tick(time.Millisecond):
					if inSession {
						channels.Data <- dataMarshallable
					}
				}
			}

			return nil
		},
	}
}
