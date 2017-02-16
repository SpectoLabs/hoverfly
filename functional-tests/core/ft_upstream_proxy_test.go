package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run two Hoverfly", func() {

	var (
		hoverflyPassThrough, hoverflyUpstream *functional_tests.Hoverfly
		hoverflyUpstreamURL                   string
	)

	BeforeEach(func() {
		hoverflyPassThrough = functional_tests.NewHoverfly()
		hoverflyUpstream = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverflyPassThrough.Stop()
		hoverflyUpstream.Stop()
	})

	Context("and configure the upstream proxy", func() {

		BeforeEach(func() {
			hoverflyUpstream.Start()
			hoverflyUpstream.SetMode("capture")
			hoverflyUpstreamURL = "localhost:" + hoverflyUpstream.GetProxyPort()

			hoverflyPassThrough.Start("-upstream-proxy=" + hoverflyUpstreamURL)
			hoverflyPassThrough.SetMode("capture")
		})

		It("Should capture the request and response", func() {

			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			resp := hoverflyPassThrough.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			expectedDestination := strings.Replace(fakeServer.URL, "http://", "", 1)

			recordsJson, err := ioutil.ReadAll(hoverflyPassThrough.GetSimulation())
			Expect(err).To(BeNil())

			Expect(recordsJson).To(ContainSubstring(expectedDestination))

			recordsJson, err = ioutil.ReadAll(hoverflyUpstream.GetSimulation())
			Expect(err).To(BeNil())

			Expect(recordsJson).To(ContainSubstring(expectedDestination))
		})
	})
})
