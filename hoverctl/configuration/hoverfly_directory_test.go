package configuration

import (
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
