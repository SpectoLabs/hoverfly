package wrapper

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

var logFile_testDirectory = "/tmp/hoverctl-tests-logfile"

func logFile_setup() {
	os.Mkdir(logFile_testDirectory, 0777)
}

func logFile_teardown() {
	os.RemoveAll(logFile_testDirectory)
}

func Test_NewLogfile(t *testing.T) {
	RegisterTestingT(t)
	logFile_setup()
	defer logFile_teardown()

	hoverflyDirectory := HoverflyDirectory{
		Path: "/test/dir",
	}

	logfile := NewLogFile(hoverflyDirectory, "1234", "4321")

	Expect(logfile.Path).To(Equal("/test/dir/hoverfly.1234.4321.log"))
	Expect(logfile.Name).To(Equal("hoverfly.1234.4321.log"))
}

func Test_Logfile_GetLogs(t *testing.T) {
	RegisterTestingT(t)
	logFile_setup()
	defer logFile_teardown()

	hoverflyDirectory := HoverflyDirectory{
		Path: logFile_testDirectory,
	}

	logfile := NewLogFile(hoverflyDirectory, "9856", "6589")

	ioutil.WriteFile(logfile.Path, []byte("testing rocks"), 0644)

	logs, err := logfile.getLogs()
	Expect(err).To(BeNil())

	Expect(logs).To(Equal("testing rocks"))
}

func Test_Logfile_GetLogs_ReturnsErrorIfLogDoesNotExist(t *testing.T) {
	RegisterTestingT(t)
	logFile_setup()
	defer logFile_teardown()

	hoverflyDirectory := HoverflyDirectory{
		Path: logFile_testDirectory,
	}

	logfile := NewLogFile(hoverflyDirectory, "9856", "6589")

	_, err := logfile.getLogs()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not open Hoverfly log file"))
}
