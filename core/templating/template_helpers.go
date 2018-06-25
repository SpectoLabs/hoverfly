package templating

import (
	"strconv"
	"time"

	"github.com/SpectoLabs/hoverfly/core/util"
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
