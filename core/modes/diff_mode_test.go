package modes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/SpectoLabs/hoverfly/core/errors"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

type hoverflyDiffStub struct{}

func (this hoverflyDiffStub) DoRequest(request *http.Request) (*http.Response, error) {
	switch request.Host {
	case "error.com":
		return nil, fmt.Errorf("Could not reach error.com")
	case "positive-match-with-same-response.com":
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("expected")),
			Header:     map[string][]string{"header": {"expected"}, "source": {"service"}},
		}, nil
	case "positive-match-with-different-response.com":
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("actual")),
			Header:     map[string][]string{"header": {"actual"}, "source": {"service"}},
		}, nil
	case "negative-match.com":
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("actual")),
			Header:     map[string][]string{"header": {"actual"}, "source": {"service"}},
		}, nil
	case "positive-match-with-different-trailers.com":
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("actual")),
			Header:     map[string][]string{"header": {"actual"}},
			Trailer:	map[string][]string{"trailer1": {"actual"}},
		}, nil
	default:
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
		}, nil
	}
}

func (this hoverflyDiffStub) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *errors.HoverflyError) {
	switch requestDetails.Destination {
	case "positive-match-with-same-response.com":
		return &models.ResponseDetails{
			Status:  200,
			Body:    "expected",
			Headers: map[string][]string{"header": {"expected"}},
		}, nil
	case "positive-match-with-different-trailers.com":
		return &models.ResponseDetails{
			Status:  200,
			Body:    "actual",
			Headers: map[string][]string{"header": {"simulated"}, "Trailer": {"trailer1"}, "trailer1": {"simulated"}},
		}, nil
	case "positive-match-with-different-response.com":
		return &models.ResponseDetails{
			Status:  200,
			Body:    "simulated",
			Headers: map[string][]string{"header": {"simulated"}, "source": {"simulation"}},
		}, nil
	default:
		return nil, &errors.HoverflyError{
			Message: "matching-error",
		}
	}
}

func (this hoverflyDiffStub) AddDiff(key v2.SimpleRequestDefinitionView, diffReport v2.DiffReport) {
}

func Test_DiffMode_WhenGivenAMatchingRequestReturningTheSameResponse(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
	}

	request := models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "positive-match-with-same-response.com",
		Path:        "/",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("expected"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"expected"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(len(unit.DiffReport.DiffEntries)).To(Equal(0))
}

