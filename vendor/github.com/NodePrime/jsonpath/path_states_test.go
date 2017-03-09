package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pathTests = []lexTest{
	{"simple root node", `$.akey`, []int{pathRoot, pathPeriod, pathKey, pathEOF}},
	{"simple current node", `@.akey`, []int{pathCurrent, pathPeriod, pathKey, pathEOF}},
	{"simple root node w/ value", `$.akey+`, []int{pathRoot, pathPeriod, pathKey, pathValue, pathEOF}},
	{"nested object", `$.akey.akey2`, []int{pathRoot, pathPeriod, pathKey, pathPeriod, pathKey, pathEOF}},
	{"nested objects", `$.akey.akey2.akey3`, []int{pathRoot, pathPeriod, pathKey, pathPeriod, pathKey, pathPeriod, pathKey, pathEOF}},
	{"quoted keys", `$.akey["akey2"].akey3`, []int{pathRoot, pathPeriod, pathKey, pathBracketLeft, pathKey, pathBracketRight, pathPeriod, pathKey, pathEOF}},
	{"wildcard key", `$.akey.*.akey3`, []int{pathRoot, pathPeriod, pathKey, pathPeriod, pathWildcard, pathPeriod, pathKey, pathEOF}},
	{"wildcard index", `$.akey[*]`, []int{pathRoot, pathPeriod, pathKey, pathBracketLeft, pathWildcard, pathBracketRight, pathEOF}},
	{"key with where expression", `$.akey?(@.ten = 5)`, []int{pathRoot, pathPeriod, pathKey, pathWhere, pathExpression, pathEOF}},
	{"bracket notation", `$["aKey"][*][32][23:42]`, []int{pathRoot, pathBracketLeft, pathKey, pathBracketRight, pathBracketLeft, pathWildcard, pathBracketRight, pathBracketLeft, pathIndex, pathBracketRight, pathBracketLeft, pathIndex, pathIndexRange, pathIndex, pathBracketRight, pathEOF}},
}

func TestValidPaths(t *testing.T) {
	as := assert.New(t)
	for _, test := range pathTests {
		lexer := NewSliceLexer([]byte(test.input), PATH)
		types := itemsToTypes(readerToArray(lexer))

		as.EqualValues(types, test.tokenTypes, "Testing of %s: \nactual\n\t%+v\nexpected\n\t%v", test.name, typesDescription(types, pathTokenNames), typesDescription(test.tokenTypes, pathTokenNames))
	}
}
