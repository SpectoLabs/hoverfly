package templating

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SpectoLabs/hoverfly/core/journal"

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
	journal            *journal.Journal
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

// Concatenates any number of arguments together as strings.
func (t templateHelpers) concat(args ...interface{}) string {
	var parts []string
	for _, arg := range args {
		// If arg is a slice, flatten it
		if s, ok := arg.([]interface{}); ok {
			for _, v := range s {
				parts = append(parts, fmt.Sprint(v))
			}
		} else {
			parts = append(parts, fmt.Sprint(arg))
		}
	}
	return strings.Join(parts, "")
}

func (t templateHelpers) isNumeric(stringToCheck string) bool {
	_, err := strconv.ParseFloat(stringToCheck, 64)
	//return fmt.Sprintf("%t", err == nil)
	return err == nil
}

func (t templateHelpers) isAlphanumeric(s string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regex.MatchString(s)
}

func (t templateHelpers) isBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

// Comparison function type that takes two float64 values and returns a boolean.
type comparator func(float64, float64) bool

// Generic comparison function that parses strings to floats and applies the given comparison function.
func compareValues(value1, value2 string, comp comparator) bool {
	num1, err := strconv.ParseFloat(value1, 64)
	if err != nil {
		return false
	}
	num2, err := strconv.ParseFloat(value2, 64)
	if err != nil {
		return false
	}
	return comp(num1, num2)
}

// Specific comparison functions using the generic compareValues function.
func isGreaterThan(value1, value2 string) bool {
	return compareValues(value1, value2, func(a, b float64) bool { return a > b })
}

func isGreaterThanOrEqual(value1, value2 string) bool {
	return compareValues(value1, value2, func(a, b float64) bool { return a >= b })
}

func isLessThan(value1, value2 string) bool {
	return compareValues(value1, value2, func(a, b float64) bool { return a < b })
}

func isLessThanOrEqual(value1, value2 string) bool {
	return compareValues(value1, value2, func(a, b float64) bool { return a <= b })
}

func (t templateHelpers) isGreaterThan(value1, value2 string) bool {
	return isGreaterThan(value1, value2)
}

func (t templateHelpers) isGreaterThanOrEqual(value1, value2 string) bool {
	return isGreaterThanOrEqual(value1, value2)
}

func (t templateHelpers) isLessThan(value1, value2 string) bool {
	return isLessThan(value1, value2)
}

func (t templateHelpers) isLessThanOrEqual(value1, value2 string) bool {
	return isLessThanOrEqual(value1, value2)
}

func (t templateHelpers) isBetween(value, min, max string) bool {
	return t.isGreaterThan(value, min) && t.isLessThan(value, max)
}

func (t templateHelpers) matchesRegex(valueToCheck, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(valueToCheck)
}

func (t templateHelpers) length(stringToCheck string) string {
	return strconv.Itoa(len(stringToCheck))
}

func (t templateHelpers) substring(str, startStr, endStr string) string {
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return ""
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		return ""
	}
	if start < 0 || end > len(str) || start > end {
		return ""
	}
	return str[start:end]
}

func (t templateHelpers) rightmostCharacters(str, countStr string) string {
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return ""
	}
	if count < 0 || count > len(str) {
		return ""
	}
	return str[len(str)-count:]
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

func (t templateHelpers) fetchSingleFieldCsv(dataSourceName, searchFieldName, searchFieldValue, returnFieldName string, options *raymond.Options) string {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if !exists {
		log.Error("could not find datasource " + dataSourceName)
		return getEvaluationString("csv", options)
	}
	source.mu.Lock()
	defer source.mu.Unlock()
	searchIndex, err := getHeaderIndex(source.Data, searchFieldName)
	if err != nil {
		log.Error(err)
		return getEvaluationString("csv", options)
	}
	returnIndex, err := getHeaderIndex(source.Data, returnFieldName)
	if err != nil {
		log.Error(err)
		return getEvaluationString("csv", options)
	}
	searchValue := getSearchFieldValue(options, searchFieldValue)
	var fallbackString string
	for i := 1; i < len(source.Data); i++ {
		record := source.Data[i]
		if strings.ToLower(record[searchIndex]) == strings.ToLower(searchValue) {
			return record[returnIndex]
		} else if record[searchIndex] == "*" {
			fallbackString = record[returnIndex]
		}
	}
	if fallbackString != "" {
		return fallbackString
	}
	return getEvaluationString("csv", options)
}

