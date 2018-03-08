package plugin

import (
	"testing"

	. "github.com/onsi/gomega"
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
