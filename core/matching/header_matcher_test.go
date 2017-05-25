package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_HeaderMatcher(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	Expect(matching.HeaderMatcher(matcherHeaders, matcherHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatcher_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": {"val1"},
		"Header2": {"VAL2"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatcher_MatchingHeadersHasMoreKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": {"val1"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
}

func Test_HeaderMatcher_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeFalse())
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
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}

func Test_HeaderMatch_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": {"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": {"val1", "val2"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders).Matched).To(BeTrue())
}
