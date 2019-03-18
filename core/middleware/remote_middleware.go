package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func (this Middleware) executeMiddlewareRemotely(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if this.Remote == "" {
		return pair, &MiddlewareError{
			Message: "Remote middleware not set",
		}
	}

	req, err := http.NewRequest("POST", this.Remote, bytes.NewBuffer(pairViewBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when building request to remote middleware")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Error when building request to remote middleware: ",
			Url:           this.Remote,
			Stdin:         string(pairViewBytes),
		}
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when communicating with remote middleware")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Error when communicating with remote middleware:",
			Url:           this.Remote,
			Stdin:         string(pairViewBytes),
		}
	}

	if resp.StatusCode != 200 {
		log.Error("Remote middleware did not process payload")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       fmt.Sprintf("Error when communicating with remote middleware: received %d", resp.StatusCode),
			Url:           this.Remote,
			Stdin:         string(pairViewBytes),
		}
	}

	returnedPairViewBytes, err := ioutil.ReadAll(resp.Body)
	if len(returnedPairViewBytes) == 0 {
		returnedPairViewBytes = []byte(" ")
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when process response from remote middleware")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Error when reading response body from remote middleware",
			Url:           this.Remote,
			Stdin:         string(pairViewBytes),
			Stdout:        string(returnedPairViewBytes),
		}
	}

	var newPairView RequestResponsePairView

	err = json.Unmarshal(returnedPairViewBytes, &newPairView)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when trying to serialize response from remote middleware")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Error when trying to serialize response from remote middleware",
			Url:           this.Remote,
			Stdin:         string(pairViewBytes),
			Stdout:        string(returnedPairViewBytes),
		}
	}
	return models.NewRequestResponsePairFromRequestResponsePairView(newPairView), nil
}
