package hoverfly_test

import (
	"bytes"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("in simulate mode", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("should not template response if templating is disabled explicitely", func() {
			hoverfly.ImportSimulation(testdata.TemplatingDisabled)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("{{ Request.QueryParam.singular }}"))
		})

		It("should not template response if templating is not explcitely enabled or disabled", func() {
			hoverfly.ImportSimulation(testdata.TemplatingDisabledByDefault)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("{{ Request.QueryParam.one }}"))
		})

		It("should template response if templating is enabled and cache template not response", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabled)

			hoverfly.WriteLogsIfError()

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("foo"))

			resp = hoverfly.Proxy(sling.New().Get("http://test-server.com?one=bar"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("bar"))
		})

		It("should be able to use state in templating", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabledWithStateInBody)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/one"))
			Expect(resp.StatusCode).To(Equal(200))

			resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/two"))
			Expect(resp.StatusCode).To(Equal(200))
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("state for eggs"))
		})

		It("should not crash when templating a response if templating variable does not exist", func() {
			hoverfly.ImportSimulation(testdata.TemplatingEnabled)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?wrong=foo"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal(""))
		})
	})

	Context("in simulate mode, template helpers", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("randomString", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomString"))
			Expect(resp.StatusCode).To(Equal(200))

			_, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
		})

		It("randomStringLength10", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomStringLength10"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(len(string(body))).To(Equal(10))
		})

		It("randomBoolean", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomBoolean"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			_, err = strconv.ParseBool(string(body))
			Expect(err).To(BeNil())
		})

		It("randomInteger", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomInteger"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			parsedInt, err := strconv.ParseInt(string(body), 10, 0)
			Expect(err).To(BeNil())

			Expect(parsedInt > 0).To(BeTrue())
		})

		It("randomIntegerRange", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomIntegerRange1-10"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			parsedInt, err := strconv.ParseInt(string(body), 10, 0)
			Expect(err).To(BeNil())
			Expect(parsedInt >= 1 && parsedInt <= 10).To(BeTrue())
		})

		It("randomFloat", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomFloat"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			parsedFloat, err := strconv.ParseFloat(string(body), 1)
			Expect(err).To(BeNil())

			Expect(parsedFloat > 0).To(BeTrue())
		})

		It("randomFloatRange", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomFloatRange1-10"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			parsedFloat, err := strconv.ParseFloat(string(body), 1)
			Expect(err).To(BeNil())

			Expect(parsedFloat >= 1 && parsedFloat <= 10).To(BeTrue())
		})

		It("randomEmail", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomEmail"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(IsEmail(string(body))).To(BeTrue())
		})

		It("randomIPv4", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomIPv4"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(IsIPv4(string(body))).To(BeTrue())
		})

		It("randomIPv6", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomIPv6"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(IsIPv6(string(body))).To(BeTrue())
		})

		It("randomUuid", func() {
			hoverfly.ImportSimulation(testdata.TemplatingHelpers)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/randomuuid"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			returnedUuid := uuid.Parse(string(body))

			Expect(returnedUuid).To(Not(BeNil()))
		})
	})

	Context("in simulate mode, template request", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("simulate")
		})

		It("Request", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/Request"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// TODO: Handle this?
			Expect(string(body)).To(ContainSubstring("{map[] [Request] http %!s(func(string, string, *raymond.Options)"))
		})

		It("Request.Body jsonpath", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)
			resp := hoverfly.Proxy(sling.New().Post("http://test-server.com/Request.Body_jsonpath").BodyJSON(map[string]string{
				"test": "value",
			}))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("value"))
		})

		It("Request.Body jsonpath incorect", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)
			resp := hoverfly.Proxy(sling.New().Post("http://test-server.com/Request.Body_jsonpath").BodyJSON(map[string]string{
				"nottest": "value",
			}))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal(""))
		})

		It("Request.Body xpath", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)
			resp := hoverfly.Proxy(sling.New().Post("http://test-server.com/Request.Body_xpath").Body(bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8"?><root><text>value</text></root>`))))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("value"))
		})

		It("Request.Body xpath incorect", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)
			resp := hoverfly.Proxy(sling.New().Post("http://test-server.com/Request.Body_xpath").Body(bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8"?><root><nottext>test</text></root>`))))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal(""))
		})

		It("Request.Method", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/Request.Method"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("GET"))
		})

		It("Request.Scheme", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/Request.Scheme"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("http"))
		})

		It("Request.Path", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/one/two/three/Request.Path"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// TODO: Handle this?
			Expect(string(body)).To(Equal("onetwothreeRequest.Path"))
		})

		It("Request.Path.[0]", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/one/two/three/Request.Path0"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// TODO: Handle this?
			Expect(string(body)).To(Equal("one"))
		})

		It("Request.QueryParam", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/Request.QueryParam?query=param"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// TODO: Handle this?
			Expect(string(body)).To(Equal("map[query:[param]]"))
		})

		It("Request.QueryParam.query", func() {
			hoverfly.ImportSimulation(testdata.TemplatingRequest)

			resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/Request.QueryParam.query?query=param"))
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(string(body)).To(Equal("param"))
		})

	})
})

// Credit: https://github.com/asaskevich/govalidator
func IsIPv4(input string) bool {
	ip := net.ParseIP(input)
	return ip != nil && strings.Contains(input, ".")
}

// Credit: https://github.com/asaskevich/govalidator
func IsIPv6(input string) bool {
	ip := net.ParseIP(input)
	return ip != nil && strings.Contains(input, ":")
}

// Credit: https://github.com/badoux/checkmai
func IsEmail(email string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return false
	}
	return true
}
