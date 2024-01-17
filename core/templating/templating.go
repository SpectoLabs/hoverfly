package templating

import (
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/util"
	"reflect"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"

	log "github.com/sirupsen/logrus"
)

const REQUEST_BODY_HELPER = "requestBody"

type TemplatingData struct {
	Request         Request
	State           map[string]string
	CurrentDateTime func(string, string, string) string
	Literals        map[string]interface{}
	Vars            map[string]interface{}
	Journal         Journal
}

type Request struct {
	QueryParam map[string][]string
	Header     map[string][]string
	Path       []string
	Scheme     string
	Body       func(queryType, query string, options *raymond.Options) interface{}
	FormData   map[string][]string
	body       string
	Method     string
	Host       string
}

type JournalEntry struct {
	requestBody  string
	responseBody string
}

type Journal struct {
	indexes []JournalIndex
}

type JournalIndex struct {
	name    string
	entries map[string]JournalEntry
}

type Templator struct {
	SupportedMethodMap map[string]interface{}
	TemplateHelper     templateHelpers
}

var helpersRegistered = false

func NewTemplator() *Templator {

	templateDataSource := NewTemplateDataSource()
	arrayData := make(map[string][]string)

	t := templateHelpers{
		now:                time.Now,
		fakerSource:        gofakeit.New(0),
		TemplateDataSource: templateDataSource,
		ArrayData:          arrayData,
	}
	helperMethodMap := make(map[string]interface{})
	helperMethodMap["now"] = t.nowHelper
	helperMethodMap["randomString"] = t.randomString
	helperMethodMap["randomStringLength"] = t.randomStringLength
	helperMethodMap["randomBoolean"] = t.randomBoolean
	helperMethodMap["randomInteger"] = t.randomInteger
	helperMethodMap["randomIntegerRange"] = t.randomIntegerRange
	helperMethodMap["randomFloat"] = t.randomFloat
	helperMethodMap["randomFloatRange"] = t.randomFloatRange
	helperMethodMap["randomEmail"] = t.randomEmail
	helperMethodMap["randomIPv4"] = t.randomIPv4
	helperMethodMap["randomIPv6"] = t.randomIPv6
	helperMethodMap["randomUuid"] = t.randomUuid
	helperMethodMap["replace"] = t.replace
	helperMethodMap["split"] = t.split
	helperMethodMap["faker"] = t.faker
	helperMethodMap["requestBody"] = t.requestBody
	helperMethodMap["csv"] = t.parseCsv
	helperMethodMap["journal"] = t.parseJournalBasedOnIndex
	helperMethodMap["sum"] = t.sum
	helperMethodMap["add"] = t.add
	helperMethodMap["subtract"] = t.subtract
	helperMethodMap["multiply"] = t.multiply
	helperMethodMap["divide"] = t.divide
	helperMethodMap["addToArray"] = t.addToArray
	helperMethodMap["getArray"] = t.getArray
	if !helpersRegistered {
		raymond.RegisterHelpers(helperMethodMap)
		helpersRegistered = true
	}

	return &Templator{
		SupportedMethodMap: helperMethodMap,
		TemplateHelper:     t,
	}
}

func (*Templator) ParseTemplate(responseBody string) (*raymond.Template, error) {

	return raymond.Parse(responseBody)
}

func (t *Templator) RenderTemplate(tpl *raymond.Template, requestDetails *models.RequestDetails, literals *models.Literals, vars *models.Variables, state map[string]string, journal *journal.Journal) (string, error) {
	if tpl == nil {
		return "", fmt.Errorf("template cannot be nil")
	}

	ctx := t.NewTemplatingData(requestDetails, literals, vars, state, journal)
	return tpl.Exec(ctx)
}

func (templator *Templator) GetSupportedMethodMap() map[string]interface{} {
	return templator.SupportedMethodMap
}

