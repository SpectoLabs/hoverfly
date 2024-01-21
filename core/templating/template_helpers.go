package templating

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/SpectoLabs/raymond"
	"github.com/pborman/uuid"

	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/icrowley/fake"
	log "github.com/sirupsen/logrus"
)

const defaultDateTimeFormat = "2006-01-02T15:04:05Z07:00"

type templateHelpers struct {
	now                func() time.Time
	fakerSource        *gofakeit.Faker
	TemplateDataSource *TemplateDataSource
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

func (t templateHelpers) requestBody(queryType, query string, options *raymond.Options) interface{} {
	body := ""
	if toMatch, exists := options.Value("request").(Request); exists {
		body = toMatch.body
	} else {
		journalToMatch := options.Value("Request").(journal.Request)
		body = journalToMatch.BodyStr
	}
	queryType = strings.ToLower(queryType)
	return util.FetchFromRequestBody(queryType, query, body)
}

func (t templateHelpers) replace(target, oldValue, newValue string) string {
	return strings.Replace(target, oldValue, newValue, -1)
}

func (t templateHelpers) split(target, separator string) []string {
	return strings.Split(target, separator)
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

func (t templateHelpers) parseCsv(dataSourceName, searchFieldName, searchFieldValue, returnFieldName string, options *raymond.Options) string {

	templateDataSources := t.TemplateDataSource.DataSources
	source, exists := templateDataSources[dataSourceName]
	if exists {
		searchIndex, err := getHeaderIndex(source.Data, searchFieldName)
		if err != nil {
			log.Error(err)
			getEvaluationString("csv", options)
		}
		returnIndex, err := getHeaderIndex(source.Data, returnFieldName)
		if err != nil {
			log.Error(err)
			return getEvaluationString("csv", options)
		}

		var fallbackString string
		searchFieldValue := getSearchFieldValue(options, searchFieldValue)
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
	return getEvaluationString("csv", options)

}

func (t templateHelpers) parseJournalBasedOnIndex(indexName, keyValue, dataSource, queryType, lookupQuery string, options *raymond.Options) interface{} {
	journalDetails := options.Value("Journal").(Journal)
	if journalEntry, err := getIndexEntry(journalDetails, indexName, keyValue); err == nil {
		if body := getBodyDataToParse(dataSource, journalEntry); body != "" {
			data := util.FetchFromRequestBody(queryType, lookupQuery, body)
			if _, ok := data.(error); ok {
				// The interface is an error
				return getEvaluationString("journal", options)
			} else {
				return data
			}
		}
	}
	return getEvaluationString("journal", options)
}

func (t templateHelpers) sum(numbers []string, format string) string {
	return sumNumbers(numbers, format)
}

func (t templateHelpers) add(val1 string, val2 string, format string) string {
	return sumNumbers([]string{val1, val2}, format)
}

func (t templateHelpers) subtract(val1 string, val2 string, format string) string {
	f1, err1 := strconv.ParseFloat(val1, 64)
	f2, err2 := strconv.ParseFloat(val2, 64)
	if err1 != nil || err2 != nil {
		return "NaN"
	}
	return formatNumber(f1-f2, format)
}

func (t templateHelpers) multiply(val1 string, val2 string, format string) string {
	f1, err1 := strconv.ParseFloat(val1, 64)
	f2, err2 := strconv.ParseFloat(val2, 64)
	if err1 != nil || err2 != nil {
		return "NaN"
	}
	return formatNumber(f1*f2, format)
}

func (t templateHelpers) divide(val1 string, val2 string, format string) string {
	f1, err1 := strconv.ParseFloat(val1, 64)
	f2, err2 := strconv.ParseFloat(val2, 64)
	if err1 != nil || err2 != nil {
		return "NaN"
	}
	return formatNumber(f1/f2, format)
}

func (t templateHelpers) addToArray(key string, value string, options *raymond.Options) string {
	arrayData := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	if array, ok := arrayData[key]; ok {
		arrayData[key] = append(array.([]string), value)
	} else {
		arrayData[key] = []string{value}
	}
	return value
}

func (t templateHelpers) getArray(key string, options *raymond.Options) []string {
	arrayData := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	if array, ok := arrayData[key]; ok {
		return array.([]string)
	} else {
		return []string{}
	}
}

func sumNumbers(numbers []string, format string) string {
	var sum float64 = 0
	for _, number := range numbers {
		value, err := strconv.ParseFloat(number, 64)
		if err != nil {
			log.Error(err)
			return "NaN"
		}
		sum += value
	}

	return formatNumber(sum, format)
}

func formatNumber(number float64, format string) string {
	if format == "" {
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", number), "0"), ".")
	}

	decimalPlaces := 0
	parts := strings.Split(format, ".")
	if len(parts) == 2 {
		decimalPlaces = len(parts[1])
	}

	multiplier := math.Pow(10, float64(decimalPlaces))
	rounded := math.Round(number*multiplier) / multiplier
	return fmt.Sprintf("%."+strconv.Itoa(decimalPlaces)+"f", rounded)
}

func getIndexEntry(journalIndexDetails Journal, indexName, indexValue string) (*JournalEntry, error) {

	for _, index := range journalIndexDetails.indexes {
		if index.name == indexName {
			if journalEntry, exists := index.entries[indexValue]; exists {
				return &journalEntry, nil
			}
		}
	}
	return nil, fmt.Errorf("no entry found for index %s", indexName)
}

func getBodyDataToParse(source string, journalEntry *JournalEntry) string {

	if strings.EqualFold(source, "request") {
		return journalEntry.requestBody
	}
	if strings.EqualFold(source, "response") {
		return journalEntry.responseBody
	}
	return ""
}

func getSearchFieldValue(options *raymond.Options, value string) string {

	if tpl, err := raymond.Parse("{{ " + value + " }}"); err == nil {
		if returnValue, err := tpl.Exec(options.Ctx()); err == nil && returnValue != "" {
			return returnValue
		}
	}
	return value
}

func getEvaluationString(helperName string, options *raymond.Options) string {

	evaluationString := "{{ " + helperName + " "
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
