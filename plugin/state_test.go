package plugin

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestState_Properties(t *testing.T) {
	RegisterTestingT(t)

	cases := []struct {
		state      State
		code       byte
		desc       string
		rawMatcher OmegaMatcher
	}{
		{NewState(0x30, "custom state"), byte(0x30), "custom state", BeNil()},
		{NewState(IdleState.Code(), "rewrote state", "with raw"), IdleState.Code(), "rewrote state", Equal("with raw")},
		{IdleState.AddRaw("new raw"), IdleState.Code(), IdleState.Desc(), Equal("new raw")},
	}

	for _, test := range cases {
		Expect(test.state.Code()).Should(Equal(test.code))
		Expect(test.state.Desc()).Should(Equal(test.desc))
		Expect(test.state.Raw()).Should(test.rawMatcher)
	}
}
func TestState_Equal(t *testing.T) {
	RegisterTestingT(t)

	custom := NewState(0x30, "custom state")
	rewrote := NewState(IdleState.Code(), "rewrote state", "with raw")
	edited := IdleState.AddRaw("new raw")

	Expect(IdleState.Equal("...")).Should(BeFalse())
	Expect(IdleState.Equal(custom)).Should(BeFalse())
	Expect(IdleState.Equal(rewrote)).Should(BeTrue()) // equality is only based on the state code. DO NOT USE AN EXISTING STATE CODE
	Expect(IdleState.Equal(edited)).Should(BeTrue())
}
func TestState_Immutability(t *testing.T) {
	RegisterTestingT(t)

	edited := IdleState.AddRaw("new raw")

	edited.code = 0x00
	Expect(edited.Code()).Should(Equal(byte(0x00)))
	Expect(edited.Code()).ShouldNot(Equal(IdleState.Code()))
}