package action_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/action"
	. "github.com/onsi/gomega"
)

func Test_SetPostServeActionMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewAction("test-callback", "python3", "dummy script", 1800)
	Expect(err).To(BeNil())

	unit := action.NewPostServeActionDetails()
	err = unit.SetAction("test-callback", newAction)

	Expect(err).To(BeNil())
	Expect(unit.Actions).To(HaveLen(1))
	Expect(unit.Actions["test-callback"].Binary).To(Equal("python3"))
	Expect(unit.Actions["test-callback"].DelayInMs).To(Equal(1800))
}

func Test_DeletePostServeActionMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewAction("test-callback", "python3", "dummy script", 1800)
	Expect(err).To(BeNil())

	unit := action.NewPostServeActionDetails()
	err = unit.SetAction("test-callback", newAction)
	Expect(err).To(BeNil())

	err = unit.DeleteAction("test-callback")
	Expect(err).To(BeNil())
	Expect(unit.Actions).To(HaveLen(0))
}
