package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_HeaderMatching(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	Expect(matching.HeaderMatching(matcherHeaders, matcherHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"Header2": {"VAL2"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_MatchingHeadersHasMoreKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": {"val1"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatching_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatching_MatchingHeadersHasLessKeysThanRequestMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"header2": {"val2"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_ShouldNotMatchUnlessAllHeadersValuesAreFound(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "nomatch"},
	}
	Expect(matching.HeaderMatching(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatching_CountsMatches_WhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	matcher := matching.HeaderMatching(
		map[string][]string{
			"header1": {"val1", "val2"},
		},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(2))

	matcher = matching.HeaderMatching(
		map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3"},
		},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(3))
}

func Test_HeaderMatching_CountsMatches_WhenThereIsNoMatch(t *testing.T) {
	RegisterTestingT(t)

	matcher := matching.HeaderMatching(
		map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3", "nomatch"},
		},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(matcher.Matched).To(BeFalse())
	Expect(matcher.MatchScore).To(Equal(3))
}

func Test_HeaderMatching_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(nil, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(0))
}
