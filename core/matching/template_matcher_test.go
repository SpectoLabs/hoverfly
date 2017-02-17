package matching

import (
	"testing"

	. "github.com/onsi/gomega"
)

// func Test_Match_EmptyTemplateShouldMatchOnAnyRequest(t *testing.T) {
// 	RegisterTestingT(t)

// 	var templates []interfaces.RequestResponsePair

// 	templates = append(templates, models.RequestResponsePair{
// 		Request: models.RequestDetails{},
// 		Response: models.ResponseDetails{
// 			Body: "test-body",
// 		},
// 	})

// 	r := models.RequestDetails{
// 		Method:      "GET",
// 		Destination: "somehost.com",
// 		Headers: map[string][]string{
// 			"sdv": []string{"ascd"},
// 		},
// 	}
// 	result, _ := Match(r, false, templates)

// 	Expect(result.Body).To(Equal("test-body"))
// }

func Test_headerMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	res := headerMatch(tmplHeaders, tmplHeaders)
	Expect(res).To(BeTrue())
}

func Test_headerMatch_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"Header2": []string{"VAL2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func Test_HeaderMatch_MatchingTemplateHasMoreHeaderKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": []string{"val1"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeFalse())
}

func Test_headerMatch_MatchingTemplateHasMoreHeaderValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeFalse())
}

func Test_headerMatch_MatchingRequestHasMoreHeaderKeysThanTemplateMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"header2": []string{"val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func Test_headerMatch_MatchingRequestHasMoreHeaderValuesThanTemplateMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}
