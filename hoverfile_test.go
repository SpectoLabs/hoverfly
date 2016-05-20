package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_NewHoverfile_CanCreateAHoverfileFromCompleteKey(t *testing.T) {
	RegisterTestingT(t)

	hoverfile, err := NewHoverfile("testvendor/testname:v1")

	Expect(err).To(BeNil())

	Expect(hoverfile.Vendor).To(Equal("testvendor"))
	Expect(hoverfile.Name).To(Equal("testname"))
	Expect(hoverfile.Version).To(Equal("v1"))
}

func Test_NewHoverfile_CanCreateAHoverfileFromDifferentCompleteKey(t *testing.T) {
	RegisterTestingT(t)

	hoverfile, err := NewHoverfile("another-vendor/test_simulation:v7")

	Expect(err).To(BeNil())

	Expect(hoverfile.Vendor).To(Equal("another-vendor"))
	Expect(hoverfile.Name).To(Equal("test_simulation"))
	Expect(hoverfile.Version).To(Equal("v7"))
}

func Test_NewHoverfile_CanCreateAHoverfileFromKey_WithNoVersion(t *testing.T) {
	RegisterTestingT(t)

	hoverfile, err := NewHoverfile("tester/tested")

	Expect(err).To(BeNil())

	Expect(hoverfile.Vendor).To(Equal("tester"))
	Expect(hoverfile.Name).To(Equal("tested"))
	Expect(hoverfile.Version).To(Equal("v1"))
}

func Test_NewHoverfile_CanCreateAHoverfileFromKey_WithNoVendor(t *testing.T) {
	RegisterTestingT(t)

	hoverfile, err := NewHoverfile("just_a-name")
	Expect(err).To(BeNil())

	Expect(hoverfile.Vendor).To(Equal(""))
	Expect(hoverfile.Name).To(Equal("just_a-name"))
	Expect(hoverfile.Version).To(Equal("v1"))
}

func Test_NewHoverfile_WontCreateAHoverfileFromKey_WithSpecialCharacters(t *testing.T) {
	RegisterTestingT(t)

	hoverfile, err := NewHoverfile("just_@-name")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in hoverfile name"))
	Expect(hoverfile).To(Equal(Hoverfile{}))

	hoverfile, err = NewHoverfile("just_\\-name")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in hoverfile name"))
	Expect(hoverfile).To(Equal(Hoverfile{}))

	hoverfile, err = NewHoverfile("just()an&simulation")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in hoverfile name"))
	Expect(hoverfile).To(Equal(Hoverfile{}))

	hoverfile, err = NewHoverfile("just()an£im%lation")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in hoverfile name"))
	Expect(hoverfile).To(Equal(Hoverfile{}))
}

func Test_Hoverfile_GetFileName(t *testing.T) {
	RegisterTestingT(t)

	hoverfile := Hoverfile{
		Vendor: "vendor",
		Name: "name",
		Version: "version",
	}

	resultFileName := hoverfile.GetFileName()
	Expect(resultFileName).To(Equal("vendor.name.version.hfile"))
}

func Test_Hoverfile_String(t *testing.T) {
	RegisterTestingT(t)

	hoverfile := Hoverfile{
		Vendor: "vendor",
		Name: "name",
		Version: "version",
	}

	resultString := hoverfile.String()
	Expect(resultString).To(Equal("vendor/name:version"))

}