package templating

import (
	"strconv"
	"time"

	"github.com/pborman/uuid"

	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/icrowley/fake"
)

type templateHelpers struct {
	now func() time.Time
}

func (t templateHelpers) iso8601DateTime() string {
	return t.now().UTC().Format("2006-01-02T15:04:05Z07:00")
}

func (t templateHelpers) iso8601DateTimePlusDays(days string) string {
	atoi, _ := strconv.Atoi(days)
	return t.now().AddDate(0, 0, atoi).UTC().Format("2006-01-02T15:04:05Z07:00")
}

func (t templateHelpers) currentDateTime(format string) string {
	formatted := t.now().UTC().Format(format)
	if formatted == format {
		return t.now().UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return formatted
}

func (t templateHelpers) currentDateTimeAdd(addTime string, format string) string {
	now := t.now()
	duration, err := ParseDuration(addTime)
	if err == nil {
		now = now.Add(duration)
	}
	formatted := now.UTC().Format(format)
	if formatted == format {
		return now.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return formatted
}

func (t templateHelpers) currentDateTimeSubtract(subtractTime string, format string) string {
	now := t.now()
	duration, err := ParseDuration(subtractTime)
	if err == nil {
		now = now.Add(-duration)
	}
	formatted := now.UTC().Format(format)
	if formatted == format {
		return now.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return formatted
}

func (t templateHelpers) randomString() string {
	return util.RandomString()
}

func (t templateHelpers) randomStringLength(length int) string {
	return util.RandomStringWithLength(length)
}

func (t templateHelpers) randomBoolean() string {
	return strconv.FormatBool(util.RandomBoolean())
}

func (t templateHelpers) randomInteger() string {
	return strconv.Itoa(util.RandomInteger())
}

func (t templateHelpers) randomIntegerRange(min, max int) string {
	return strconv.Itoa(util.RandomIntegerRange(min, max))
}

func (t templateHelpers) randomFloat() string {
	return strconv.FormatFloat(util.RandomFloat(), 'f', 6, 64)
}

func (t templateHelpers) randomFloatRange(min, max float64) string {
	return strconv.FormatFloat(util.RandomFloatRange(min, max), 'f', 6, 64)
}

func (t templateHelpers) randomEmail() string {
	return fake.EmailAddress()
}

func (t templateHelpers) randomIPv4() string {
	return fake.IPv4()
}

func (t templateHelpers) randomIPv6() string {
	return fake.IPv6()
}

func (t templateHelpers) randomUuid() string {
	return uuid.New()
}
