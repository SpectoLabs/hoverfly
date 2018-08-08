package api_test

import (
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("When I run Hoverfly with auth", func() {

	var (
		hoverfly *functional_tests.Hoverfly

		username = "ft_user"
		password = "ft_password"
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	Context("Using a password  provided via command line", func() {

		BeforeEach(func() {
			hoverfly.Start("-auth", "-username", username, "-password", password)
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("Using /api/token-auth", func() {

			It("should return a 200 with correct username and password", func() {
				request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/token-auth").BodyJSON(backends.User{
					Username: username,
					Password: password,
				})

				response := functional_tests.DoRequest(request)
				Expect(response.StatusCode).To(Equal(200))
			})

			It("should return a 401 with incorrect password", func() {
				request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/token-auth").BodyJSON(backends.User{
					Username: username,
					Password: "wfewrrw",
				})

				response := functional_tests.DoRequest(request)
				Expect(response.StatusCode).To(Equal(401))
			})
		})
	})

	Context("Using a password hash provided via command line", func() {

		BeforeEach(func() {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
			hoverfly.Start("-auth", "-username", username, "-password-hash", string(hashedPassword))
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("Using /api/token-auth", func() {

			It("should return a 200 with correct username and password", func() {
				request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/token-auth").BodyJSON(backends.User{
					Username: username,
					Password: password,
				})

				response := functional_tests.DoRequest(request)
				Expect(response.StatusCode).To(Equal(200))
			})

			It("should return a 401 with incorrect password", func() {
				request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/token-auth").BodyJSON(backends.User{
					Username: username,
					Password: "wfewrrw",
				})

				response := functional_tests.DoRequest(request)
				Expect(response.StatusCode).To(Equal(401))
			})
		})
	})

	Context("Using a password  provided via command line", func() {

		BeforeEach(func() {
			hoverfly.Start("-auth", "-username", username, "-password", password)
		})

		AfterEach(func() {
			hoverfly.Stop()
		})
		It("should return a 429 after three incorrect attempts", func() {
			request := sling.New().Post("http://localhost:" + hoverfly.GetAdminPort() + "/api/token-auth").BodyJSON(backends.User{
				Username: username,
				Password: "wfewrrw",
			})

			response := functional_tests.DoRequest(request)
			Expect(response.StatusCode).To(Equal(401))

			response = functional_tests.DoRequest(request)
			Expect(response.StatusCode).To(Equal(401))

			response = functional_tests.DoRequest(request)
			Expect(response.StatusCode).To(Equal(401))

			response = functional_tests.DoRequest(request)
			Expect(response.StatusCode).To(Equal(429))
		})

	})
})
