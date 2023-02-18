package templating

import (
	"fmt"
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
}

type Request struct {
	QueryParam map[string][]string
	Header     map[string][]string
	Path       []string
	Scheme     string
	Body       func(queryType, query string, options *raymond.Options) string
	FormData   map[string][]string
	body       string
	Method     string
}

type Templator struct {
	SupportedMethodMap map[string]interface{}
}

var helpersRegistered = false

func NewTemplator() *Templator {
	t := templateHelpers{
		now:         time.Now,
		fakerSource: gofakeit.New(0),
	}
	helperMethodMap := make(map[string]interface{})
	if !helpersRegistered {
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
		helperMethodMap["faker"] = t.faker
		helperMethodMap["requestBody"] = t.requestBody

		raymond.RegisterHelpers(helperMethodMap)
		helpersRegistered = true
	}

	return &Templator{
		SupportedMethodMap: helperMethodMap,
	}
}

func (*Templator) ParseTemplate(responseBody string) (*raymond.Template, error) {

	return raymond.Parse(responseBody)
}

func (t *Templator) RenderTemplate(tpl *raymond.Template, requestDetails *models.RequestDetails, literals *models.Literals, vars *models.Variables, state map[string]string) (string, error) {
	if tpl == nil {
		return "", fmt.Errorf("template cannot be nil")
	}

	ctx := t.NewTemplatingData(requestDetails, literals, vars, state)
	return tpl.Exec(ctx)
}

func (templator *Templator) GetSupportedMethodMap() map[string]interface{} {
	return templator.SupportedMethodMap
}

func (t *Templator) NewTemplatingData(requestDetails *models.RequestDetails, literals *models.Literals, vars *models.Variables, state map[string]string) *TemplatingData {

	literalMap := make(map[string]interface{})
	if literals != nil {
		for _, literal := range *literals {
			literalMap[literal.Name] = literal.Value
		}
	}

	variableMap := t.getVariables(vars, requestDetails)

	return &TemplatingData{
		Request: Request{
			Path:       strings.Split(requestDetails.Path, "/")[1:],
			QueryParam: requestDetails.Query,
			Header:     requestDetails.Headers,
			Scheme:     requestDetails.Scheme,
			Body:       templateHelpers{}.requestBody,
			FormData:   requestDetails.FormData,
			body:       requestDetails.Body,
			Method:     requestDetails.Method,
		},
		Literals: literalMap,
		Vars:     variableMap,
		State:    state,
		CurrentDateTime: func(a1, a2, a3 string) string {
			return a1 + " " + a2 + " " + a3
		},
	}

}

func (t *Templator) getVariables(vars *models.Variables, requestDetails *models.RequestDetails) map[string]interface{} {
	variableMap := make(map[string]interface{})
	if vars != nil {
		for _, variable := range *vars {
			if variable.Function == REQUEST_BODY_HELPER {
				variableMap[variable.Name] = getDataFromRequestBody(variable, requestDetails.Body)
			} else {
				variableMap[variable.Name] = t.callHelper(variable)
			}
		}
	}

	return variableMap
}

func getDataFromRequestBody(variable models.Variable, body string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Error("panic occurred:", err)
		}
	}()
	return fetchFromRequestBody(variable.Arguments[0].(string), variable.Arguments[1].(string), body)
}

func (t *Templator) callHelper(variable models.Variable) interface{} {

	defer func() {
		if rec := recover(); rec != nil {
			log.Error("panic occurred:", rec)
		}
	}()
	function := reflect.ValueOf(t.SupportedMethodMap[variable.Function])
	functionType := function.Type()
	arguments := make([]reflect.Value, functionType.NumIn())
	for i := range arguments {
		// validate the type of argument - as of now just passing string, int, float, so just converted those
		if functionType.In(i).Kind() == reflect.String {
			arguments[i] = reflect.ValueOf(variable.Arguments[i].(string))
		} else if functionType.In(i).Kind() == reflect.Int {
			arguments[i] = reflect.ValueOf(int(variable.Arguments[i].(float64)))
		} else if functionType.In(i).Kind() == reflect.Float64 {
			arguments[i] = reflect.ValueOf(variable.Arguments[i].(float64))
		}
	}
	return function.Call(arguments)[0]
}