func (t templateHelpers) fetchMatchingRowsCsv(dataSourceName string, searchFieldName string, searchFieldValue string) []RowMap {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if !exists {
		log.Error("could not find datasource " + dataSourceName)
		return []RowMap{}
	}
	if len(source.Data) < 1 {
		log.Error("no data available in datasource " + dataSourceName)
		return []RowMap{}
	}
	source.mu.Lock()
	defer source.mu.Unlock()

	headers := source.Data[0]
	fieldIndex := -1
	for i, header := range headers {
		if header == searchFieldName {
			fieldIndex = i
			break
		}
	}
	if fieldIndex == -1 {
		log.Error("could not find search field name " + searchFieldName)
		return []RowMap{}
	}

	var result []RowMap
	for _, row := range source.Data[1:] {
		if fieldIndex < len(row) && row[fieldIndex] == searchFieldValue {
			rowMap := make(RowMap)
			for i, cell := range row {
				if i < len(headers) {
					rowMap[headers[i]] = cell
				}
			}
			result = append(result, rowMap)
		}
	}
	return result
}

func (t templateHelpers) csvAsArray(dataSourceName string) [][]string {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if exists {
		source.mu.Lock()
		defer source.mu.Unlock()
		return source.Data
	} else {
		log.Error("could not find datasource " + dataSourceName)
		return [][]string{}
	}
}

func (t templateHelpers) csvAsMap(dataSourceName string) []RowMap {

	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if !exists {
		log.Error("could not find datasource " + dataSourceName)
		return []RowMap{}
	}
	source.mu.Lock()
	defer source.mu.Unlock()
	if len(source.Data) < 1 {
		log.Error("no data available in datasource " + dataSourceName)
		return []RowMap{}
	}
	headers := source.Data[0]
	var result []RowMap
	for _, row := range source.Data[1:] {
		rowMap := make(RowMap)
		for i, cell := range row {
			if i < len(headers) {
				rowMap[headers[i]] = cell
			}
		}
		result = append(result, rowMap)
	}
	return result
}

func (t templateHelpers) csvAddRow(dataSourceName string, newRow []string) string {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if exists {
		source.mu.Lock()
		defer source.mu.Unlock()
		source.Data = append(source.Data, newRow)
	} else {
		log.Error("could not find datasource " + dataSourceName)
	}
	return ""
}

func (t templateHelpers) csvDeleteRows(dataSourceName, searchFieldName, searchFieldValue string, output bool) string {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if !exists {
		log.Error("could not find datasource " + dataSourceName)
		return ""
	}
	source.mu.Lock()
	defer source.mu.Unlock()
	if len(source.Data) == 0 {
		log.Error("datasource " + dataSourceName + " is empty")
		return ""
	}
	header := source.Data[0]
	fieldIndex := -1
	for i, fieldName := range header {
		if fieldName == searchFieldName {
			fieldIndex = i
			break
		}
	}
	if fieldIndex == -1 {
		log.Error("could not find field name " + searchFieldName + " in datasource " + dataSourceName)
		return ""
	}
	filteredData := [][]string{header}
	rowsDeleted := 0
	for _, row := range source.Data[1:] {
		if row[fieldIndex] != searchFieldValue {
			filteredData = append(filteredData, row)
		} else {
			rowsDeleted++
		}
	}
	source.Data = filteredData
	if output {
		return fmt.Sprintf("%d", rowsDeleted)
	}
	return ""
}

func (t templateHelpers) csvCountRows(dataSourceName string) string {
	source, exists := t.TemplateDataSource.GetDataSource(dataSourceName)
	if !exists {
		log.Error("could not find datasource " + dataSourceName)
		return ""
	}
	source.mu.Lock()
	defer source.mu.Unlock()
	if len(source.Data) == 0 {
		return "0"
	}
	numRows := len(source.Data) - 1 // The number of rows is len(source.Data) - 1 (subtracting 1 for the header row)
	return fmt.Sprintf("%d", numRows)
}

func (t templateHelpers) csvSqlCommand(commandString string) []RowMap {

	// Parse the command string to get the Sql command
	command, err := parseSqlCommand(strings.TrimSpace(commandString), t.TemplateDataSource)
	if err != nil {
		log.Error("Error parsing sql command:", err)
		return []RowMap{}
	}

	// Find the data source by name
	source, exists := t.TemplateDataSource.GetDataSource(command.DataSourceName)
	if !exists {
		log.Error("Could not find datasource " + command.DataSourceName)
		return []RowMap{}
	}

	source.mu.Lock()
	defer source.mu.Unlock()

	var results []RowMap

	switch command.Type {
	case "SELECT":
		results = executeSqlSelectQuery(&source.Data, command)
	case "UPDATE":
		rowsAffected := executeSqlUpdateCommand(&source.Data, command)
		log.Debug(strconv.Itoa(rowsAffected) + " rows affected by " + commandString)
		return nil
	case "DELETE":
		rowsAffected := executeSqlDeleteCommand(&source.Data, command)
		log.Debug(strconv.Itoa(rowsAffected) + " rows affected by " + commandString)
		return nil
	default:
		log.Error(fmt.Errorf("unsupported query type %s", command.Type))
		return nil
	}

	return results
}

