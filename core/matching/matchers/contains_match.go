package matchers

import (
	"reflect"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var Contains = "contains"

func ContainsMatch(match interface{}, toMatch string) bool {
	val := reflect.ValueOf(match)
	if val.Kind() != reflect.Slice {
		return false
	}
	var matchStringArr []string
	for i := 0; i < val.Len(); i++ {
		currentValue := val.Index(i)
		if currentValue.Kind() == reflect.Interface {
			matchStringArr = append(matchStringArr, currentValue.Elem().String())
		} else {
			matchStringArr = append(matchStringArr, currentValue.String())
		}

	}
	toMatchArr := strings.Split(toMatch, ";")
	return util.Contains(matchStringArr, toMatchArr)
}
