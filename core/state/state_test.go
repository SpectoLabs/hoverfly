package state_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

func Test_NewStateFromState_ReturnsEmptyState(t *testing.T) {
	RegisterTestingT(t)

	b := state.NewStateFromState(map[string]string{})
	Expect(b).ToNot(BeNil())
	Expect(b.State).To(Equal(map[string]string{}))
}

func Test_NewStateFromState_InitializesStateForSequenceKeys(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewStateFromState(map[string]string{
		"test":       "test",
		"sequence:0": "true",
		"sequence:1": "true",
	})
	Expect(s).ToNot(BeNil())
	Expect(s.State).To(Equal(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
	}))
}

func Test_State_GetNewSequenceKey_ReturnsFirstFreeInSequence(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:0"))

	s.SetState(map[string]string{
		"sequence:0": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:1"))

	s.SetState(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:2"))

	s.SetState(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
		"sequence:2": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:3"))
}
