package configuration

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	. "github.com/onsi/gomega"
)

var hoverflydirectoryTestdirectory = "/tmp/hoverctl-hoverfly-directory-test"

func hoverflyDirectory_setup() {
	os.Mkdir(hoverflydirectoryTestdirectory, 0777)
}

func hoverflyDirectory_teardown() {
	os.RemoveAll(hoverflydirectoryTestdirectory)
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

	result := createHoverflyDirectory(hoverflydirectoryTestdirectory)

	Expect(result).To(Equal(hoverflydirectoryTestdirectory + "/.hoverfly"))

	fileInfo, _ := os.Stat(result)

	Expect(fileInfo.IsDir()).To(BeTrue())

	hoverflyDirectory_teardown()
}