func (t templateHelpers) parseJournalBasedOnIndex(indexName, keyValue, dataSource, queryType, lookupQuery string, options *raymond.Options) interface{} {
	if journalEntry, err := getIndexEntry(t.journal, indexName, keyValue); err == nil {
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

func (t templateHelpers) hasJournalKey(indexName, keyValue string) bool {

	journalEntry, _ := getIndexEntry(t.journal, indexName, keyValue)

	return journalEntry != nil
}

func (t templateHelpers) setStatusCode(statusCode string, options *raymond.Options) string {
	intStatusCode, err := strconv.Atoi(statusCode)
	if err != nil {
		log.Error("status code is not a valid integer")
		return ""
	}

	if intStatusCode < 100 || intStatusCode > 599 {
		log.Error("status code is not valid")
		return ""
	}

	internalVars := options.ValueFromAllCtx("InternalVars").(map[string]interface{})
	internalVars["statusCode"] = intStatusCode
	return ""
}

func (t templateHelpers) setHeader(headerName string, headerValue string, options *raymond.Options) string {
	if headerName == "" {
		log.Error("header name cannot be empty")
		return ""
	}
	internalVars := options.ValueFromAllCtx("InternalVars").(map[string]interface{})
	var headers map[string][]string
	if h, ok := internalVars["setHeaders"]; ok {
		headers = h.(map[string][]string)
	} else {
		headers = make(map[string][]string)
	}
	// Replace or add the header
	headers[headerName] = []string{headerValue}
	internalVars["setHeaders"] = headers
	return ""
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

func (t templateHelpers) addToArray(key string, value string, output bool, options *raymond.Options) string {
	arrayData := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	if array, ok := arrayData[key]; ok {
		arrayData[key] = append(array.([]string), value)
	} else {
		arrayData[key] = []string{value}
	}

	if output {
		return value
	} else {
		return ""
	}
}

// Initializes (clears) an array in the template context under the given key.
func (t templateHelpers) initArray(key string, options *raymond.Options) string {
	arrayData := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	arrayData[key] = []string{}
	return ""
}

func (t templateHelpers) getArray(key string, options *raymond.Options) []string {
	arrayData := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	if array, ok := arrayData[key]; ok {
		return array.([]string)
	} else {
		return []string{}
	}
}

func (t templateHelpers) putValue(key string, value string, output bool, options *raymond.Options) string {
	kvs := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	kvs[key] = value
	if output {
		return value
	} else {
		return ""
	}
}

func (t templateHelpers) getValue(key string, options *raymond.Options) string {
	kvs := options.ValueFromAllCtx("Kvs").(map[string]interface{})
	value, exits := kvs[key]

	if exits {
		return value.(string)
	} else {
		return ""
	}
}

// jsonFromJWT extracts data from a JWT using a JSONPath query.
// Returns string, []interface{} (for arrays), or "" on errors/not found.
func (t templateHelpers) jsonFromJWT(path string, token string) interface{} {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	low := strings.ToLower(token)
	if strings.HasPrefix(low, "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Error("invalid jwt token (segment count) for jsonFromJWT")
		return ""
	}

 	composite := make(map[string]interface{})
    if h, ok := decodeJWTSegment(parts[0]); ok {
        composite["header"] = h
    }
    if p, ok := decodeJWTSegment(parts[1]); ok {
        composite["payload"] = p
    }

	jsonBytes, err := json.Marshal(composite)
	if err != nil {
		log.Error("error marshaling jwt composite: ", err)
		return ""
	}

	result := util.FetchFromRequestBody("jsonpath", path, string(jsonBytes))
	switch v := result.(type) {
	case []interface{}:
		return v
	case string:
		if v == "" {
			return ""
		}
		return v
	default:
		return ""
	}
}

// decodeJWTSegment decodes a single JWT segment (header or payload) into a generic map/array.
// It returns the decoded JSON value and a boolean indicating success.
func decodeJWTSegment(seg string) (interface{}, bool) {
    b, err := base64.RawURLEncoding.DecodeString(seg)
    if err != nil {
        log.Error("error decoding jwt segment: ", err)
        return nil, false
    }
    var v interface{}
    if err := json.Unmarshal(b, &v); err != nil {
        log.Error("error unmarshalling jwt segment: ", err)
        return nil, false
    }
    return v, true
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

func getIndexEntry(journal *journal.Journal, indexName, indexValue string) (*journal.JournalEntry, error) {

	for _, index := range journal.Indexes {
		if index.Name == indexName {
			if journalEntry, exists := index.Entries[indexValue]; exists {
				return journalEntry, nil
			}
		}
	}
	return nil, fmt.Errorf("no entry found for index %s", indexName)
}

func getBodyDataToParse(source string, journalEntry *journal.JournalEntry) string {

	if strings.EqualFold(source, "request") {
		return journalEntry.Request.Body
	}
	if strings.EqualFold(source, "response") {
		return journalEntry.Response.Body
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
		evaluationString = evaluationString + fmt.Sprint(params) + ` `
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
