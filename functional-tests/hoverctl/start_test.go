package hoverctl_suite

import (
	"io/ioutil"
	"strconv"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("hoverctl `start`", func() {

	AfterEach(func() {
		output := functional_tests.Run(hoverctlBinary, "targets")
		functional_tests.KillHoverflyTargets(output)
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

		It("should create a default target", func() {
			functional_tests.Run(hoverctlBinary, "start")

			output := functional_tests.Run(hoverctlBinary, "targets")

			targets := functional_tests.TableToSliceMapStringString(output)
			Expect(targets).To(HaveKey("default"))
			Expect(targets["default"]).To(HaveKeyWithValue("TARGET NAME", "default"))
			Expect(targets["default"]).To(HaveKeyWithValue("ADMIN PORT", "8888"))
			Expect(targets["default"]).To(HaveKeyWithValue("PROXY PORT", "8500"))

			Expect(targets["default"]).To(HaveKey("PID"))
			Expect(strconv.Atoi(targets["default"]["PID"])).To(BeNumerically(">", 1))
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

			output = functional_tests.Run(hoverctlBinary, "logs")

			Expect(output).To(ContainSubstring("Default keys have been overwritten"))
		})

		It("should start an instance of  hoverfly with tls verification turned off", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--disable-tls")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs")

			Expect(output).To(ContainSubstring("tls certificate verification is now turned off!"))
		})

		It("should start an instance of hoverfly with persisted cache", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "--cache", "testdata-gen/cache.db", "-v")

			Expect(output).To(ContainSubstring("Hoverfly is now running"))

			output = functional_tests.Run(hoverctlBinary, "logs")
			Expect(output).To(ContainSubstring("Creating bolt db backend."))
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

			Expect(ioutil.ReadAll(response.Body)).To(ContainSubstring(`"upstream-proxy":"http://hoverfly.io:8080"`))
		})
	})

	Context("with a target", func() {
		It("should start an instance of Hoverfly with the target configuration", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create",
				"--target", "test-target",
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

		It("should update the target with a pid", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create",
				"--target", "test-target",
				"--admin-port", "4567",
				"--proxy-port", "4342",
			)

			functional_tests.Run(hoverctlBinary, "start", "--target", "test-target")

			output := functional_tests.Run(hoverctlBinary, "targets")

			targets := functional_tests.TableToSliceMapStringString(output)
			Expect(targets["test-target"]).To(HaveKey("PID"))
			Expect(strconv.Atoi(targets["test-target"]["PID"])).To(BeNumerically(">", 1))
		})
	})

	// Context("with a target that doesn't exist", func() {
	// 	It("should error", func() {
	// 		output := functional_tests.Run(hoverctlBinary, "start", "--target", "test-target")

	// 		Expect(output).To(ContainSubstring("test-target is not a target"))
	// 		Expect(output).To(ContainSubstring("Run `hoverctl start --new-target test-target`"))
	// 	})
	// })

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
	})

	Context("with a target that has already been started", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "start", "--new-target", "started")
		})

		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "start", "-t", "started")

			Expect(output).To(ContainSubstring("Hoverfly is already running"))
		})
	})

	Context("with port conflicts", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "start", "--new-target", "admin-conflict")
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
})
