package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_HeaderMatching(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, matcherHeaders).Matched).To(BeTrue())
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
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeTrue())
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
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatching_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1"},
	}
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeFalse())
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
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatching_ShouldNotMatchUnlessAllHeadersValuesAreFound(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "nomatch"},
	}
	Expect(matching.HeaderMatching(models.RequestMatcher{
		Headers: matcherHeaders,
	}, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatching_CountsMatches_WhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	match := matching.HeaderMatching(models.RequestMatcher{
		Headers: map[string][]string{
			"header1": {"val1", "val2"},
		},
	},

		map[string][]string{
			"header1":     {"val1", "val2"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(2))

	match = matching.HeaderMatching(models.RequestMatcher{
		Headers: map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3"},
		},
	},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(3))
}

func Test_HeaderMatching_CountsMatches_WhenThereIsNoMatch(t *testing.T) {
	RegisterTestingT(t)

	match := matching.HeaderMatching(models.RequestMatcher{
		Headers: map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3", "nomatch"},
		},
	},
		map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		})

	Expect(match.Matched).To(BeFalse())
	Expect(match.MatchScore).To(Equal(3))
}

func Test_HeaderMatching_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	match := matching.ScoredFieldMatcher(nil, `testtesttest`)

	Expect(match.Matched).To(BeTrue())
	Expect(match.MatchScore).To(Equal(0))
}
