package api_test

import (
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly/mode", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("Should get the mode", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"simulate","arguments":{"matchingStrategy":"strongest"}}`)))
		})
	})

	Context("PUT", func() {

		It("Should should get and set capture mode", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"capture"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture","arguments":{}}`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture","arguments":{}}`)))
		})

		It("Should should get and set modify mode", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"modify"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"modify","arguments":{}}`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"modify","arguments":{}}`)))
		})

		It("Should should get and set simulate mode", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"simulate"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"simulate","arguments":{"matchingStrategy":"strongest"}}`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"simulate","arguments":{"matchingStrategy":"strongest"}}`)))
		})

		It("Should should get and set synthesize mode", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"synthesize"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"synthesize","arguments":{}}`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"synthesize","arguments":{}}`)))
		})

		It("Should should get and set spy mode", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"spy"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"spy","arguments":{"matchingStrategy":"strongest"}}`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"spy","arguments":{"matchingStrategy":"strongest"}}`)))
		})

		It("Should error when header arguments use an asterisk and a header", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.BodyJSON(map[string]interface{}{
				"mode": "mode",
				"arguments": map[string][]string{
					"headersWhitelist": {"*", "Content-Type"},
				},
			})
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(400))
			errorJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			Expect(string(errorJson)).To(Equal(`{"error":"Not a valid mode"}`))

		})

		It("Should error when setting an invalid matching strategy", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`
			{
				"mode" : "simulate",
				"arguments" : {
					"matchingStrategy" : "INVALID"
				}
			}
			`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(400))
			errorJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			Expect(string(errorJson)).To(Equal(`{"error":"Only matching strategy of 'first' or 'strongest' is permitted"}`))

		})
	})
})
