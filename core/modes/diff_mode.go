package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/errors"

	log "github.com/sirupsen/logrus"

	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/dsnet/compress/brotli"
)

type HoverflyDiff interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *errors.HoverflyError)
	DoRequest(*http.Request) (*http.Response, error)
	AddDiff(requestView v2.SimpleRequestDefinitionView, diffReport v2.DiffReport)
}

type DiffMode struct {
	Hoverfly   HoverflyDiff
	DiffReport v2.DiffReport
	Arguments  ModeArguments
}

func (this *DiffMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Diff,
		Arguments: v2.ModeArgumentsView{
			Headers:          this.Arguments.Headers,
			MatchingStrategy: this.Arguments.MatchingStrategy,
			Stateful:         this.Arguments.Stateful,
		},
	}
}

func (this *DiffMode) SetArguments(arguments ModeArguments) {
	this.Arguments = arguments
}

//TODO: We should only need one of these two parameters
func (this *DiffMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	this.DiffReport = v2.DiffReport{Timestamp: time.Now().Format(time.RFC3339)}

	actualPair := models.RequestResponsePair{
		Request: details,
	}

	simResponse, simRespErr := this.Hoverfly.GetResponse(details)

	log.Info("Going to call real server")
	modifiedRequest, err := ReconstructRequest(actualPair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &actualPair, "There was an error when reconstructing the request.", Diff)
	}

	actualResponse, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &actualPair, "There was an error when forwarding the request to the intended destination", Diff)
	}

	if this.Arguments.Headers == nil {
		this.Arguments.Headers = []string{}
	}

	if simRespErr == nil {
		respBody, _ := util.GetResponseBody(actualResponse)
		respHeaders := util.GetResponseHeaders(actualResponse)

		actualResponseDetails := &models.ResponseDetails{
			Status:  actualResponse.StatusCode,
			Body:    string(respBody),
			Headers: respHeaders,
		}

		this.diffResponse(simResponse, actualResponseDetails, this.Arguments.Headers)
		this.Hoverfly.AddDiff(v2.SimpleRequestDefinitionView{
			Method: modifiedRequest.Method,
			Host:   modifiedRequest.URL.Host,
			Path:   modifiedRequest.URL.Path,
			Query:  modifiedRequest.URL.RawQuery,
		}, this.DiffReport)
	} else {
		log.WithFields(log.Fields{
			"mode":   Diff,
			"method": modifiedRequest.Method,
			"url":    modifiedRequest.URL,
		}).Info("There was no simulation matched for the request")
	}

	return actualResponse, nil
}

func (this *DiffMode) diffResponse(expected *models.ResponseDetails, actual *models.ResponseDetails, headersBlacklist []string) {
	if expected.Status != 0 && expected.Status != actual.Status {
		this.addEntry("status", expected.Status, actual.Status)
	}
	this.headerDiff(expected.Headers, actual.Headers, headersBlacklist)
	this.bodyDiff(expected, actual)
}

func (this *DiffMode) addEntry(parameterName string, expected interface{}, actual interface{}) {
	this.DiffReport.DiffEntries = append(this.DiffReport.DiffEntries,
		v2.DiffReportEntry{
			Field:    parameterName,
			Expected: nullOrValue(expected),
			Actual:   nullOrValue(actual),
		})
}

func nullOrValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprint(value)
}

func (this *DiffMode) headerDiff(expected map[string][]string, actual map[string][]string, headersBlacklist []string) bool {
	same := true
	for k := range expected {
		shouldContinue := false
		for _, header := range headersBlacklist {
			if k == header || header == "*" {
				shouldContinue = true
			}
		}
		if shouldContinue {
			continue
		}
		if _, ok := actual[k]; !ok {
			this.addEntry("header/"+k, expected[k], nil)
			same = false
		} else if !reflect.DeepEqual(expected[k], actual[k]) {
			this.addEntry("header/"+k, expected[k], actual[k])
			same = false
		}

	}
	return same
}

func (this *DiffMode) bodyDiff(expected *models.ResponseDetails, actual *models.ResponseDetails) bool {
	var expectedJson, actualJson interface{}

	err := unmarshalResponseToInterface(expected, &expectedJson)
	if err != nil {
		return this.doDeepEqual(expected.Body, actual.Body)
	}

	err = unmarshalResponseToInterface(actual, &actualJson)
	if err != nil {
		return this.doDeepEqual(expected.Body, actual.Body)
	}

	return this.JsonDiff("body", expectedJson.(map[string]interface{}), actualJson.(map[string]interface{}))
}

func (this *DiffMode) doDeepEqual(expected string, actual string) bool {
	if !reflect.DeepEqual(expected, actual) {
		this.addEntry("body", expected, actual)
		return false
	}
	return true
}

func unmarshalResponseToInterface(response *models.ResponseDetails, output interface{}) error {

	body := []byte(response.Body)

	encodings := response.Headers["Content-Encoding"]
	decompressedBody, err := decompress(body, encodings)
	if err != nil {
		return fmt.Errorf("It wasn't possible to decompress the response body: %s ", err)
	} else {
		body = decompressedBody
	}

	for i, ch := range body {
		switch {
		case ch == '\r':
			body[i] = ' '
		case ch == '\n':
			body[i] = ' '
		case ch == '\t':
			body[i] = ' '
		case ch == '\\':
			body[i] = ' '
		case ch == '\'':
			body[i] = '"'
		}
	}
	err = json.Unmarshal(body, &output)
	return err
}

func decompress(body []byte, encodings []string) ([]byte, error) {
	var err error
	var reader io.ReadCloser
	if len(encodings) > 0 {
		for index := range encodings {
			switch encodings[index] {
			case "gzip":
				reader, err = gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(body)))
				if err != nil {
					return body, err
				}
				body, err = ioutil.ReadAll(reader)
				if err != nil {
					return body, err
				}

			case "br":
				reader, err = brotli.NewReader(ioutil.NopCloser(bytes.NewBuffer(body)), &brotli.ReaderConfig{})
				if err != nil {
					return body, err
				}
				body, err = ioutil.ReadAll(reader)
				if err != nil {
					return body, err
				}

			case "deflate":
				reader = flate.NewReader(ioutil.NopCloser(bytes.NewBuffer(body)))
				body, err = ioutil.ReadAll(reader)
				if err != nil {
					return body, err
				}
			}
			reader.Close()
		}
	}
	return body, err
}

func (this *DiffMode) JsonDiff(prefix string, expected map[string]interface{}, actual map[string]interface{}) bool {
	same := true
	for k := range expected {
		param := prefix + "/" + k
		if _, ok := actual[k]; !ok {
			this.addEntry(param, expected[k], nil)
			same = false
		} else if reflect.TypeOf(expected[k]) != reflect.TypeOf(actual[k]) {
			this.addEntry(param, expected[k], actual[k])
			same = false
		} else {
			switch expected[k].(type) {
			default:
				if expected[k] != actual[k] {
					this.addEntry(param, expected[k], actual[k])
					same = false
				}
			case map[string]interface{}:
				if !this.JsonDiff(param, expected[k].(map[string]interface{}), actual[k].(map[string]interface{})) {
					same = false
				}
			case []interface{}:
				if !reflect.DeepEqual(expected[k], actual[k]) {
					this.addEntry(param, expected[k], actual[k])
					same = false
				}
			}
		}
	}

	return same
}
