package hoverfly_test

import (
	"io/ioutil"
	"net/http"

	"encoding/base64"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly

		username = "ft_user"
		password = "ft_password"
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	Context("with auth turned on", func() {

		BeforeEach(func() {
			hoverfly.Start("-auth", "-username", username, "-password", password)
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("Using HTTP Proxy authentication config values", func() {

			It("should return a 407 when trying to proxy without auth credentials", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when trying to proxy with incorrect auth credentials", func() {
				resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), "incorrect", "incorrect")
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 502 (no match in simulate mode) when trying to proxy with auth credentials", func() {
				resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), username, password)
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})

		Context("Using the `Proxy-Authorization` header over HTTP", func() {

			It("should not return a 407 (no match in simulate mode) when using Basic", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Basic with an incorrect base64 encoded credentials", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should not return a 407 (no match in simulate mode) when using Bearer", func() {
				token := hoverfly.GetAPIToken(username, password)

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Bearer "+token))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Bearer with an incorrect token", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Bearer ewGvdww.wRgFhE34.token"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `Proxy-Authorization` header over HTTPS", func() {

			It("should not return a 407 (no match in simulate mode) when using Basic", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Basic with an incorrect base64 encoded credentials", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should not return a 407 (no match in simulate mode) when using Bearer", func() {
				token := hoverfly.GetAPIToken(username, password)

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Bearer "+token))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Bearer with an incorrect token", func() {
				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Bearer ewGvdww.wRgFhE34.token"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `X-HOVERFLY-AUTHORIZATION` header`", func() {
			It("should return a 407 as `X-HOVERFLY-AUTHORIZATION` header is disabled", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})
	})

	Context("with auth turned on and using a boltdb", func() {

		BeforeEach(func() {
			hoverfly.Start("-db", "boltdb", "-auth", "-username", username, "-password", password)
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("Using HTTP Proxy authentication config values", func() {

			It("should return a 407 when trying to proxy without auth credentials", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when trying to proxy with incorrect auth credentials", func() {
				resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), "incorrect", "incorrect")
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 502 (no match in simulate mode) when trying to proxy with auth credentials", func() {
				resp := hoverfly.ProxyWithAuth(sling.New().Get("http://hoverfly.io"), username, password)
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})

		Context("Using the `Proxy-Authorization header over HTTP`", func() {

			It("should not return a 407 (no match in simulate mode) when using Basic", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Basic with an incorrect base64 encoded credentials", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should not return a 407 (no match in simulate mode) when using Bearer", func() {
				token := hoverfly.GetAPIToken(username, password)

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Bearer "+token))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Bearer with an incorrect token", func() {
				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("Proxy-Authorization", "Bearer ewGvdww.wRgFhE34.token"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `Proxy-Authorization` header over HTTPS", func() {

			It("should not return a 407 (no match in simulate mode) when using Basic", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Basic with an incorrect base64 encoded credentials", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should not return a 407 (no match in simulate mode) when using Bearer", func() {
				token := hoverfly.GetAPIToken(username, password)

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Bearer "+token))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Bearer with an incorrect token", func() {
				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Bearer ewGvdww.wRgFhE34.token"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `X-HOVERFLY-AUTHORIZATION` header`", func() {
			It("should return a 407 as `X-HOVERFLY-AUTHORIZATION` header is disabled", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("http://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})
	})

	Context("with auth turned on using `header-auth`", func() {

		BeforeEach(func() {
			hoverfly.Start("-auth", "-username", username, "-password", password, "-proxy-auth", "header-auth")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("Using the `X-HOVERFLY-AUTHORIZATION` header over HTTP", func() {

			It("should return a 407 as HTTP is disabled", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `X-HOVERFLY-AUTHORIZATION` header over HTTPS", func() {

			It("should not return a 407 (no match in simulate mode) when using Basic", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Basic "+base64Encoded))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Basic with an incorrect base64 encoded credentials", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":incorect"))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})

			It("should not return a 407 (no match in simulate mode) when using Bearer", func() {
				token := hoverfly.GetAPIToken(username, password)

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Bearer "+token))
				Expect(resp.StatusCode).ToNot(Equal(http.StatusProxyAuthRequired))
			})

			It("should return a 407 when using Bearer with an incorrect token", func() {
				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("X-HOVERFLY-AUTHORIZATION", "Bearer ewGvdww.wRgFhE34.token"))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))
			})
		})

		Context("Using the `Proxy-Authorization` header`", func() {
			It("should return a 407 as `Proxy-Authorization` header is disabled", func() {
				base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

				resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
				Expect(resp.StatusCode).To(Equal(http.StatusProxyAuthRequired))

				responseBody, _ := ioutil.ReadAll(resp.Body)
				Expect(string(responseBody)).To(Equal("407 `Proxy-Authorization` header is disabled, use `X-HOVERFLY-AUTHORIZATION` instead"))
			})
		})
	})

	Context("with auth turned on and using https", func() {

		BeforeEach(func() {
			hoverfly.Start("-auth", "-username", username, "-password", password)
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		// This test fails, if you try to manually execute it using the command below
		// you will get an error, making me believe this isn't supported

		// curl -proxy benji:password@localhost:8500 https://hoverfly.io

		// It("should return a 502 (no match in simulate mode) when trying to proxy with auth credentials", func() {
		// 	resp := hoverfly.ProxyWithAuth(sling.New().Get("https://hoverfly.io"), username, password)
		// 	Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
		// })

		It("should return a 502 (no match in simulate mode) when using Basic Proxy-Authorization", func() {
			base64Encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

			resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Basic "+base64Encoded))
			Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
		})

		It("should return a 502 (no match in simulate mode) when using Bearer Proxy-Authorization", func() {
			token := hoverfly.GetAPIToken(username, password)

			resp := hoverfly.Proxy(sling.New().Get("https://hoverfly.io").Add("Proxy-Authorization", "Bearer "+token))
			Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
		})
	})
})
