package matchers

import (
	"bytes"
	"encoding/json"
	"reflect"
)

var Json = "json"

func JsonMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	if matchString == toMatch {
		return true
	}
	var matchingObject interface{}
	d := json.NewDecoder(bytes.NewBuffer([]byte(matchString)))
	d.UseNumber()
	err := d.Decode(&matchingObject)
	if err != nil {
		return false
	}

	var toMatchObject interface{}
	d = json.NewDecoder(bytes.NewBuffer([]byte(toMatch)))
	d.UseNumber()
	err = d.Decode(&toMatchObject)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(matchingObject, toMatchObject)
}
