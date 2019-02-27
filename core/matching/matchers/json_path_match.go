package matchers

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/third_party/forked/golang/template"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/util/jsonpath"
)

var JsonPath = "jsonpath"

func JsonPathMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	matchString = prepareJsonPathQuery(matchString)
	returnedString, err := JsonPathExecution(matchString, toMatch)
	if err != nil {
		return
	}

	return true, returnedString
}

func JsonPathExecution(matchString, toMatch string) (string, error) {
	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchString, err.Error())
		return "", err
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(toMatch), &jsonData); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return "", err
	}
	
	fullResults, err := jsonPath.FindResults(jsonData)
	if err != nil {
		log.Warnf("Json path match for `%s` failed: %s", matchString, err.Error())
		return "", err
	}

	// TODO handle multiple matching elements
	var foundEl interface{}
	var ok bool
	for _, results := range fullResults {
		for _, r := range results {
			foundEl, ok = template.PrintableValue(r)
			if !ok {
				log.Errorf("Failed to get json path match results: %s", err.Error())
				return "", err
			}
		}
	}

	textEl, ok := foundEl.(string)
	if ok {
		return textEl, nil
	}

	found, err := json.Marshal(foundEl)
	if err != nil {
		log.Errorf("Failed to marshal json path match results: %s", err.Error())
		return "", err
	}
	return string(found), nil
}

func prepareJsonPathQuery(query string) string {
	if string(query[0:1]) != "{" && string(query[len(query)-1:]) != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}
