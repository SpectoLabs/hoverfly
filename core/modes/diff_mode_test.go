package modes

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"bytes"
	"fmt"
	"encoding/json"
)

type hoverflyDiffStub struct{}

func (this hoverflyDiffStub) DoRequest(request *http.Request) (*http.Response, error) {
	response := &http.Response{}
	if request.Host == "error.com" {
		return nil, errors.New("Could not reach error.com")
	}

	if request.Host == "positive-match-with-same-response.com" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("expected")),
			Header:     map[string][]string{"header": {"expected"}, "source": {"service"}},
		}, nil
	} else if request.Host == "positive-match-with-different-response.com" || request.Host == "negative-match.com" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("actual")),
			Header:     map[string][]string{"header": {"actual"}, "source": {"service"}},
		}, nil
	}

	response.StatusCode = 200
	response.Body = ioutil.NopCloser(bytes.NewBufferString("test"))

	return response, nil
}

func (this hoverflyDiffStub) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError) {
	if requestDetails.Destination == "positive-match-with-same-response.com" {
		return &models.ResponseDetails{
			Status:  200,
			Body:    "expected",
			Headers: map[string][]string{"header": {"expected"}},
		}, nil
	} else if requestDetails.Destination == "positive-match-with-different-response.com" {
		return &models.ResponseDetails{
			Status:  200,
			Body:    "simulated",
			Headers: map[string][]string{"header": {"simulated"}, "source": {"simulation"}},
		}, nil
	} else {
		return nil, &matching.MatchingError{
			Description: "matching-error",
			StatusCode:  500,
		}
	}
}

func Test_DiffMode_WhenGivenAMatchingRequestReturningTheSameResponse(t *testing.T) {
	RegisterTestingT(t)

	//given
	unit := &DiffMode{
		Hoverfly: hoverflyDiffStub{},
	}

	request := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match-with-same-response.com",
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
	message := unit.GetMessage().DiffMessage
	Expect(message.String()).To(Equal(""))
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
	verifyParamsDiffAreInMessage(unit.GetMessage(),
		paramDiff{"header/source", "[simulation]", "[service]"},
		paramDiff{"header/header", "[simulated]", "[actual]"},
		paramDiff{"body", "simulated", "actual"})
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
	message := unit.GetMessage().DiffMessage
	Expect(message.String()).To(Equal(""))
}

type paramDiff struct {
	param    string
	expected string
	actual   string
}

func TestJsonDiff_WhenDifferentThenCreatesErrorMessage(t *testing.T) {
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

	diffMessage := DiffErrorMessage{}

	// when
	result := JsonDiff(&diffMessage, "test", jsonExpected.(map[string]interface{}), jsonActual.(map[string]interface{}))

	// then
	Expect(result).To(Equal(false))
	Expect(diffMessage.Counter).To(Equal(6))
	Expect(diffMessage.DiffMessage.String()).NotTo(Equal(""))
	verifyParamsDiffAreInMessage(diffMessage,
		paramDiff{"test/foo", "bar", "baz"},
		paramDiff{"test/fooInt", "1", "2"},
		paramDiff{"test/fooDouble", "0", "0.1"},
		paramDiff{"test/fooBool", "true", "false"},
		paramDiff{"test/fooBool", "true", "false"},
		paramDiff{"test/anotherExpFoo", "foo", "undefined"},
		paramDiff{"test/nested/baz", "boo", "bar"})
}

func TestJsonDiff_WhenExpectedEmptyThenReturnsTrue(t *testing.T) {
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

	diffMessage := DiffErrorMessage{}

	// when
	result := JsonDiff(&diffMessage, "test", jsonExpected.(map[string]interface{}), jsonActual.(map[string]interface{}))

	// then
	Expect(result).To(Equal(true))
	Expect(diffMessage.Counter).To(Equal(0))
	Expect(diffMessage.DiffMessage.String()).To(Equal(""))
}

func TestJsonDiff_WhenSameThenReturnsTrue(t *testing.T) {
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

	diffMessage := DiffErrorMessage{}

	// when
	result := JsonDiff(&diffMessage, "test", jsonObject.(map[string]interface{}), jsonObject.(map[string]interface{}))

	// then
	Expect(result).To(Equal(true))
	Expect(diffMessage.Counter).To(Equal(0))
	Expect(diffMessage.DiffMessage.String()).To(Equal(""))
}

func verifyParamsDiffAreInMessage(diffMessage DiffErrorMessage, params ...paramDiff) {
	for _, param := range params {
		Expect(diffMessage.DiffMessage.String()).To(
			ContainSubstring(fmt.Sprintf(errorMsgTemplate, param.param, param.expected, param.actual)))
	}
}
