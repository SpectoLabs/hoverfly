package matchers

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/jsonpath"
)

func IdentityValueGenerator(match interface{}, toMatch string) string {

	return toMatch
}

func JsonPathMatcherValueGenerator(match interface{}, toMatch string) string {

	matchString := prepareJsonPathQuery(match.(string))

	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchString, err.Error())
		return ""
	}

	var data interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return ""
	}

	results, err := jsonPath.FindResults(data)
	if err != nil {
		log.Errorf("err to execute json path match: %s", err.Error())
		return ""
	}

	return getResult(results)
}

func getResult(results [][]reflect.Value) string {
	var allResultsInterface []interface{}
	for _, result := range results {
		for _, singleResult := range result {
			allResultsInterface = append(allResultsInterface, singleResult.Interface())
		}
	}
	if len(allResultsInterface) == 1 {
		if _, ok := allResultsInterface[0].(string); ok {
			return allResultsInterface[0].(string)
		}
		bytes, err := json.Marshal(allResultsInterface[0])
		if err != nil {
			log.Errorf("err to marshal %s", err.Error())
			return ""
		}
		return string(bytes)
	}

	bytes, err := json.Marshal(allResultsInterface)
	if err != nil {
		log.Errorf("err to marshal %s", err.Error())
		return ""
	}
	return string(bytes)
}

func XPathMatchValueGenerator(match interface{}, toMatch string) string {

	results, err := XpathExecution(match.(string), toMatch)
	if err != nil {
		log.Errorf("Failed to generate xpath value: %s", err.Error())
		return ""
	}
	return results.String()
}

func JwtMatchValueGenerator(match interface{}, toMatch string) string {

	if jwt, err := ParseJWT(toMatch); err == nil {
		return jwt
	} else {
		log.Errorf("Failed to parse JWT %s", err.Error())
		return ""
	}
}
