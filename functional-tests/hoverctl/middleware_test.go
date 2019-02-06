package hoverctl_suite

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

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
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
				"\n" +
				"\n..."))
		})

		It("I can get the hoverfly's middleware with script less than 5 lines", func() {
			hoverfly.SetMiddleware("node", "#!/usr/bin/env node\nprocess.stdin.resume();process.stdin.setEncoding('utf8');\nprocess.stdin.on('data', function(data) {var parsed_json = JSON.parse(data);process.stdout.write(JSON.stringify(parsed_json));});")
			output := functional_tests.Run(hoverctlBinary, "middleware")

			Expect(output).To(ContainSubstring("Hoverfly middleware configuration is currently set to"))
			Expect(output).To(ContainSubstring("Binary: node"))
			Expect(output).To(ContainSubstring("Script: #!/usr/bin/env node" +
				"\nprocess.stdin.resume();process.stdin.setEncoding('utf8');" +
				"\nprocess.stdin.on('data', function(data) {var parsed_json = JSON.parse(data);process.stdout.write(JSON.stringify(parsed_json));});"))
		})

		It("I can set the hoverfly's middleware with a binary and a script", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--binary", "python", "--script", "testdata/add_random_delay.py")

			Expect(output).To(ContainSubstring("Testing middleware against Hoverfly..."))
			Expect(output).To(ContainSubstring("Hoverfly middleware configuration has been set to"))
			Expect(output).To(ContainSubstring("Binary: python"))
			Expect(output).To(ContainSubstring("Script: #!/usr/bin/env python" +
				"\nimport sys" +
				"\nimport logging" +
				"\nimport random" +
				"\nfrom time import sleep" +
				"\n..."))
		})

		It("prints the full script when in verbose", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "-v", "--binary", "python", "--script", "testdata/add_random_delay.py")

			Expect(output).To(ContainSubstring("Testing middleware against Hoverfly..."))
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
				var newPairView v2.RequestResponsePairViewV1

				json.Unmarshal(body, &newPairView)

				newPairView.Response.Body = "modified body"

				pairViewBytes, _ := json.Marshal(newPairView)
				w.Write(pairViewBytes)
			}))
			defer middlewareServer.Close()

			output := functional_tests.Run(hoverctlBinary, "middleware", "--remote", middlewareServer.URL)

			Expect(output).To(ContainSubstring("Testing middleware against Hoverfly..."))
			Expect(output).To(ContainSubstring("Hoverfly middleware configuration has been set to"))
			Expect(output).To(ContainSubstring("Remote: " + middlewareServer.URL))
		})

		It("I cannae set the hoverfly's middleware when specifying non-existing file", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--binary", "python", "--script", "testdata/not_a_real_file.fake")

			Expect(output).To(ContainSubstring("File not found: testdata/not_a_real_file.fake"))
		})

		It("I cannae set the hoverfly's middleware when specifying non-existing remote middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--remote", "http://specto.io/404/nothere")

			Expect(output).To(ContainSubstring("Could not set middleware, it may have failed the test"))
			Expect(output).To(ContainSubstring("Error when communicating with remote middleware: received 404"))
			Expect(output).To(ContainSubstring("URL: http://specto.io/404/nothere"))
			Expect(output).To(ContainSubstring("STDIN:"))
			Expect(output).To(ContainSubstring(`{"response":{"status":200,"body":"ok","encodedBody":false,"headers":{"test_header":["true"]}},"request":{"path":"/","method":"GET","destination":"www.test.com","scheme":"","query":"","body":"","headers":{"test_header":["true"]}}}`))
		})

		It("When I use the verbose flag, I see that notpython is not an executable", func() {
			output := functional_tests.Run(hoverctlBinary, "-v", "middleware", "--binary", "notpython", "--script", "testdata/add_random_delay.py")

			Expect(output).To(ContainSubstring("Could not set middleware"))
			Expect(output).To(ContainSubstring(`exec: "notpython": executable file not found in $PATH`))
		})

		It("When I use the verbose flag, I see that http://wqrwwewf.wewefwef.specto is not an executable remote middleware", func() {
			output := functional_tests.Run(hoverctlBinary, "-v", "middleware", "--remote", "http://wqrwwewf.wewefwef.specto")

			Expect(output).To(ContainSubstring("Could not set middleware, it may have failed the test"))
			Expect(output).To(ContainSubstring(`Error when communicating with remote middleware:`))
			Expect(output).To(MatchRegexp("Post http://wqrwwewf.wewefwef.specto: dial tcp: lookup wqrwwewf.wewefwef.specto on|no such host"))
			Expect(output).To(ContainSubstring("URL: http://wqrwwewf.wewefwef.specto"))
			Expect(output).To(ContainSubstring("STDIN:"))
			Expect(output).To(ContainSubstring(`{"response":{"status":200,"body":"ok","encodedBody":false,"headers":{"test_header":["true"]}},"request":{"path":"/","method":"GET","destination":"www.test.com","scheme":"","query":"","body":"","headers":{"test_header":["true"]}}}`))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "middleware", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})
})
