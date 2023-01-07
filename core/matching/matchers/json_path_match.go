package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/jsonpath"
)

var JsonPath = "jsonpath"

func JsonPathMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchString, ok := match.(string)
	if !ok {
		return "", false
	}

	matchString = prepareJsonPathQuery(matchString)
	returnedString, err := jsonPathExecutionReturningJsonString(matchString, toMatch)
	if err != nil || returnedString == matchString {
		return "", false
	}

	return returnedString, true
}

func JsonPathExecution(matchString, toMatch string) (string, error) {
	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchString, err.Error())
		return "", err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return "", err
	}

	buf := new(bytes.Buffer)

	err = jsonPath.Execute(buf, data)
	if err != nil {
		log.Errorf("err to execute json path match: %s", err.Error())
		return "", err
	}

	return buf.String(), nil
}

func jsonPathExecutionReturningJsonString(matchString, toMatch string) (string, error) {
	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchString, err.Error())
		return "", err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return "", err
	}

	results, err := jsonPath.FindResults(data)
	if err != nil {
		log.Errorf("err to execute json path match: %s", err.Error())
		return "", err
	}

	return getResult(results)
}

func getResult(results [][]reflect.Value) (string, error) {
	var allResultsInterface []interface{}
	for _, result := range results {
		for _, singleResult := range result {
			allResultsInterface = append(allResultsInterface, singleResult.Interface())
		}
	}
	if len(allResultsInterface) == 1 {
		if _, ok := allResultsInterface[0].(string); ok {
			return allResultsInterface[0].(string), nil
		}
		bytes, err := json.Marshal(allResultsInterface[0])
		if err != nil {
			log.Errorf("err to marshal %s", err.Error())
			return "", err
		}
		return string(bytes), nil
	}

	bytes, err := json.Marshal(allResultsInterface)
	if err != nil {
		log.Errorf("err to marshal %s", err.Error())
		return "", err
	}
	return string(bytes), nil
}

func prepareJsonPathQuery(query string) string {
	if string(query[0:1]) != "{" && string(query[len(query)-1:]) != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}
