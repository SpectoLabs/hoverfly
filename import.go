package hoverfly

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"net/http"
	"encoding/base64"
)

// Import is a function that based on input decides whether it is a local resource or whether
// it should fetch it from remote server. It then imports given payload into the database
// or returns an error
func (d *Hoverfly) Import(uri string) error {

	// assuming file URI is URL:
	if isURL(uri) {
		log.WithFields(log.Fields{
			"isURL":      isURL(uri),
			"importFrom": uri,
		}).Info("URL")
		return d.ImportFromURL(uri)
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
		return d.ImportFromDisk(uri)
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
func (d *Hoverfly) ImportFromDisk(path string) error {
	payloadsFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Got error while opening payloads file, error %s", err.Error())
	}

	var requests recordedRequests

	jsonParser := json.NewDecoder(payloadsFile)
	if err = jsonParser.Decode(&requests); err != nil {
		return fmt.Errorf("Got error while parsing payloads file, error %s", err.Error())
	}

	return d.ImportPayloads(requests.Data)
}

// ImportFromURL - takes one string value and tries connect to a remote server, then parse response body into
// recordedRequests structure (which is default format in which Hoverfly exports captured requests) and
// imports those requests into the database
func (d *Hoverfly) ImportFromURL(url string) error {

	resp, err := d.HTTP.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to fetch given URL, error %s", err.Error())
	}

	var requests recordedRequests

	jsonParser := json.NewDecoder(resp.Body)
	if err = jsonParser.Decode(&requests); err != nil {
		return fmt.Errorf("Got error while parsing payloads, error %s", err.Error())
	}

	return d.ImportPayloads(requests.Data)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

// ImportPayloads - a function to save given payloads into the database.
func (d *Hoverfly) ImportPayloads(payloads []PayloadView) error {
	if len(payloads) > 0 {
		success := 0
		failed := 0
		for _, payloadView := range payloads {

			// Decode base64 if body is encoded

			if payloadView.Response.EncodedBody {
				decodedBody, err := base64.StdEncoding.DecodeString(payloadView.Response.Body)
				if err != nil {
					log.Fatal("error:", err)
				}

				payloadView.Response.Body = string(decodedBody)
			}

			// Convert PayloadView back to Payload for internal storage
			pl := payloadView.ConvertToPayload()

			if len(pl.Request.Headers) == 0 {
				pl.Request.Headers = make(map[string][]string)
			}

			if _, present := pl.Request.Headers["Content-Type"]; !present {
				// sniffing content types
				if isJSON(pl.Request.Body) {
					pl.Request.Headers["Content-Type"] = []string{"application/json"}
				} else {
					ct := http.DetectContentType([]byte(pl.Request.Body))
					pl.Request.Headers["Content-Type"] = []string{ct}
				}
			}

			// recalculating request hash and storing it in database
			r := RequestContainer{Details: pl.Request, Minifier: d.MIN}
			key := r.Hash()

			// regenerating key
			pl.ID = key

			bts, err := pl.Encode()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error("Failed to encode payload")
				failed++
			} else {
				// hook
				var en Entry
				en.ActionType = ActionTypeRequestCaptured
				en.Message = "imported"
				en.Time = time.Now()
				en.Data = bts

				if err := d.Hooks.Fire(ActionTypeRequestCaptured, &en); err != nil {
					log.WithFields(log.Fields{
						"error":      err.Error(),
						"message":    en.Message,
						"actionType": ActionTypeRequestCaptured,
					}).Error("failed to fire hook")
				}

				d.RequestCache.Set([]byte(key), bts)
				if err == nil {
					success++
				} else {
					failed++
				}
			}
		}
		log.WithFields(log.Fields{
			"total":      len(payloads),
			"successful": success,
			"failed":     failed,
		}).Info("payloads imported")
		return nil
	}
	return fmt.Errorf("Bad request. Nothing to import!")
}
