package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/util/jsonpath"
)

func JsonPathMatch(matchingString string, toMatch string) bool {
	matchingString = prepareJsonPathQuery(matchingString)

	jsonPath := jsonpath.New("")

	err := jsonPath.Parse(matchingString)
	if err != nil {
		log.Errorf("Failed to parse json path query %s: %s", matchingString, err.Error())
		return false
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(toMatch), &data); err != nil {
		log.Errorf("Failed to unmarshal body to JSON: %s", err.Error())
		return false
	}

	buf := new(bytes.Buffer)

	err = jsonPath.Execute(buf, data)
	if err != nil {
		log.Errorf("Failed to execute json path match: %s", err.Error())
		return false
	}

	returnedString := buf.String()
	if returnedString == matchingString {
		return false
	}

	return true
}

func prepareJsonPathQuery(query string) string {
	if string(query[0:1]) != "{" && string(query[len(query)-1:]) != "}" {
		query = fmt.Sprintf("{%s}", query)
	}

	return query
}
