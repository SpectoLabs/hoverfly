package main

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/mitchellh/go-homedir"
)


func Test_getHoverflyDirectory(t *testing.T) {
	RegisterTestingT(t)

	result := getHoverflyDirectory("/test/dir/config.yaml")

	Expect(result).To(Equal("/test/dir"))
}

func Test_getHoverflyDirectory_ReturnsHomeDirIfConfigUriIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	result := getHoverflyDirectory("")

	homeDir, _ := homedir.Dir()

	Expect(result).To(Equal(homeDir))
}
