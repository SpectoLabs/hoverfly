package templating

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func testNow() time.Time {
	parsedTime, _ := time.Parse("2006-01-02T15:04:05Z", "2018-01-01T00:00:00Z")
	return parsedTime
}

func Test_now_withEmptyOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_now_withEmptyOffsetAndUnixFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "unix")).To(Equal("1514764800"))
}

func Test_now_withEmptyOffsetAndUnixMillisFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "epoch")).To(Equal("1514764800000"))
}

func Test_now_withEmptyOffsetAndCustomFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:00 UTC 2018"))
}

func Test_now_withPositiveOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("1d", "")).To(Equal("2018-01-02T00:00:00Z"))
}

func Test_now_withNegativeOffsetAndEmptyFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("-1d", "")).To(Equal("2017-12-31T00:00:00Z"))
}

func Test_now_withInvalidOffset(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("cat", "")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_now_withInvalidFormat(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{now: testNow}

	Expect(unit.nowHelper("", "dog")).To(Equal("dog"))
}

func Test_replace(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.replace("oink, oink, oink", "oink", "moo")).To(Equal("moo, moo, moo"))
}

func Test_faker(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{}

	Expect(unit.faker("JobTitle")[0].String()).To(Not(BeEmpty()))
}