func Test_DiffMode_WhenGivenAMatchingRequestReturningDifferentResponse(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-different-response.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("actual"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(unit.DiffReport.DiffEntries).To(ConsistOf(
		v2.DiffReportEntry{Field: "header/source", Expected: "[simulation]", Actual: "[service]"},
		v2.DiffReportEntry{Field: "header/header", Expected: "[simulated]", Actual: "[actual]"},
		v2.DiffReportEntry{Field: "body", Expected: "simulated", Actual: "actual"}))
}

func Test_DiffMode_IncludeResponseTrailerForDiffing(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-different-trailers.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	Expect(len(response.Header)).To(Equal(1))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(unit.DiffReport.DiffEntries).To(ConsistOf(
		v2.DiffReportEntry{Field: "header/header", Expected: "[simulated]", Actual: "[actual]"},
		v2.DiffReportEntry{Field: "header/trailer1", Expected: "[simulated]", Actual: "[actual]"},
		))
}

func Test_DiffMode_BlacklistAllHeaders(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
		Arguments: ModeArguments{
			Headers: []string{
				"*",
			},
		},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-different-response.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("actual"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(unit.DiffReport.DiffEntries).To(ConsistOf(
		v2.DiffReportEntry{"body", "simulated", "actual"}))
}

func Test_DiffMode_BlacklistOneHeader(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
		Arguments: ModeArguments{
			Headers: []string{
				"header",
			},
		},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-different-response.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("actual"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(unit.DiffReport.DiffEntries).To(ConsistOf(
		v2.DiffReportEntry{"header/source", "[simulation]", "[service]"},
		v2.DiffReportEntry{"body", "simulated", "actual"}))
}

func Test_DiffMode_BlacklistlistTwoHeaders(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
		Arguments: ModeArguments{
			Headers: []string{
				"header", "source",
			},
		},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-different-response.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("actual"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(unit.DiffReport.DiffEntries).To(ConsistOf(
		v2.DiffReportEntry{"body", "simulated", "actual"}))
}

func Test_DiffMode_WhenGivenANonMatchingRequestDiffIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Method:      "GET",
		Destination: "negative-match.com",
	}

	// when
	response, err := unit.Process(nil, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("actual"))
	Expect(len(response.Header)).To(Equal(2))
	Expect(response.Header["header"]).To(Equal([]string{"actual"}))
	Expect(response.Header["source"]).To(Equal([]string{"service"}))
	Expect(unit.DiffReport.DiffEntries).To(BeEmpty())
}

func Test_JsonDiff_WhenDifferentThenCreatesErrorMessage(t *testing.T) {
	RegisterTestingT(t)

	//when
	expected := []byte(`{
	"foo": "bar",
	"fooInt": 1,
	"fooDouble": 0,
	"fooBool": true,
	"anotherExpFoo": "foo",
	"nested": {
		"baz": "boo"
	}}`)
	actual := []byte(`{
	"foo": "baz",
	"fooInt": 2,
	"fooDouble": 0.1,
	"fooBool": false,
	"anotherActFoo": "foo",
	"nested": {
		"baz": "bar"
	}}`)

	var jsonExpected interface{}
	var jsonActual interface{}
	_ = json.Unmarshal(expected, &jsonExpected)
	_ = json.Unmarshal(actual, &jsonActual)

	diffMode := DiffMode{DiffReport: v2.DiffReport{}}

	// when
	result := diffMode.JsonDiff("test", jsonExpected.(map[string]interface{}), jsonActual.(map[string]interface{}))

	// then
	Expect(result).To(Equal(false))
	Expect(len(diffMode.DiffReport.DiffEntries)).To(Equal(6))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/foo", "bar", "baz"}))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/fooInt", "1", "2"}))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/fooDouble", "0", "0.1"}))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/fooBool", "true", "false"}))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/anotherExpFoo", "foo", "null"}))
	Expect(diffMode.DiffReport.DiffEntries).To(ContainElement(
		v2.DiffReportEntry{"test/nested/baz", "boo", "bar"}))
}

func Test_JsonDiff_WhenExpectedEmptyThenReturnsTrue(t *testing.T) {
	RegisterTestingT(t)

	// given
	actual := []byte(`{
	"foo": "bar",
	"bar": {
		"baz": "xyzzy"
	},
	"xyzzy": [1,2]
	}`)

	expected := []byte(`{}`)

	var jsonActual interface{}
	_ = json.Unmarshal(actual, &jsonActual)
	var jsonExpected interface{}
	_ = json.Unmarshal(expected, &jsonExpected)

	diffMode := DiffMode{DiffReport: v2.DiffReport{}}

	// when
	result := diffMode.JsonDiff("test", jsonExpected.(map[string]interface{}), jsonActual.(map[string]interface{}))

	// then
	Expect(result).To(Equal(true))
	Expect(len(diffMode.DiffReport.DiffEntries)).To(Equal(0))
}

func Test_JsonDiff_WhenSameThenReturnsTrue(t *testing.T) {
	RegisterTestingT(t)

	// given
	data := []byte(`{
	"foo": "bar",
	"bar": {
		"baz": "xyzzy"
	},
	"xyzzy": [1,2]
	}`)

	var jsonObject interface{}
	_ = json.Unmarshal(data, &jsonObject)

	diffMode := DiffMode{DiffReport: v2.DiffReport{}}

	// when
	result := diffMode.JsonDiff("test", jsonObject.(map[string]interface{}), jsonObject.(map[string]interface{}))

	// then
	Expect(result).To(Equal(true))
	Expect(len(diffMode.DiffReport.DiffEntries)).To(Equal(0))
}
