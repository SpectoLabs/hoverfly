package hoverfly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

const pythonMiddlewareBasic = "import sys\nprint(sys.stdin.readlines()[0])"

const pythonModifyResponse = "#!/usr/bin/env python\n" +
	"import sys\n" +
	"import json\n" +

	"def main():\n" +
	"	data = sys.stdin.readlines()\n" +
	"	payload = data[0]\n" +

	"	payload_dict = json.loads(payload)\n" +

	"	payload_dict['response']['status'] = 201\n" +
	"	payload_dict['response']['body'] = \"body was replaced by middleware\"\n" +

	"	print(json.dumps(payload_dict))\n" +

	"if __name__ == \"__main__\":\n" +
	"	main()\n"

const pythonReflectBody = "#!/usr/bin/env python\n" +
	"import sys\n" +
	"import json\n" +

	"def main():\n" +
	"	data = sys.stdin.readlines()\n" +
	"	payload = data[0]\n" +

	"	payload_dict = json.loads(payload)\n" +

	"	payload_dict['response']['status'] = 201\n" +
	"	payload_dict['response']['body'] = payload_dict['request']['body']\n" +

	"	print(json.dumps(payload_dict))\n" +

	"if __name__ == \"__main__\":\n" +
	"	main()\n"

const pythonMiddlewareBad = "this shouldn't work"

const rubyEcho = "#!/usr/bin/env ruby\n" +
	"# encoding: utf-8\n" +
	"while payload = STDIN.gets\n" +
	"  next unless payload\n" +
	"\n" +
	"  STDOUT.puts payload\n" +
	"\n" +
	"  STDERR.puts \"Payload data: #{payload}\"\n" +
	"\n" +
	"end"

func processHandlerOkayButNoResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func processHandlerNotOkay(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func Test_ConvertToNewMiddleware_WillCreateAMiddlewareObjectFromAUrlString(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkayButNoResponse).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	unit, err := ConvertToNewMiddleware(server.URL + "/process")
	Expect(err).To(BeNil())

	Expect(unit.Binary).To(Equal(""))
	Expect(unit.Script).To(BeNil())
	Expect(unit.Remote).To(Equal(server.URL + "/process"))
}

func Test_ConvertToNewMiddleware_ReturnsErrorIfUrlIsUnaccessible(t *testing.T) {
	RegisterTestingT(t)

	unit, err := ConvertToNewMiddleware("http://specto.io/404/process")
	Expect(err).ToNot(BeNil())

	Expect(unit).To(BeNil())
}

func Test_ConvertToNewMiddleware_WillCreateAMiddlewareObjectFromASingleBinary(t *testing.T) {
	RegisterTestingT(t)

	unit, err := ConvertToNewMiddleware("cat")
	Expect(err).To(BeNil())

	Expect(unit.Binary).To(Equal("cat"))
	Expect(unit.Script).To(BeNil())
	Expect(unit.Remote).To(Equal(""))
}

func Test_ConvertToNewMiddleware_ReturnsErrorIfBinaryIsUnaccessible(t *testing.T) {
	RegisterTestingT(t)

	unit, err := ConvertToNewMiddleware("fake-binary")
	Expect(err).ToNot(BeNil())

	Expect(unit).To(BeNil())
}

func Test_ConvertToNewMiddleware_WillCreateAMiddlewareObjectFromASingleBinaryAndScript(t *testing.T) {
	RegisterTestingT(t)

	unit, err := ConvertToNewMiddleware("python examples/middleware/reflect_body/reflect_body.py")
	Expect(err).To(BeNil())

	Expect(unit.Binary).To(Equal("python"))
	Expect(unit.Script).ToNot(BeNil())

	script, err := unit.GetScript()
	Expect(err).To(BeNil())

	Expect(script).ToNot(BeNil())

	Expect(unit.Remote).To(Equal(""))
}

func TestChangeBodyMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("body was replaced by middleware"))
}

func TestMalformedRequestResponsePairWithMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	malformedPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("ruby")
	Expect(err).To(BeNil())

	err = unit.SetScript(rubyEcho)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(malformedPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("original body"))
}

func TestReflectBody(t *testing.T) {
	RegisterTestingT(t)

	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	originalPair := models.RequestResponsePair{Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal(req.Body))
	Expect(newPair.Request.Method).To(Equal(req.Method))
	Expect(newPair.Request.Destination).To(Equal(req.Destination))
}

func processHandlerOkay(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var newPairView v1.RequestResponsePairView

	json.Unmarshal(body, &newPairView)

	newPairView.Response.Body = "You got straight up messed with"

	pairViewBytes, _ := json.Marshal(newPairView)
	w.Write(pairViewBytes)
}

