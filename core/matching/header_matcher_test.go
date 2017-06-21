package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_CountlessHeaderMatcher(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	Expect(matching.CountlessHeaderMatcher(matcherHeaders, matcherHeaders).Matched).To(BeTrue())
}

func Test_CountlessHeaderMatcher_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"Header2": {"VAL2"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_CountlessHeaderMatcher_MatchingHeadersHasMoreKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": {"val1"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_CountlessHeaderMatcher_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatch_MatchingHeadersHasLessKeysThanRequestMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"header2": {"val2"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatch_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatch_ShouldNotMatchUnlessAllHeadersValuesAreFound(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "nomatch"},
	}
	Expect(matching.CountlessHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_CountingHeaderMatcher(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	Expect(matching.CountingHeaderMatcher(matcherHeaders, matcherHeaders).Matched).To(BeTrue())
}

func Test_CountingHeaderMatcher_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"Header2": {"VAL2"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_CountingHeaderMatcher_MatchingHeadersHasMoreKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": {"val1"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_CountingHeaderMatcher_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_CountingHeaderMatcher_MatchingHeadersHasLessKeysThanRequestMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"header2": {"val2"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_CountingHeaderMatcher_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_CountingHeaderMatcher_ShouldNotMatchUnlessAllHeadersValuesAreFound(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "nomatch"},
	}
	Expect(matching.CountingHeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_CountingHeaderMatcher_CountsMatches_WhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	matcher := matching.CountingHeaderMatcher(
		map[string][]string{
			"header1": {"val1", "val2"},
		},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(2))

	matcher = matching.CountingHeaderMatcher(
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

func Test_CountingHeaderMatcher_CountsMatches_WhenThereIsNoMatch(t *testing.T) {
	RegisterTestingT(t)

	matcher := matching.CountingHeaderMatcher(
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

func Test_CountingHeaderMatcher_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(nil, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(0))
}
