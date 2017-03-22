package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_HeaderMatcher(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	Expect(matching.HeaderMatcher(matcherHeaders, matcherHeaders)).To(BeTrue())
}

func Test_HeaderMatcher_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"Header2": []string{"VAL2"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders)).To(BeTrue())
}

func Test_HeaderMatcher_MatchingHeadersHasMoreKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": []string{"val1"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders)).To(BeFalse())
}

func Test_HeaderMatcher_MatchingHeadersHasMoreValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders)).To(BeFalse())
}

func Test_HeaderMatch_MatchingHeadersHasLessKeysThanRequestMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"header2": []string{"val2"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders)).To(BeTrue())
}

func Test_HeaderMatch_RequestHeadersContainsAllOFMatcherHeaders(t *testing.T) {
	RegisterTestingT(t)

	matcherHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	Expect(matching.HeaderMatcher(matcherHeaders, reqHeaders)).To(BeTrue())
}
