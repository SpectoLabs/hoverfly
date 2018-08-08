package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly/diff", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("capture")
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("Should get empty diff", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/diff")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			diffJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(diffJson).To(Equal([]byte(`{"diff":null}`)))
		})

		It("Diffs should remain empty if response is unchanged", func() {
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
			}))

			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			hoverfly.SetMode("diff")

			resp = hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/diff")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			diffJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(diffJson).To(Equal([]byte(`{"diff":null}`)))
		})

		It("Should be diffs if response is changed", func() {
			contentType := "text/plain"
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", contentType)
			}))

			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?test=one"))
			Expect(resp.StatusCode).To(Equal(200))
			contentType = "application/json"
			hoverfly.SetMode("diff")

			resp = hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?test=one"))
			Expect(resp.StatusCode).To(Equal(200))

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/diff")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			var diffs v2.DiffView
			functional_tests.UnmarshalFromResponse(res, &diffs)

			Expect(diffs.Diff).To(HaveLen(1))
			Expect(diffs.Diff[0].Request).To(Equal(v2.SimpleRequestDefinitionView{
				Method: "GET",
				Host:   strings.Replace(fakeServer.URL, "http://", "", 1),
				Path:   "/",
				Query:  "test=one",
			}))
			Expect(diffs.Diff[0].DiffReport).To(HaveLen(1))
			Expect(diffs.Diff[0].DiffReport[0].DiffEntries).To(HaveLen(1))
			Expect(diffs.Diff[0].DiffReport[0].DiffEntries[0]).To(Equal(v2.DiffReportEntry{
				Field:    "header/Content-Type",
				Expected: "[text/plain]",
				Actual:   "[application/json]",
			}))
		})
	})

	Context("DELETE", func() {

		It("Should delete diffs", func() {
			contentType := "text/plain"
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", contentType)
			}))

			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?test=one"))
			Expect(resp.StatusCode).To(Equal(200))
			contentType = "application/json"
			hoverfly.SetMode("diff")

			resp = hoverfly.Proxy(sling.New().Get(fakeServer.URL + "?test=one"))
			Expect(resp.StatusCode).To(Equal(200))

			req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/diff")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/diff")
			res = functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			var diffs v2.DiffView
			functional_tests.UnmarshalFromResponse(res, &diffs)

			Expect(diffs.Diff).To(BeNil())
		})
	})
})
