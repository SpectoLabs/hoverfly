package matching

import (
	"testing"
	. "github.com/onsi/gomega"
)

func Test_UnscoredShouldMatchIfBothCurrentAndRequiredStateAreNil(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(nil, nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldMatchIfCurrentStateIsNilAndRequiredStateIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(nil, make(map[string]string))

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldMatchIfCurrentStateIEmptyAndRequiredStateIsNil(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(make(map[string]string), nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldNotMatchIfRequiredStateLengthIsGreaterThanActualStateLength(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(make(map[string]string), map[string]string{"foo": "bar"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldNotMatchIfLengthsAreTheSameButKeysAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"adasd": "bar", "sadsad": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldNotMatchIfKeysAreTheSameButValuesAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "adsad", "cheese": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredShouldMatchIsKeysAndValuesAreTheSame(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "bar", "cheese": "ham"})

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}