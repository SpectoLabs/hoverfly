 package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

var jsonMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithJSON",
		match:   `{"test":{"json":true,"minified":true}}`,
		toMatch: `{"test":{"json":true,"minified":true}}`,
		matched: true,
		result:  `{"test":{"json":true,"minified":true}}`,
	},
	{
		name:    "MatchesTrueWithJSON_InADifferentOrder",
		match:   `{"test":{"minified":true, "json":true}}`,
		toMatch: `{"test":{"json":true,"minified":true}}`,
		matched: true,
		result:  `{"test":{"json":true,"minified":true}}`,
	},
	{
		name:      "MatchesTrueWithUnminifiedJSON",
		match:     `{"test":{"json":true,"minified":true}}`,
		toMatch:   `{
		"test": {
			"json": true,
			"minified": true
		}
	}`,
		matched: true,
		result:    `{
		"test": {
			"json": true,
			"minified": true
		}
	}`,
	},
	{
		name:      "MatchesFalseReturnsEmptyString",
		match:     `{"test":{"json":true,"minified":true}}`,
		toMatch:   `{
		"test": {
			"json": false,
			"minified": false
		}
	}`,
		matched: false,
		result:  "",
	},
	{
		name:      "MatchesFalseWithInvalidJSONAsMatcher",
		match:     `"test":"json":true,"minified"`,
		toMatch:   `{
		"test": {
			"json": true,
			"minified": true
		}
	}`,
		matched: false,
		result:  "",
	},
	{
		name:      "MatchesFalseWithInvalidJSON",
		match:     `{"test":{"json":true,"minified":true}}`,
		toMatch:   `{
		"test": {
			"json": true,
			"minified": 
		}
	}`,
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithTwoEmptyString",
		match:   "",
		toMatch: "",
		matched: true,
		result:  "",
	},
	{
		name:      "MatchesFalseAgainstEmptyString",
		match:     `{
		"test": {
			"json": true,
			"minified": 
		}
	}`,
		toMatch: "",
		matched: false,
		result:  "",
	},
	{
		name:      "MatchesFalseWithEmptyString",
		match:     "",
		toMatch:   `{
		"test": {
			"json": true,
			"minified": 
		}
	}`,
		matched: false,
		result:  "",
	},
}

func Test_JsonMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range jsonMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.JsonMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}

