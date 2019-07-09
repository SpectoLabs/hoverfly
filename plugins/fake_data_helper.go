package main

import (
	"github.com/icrowley/fake"
)


type FakeDataHelper struct {

}
func (t FakeDataHelper) Helper() string {
	return fake.FullName()
}

func (t FakeDataHelper) Name() string {
	return "fake"
}

var HelperDef FakeDataHelper