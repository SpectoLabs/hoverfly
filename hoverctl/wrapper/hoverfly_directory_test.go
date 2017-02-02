package wrapper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	. "github.com/onsi/gomega"
)

var hoverflyDirectory_testDirectory = "/tmp/hoverctl-hoverfly-directory-test"

func hoverflyDirectory_setup() {
	os.Mkdir(hoverflyDirectory_testDirectory, 0777)
}

func hoverflyDirectory_teardown() {
	os.RemoveAll(hoverflyDirectory_testDirectory)
}

func Test_getHomeDir_ReturnsHomeDirectoryOfSystem(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	result := getHomeDirectory()
	homeDirectory, _ := homedir.Dir()

	Expect(result).To(Equal(homeDirectory))

	hoverflyDirectory_teardown()
}

func Test_createHoverflyDirectory_CreatesADirectory(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	result := createHoverflyDirectory(hoverflyDirectory_testDirectory)

	Expect(result).To(Equal(hoverflyDirectory_testDirectory + "/.hoverfly"))

	fileInfo, _ := os.Stat(result)

	Expect(fileInfo.IsDir()).To(BeTrue())

	hoverflyDirectory_teardown()
}

func Test_GetPid_ReturnsAnIntOfTheFileContents(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	ioutil.WriteFile(hoverflyDirectory_testDirectory+"/hoverfly.9654.4332.pid", []byte("5555"), 0644)

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result, err := hoverflyDirectory.GetPid("9654", "4332")

	Expect(result).To(Equal(5555))
	Expect(err).To(BeNil())

	hoverflyDirectory_teardown()
}

func Test_GetPid_ReturnsZeroIfFileIsNotFound(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result, err := hoverflyDirectory.GetPid("8321", "3342")

	Expect(result).To(Equal(0))
	Expect(err).To(BeNil())

	hoverflyDirectory_teardown()
}

func Test_WritePid(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result := hoverflyDirectory.WritePid("8765", "4567", 16)

	Expect(result).To(BeNil())

	data, err := ioutil.ReadFile(hoverflyDirectory_testDirectory + "/hoverfly.8765.4567.pid")
	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal("16"))

	hoverflyDirectory_teardown()
}

func Test_WritePid_WhenPidExists_WontOverwrite(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	ioutil.WriteFile(hoverflyDirectory_testDirectory+"/hoverfly.4321.6654.pid", []byte("untouched"), 0644)

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result := hoverflyDirectory.WritePid("4321", "6654", 72)

	Expect(result).ToNot(BeNil())

	data, err := ioutil.ReadFile(hoverflyDirectory_testDirectory + "/hoverfly.4321.6654.pid")
	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal("untouched"))

	hoverflyDirectory_teardown()
}

func Test_DeletePid_WhenPidExists(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	ioutil.WriteFile(hoverflyDirectory_testDirectory+"/hoverfly.4332.3342.pid", []byte("1234"), 0644)

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result := hoverflyDirectory.DeletePid("4332", "3342")

	Expect(result).To(BeNil())

	_, err := os.Stat(hoverflyDirectory_testDirectory + "/hoverfly.4332.3342.pid")

	Expect(os.IsExist(err)).To(BeFalse())

	hoverflyDirectory_teardown()
}

func Test_DeletePid_WhenPidDoesNotExist(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	hoverflyDirectory := HoverflyDirectory{
		Path: hoverflyDirectory_testDirectory,
	}

	result := hoverflyDirectory.DeletePid("4332", "3342")

	Expect(result).ToNot(BeNil())

	hoverflyDirectory_teardown()
}

func Test_buildPidFilePath(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	hoverflyDirectory := HoverflyDirectory{
		Path: "/test",
	}

	result := hoverflyDirectory.buildPidFilePath("1234", "5678")

	Expect(result).To(Equal("/test/hoverfly.1234.5678.pid"))

	hoverflyDirectory_teardown()
}

func Test_buildPidFilePath_WithDifferentData(t *testing.T) {
	RegisterTestingT(t)
	hoverflyDirectory_setup()

	hoverflyDirectory := HoverflyDirectory{
		Path: "/another/test",
	}

	result := hoverflyDirectory.buildPidFilePath("4321", "8765")

	Expect(result).To(Equal("/another/test/hoverfly.4321.8765.pid"))

	hoverflyDirectory_teardown()
}
