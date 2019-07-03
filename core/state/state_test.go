package state_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

func Test_InitializeSequences_InitializesStateForSequenceKeys(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	s.InitializeSequences(map[string]string{
		"test":       "test",
		"sequence:0": "true",
		"sequence:1": "true",
	})
	Expect(s.State).To(Equal(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
	}))
}

func Test_InitializeSequences_ResetExistingSequenceKeys(t *testing.T) {
	RegisterTestingT(t)

	s := state.NewState()
	s.SetState(map[string]string{
		"sequence:0": "2",
		"sequence:1": "3",
		"sequence:2": "4",
	})
	s.InitializeSequences(map[string]string{
		"test":       "test",
		"sequence:0": "true",
		"sequence:1": "true",
	})
	Expect(s.State).To(Equal(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
		"sequence:2": "4",
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
	Eventually(func() string {
			val, _ := s.GetState("test1")
			return val
	}).Should(Equal("1"))
	Eventually(func() string {
			val, _ := s.GetState("test2")
			return val
	}).Should(Equal("2"))
	Eventually(func() string {
			val, _ := s.GetState("test3")
			return val
	}).Should(Equal("3"))
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