func TestExecuteMiddlewareRemotely(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkay).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{}

	err := unit.SetRemote(server.URL + "/process")
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).To(BeNil())

	Expect(newPair).ToNot(Equal(originalPair))
	Expect(newPair.Response.Body).To(Equal("You got straight up messed with"))
}

func Test_Middleware_executeMiddlewareRemotely_ReturnsErrorIfDoesntGetA200_AndSameRequestResponsePairs(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerNotOkay).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{}

	err := unit.SetRemote(server.URL + "/process")
	Expect(err).ToNot(BeNil())

	newPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Error when communicating with remote middleware"))

	Expect(newPair).To(Equal(originalPair))
}

func Test_Middleware_executeMiddlewareRemotely_ReturnsErrorIfNoRequestResponsePairOnResponse_TheUntouchedPairIsReturned(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkayButNoResponse).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{}

	err := unit.SetRemote(server.URL + "/process")
	Expect(err).To(BeNil())

	untouchedPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("unexpected end of JSON input"))

	Expect(untouchedPair).To(Equal(originalPair))
}

func Test_Middleware_SetBinary_SetsBinaryIfItCanRunIt(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetBinary("go")
	Expect(err).To(BeNil())
	Expect(unit.Binary).To(Equal("go"))
}

func Test_Middleware_SetBinary_DoesNotSetIfCantRun(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetBinary("|{}|")
	Expect(err).ToNot(BeNil())
	Expect(unit.Binary).To(Equal(""))
}

func Test_Middleware_SetBinary_SetsStringToEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Binary: "test",
	}

	err := unit.SetBinary("")
	Expect(err).To(BeNil())
	Expect(unit.Binary).To(Equal(""))
}

func Test_Middleware_SetScript_WritesScriptToFile(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript("just a test")
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	fileContents, err := ioutil.ReadFile(unit.Script.Name())
	Expect(err).To(BeNil())

	Expect(string(fileContents)).To(Equal("just a test"))
}

func Test_Middleware_SetScript_DeletesPreviousScript(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript("just a test")
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	firstScript := unit.Script

	err = unit.SetScript("just a test 2")
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	_, err = ioutil.ReadFile(firstScript.Name())
	Expect(err).ToNot(BeNil())
}

func Test_Middleware_GetScript_GetsScript(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript("just a test")
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	result, err := unit.GetScript()
	Expect(err).To(BeNil())
	Expect(result).To(Equal("just a test"))
}

func Test_Middleware_GetScript_DoesNotErrorIfScriptIsNotSet(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	result, err := unit.GetScript()
	Expect(err).To(BeNil())
	Expect(result).To(Equal(""))
}

func Test_Middleware_DeleteScripts_WillDeleteScriptAndSetScriptToNil(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript("just a test")
	Expect(err).To(BeNil())

	firstScript := unit.Script

	err = unit.DeleteScripts(path.Join(os.TempDir(), "hoverfly"))
	Expect(err).To(BeNil())
	Expect(unit.Script).To(BeNil())

	_, err = ioutil.ReadFile(firstScript.Name())
	Expect(err).ToNot(BeNil())
}

func Test_Middleware_DeleteScripts_WillDeletePreviousScripts(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript("just a test")
	Expect(err).To(BeNil())

	firstScript := unit.Script

	err = ioutil.WriteFile(path.Join(os.TempDir(), "hoverfly", "test"), []byte("test"), 0644)
	Expect(err).To(BeNil())

	err = unit.DeleteScripts(path.Join(os.TempDir(), "hoverfly"))
	Expect(err).To(BeNil())
	Expect(unit.Script).To(BeNil())

	_, err = ioutil.ReadFile(firstScript.Name())
	Expect(err).ToNot(BeNil())

	_, err = ioutil.ReadFile(path.Join(os.TempDir(), "hoverfly", "test"))
	Expect(err).ToNot(BeNil())

}

func Test_Middleware_DeleteScripts_DoesNotErrorIfNoScriptWasSet(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.DeleteScripts(path.Join(os.TempDir(), "hoverfly"))
	Expect(err).To(BeNil())
}

func Test_Middleware_SetScript_WritesMultiLineStringScriptToFile(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetScript(rubyEcho)
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	fileContents, err := ioutil.ReadFile(unit.Script.Name())
	Expect(err).To(BeNil())

	Expect(string(fileContents)).To(Equal(rubyEcho))
}

