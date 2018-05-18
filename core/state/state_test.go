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

	b := state.NewStateFromState(map[string]string{
		"test":       "test",
		"sequence:0": "true",
		"sequence:1": "true",
	})
	Expect(b).ToNot(BeNil())
	Expect(b.State).To(Equal(map[string]string{
		"sequence:0": "1",
		"sequence:1": "1",
	}))
}
