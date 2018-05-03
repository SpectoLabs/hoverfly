package hoverfly_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func checkHeadersHttpMiddleware(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newPairView v2.RequestResponsePairViewV1

	json.Unmarshal(body, &newPairView)

	if _, present := newPairView.Request.Headers["New-Header"]; present {
		newPairView.Response.Body = "New-Header present"
	} else {
		newPairView.Response.Body = "New-Header not present"
	}

	pairViewBytes, _ := json.Marshal(newPairView)
	w.Write(pairViewBytes)
}

var server *httptest.Server

var _ = Describe("Running Hoverfly with middleware", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
		server.Close()
	})

	Context("in simulate mode", func() {

		BeforeEach(func() {
			muxRouter := mux.NewRouter()
			muxRouter.HandleFunc("/process", checkHeadersHttpMiddleware).Methods("POST")
			server = httptest.NewServer(muxRouter)

			hoverfly.Start("-middleware", server.URL+"/process")
			hoverfly.ImportSimulation(testdata.JsonGetAndPost)
			hoverfly.SetMode("simulate")
		})

		It("the middleware should recieve the request made instead of the request stored in the cache", func() {
			slingRequest := sling.New().Get("http://destination1/path1").Add("New-Header", "here")

			resp := hoverfly.Proxy(slingRequest)

			body, _ := ioutil.ReadAll(resp.Body)

			Expect(string(body)).To(Equal("New-Header present"))

		})
	})
})
