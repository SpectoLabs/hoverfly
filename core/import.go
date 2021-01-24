package hoverfly

import (
	"encoding/json"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/delay"
	"github.com/SpectoLabs/hoverfly/core/state"
	"github.com/SpectoLabs/hoverfly/core/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	log "github.com/sirupsen/logrus"
)

// Import is a function that based on input decides whether it is a local resource or whether
// it should fetch it from remote server. It then imports given payload into the database
// or returns an error
func (hf *Hoverfly) Import(uri string) error {

	// assuming file URI is URL:
	if util.IsURL(uri) {
		log.WithFields(log.Fields{
			"isURL":      util.IsURL(uri),
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

	var simulation v2.SimulationViewV5

	body, err := ioutil.ReadAll(pairsFile)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	err = json.Unmarshal(body, &simulation)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	return hf.PutSimulation(simulation).GetError()
}

// ImportFromURL - takes one string value and tries connect to a remote server, then parse response body into
// recordedRequests structure (which is default format in which Hoverfly exports captured requests) and
// imports those requests into the database
func (hf *Hoverfly) ImportFromURL(url string) error {
	resp, err := http.DefaultClient.Get(url)

	if err != nil {
		return fmt.Errorf("Failed to fetch given URL, error %s", err.Error())
	}

	var simulation v2.SimulationViewV5

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	err = json.Unmarshal(body, &simulation)
	if err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	return hf.PutSimulation(simulation).GetError()
}

// importRequestResponsePairViews - a function to save given pairs into the database.
func (hf *Hoverfly) importRequestResponsePairViews(pairViews []v2.RequestMatcherResponsePairViewV5) v2.SimulationImportResult {
	importResult := v2.SimulationImportResult{}
	initialStates := map[string]string{}
	if len(pairViews) > 0 {
		success := 0
		failed := 0
		for i, pairView := range pairViews {

			pair := models.NewRequestMatcherResponsePairFromView(&pairView)

			if pairView.Response.LogNormalDelay != nil {
				d := *pairView.Response.LogNormalDelay
				if err := delay.ValidateLogNormalDelayOptions(d.Min, d.Max, d.Mean, d.Median); err != nil {
					failed++
					importResult.SetError(err)
					break
				}
			}

			var isPairAdded bool
			if hf.Cfg.NoImportCheck {
				hf.Simulation.AddPairWithoutCheck(pair)
				isPairAdded = true
			} else {
				isPairAdded = hf.Simulation.AddPair(pair)
			}

			if isPairAdded {
				for k, v := range pair.RequestMatcher.RequiresState {
					initialStates[k] = v
				}
				success++
			} else {
				importResult.AddPairIgnoredWarning(i)
			}

			if pairView.RequestMatcher.DeprecatedQuery != nil && len(pairView.RequestMatcher.DeprecatedQuery) != 0 {
				importResult.AddDeprecatedQueryWarning(i)
			}

			if len(pairView.Response.Headers["Content-Length"]) > 0 && len(pairView.Response.Headers["Transfer-Encoding"]) > 0 {
				importResult.AddContentLengthAndTransferEncodingWarning(i)
			}

			if len(pairView.Response.Headers["Content-Length"]) > 0 {
				contentLength, err := strconv.Atoi(pairView.Response.Headers["Content-Length"][0])
				if err == nil && contentLength != len(pair.Response.Body) {
					importResult.AddContentLengthMismatchWarning(i)
				}
			}

			continue
		}

		if hf.state == nil {
			hf.state = state.NewState()
		}
		hf.state.InitializeSequences(initialStates)

		log.WithFields(log.Fields{
			"total":      len(pairViews),
			"successful": success,
			"failed":     failed,
		}).Info("payloads imported")
		return importResult
	}

	return importResult
}
