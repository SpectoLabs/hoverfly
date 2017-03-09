package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonTests = []lexTest{
	{"empty object", `{}`, []int{jsonBraceLeft, jsonBraceRight, jsonEOF}},
	{"empty array", `[]`, []int{jsonBracketLeft, jsonBracketRight, jsonEOF}},
	{"key string", `{"key" :"value"}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonString, jsonBraceRight, jsonEOF}},
	{"multiple pairs", `{"key" :"value","key2" :"value"}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonString, jsonComma, jsonKey, jsonColon, jsonString, jsonBraceRight, jsonEOF}},
	{"key number", `{"key" : 12.34e+56}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonNumber, jsonBraceRight, jsonEOF}},
	{"key true", `{"key" :true}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBool, jsonBraceRight, jsonEOF}},
	{"key false", `{"key" :false}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBool, jsonBraceRight, jsonEOF}},
	{"key null", `{"key" :null}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonNull, jsonBraceRight, jsonEOF}},
	{"key arrayOf number", `{"key" :[23]}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonNumber, jsonBracketRight, jsonBraceRight, jsonEOF}},
	{"key array", `{"key" :[23,"45",67]}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonNumber, jsonComma, jsonString, jsonComma, jsonNumber, jsonBracketRight, jsonBraceRight, jsonEOF}},
	{"key array", `{"key" :["45",{}]}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonString, jsonComma, jsonBraceLeft, jsonBraceRight, jsonBracketRight, jsonBraceRight, jsonEOF}},
	{"key nestedObject", `{"key" :{"innerkey":"value"}}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBraceLeft, jsonKey, jsonColon, jsonString, jsonBraceRight, jsonBraceRight, jsonEOF}},
	{"key nestedArray", `[1,["a","b"]]`, []int{jsonBracketLeft, jsonNumber, jsonComma, jsonBracketLeft, jsonString, jsonComma, jsonString, jsonBracketRight, jsonBracketRight, jsonEOF}},
}

func TestValidJson(t *testing.T) {
	as := assert.New(t)

	for _, test := range jsonTests {
		lexer := NewSliceLexer([]byte(test.input), JSON)
		types := itemsToTypes(readerToArray(lexer))

		as.EqualValues(types, test.tokenTypes, "Testing of %q: \nactual\n\t%+v\nexpected\n\t%v", test.name, typesDescription(types, jsonTokenNames), typesDescription(test.tokenTypes, jsonTokenNames))
	}
}

var errorJsonTests = []lexTest{
	{"Missing end brace", `{`, []int{jsonBraceLeft, jsonError}},
	{"Missing start brace", `}`, []int{jsonError}},
	{"Missing key start quote", `{key":true}`, []int{jsonBraceLeft, jsonError}},
	{"Missing key end quote", `{"key:true}`, []int{jsonBraceLeft, jsonError}},
	{"Missing colon", `{"key"true}`, []int{jsonBraceLeft, jsonKey, jsonError}},
	{"Missing value", `{"key":}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonError}},
	{"Missing string start quote", `{"key":test"}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonError}},
	{"Missing embedded array bracket", `{"key":[}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonError}},
	{"Missing values in array", `{"key":[,]`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonError}},
	{"Missing value after comma", `{"key":[343,]}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonNumber, jsonComma, jsonError}},
	{"Missing comma in array", `{"key":[234 424]}`, []int{jsonBraceLeft, jsonKey, jsonColon, jsonBracketLeft, jsonNumber, jsonError}},
}

func TestMalformedJson(t *testing.T) {
	as := assert.New(t)

	for _, test := range errorJsonTests {
		lexer := NewSliceLexer([]byte(test.input), JSON)
		types := itemsToTypes(readerToArray(lexer))

		as.EqualValues(types, test.tokenTypes, "Testing of %q: \nactual\n\t%+v\nexpected\n\t%v", test.name, typesDescription(types, jsonTokenNames), typesDescription(test.tokenTypes, jsonTokenNames))
	}
}

func itemsToTypes(items []Item) []int {
	types := make([]int, len(items))
	for i, item := range items {
		types[i] = item.typ
	}
	return types
}

// func TestEarlyTerminationForJSON(t *testing.T) {
// 	as := assert.New(t)
// 	wg := sync.WaitGroup{}

// 	lexer := NewSliceLexer(`{"key":"value", "key2":{"ikey":3}, "key3":[1,2,3,4]}`)
// 	wg.Add(1)
// 	go func() {
// 		lexer.Run(JSON)
// 		wg.Done()
// 	}()

// 	// Pop a few items
// 	<-lexer.items
// 	<-lexer.items
// 	// Kill command
// 	close(lexer.kill)

// 	wg.Wait()
// 	remainingItems := readerToArray(lexer.items)
// 	// TODO: Occasionally fails - rethink this
// 	_ = as
// 	_ = remainingItems
// 	// as.True(len(remainingItems) <= bufferSize, "Count of remaining items should be less than buffer size: %d", len(remainingItems))
// }

var examples = []string{
	`{"items":[
	  {
	    "name": "example document for wicked fast parsing of huge json docs",
	    "integer": 123,
	    "totally sweet scientific notation": -123.123e-2,
	    "unicode? you betcha!": "ú™£¢∞§\u2665",
	    "zero character": "0",
	    "null is boring": null
	  },
	  {
	    "name": "another object",
	    "cooler than first object?": true,
	    "nested object": {
	      "nested object?": true,
	      "is nested array the same combination i have on my luggage?": true,
	      "nested array": [1,2,3,4,5]
	    },
	    "false": false
	  }
]}`,
	`{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "Product set",
    "type": "array",
    "items": {
        "title": "Product",
        "type": "object",
        "properties": {
            "id": {
                "description": "The unique identifier for a product",
                "type": "number"
            },
            "name": {
                "type": "string"
            },
            "price": {
                "type": "number",
                "minimum": 0,
                "exclusiveMinimum": true
            },
            "tags": {
                "type": "array",
                "items": {
                    "type": "string"
                },
                "minItems": 1,
                "uniqueItems": true
            },
            "dimensions": {
                "type": "object",
                "properties": {
                    "length": {"type": "number"},
                    "width": {"type": "number"},
                    "height": {"type": "number"}
                },
                "required": ["length", "width", "height"]
            },
            "warehouseLocation": {
                "description": "Coordinates of the warehouse with the product",
                "$ref": "http://json-schema.org/geo"
            }
        },
        "required": ["id", "name", "price"]
    }
   }`,
}

func TestMixedCaseJson(t *testing.T) {
	as := assert.New(t)
	for _, json := range examples {
		lexer := NewSliceLexer([]byte(json), JSON)
		items := readerToArray(lexer)

		for _, i := range items {
			as.False(i.typ == jsonError, "Found error while parsing: %q", i.val)
		}
	}
}
