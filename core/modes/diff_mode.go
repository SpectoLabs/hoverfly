package modes

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"reflect"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/util"
	"fmt"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"io"
	"github.com/dsnet/compress/brotli"
	"compress/flate"
)

var DiffErrorMsg DiffErrorMessage

type HoverflyDiff interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, error)
}

type DiffMode struct {
	Hoverfly         HoverflyDiff
	MatchingStrategy string
	DiffErrorMessage DiffErrorMessage
}

func (this DiffMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Diff,
		Arguments: v2.ModeArgumentsView{
			MatchingStrategy: &this.MatchingStrategy,
		},
	}
}

func (this DiffMode) GetMessage() DiffErrorMessage {
	return DiffErrorMsg
}

func (this DiffMode) SetArguments(arguments ModeArguments) {
	if arguments.MatchingStrategy == nil {
		this.MatchingStrategy = "strongest"
	} else {
		this.MatchingStrategy = *arguments.MatchingStrategy
	}
}

//TODO: We should only need one of these two parameters
func (this DiffMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	DiffErrorMsg = DiffErrorMessage{}

	actualPair := models.RequestResponsePair{
		Request: details,
	}

	simResponse, simRespErr := this.Hoverfly.GetResponse(details)

	log.Info("Going to call real server")
	modifiedRequest, err := ReconstructRequestForPassThrough(actualPair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &actualPair, "There was an error when reconstructing the request.", Diff)
	}

	actualResponse, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &actualPair, "There was an error when forwarding the request to the intended destination", Diff)
	}

	if simRespErr == nil {
		respBody, _ := util.GetResponseBody(actualResponse)

		actualResponseDetails := &models.ResponseDetails{
			Status:  actualResponse.StatusCode,
			Body:    string(respBody),
			Headers: actualResponse.Header,
		}

		diffResponse(simResponse, actualResponseDetails)
	} else {
		log.WithFields(log.Fields{
			"mode":   Diff,
			"method": request.Method,
			"url":    request.URL,
		}).Info("There was no simulation matched for the request")
	}

	return actualResponse, nil
}

func diffResponse(expected *models.ResponseDetails, actual *models.ResponseDetails) {
	if expected.Status != 0 && expected.Status != actual.Status {
		DiffErrorMsg.write("status", expected.Status, actual.Status)
	}
	headerDiff(&DiffErrorMsg, expected.Headers, actual.Headers)
	bodyDiff(&DiffErrorMsg, expected, actual)
}

type DiffErrorMessage struct {
	DiffMessage bytes.Buffer
	counter     int
}

func (message *DiffErrorMessage) GetErrorMessage() string {
	return message.DiffMessage.String()
}

func (message *DiffErrorMessage) write(parameterName string, expected interface{}, actual interface{}) {
	message.DiffMessage.WriteString(
		fmt.Sprintf("(%d)The \"%s\" parameter is not same - the expected value was [%s], but the actual one [%s]\n",
			message.counter, parameterName, expected, actual))
	message.counter++
}

func headerDiff(message *DiffErrorMessage, expected map[string][]string, actual map[string][]string) bool {
	same := true
	for k := range expected {
		if _, ok := actual[k]; !ok {
			message.write("header/"+k, expected[k], "undefined")
			same = false
		} else
		if !reflect.DeepEqual(expected[k], actual[k]) {
			message.write("header/"+k, expected[k], actual[k])
			same = false
		}

	}
	return same
}

func bodyDiff(message *DiffErrorMessage, expected *models.ResponseDetails, actual *models.ResponseDetails) bool {
	var expectedJson, actualJson interface{}

	err := unmarshalResponseToInterface(expected, &expectedJson)
	if err != nil {
		return doDeepEqual(message, expected.Body, actual.Body)
	}

	err = unmarshalResponseToInterface(actual, &actualJson)
	if err != nil {
		return doDeepEqual(message, expected.Body, actual.Body)
	}

	return JsonDiff(message, "body", expectedJson.(map[string]interface{}), actualJson.(map[string]interface{}))
}

func doDeepEqual(message *DiffErrorMessage, expected string, actual string) bool {
	if !reflect.DeepEqual(expected, actual) {
		message.write("body", expected, actual)
		return false
	}
	return true
}

func unmarshalResponseToInterface(response *models.ResponseDetails, output interface{}) error {

	body := []byte(response.Body)

	encodings := response.Headers["Content-Encoding"]
	decompressedBody, err := decompress(body, encodings)
	if err != nil {
		fmt.Errorf("It wasn't possible to decompress the response body: %s ", err)
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

func JsonDiff(message *DiffErrorMessage, prefix string, expected map[string]interface{}, actual map[string]interface{}) bool {
	same := true
	for k := range expected {
		param := prefix + "/" + k
		if _, ok := actual[k]; !ok {
			message.write(param, expected[k], "undefined")
			same = false
		} else if reflect.TypeOf(expected[k]) != reflect.TypeOf(actual[k]) {
			message.write(param, expected[k], actual[k])
			same = false
		} else {
			switch expected[k].(type) {
			default:
				if expected[k] != actual[k] {
					message.write(param, expected[k], actual[k])
					same = false
				}
			case map[string]interface{}:
				if !JsonDiff(message, param, expected[k].(map[string]interface{}), actual[k].(map[string]interface{})) {
					same = false
				}
			case []interface{}:
				if !reflect.DeepEqual(expected[k], actual[k]) {
					message.write(param, expected[k], actual[k])
					same = false
				}
			}
		}
	}

	return same
}
