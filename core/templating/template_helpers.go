package templating

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aymerick/raymond"
	"github.com/pborman/uuid"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/icrowley/fake"
	log "github.com/sirupsen/logrus"
)

const defaultDateTimeFormat = "2006-01-02T15:04:05Z07:00"

type templateHelpers struct {
	now func() time.Time
}

func (t templateHelpers) iso8601DateTime() string {
	return t.now().UTC().Format(defaultDateTimeFormat)
}

func (t templateHelpers) iso8601DateTimePlusDays(days string) string {
	atoi, _ := strconv.Atoi(days)
	return t.now().AddDate(0, 0, atoi).UTC().Format(defaultDateTimeFormat)
}

func (t templateHelpers) currentDateTime(format string) string {
	return t.now().UTC().Format(format)
}

func (t templateHelpers) currentDateTimeAdd(addTime string, format string) string {
	now := t.now()
	duration, err := ParseDuration(addTime)
	if err == nil {
		now = now.Add(duration)
	}
	return now.UTC().Format(format)
}

func (t templateHelpers) currentDateTimeSubtract(subtractTime string, format string) string {
	now := t.now()
	duration, err := ParseDuration(subtractTime)
	if err == nil {
		now = now.Add(-duration)
	}
	return now.UTC().Format(format)
}

func (t templateHelpers) nowHelper(offset string, format string) string {
	now := t.now()
	if offset != "" {
		duration, err := ParseDuration(offset)
		if err == nil {
			now = now.Add(duration)
		}
	}

	var formatted string
	if format == "" {
		formatted = now.UTC().Format(defaultDateTimeFormat)
	} else if format == "unix" {
		formatted = strconv.FormatInt(now.Unix(), 10)
	} else if format == "epoch" {
		formatted = strconv.FormatInt(now.UnixNano()/1000000, 10)
	} else {
		formatted = now.UTC().Format(format)
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

func (t templateHelpers) requestBody(queryType, query string, options *raymond.Options) string {
	toMatch := options.Value("request").(Request).body
	queryType = strings.ToLower(queryType)
	if queryType == "jsonpath" {
		return t.jsonPath(query, toMatch)
	} else if queryType == "xpath" {
		return t.xPath(query, toMatch)
	}
	log.Errorf("Unknown query type \"%s\" for templating Request.Body", queryType)
	return ""
}

func (t templateHelpers) jsonPath(query, toMatch string) string {
	query = prepareJsonPathQuery(query)

	result, err := matchers.JsonPathExecution(query, toMatch)
	if err != nil {
		return ""
	}
	return result
}

func (t templateHelpers) xPath(query, toMatch string) string {
	result, err := matchers.XpathExecution(query, toMatch)
	if err != nil {
		return ""
	}
	return result.String()
}

func (t templateHelpers) replace(target, oldValue, newValue string) string {
	return strings.Replace(target, oldValue, newValue, -1)
}

func prepareJsonPathQuery(query string) string {
	if query[0:1] != "{" && query[len(query)-1:] != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}