func Test_Middleware_Execute_RunsMiddlewareCorrectly(t *testing.T) {
	RegisterTestingT(t)

	binary := "python"
	script := "#!/usr/bin/env python\n" +
		"import sys\n" +
		"import json\n" +
		"\n" +
		"def main():\n" +
		"	data = sys.stdin.readlines()\n" +
		"	payload = data[0]\n" +
		"\n" +
		"	payload_dict = json.loads(payload)\n" +
		"\n" +
		"	payload_dict['response']['status'] = 200" +
		"\n" +
		"	print(json.dumps(payload_dict))\n" +
		"\n" +
		"if __name__ == \"__main__\":\n" +
		"	main()"

	unit := Middleware{}

	err := unit.SetScript(script)
	Expect(err).To(BeNil())
	Expect(unit.Script).ToNot(BeNil())

	err = unit.SetBinary(binary)
	Expect(err).To(BeNil())
	Expect(unit.Binary).To(Equal(binary))

	resp := models.ResponseDetails{Status: 0, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	resultPair, err := unit.Execute(originalPair)
	Expect(err).To(BeNil())

	Expect(resultPair.Response.Status).To(Equal(200))
}

func Test_Middleware_Execute_WillErrorIfMiddlewareHasNotBeenCorrectlySet(t *testing.T) {

	RegisterTestingT(t)

	unit := Middleware{}

	_, err := unit.Execute(models.RequestResponsePair{})
	Expect(err).ToNot(BeNil())

	Expect(err.Error()).To(Equal("Cannot execute middleware as middleware has not been correctly set"))
}

func Test_Middleware_SetRemote_CanSetRemote(t *testing.T) {
	RegisterTestingT(t)

	remoteMiddleware := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var newPairView v2.RequestResponsePairView

		json.Unmarshal(body, &newPairView)

		newPairView.Response.Body = "modified body"

		pairViewBytes, _ := json.Marshal(newPairView)
		w.Write(pairViewBytes)
	}))
	defer remoteMiddleware.Close()

	unit := Middleware{}

	err := unit.SetRemote(remoteMiddleware.URL)
	Expect(err).To(BeNil())

	Expect(unit.Remote).To(Equal(remoteMiddleware.URL))
}

func Test_Middleware_SetRemote_WontSetRemoteIfRemoteDoesntExist(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetRemote("http://www.specto.io/madeupmiddlewareendpoint")
	Expect(err).ToNot(BeNil())

	Expect(unit.Remote).To(Equal(""))
}

func Test_Middleware_SetRemote_WontSetRemoteIfCouldNeverExist(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetRemote("[]somemadeupwebsite*&*^&$%^")
	Expect(err).ToNot(BeNil())

	Expect(unit.Remote).To(Equal(""))
}

func Test_Middleware_SetRemote_CanBeSetToEmptyStringWithoutError(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{}

	err := unit.SetRemote("")
	Expect(err).To(BeNil())

	Expect(unit.Remote).To(Equal(""))
}

func Test_Middleware_Execute_RunsRemoteMiddlewareCorrectly(t *testing.T) {
	RegisterTestingT(t)

	middlewareServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var newPairView v2.RequestResponsePairView

		json.Unmarshal(body, &newPairView)

		newPairView.Response.Body = "modified body"

		pairViewBytes, _ := json.Marshal(newPairView)
		w.Write(pairViewBytes)
	}))
	defer middlewareServer.Close()

	unit := Middleware{}
	err := unit.SetRemote(middlewareServer.URL)
	Expect(err).To(BeNil())

	resp := models.ResponseDetails{Status: 0, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	resultPair, err := unit.Execute(originalPair)
	Expect(err).To(BeNil())

	Expect(resultPair.Response.Body).To(Equal("modified body"))
}

func Test_Middleware_IsSet_WillSayItsSetIfARemoteIsDefined(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Remote: "test-remote",
	}

	Expect(unit.IsSet()).To(BeTrue())
}

func Test_Middleware_IsSet_WillSayItsSetIfABinaryIsDefined(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Binary: "test-binary",
	}

	Expect(unit.IsSet()).To(BeTrue())
}

func Test_Middleware_IsSet_WillSayItsNotSetIfAOnlyAScriptIsDefined(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Script: os.NewFile(0, "testfile.txt"),
	}

	Expect(unit.IsSet()).To(BeFalse())
}

func Test_Middleware_toString_WillProduceAStringRepresentationOfMiddlewareThatUsesRemote(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Remote: "test-remote",
	}

	Expect(unit.toString()).To(Equal("test-remote"))
}

func Test_Middleware_toString_WillProduceAStringRepresentationOfMiddlewareThatUsesBinary(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Binary: "test-binary",
	}

	Expect(unit.toString()).To(Equal("test-binary"))
}

func Test_Middleware_toString_WillProduceAStringRepresentationOfMiddlewareThatUsesBinaryAndScript(t *testing.T) {
	RegisterTestingT(t)

	unit := Middleware{
		Binary: "test-binary",
		Script: os.NewFile(0, "testfile.txt"),
	}

	Expect(unit.toString()).To(Equal("test-binary testfile.txt"))
}
