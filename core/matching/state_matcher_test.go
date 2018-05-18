package matching

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

func Test_StateMatcher_houldMatchIfBothCurrentAndRequiredStateAreNil(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(nil, nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.Score).To(Equal(0))
}

func Test_StateMatcher_ShouldMatchIfCurrentStateIsNilAndRequiredStateIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(nil, make(map[string]string))

	Expect(match.Matched).To(BeTrue())
	Expect(match.Score).To(Equal(0))
}

func Test_StateMatcher_ShouldMatchIfCurrentStateIEmptyAndRequiredStateIsNil(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(&state.State{State: make(map[string]string)}, nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.Score).To(Equal(0))
}

func Test_StateMatcher_ShouldNotMatchIfRequiredStateLengthIsGreaterThanActualStateLength(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(&state.State{State: make(map[string]string)}, map[string]string{"foo": "bar"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.Score).To(Equal(0))
}

func Test_StateMatcher_ShouldNotMatchIfLengthsAreTheSameButKeysAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(
		&state.State{State: map[string]string{"foo": "bar", "cheese": "ham"}},
		map[string]string{"adasd": "bar", "sadsad": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.Score).To(Equal(0))
}

func Test_StateMatcher_ShouldNotMatchIfKeysAreTheSameButValuesAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(
		&state.State{State: map[string]string{"foo": "bar", "cheese": "ham"}},
		map[string]string{"foo": "adsad", "cheese": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.Score).To(Equal(1))
}

func Test_StateMatcher_ShouldMatchIsKeysAndValuesAreTheSame(t *testing.T) {
	RegisterTestingT(t)

	match := StateMatcher(
		&state.State{State: map[string]string{"foo": "bar", "cheese": "ham"}},
		map[string]string{"foo": "bar", "cheese": "ham"})

	Expect(match.Matched).To(BeTrue())
	Expect(match.Score).To(Equal(2))
}
