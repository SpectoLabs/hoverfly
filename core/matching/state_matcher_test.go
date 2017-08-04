package matching

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_UnscoredStateMatcher_ShouldMatchIfBothCurrentAndRequiredStateAreNil(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(nil, nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldMatchIfCurrentStateIsNilAndRequiredStateIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(nil, make(map[string]string))

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldMatchIfCurrentStateIEmptyAndRequiredStateIsNil(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(make(map[string]string), nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldNotMatchIfRequiredStateLengthIsGreaterThanActualStateLength(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(make(map[string]string), map[string]string{"foo": "bar"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldNotMatchIfLengthsAreTheSameButKeysAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"adasd": "bar", "sadsad": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldNotMatchIfKeysAreTheSameButValuesAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "adsad", "cheese": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_UnscoredStateMatcher_ShouldMatchIsKeysAndValuesAreTheSame(t *testing.T) {
	RegisterTestingT(t)

	match := UnscoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "bar", "cheese": "ham"})

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_houldMatchIfBothCurrentAndRequiredStateAreNil(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(nil, nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_ShouldMatchIfCurrentStateIsNilAndRequiredStateIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(nil, make(map[string]string))

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_ShouldMatchIfCurrentStateIEmptyAndRequiredStateIsNil(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(make(map[string]string), nil)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_ShouldNotMatchIfRequiredStateLengthIsGreaterThanActualStateLength(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(make(map[string]string), map[string]string{"foo": "bar"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_ShouldNotMatchIfLengthsAreTheSameButKeysAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"adasd": "bar", "sadsad": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(0))
}

func Test_ScoredStateMatcher_ShouldNotMatchIfKeysAreTheSameButValuesAreDifferent(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "adsad", "cheese": "ham"})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(1))
}

func Test_ScoredStateMatcher_ShouldMatchIsKeysAndValuesAreTheSame(t *testing.T) {
	RegisterTestingT(t)

	match := ScoredStateMatcher(
		map[string]string{"foo": "bar", "cheese": "ham"},
		map[string]string{"foo": "bar", "cheese": "ham"})

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(2))
}
