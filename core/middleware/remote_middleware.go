package middleware

import (
	"bytes"
	"encoding/json"

	"errors"
	"io/ioutil"
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func (this Middleware) executeMiddlewareRemotely(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if this.Remote == "" {
		return pair, fmt.Errorf("Error when communicating with remote middleware")
	}

	req, err := http.NewRequest("POST", this.Remote, bytes.NewBuffer(pairViewBytes))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when building request to remote middleware")
		return pair, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when communicating with remote middleware")
		return pair, err
	}

	if resp.StatusCode != 200 {
		log.Error("Remote middleware did not process payload")
		return pair, errors.New("Error when communicating with remote middleware")
	}

	returnedPairViewBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when process response from remote middleware")
		return pair, err
	}

	var newPairView RequestResponsePairView

	err = json.Unmarshal(returnedPairViewBytes, &newPairView)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when trying to serialize response from remote middleware")
		return pair, err
	}
	return models.NewRequestResponsePairFromRequestResponsePairView(newPairView), nil
}
