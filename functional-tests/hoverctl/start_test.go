package hoverctl_suite

import (
	"io/ioutil"
	"strconv"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("hoverctl `start`", func() {

	AfterEach(func() {
		output := functional_tests.Run(hoverctlBinary, "targets")
		KillHoverflyTargets(output)
	})

	Context("without a target", func() {
		It("should start an instance of hoverfly with default configuration", func() {
			output := functional_tests.Run(hoverctlBinary, "start")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring("8888"))

			Expect(output).To(ContainSubstring("proxy-port"))
			Expect(output).To(ContainSubstring("8500"))
		})

		It("should create a local target", func() {
			functional_tests.Run(hoverctlBinary, "start")

			output := functional_tests.Run(hoverctlBinary, "targets")

			targets := functional_tests.TableToSliceMapStringString(output)
			Expect(targets).To(HaveKey("local"))
			Expect(targets["local"]).To(HaveKeyWithValue("TARGET NAME", "local"))
			Expect(targets["local"]).To(HaveKeyWithValue("ADMIN PORT", "8888"))
			Expect(targets["local"]).To(HaveKeyWithValue("PROXY PORT", "8500"))

		})
	})

	Context("with default target, but providing configuration flags", func() {
		It("should start an instance of hoverfly setting the admin port based on the flag", func() {
			randomAdminPort := strconv.Itoa(freeport.GetPort())

			output := functional_tests.Run(hoverctlBinary, "start", "--admin-port", randomAdminPort)

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring(randomAdminPort))

			Expect(output).To(ContainSubstring("proxy-port"))
			Expect(output).To(ContainSubstring("8500"))
		})

		It("should start an instance of hoverfly setting the proxy port based on the flag", func() {
			randomProxyPort := strconv.Itoa(freeport.GetPort())

			output := functional_tests.Run(hoverctlBinary, "start", "--proxy-port", randomProxyPort)

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring("8888"))

			Expect(output).To(ContainSubstring("proxy-port"))
			Expect(output).To(ContainSubstring(randomProxyPort))
		})

		It("should start an instance of hoverfly setting the admin and proxy port based on the flags", func() {
			randomAdminPort := strconv.Itoa(freeport.GetPort())
			randomProxyPort := strconv.Itoa(freeport.GetPort())

			output := functional_tests.Run(hoverctlBinary, "start", "--admin-port", randomAdminPort, "--proxy-port", randomProxyPort)

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring(randomAdminPort))

			Expect(output).To(ContainSubstring("proxy-port"))
			Expect(output).To(ContainSubstring(randomProxyPort))
		})

		It("should start an instance of hoverfly as a webserver", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "webserver")

			Expect(output).To(ContainSubstring("Hoverfly is now running as a webserver"))

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring("8888"))

			Expect(output).To(ContainSubstring("webserver-port"))
			Expect(output).To(ContainSubstring("8500"))

			request := sling.New().Get("http://localhost:8500")
			response := functional_tests.DoRequest(request)

			Expect(ioutil.ReadAll(response.Body)).ToNot(ContainSubstring("This is a proxy server"))
		})

		It("should start an instance of hoverfly with custom certificate and key", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--certificate", "testdata/cert.pem", "--key", "testdata/key.pem")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs", "--json")

			Expect(output).To(ContainSubstring("Default keys have been overwritten"))
		})

		It("should start an instance of  hoverfly with tls verification turned off", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--disable-tls")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs", "--json")

			Expect(output).To(ContainSubstring("TLS certificate verification has been disabled"))
		})

		It("should start an instance of hoverfly with persisted cache", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--cache", "testdata-gen/cache.db", "-v")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs", "--json")
			Expect(output).To(ContainSubstring("Using boltdb backend"))
		})

		It("should start and instance of hoverfly without a cache", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--disable-cache")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			response := functional_tests.DoRequest(sling.New().Get("http://localhost:8888/api/v2/cache"))
			Expect(response.StatusCode).To(Equal(500))

			Expect(ioutil.ReadAll(response.Body)).To(ContainSubstring(`{"error":"No cache set"}`))
		})

		It("starts hoverfly with upstream-proxy for environments with a proxy set up already", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--upstream-proxy", "hoverfly.io:8080", "-v")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			response := functional_tests.DoRequest(sling.New().Get("http://localhost:8888/api/v2/hoverfly"))
			Expect(response.StatusCode).To(Equal(200))

			Expect(ioutil.ReadAll(response.Body)).To(ContainSubstring(`"upstreamProxy":"http://hoverfly.io:8080"`))
		})

		It("should start an instance of hoverfly with HTTPS only", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--https-only")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs", "--json")
			Expect(output).To(ContainSubstring("Disabling HTTP"))
		})

		It("should start hoverfly with authentication", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--auth", "--username", functional_tests.HoverflyUsername, "--password", functional_tests.HoverflyPassword)

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			response := functional_tests.DoRequest(sling.New().Get("http://localhost:8888/api/v2/hoverfly"))
			Expect(response.StatusCode).To(Equal(401))

			response = functional_tests.DoRequest(sling.New().Post("http://localhost:8888/api/token-auth").BodyJSON(backends.User{
				Username: functional_tests.HoverflyUsername,
				Password: functional_tests.HoverflyPassword,
			}))

			Expect(response.StatusCode).To(Equal(200))
		})
	})

	Context("with a target", func() {
		It("should start an instance of Hoverfly with the target configuration", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-target",
				"--admin-port", "1234",
				"--proxy-port", "8765",
			)

			output := functional_tests.Run(hoverctlBinary, "start", "--target", "test-target")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			Expect(output).To(ContainSubstring("admin-port"))
			Expect(output).To(ContainSubstring("123"))

			Expect(output).To(ContainSubstring("proxy-port"))
			Expect(output).To(ContainSubstring("8765"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})

	Context("with a target with a remote url", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "remote", "--host", "hoverfly.io")
		})
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--target", "remote")

			Expect(output).To(ContainSubstring("Unable to start an instance of Hoverfly on a remote host (remote host: hoverfly.io)"))
			Expect(output).To(ContainSubstring("Run `hoverctl start --new-target <name>`"))
		})
	})

	Context("with a target with a running hoverfly", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "start", "--new-target", "running")
		})
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--target", "running")

			Expect(output).To(ContainSubstring("Target Hoverfly is already running"))
			Expect(output).To(ContainSubstring("Run `hoverctl stop -t running` to stop it"))
		})
	})

	Context("with --new-target flag", func() {

		It("should create a target using the flags to configure the target", func() {
			randomAdminPort := strconv.Itoa(freeport.GetPort())
			randomProxyPort := strconv.Itoa(freeport.GetPort())

			functional_tests.Run(hoverctlBinary, "start",
				"--new-target", "notdefault",
				"--admin-port", randomAdminPort,
				"--proxy-port", randomProxyPort,
			)

			output := functional_tests.Run(hoverctlBinary, "targets")

			targets := functional_tests.TableToSliceMapStringString(output)
			Expect(targets).To(HaveKey("notdefault"))
			Expect(targets["notdefault"]).To(HaveKeyWithValue("TARGET NAME", "notdefault"))
			Expect(targets["notdefault"]).To(HaveKeyWithValue("ADMIN PORT", randomAdminPort))
			Expect(targets["notdefault"]).To(HaveKeyWithValue("PROXY PORT", randomProxyPort))
		})

		It("should error when trying to create a target that already exists", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "exists")

			output := functional_tests.Run(hoverctlBinary, "start",
				"--new-target", "exists",
			)

			Expect(output).To(ContainSubstring("Target exists already exists"))
			Expect(output).To(ContainSubstring("Use a different target name or run `hoverctl targets update exists`"))
		})
	})

	Context("with a target that has already been started", func() {
		BeforeEach(func() {

		})

		It("should error", func() {
			functional_tests.Run(hoverctlBinary, "start", "--new-target", "started")
			output := functional_tests.Run(hoverctlBinary, "start", "-t", "started")

			Expect(output).To(ContainSubstring("Hoverfly is already running"))
		})
	})

	Context("with port conflicts", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "start", "--new-target", "admin-conflict", "--admin-port", "8500", "--proxy-port", "8888")
		})

		It("errors when the admin port is the same", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--proxy-port", "1234")

			Expect(output).To(ContainSubstring("Could not start Hoverfly"))
			Expect(output).To(ContainSubstring("Port 8888 was not free"))
		})

		It("errors when the proxy port is the same", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--admin-port", "1234")

			Expect(output).To(ContainSubstring("Could not start Hoverfly"))
			Expect(output).To(ContainSubstring("Port 8500 was not free"))
		})
	})

	Context("with pac-file", func() {
		BeforeEach(func() {
		})

		It("starts with pac file when defined", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--pac-file", "testdata/test.pac")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			response := functional_tests.DoRequest(sling.New().Get("http://localhost:8888/api/v2/hoverfly/pac"))
			Expect(response.StatusCode).To(Equal(200))
			responseBody, err := ioutil.ReadAll(response.Body)
			Expect(err).To(BeNil())
			Expect(string(responseBody)).To(ContainSubstring(`function FindProxyForURL(url, host) {`))
		})

		It("errors when pac file not found", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--pac-file", "unknown.pac")

			Expect(output).To(ContainSubstring("File not found: unknown.pac"))
		})
	})
})
