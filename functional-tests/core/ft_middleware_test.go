package hoverfly_test

import (
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/SpectoLabs/hoverfly/core/models"
	"encoding/json"
	"bytes"
)

func checkHeadersHttpMiddleware(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newPayloadView models.PayloadView

	json.Unmarshal(body, &newPayloadView)

	if _, present := newPayloadView.Request.Headers["New-Header"]; present {
		newPayloadView.Response.Body = "New-Header present"
	} else {
		newPayloadView.Response.Body = "New-Header not present"
	}

	bts, _ := json.Marshal(newPayloadView)
	w.Write(bts)
}

var server *httptest.Server

var _ = Describe("Running Hoverfly with middleware", func() {


	Context("in simulate mode", func() {

		BeforeEach(func() {
			muxRouter := mux.NewRouter()
			muxRouter.HandleFunc("/process", checkHeadersHttpMiddleware).Methods("POST")
			server = httptest.NewServer(muxRouter)

			hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, server.URL + "/process")

			jsonPayload := bytes.NewBufferString(`{"data":[{"request": {"path": "/path1", "method": "GET", "destination": "destination1", "scheme": "http", "query": "", "body": "", "headers": {}}, "response": {"status": 200, "encodedBody": false, "body": "body1", "headers": {"Header": ["value1"]}}}]}`)
			ImportHoverflyRecords(jsonPayload)

			SetHoverflyMode("simulate")
		})

		AfterEach(func() {
			server.Close()
			stopHoverfly()
		})

		It("the middleware should recieve the request made instead of the request stored in the cache", func() {
			slingRequest := sling.New().Get("http://destination1/path1").Add("New-Header", "true")

			resp := DoRequestThroughProxy(slingRequest)

			body, _ := ioutil.ReadAll(resp.Body)

			Expect(string(body)).To(Equal("New-Header present"))

		})
	})
})
