package hoverfly

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
)

// Import is a function that based on input decides whether it is a local resource or whether
// it should fetch it from remote server. It then imports given payload into the database
// or returns an error
func (hf *Hoverfly) Import(uri string) error {

	// assuming file URI is URL:
	if isURL(uri) {
		log.WithFields(log.Fields{
			"isURL":      isURL(uri),
			"importFrom": uri,
		}).Info("URL")
		return hf.ImportFromURL(uri)
	}
	// assuming file URI is disk location
	ext := path.Ext(uri)
	if ext != ".json" {
		return fmt.Errorf("Failed to import payloads, only JSON files are acceppted. Given file: %s", uri)
	}
	// checking whether it exists
	exists, err := exists(uri)
	if err != nil {
		return fmt.Errorf("Failed to import payloads from %s. Got error: %s", uri, err.Error())
	}
	if exists {
		// file is JSON and it exist
		return hf.ImportFromDisk(uri)
	}
	return fmt.Errorf("Failed to import payloads, given file '%s' does not exist", uri)
}

// URL is regexp to match http urls
const URL string = `^((ftp|https?):\/\/)(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`

var rxURL = regexp.MustCompile(URL)

func isURL(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return rxURL.MatchString(str)

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ImportFromDisk - takes one string value and tries to open a file, then parse it into recordedRequests structure
// (which is default format in which Hoverfly exports captured requests) and imports those requests into the database
func (hf *Hoverfly) ImportFromDisk(path string) error {
	pairsFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Got error while opening payloads file, error %s", err.Error())
	}

	var requests v1.RequestResponsePairPayload

	body, err := ioutil.ReadAll(pairsFile)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	err = json.Unmarshal(body, &requests)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	requestResponsePairViews := make([]interfaces.RequestResponsePair, len(requests.Data))
	for i, v := range requests.Data {
		requestResponsePairViews[i] = v
	}

	return hf.ImportRequestResponsePairViews(requestResponsePairViews)
}

// ImportFromURL - takes one string value and tries connect to a remote server, then parse response body into
// recordedRequests structure (which is default format in which Hoverfly exports captured requests) and
// imports those requests into the database
func (hf *Hoverfly) ImportFromURL(url string) error {
	resp, err := http.DefaultClient.Get(url)

	if err != nil {
		return fmt.Errorf("Failed to fetch given URL, error %s", err.Error())
	}

	var requests v1.RequestResponsePairPayload

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	err = json.Unmarshal(body, &requests)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	requestResponsePairViews := make([]interfaces.RequestResponsePair, len(requests.Data))
	for i, v := range requests.Data {
		requestResponsePairViews[i] = v
	}

	return hf.ImportRequestResponsePairViews(requestResponsePairViews)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

// ImportRequestResponsePairViews - a function to save given pairs into the database.
func (hf *Hoverfly) ImportRequestResponsePairViews(pairViews []interfaces.RequestResponsePair) error {
	if len(pairViews) > 0 {
		success := 0
		failed := 0
		for _, pairView := range pairViews {

			// Convert PayloadView back to Payload for internal storage
			pair := models.NewRequestResponsePairFromRequestResponsePairView(pairView)

			if pairView.GetRequest().GetRequestType() != nil && *pairView.GetRequest().GetRequestType() == *StringToPointer("template") {
				responseDetails := models.NewResponseDetailsFromResponse(pairView.GetResponse())

				requestTemplate := models.RequestTemplate{
					Path:        pairView.GetRequest().GetPath(),
					Method:      pairView.GetRequest().GetMethod(),
					Destination: pairView.GetRequest().GetDestination(),
					Scheme:      pairView.GetRequest().GetScheme(),
					Query:       pairView.GetRequest().GetQuery(),
					Body:        pairView.GetRequest().GetBody(),
					Headers:     pairView.GetRequest().GetHeaders(),
				}

				requestTemplateResponsePair := models.RequestTemplateResponsePair{
					RequestTemplate: requestTemplate,
					Response:        responseDetails,
				}

				hf.Simulation.Templates = append(hf.Simulation.Templates, requestTemplateResponsePair)
				success++
				continue
			}

			if len(pair.Request.Headers) == 0 {
				pair.Request.Headers = make(map[string][]string)
			}

			if _, present := pair.Request.Headers["Content-Type"]; !present {
				// sniffing content types
				if isJSON(pair.Request.Body) {
					pair.Request.Headers["Content-Type"] = []string{"application/json"}
				} else {
					ct := http.DetectContentType([]byte(pair.Request.Body))
					pair.Request.Headers["Content-Type"] = []string{ct}
				}
			}

			err := hf.RequestMatcher.SaveRequestResponsePair(&pair)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("Failed to save payload")
			}

			if err == nil {
				success++
			} else {
				failed++
			}
		}
		log.WithFields(log.Fields{
			"total":      len(pairViews),
			"successful": success,
			"failed":     failed,
		}).Info("payloads imported")
		return nil
	}
	return nil
}