func (t *Templator) NewTemplatingData(requestDetails *models.RequestDetails, literals *models.Literals, vars *models.Variables, state map[string]string, journal *journal.Journal) *TemplatingData {

	literalMap := make(map[string]interface{})
	if literals != nil {
		for _, literal := range *literals {
			literalMap[literal.Name] = literal.Value
		}
	}

	variableMap := t.getVariables(vars, requestDetails)
	templateJournal := Journal{}
	if journal != nil {

		indexes := make([]JournalIndex, len(journal.Indexes))
		for _, index := range journal.Indexes {

			journalIndexEntries := make(map[string]JournalEntry)
			for indexKey, entry := range index.Entries {

				journalEntry := JournalEntry{
					requestBody:  entry.Request.Body,
					responseBody: entry.Response.Body,
				}
				journalIndexEntries[indexKey] = journalEntry
			}
			journalIndex := JournalIndex{
				name:    index.Name,
				entries: journalIndexEntries,
			}
			indexes = append(indexes, journalIndex)
		}
		templateJournal = Journal{
			indexes: indexes,
		}

	}

	return &TemplatingData{
		Request:  getRequest(requestDetails),
		Literals: literalMap,
		Vars:     variableMap,
		State:    state,
		Journal:  templateJournal,
		CurrentDateTime: func(a1, a2, a3 string) string {
			return a1 + " " + a2 + " " + a3
		},
	}

}

func getRequest(requestDetails *models.RequestDetails) Request {
	return Request{
		Path:       strings.Split(requestDetails.Path, "/")[1:],
		QueryParam: requestDetails.Query,
		Header:     requestDetails.Headers,
		Scheme:     requestDetails.Scheme,
		Body:       templateHelpers{}.requestBody,
		FormData:   requestDetails.FormData,
		body:       requestDetails.Body,
		Method:     requestDetails.Method,
		Host:       requestDetails.Destination,
	}
}

func (t *Templator) getVariables(vars *models.Variables, requestDetails *models.RequestDetails) map[string]interface{} {
	variableMap := make(map[string]interface{})
	if vars != nil {
		for _, variable := range *vars {
			if variable.Function == REQUEST_BODY_HELPER {
				variableMap[variable.Name] = getDataFromRequestBody(variable, requestDetails.Body)
			} else {
				resultValue := t.callHelper(variable, requestDetails)
				if resultValue != nil {
					variableMap[variable.Name] = resultValue.(reflect.Value).Interface()
				}
			}
		}
	}

	return variableMap
}

func getDataFromRequestBody(variable models.Variable, body string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Error("panic occurred:", err)
		}
	}()
	return util.FetchFromRequestBody(variable.Arguments[0].(string), variable.Arguments[1].(string), body)
}

/*
*
This method is basically invoking helper function via reflection and returning the value.
Disclaimer: we cannot use helper functions that are taking *raymond.Options as an argument
*/
func (t *Templator) callHelper(variable models.Variable, requestDetails *models.RequestDetails) interface{} {

	defer func() {
		if rec := recover(); rec != nil {
			log.Error("panic occurred:", rec)
		}
	}()
	function := reflect.ValueOf(t.SupportedMethodMap[variable.Function])
	functionType := function.Type()
	arguments := make([]reflect.Value, functionType.NumIn())
	for i := range arguments {
		if functionType.In(i).Kind() == reflect.String {
			arguments[i] = reflect.ValueOf(parseValidRequestTemplate(variable.Arguments[i].(string), requestDetails))
		} else if functionType.In(i).Kind() == reflect.Int {
			arguments[i] = reflect.ValueOf(int(variable.Arguments[i].(float64)))
		} else if functionType.In(i).Kind() == reflect.Float64 {
			arguments[i] = reflect.ValueOf(variable.Arguments[i].(float64))
		}
	}
	return function.Call(arguments)[0]
}

func parseValidRequestTemplate(source string, details *models.RequestDetails) string {

	if tpl, err := raymond.Parse("{{" + source + "}}"); err == nil {
		ctx := &TemplatingData{Request: getRequest(details)}
		if parsedValue, execErr := tpl.Exec(ctx); execErr == nil {
			return parsedValue
		}
	}
	return source
}
