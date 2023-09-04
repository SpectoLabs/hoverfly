package templating

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aymerick/raymond"
	"github.com/pborman/uuid"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/icrowley/fake"
	log "github.com/sirupsen/logrus"
)

const defaultDateTimeFormat = "2006-01-02T15:04:05Z07:00"

type templateHelpers struct {
	now         func() time.Time
	fakerSource *gofakeit.Faker
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
	return fetchFromRequestBody(queryType, query, toMatch)
}

func (t templateHelpers) parseCsv(dataSourceName, searchFieldName, searchFieldValue, returnFieldName string, options *raymond.Options) string {

	templateDataSources := options.Value("TemplateDataSources").(map[string]*DataSource)
	source, exists := templateDataSources[dataSourceName]
	if exists {
		searchIndex, err := getHeaderIndex(source.Data, searchFieldName)
		if err != nil {
			log.Error(err)
			getCSVEvaluationString(options)
		}
		returnIndex, err := getHeaderIndex(source.Data, returnFieldName)
		if err != nil {
			log.Error(err)
			return getCSVEvaluationString(options)
		}
		var fallbackString string
		for i := 1; i < len(source.Data); i++ {
			record := source.Data[i]
			if strings.ToLower(record[searchIndex]) == strings.ToLower(searchFieldValue) {
				return record[returnIndex]
			} else if record[searchIndex] == "*" {
				fallbackString = record[returnIndex]
			}
		}

		if fallbackString != "" {
			return fallbackString
		}

	}
	return getCSVEvaluationString(options)

}

func getCSVEvaluationString(options *raymond.Options) string {

	evaluationString := "{{ csv "
	for _, params := range options.Params() {
		evaluationString = evaluationString + params.(string) + ` `
	}
	return evaluationString + "}}"
}

func getHeaderIndex(data [][]string, headerName string) (int, error) {

	if len(data) == 0 {
		return -1, fmt.Errorf("empty file provided")
	}
	headerRecord := data[0]
	for index, fieldName := range headerRecord {
		if strings.ToLower(fieldName) == strings.ToLower(headerName) {
			return index, nil
		}
	}
	return -1, fmt.Errorf("search field %s does not found", headerName)
}

func fetchFromRequestBody(queryType, query, toMatch string) string {

	if queryType == "jsonpath" {
		return jsonPath(query, toMatch)
	} else if queryType == "xpath" {
		return xPath(query, toMatch)
	}
	log.Errorf("Unknown query type \"%s\" for templating Request.Body", queryType)
	return ""
}

func jsonPath(query, toMatch string) string {
	query = prepareJsonPathQuery(query)

	result, err := matchers.JsonPathExecution(query, toMatch)
	if err != nil {
		return ""
	}
	return result
}

func xPath(query, toMatch string) string {
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

func (t templateHelpers) faker(fakerType string) []reflect.Value {

	if t.fakerSource == nil {
		t.fakerSource = gofakeit.New(0)
	}
	if reflect.ValueOf(t.fakerSource).MethodByName(fakerType).IsValid() {
		return reflect.ValueOf(t.fakerSource).MethodByName(fakerType).Call([]reflect.Value{})
	}
	return []reflect.Value{}
}
