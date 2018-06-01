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

func Test_GetState(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	s.State = map[string]string{
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	Expect(s).ToNot(BeNil())
	Expect(s.GetState("test1")).To(Equal("1"))
	Expect(s.GetState("test2")).To(Equal("2"))
	Expect(s.GetState("test3")).To(Equal("3"))
}

func Test_PatchState(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	s.State = map[string]string{
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	s.PatchState(map[string]string{
		"test1": "modified",
	})

	Expect(s).ToNot(BeNil())
	Expect(s.State).To(Equal(map[string]string{
		"test1": "modified",
		"test2": "2",
		"test3": "3",
	}))

	s.PatchState(map[string]string{
		"test2": "modified",
		"test3": "modified",
	})

	Expect(s).ToNot(BeNil())
	Expect(s.State).To(Equal(map[string]string{
		"test1": "modified",
		"test2": "modified",
		"test3": "modified",
	}))
}

func Test_RemoveState(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	s.State = map[string]string{
		"test1": "1",
		"test2": "2",
		"test3": "3",
	}

	s.RemoveState([]string{"test1"})

	Expect(s).ToNot(BeNil())
	Expect(s.State).To(Equal(map[string]string{
		"test2": "2",
		"test3": "3",
	}))

	s.RemoveState([]string{"test2", "test3"})

	Expect(s).ToNot(BeNil())
	Expect(s.State).To(Equal(map[string]string{}))
}

func Test_State_GetNewSequenceKey_ReturnsFirstFreeInSequence(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:1"))

	s.SetState(map[string]string{
		"sequence:1": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:2"))

	s.SetState(map[string]string{
		"sequence:1": "1",
		"sequence:2": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:3"))

	s.SetState(map[string]string{
		"sequence:1": "1",
		"sequence:2": "1",
		"sequence:3": "1",
	})
	Expect(s.GetNewSequenceKey()).To(Equal("sequence:4"))
}
