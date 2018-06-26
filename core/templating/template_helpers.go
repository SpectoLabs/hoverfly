package templating

import (
	"strconv"
	"time"

	"github.com/satori/go.uuid"

	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/icrowley/fake"
)

func iso8601DateTime() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
}

func iso8601DateTimePlusDays(days string) string {
	atoi, _ := strconv.Atoi(days)
	return time.Now().AddDate(0, 0, atoi).UTC().Format("2006-01-02T15:04:05Z07:00")
}

func randomString() string {
	return util.RandomString()
}

func randomStringLength(length int) string {
	return util.RandomStringWithLength(length)
}

func randomBoolean() string {
	return strconv.FormatBool(util.RandomBoolean())
}

func randomInteger() string {
	return strconv.Itoa(util.RandomInteger())
}

func randomIntegerRange(min, max int) string {
	return strconv.Itoa(util.RandomIntegerRange(min, max))
}

func randomFloat() string {
	return strconv.FormatFloat(util.RandomFloat(), 'f', 6, 64)
}

func randomFloatRange(min, max float64) string {
	return strconv.FormatFloat(util.RandomFloatRange(min, max), 'f', 6, 64)
}

func randomEmail() string {
	return fake.EmailAddress()
}

func randomIPv4() string {
	return fake.IPv4()
}

func randomIPv6() string {
	return fake.IPv6()
}

func randomUuid() string {
	return uuid.NewV4().String()
}
