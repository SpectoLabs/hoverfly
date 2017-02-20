package hoverctl_end_to_end

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Describe("with a running hoverfly which has middleware configured", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
			hoverfly.SetMiddleware("ruby", "#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend")

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("I can get the hoverfly's middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware")

			Expect(output).To(ContainSubstring("Hoverfly middleware configuration is currently set to"))
			Expect(output).To(ContainSubstring("Binary: ruby"))
			Expect(output).To(ContainSubstring("Script: #!/usr/bin/env ruby" +
				"\n# encoding: utf-8" +
				"\nwhile payload = STDIN.gets" +
				"\nnext unless payload" +

				"\n\nSTDOUT.puts payload" +
				"\nend"))
		})

		It("I can set the hoverfly's middleware with a binary and a script", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--binary", "python", "--script", "testdata/add_random_delay.py")

			Expect(output).To(ContainSubstring("Hoverfly middleware configuration has been set to"))
			Expect(output).To(ContainSubstring("Binary: python"))
			Expect(output).To(ContainSubstring("Script: #!/usr/bin/env python" +
				"\nimport sys" +
				"\nimport logging" +
				"\nimport random" +
				"\nfrom time import sleep" +
				"\n\n" +
				"logging.basicConfig(filename='random_delay_middleware.log', level=logging.DEBUG)" +
				"\nlogging.debug('Random delay middleware is called')" +
				"\n\n" +
				"# set delay to random value less than one second" +
				"\n\n" +
				"SLEEP_SECS = random.random()" +
				"\n\n" +
				"def main():" +
				"\n\n    data = sys.stdin.readlines()" +
				"\n    # this is a json string in one line so we are interested in that one line" +
				"\n    payload = data[0]" +
				"\n    logging.debug(\"sleeping for %s seconds\" % SLEEP_SECS)" +
				"\n    sleep(SLEEP_SECS)\n\n\n    # do not modifying payload, returning same one" +
				"\n    print(payload)" +
				"\n\nif __name__ == \"__main__\":" +
				"\n    main()"))
		})

		It("I can set the hoverfly's middleware with a remote", func() {

			middlewareServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				var newPairView v2.RequestResponsePairView

				json.Unmarshal(body, &newPairView)

				newPairView.Response.Body = "modified body"

				pairViewBytes, _ := json.Marshal(newPairView)
				w.Write(pairViewBytes)
			}))
			defer middlewareServer.Close()

			output := functional_tests.Run(hoverctlBinary, "middleware", "--remote", middlewareServer.URL)

			Expect(output).To(ContainSubstring("Hoverfly middleware configuration has been set to"))
			Expect(output).To(ContainSubstring("Remote: " + middlewareServer.URL))
		})

		It("I cannae set the hoverfly's middleware when specifying non-existing file", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--binary", "python", "--script", "testdata/not_a_real_file.fake")

			Expect(output).To(ContainSubstring("File not found: testdata/not_a_real_file.fake"))
		})

		It("I cannae set the hoverfly's middleware when specifying non-existing remote middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--remote", "http://specto.io/404/nothere")

			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
		})

		It("When I use the verbose flag, I see that notpython is not an executable", func() {
			output := functional_tests.Run(hoverctlBinary, "-v", "middleware", "--binary", "notpython", "--script", "testdata/add_random_delay.py")

			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
			Expect(output).To(ContainSubstring(`Invalid middleware: exec: \"notpython\": executable file not found in $PATH`))
		})

		It("When I use the verbose flag, I see that http://wqrwwewf.wewefwef.specto is not an executable remote middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "-v", "middleware", "--remote", "http://wqrwwewf.wewefwef.specto")

			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
			Expect(output).To(ContainSubstring(`Invalid middleware: Could not reach remote middleware`))
		})
	})
})
