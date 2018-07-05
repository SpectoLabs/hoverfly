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

func Test_iso8601DateTime(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.iso8601DateTime()).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_iso8601DateTimePlusDays(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.iso8601DateTimePlusDays("14")).To(Equal("2018-01-15T00:00:00Z"))
}

func Test_iso8601DateTimePlusDays_failure(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.iso8601DateTimePlusDays("cat")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_currentDateTime(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTime("Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:00 UTC 2018"))
}

func Test_currentDateTime_failure(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTime("cat")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_currentDateTimeAdd(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTimeAdd("1s", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:01 UTC 2018"))
	Expect(unit.currentDateTimeAdd("2m", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:02:00 UTC 2018"))
	Expect(unit.currentDateTimeAdd("3h", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 03:00:00 UTC 2018"))
	Expect(unit.currentDateTimeAdd("4d", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Fri Jan 5 00:00:00 UTC 2018"))
	Expect(unit.currentDateTimeAdd("5y", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Sat Dec 31 00:00:00 UTC 2022"))
	Expect(unit.currentDateTimeAdd("1y2d3h4m5s", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Thu Jan 3 03:04:05 UTC 2019"))
}

func Test_currentDateTimeAdd_failure(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTimeAdd("cat", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:00 UTC 2018"))
	Expect(unit.currentDateTimeAdd("1s", "cat")).To(Equal("2018-01-01T00:00:01Z"))
	Expect(unit.currentDateTimeAdd("cat", "cat")).To(Equal("2018-01-01T00:00:00Z"))
}

func Test_currentDateTimeSubtract(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTimeSubtract("1s", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Sun Dec 31 23:59:59 UTC 2017"))
	Expect(unit.currentDateTimeSubtract("2m", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Sun Dec 31 23:58:00 UTC 2017"))
	Expect(unit.currentDateTimeSubtract("3h", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Sun Dec 31 21:00:00 UTC 2017"))
	Expect(unit.currentDateTimeSubtract("4d", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Thu Dec 28 00:00:00 UTC 2017"))
	Expect(unit.currentDateTimeSubtract("5y", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Wed Jan 2 00:00:00 UTC 2013"))
	Expect(unit.currentDateTimeSubtract("1y2d3h4m5s", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Thu Dec 29 20:55:55 UTC 2016"))
}

func Test_currentDateTimeSubtract_failure(t *testing.T) {
	RegisterTestingT(t)

	unit := templateHelpers{
		now: testNow,
	}

	Expect(unit.currentDateTimeSubtract("cat", "Mon Jan 2 15:04:05 MST 2006")).To(Equal("Mon Jan 1 00:00:00 UTC 2018"))
	Expect(unit.currentDateTimeSubtract("1s", "cat")).To(Equal("2017-12-31T23:59:59Z"))
	Expect(unit.currentDateTimeSubtract("cat", "cat")).To(Equal("2018-01-01T00:00:00Z"))
}
