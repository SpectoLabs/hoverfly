package templating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aymerick/raymond"
	"github.com/pborman/uuid"
	"k8s.io/client-go/util/jsonpath"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	log "github.com/Sirupsen/logrus"
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

	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(query)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", query, err.Error())
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return ""
	}

	buf := new(bytes.Buffer)

	err = jsonPath.Execute(buf, data)
	if err != nil {
		log.Errorf("Failed to execute json path match: %s", err.Error())
		return ""
	}

	return buf.String()
}

func (t templateHelpers) xPath(query, toMatch string) string {
	xpathRule, err := goxpath.Parse(query)
	if err != nil {
		log.Errorf("Failed to parse xpath query %s: %s", query, err.Error())
		return ""
	}

	xTree, err := xmltree.ParseXML(bytes.NewBufferString(toMatch))
	if err != nil {
		log.Errorf("Failed to load XML tree: %s", err.Error())
		return ""
	}

	results, err := xpathRule.ExecNode(xTree)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return ""
	}
	return results.String()
}

func prepareJsonPathQuery(query string) string {
	if string(query[0:1]) != "{" && string(query[len(query)-1:]) != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}
