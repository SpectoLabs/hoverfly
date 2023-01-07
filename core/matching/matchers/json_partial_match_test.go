package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonPartialMatch_MatchesTrueWithEqualsJSON(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	{
        "name": "Object 2",
        "set": false,
        "age": 400
    }`, `{
"objects": [
    {
        "name": "Object 1",
        "set": true
    },{
        "name": "Object 2",
        "set": false,
        "age": 400
    }]
}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithNotOrderedJSON(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPartialMatch(`{"test":{"minified":true,"json":true}}`, `{"test":{"json":true,"minified":true}}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{"test":{"json":true,"minified":true}}`))
}

func Test_JsonPartialMatch_MatchesTrueWithAbsentNode(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPartialMatch(`{"test":{"minified":true}}`, `{"test":{"json":true,"minified":true}}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{"test":{"json":true,"minified":true}}`))
}

func Test_JsonPartialMatch_MatchesTrueWithAbsentObject(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPartialMatch(`{"test":{"minified":true}}`, `{"test":{"json":true,"minified":true,"someObject":{"fieldA":"valueA"}}}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{"test":{"json":true,"minified":true,"someObject":{"fieldA":"valueA"}}}`))
}

func Test_JsonPartialMatch_MatchesFalseWithAbsentNode(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"minified":true}}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesFalseWithAbsentObject(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true,"someObject":{"fieldA":"valueA"}}}`, `{"test":{"minified":true}}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesTrueEmptyJson(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPartialMatch(`{}`, `{}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{}`))
}

func Test_JsonPartialMatch_MatchesFalseInvalidJson(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":}}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesTrueDeep(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(
		`{
  "fieldA": "valueA"
}`,
		`{
	"test": {
		"json": true,
		"minified": true,
		"someObject": {
			"fieldA": "valueA"
		}
}}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueDeepArrayInside(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(
		`{
  "NAME": "79684881033",
  "REDIRECT_NUMBER": "79684881033"
}`,
		`{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "redirect_type": 1,
    "followme_struct": [
      3,
      [
        {
          "I_FOLLOW_ORDER": "1",
          "ACTIVE": "Y",
          "NAME": "79684881033",
          "REDIRECT_NUMBER": "79684881033",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        },
        {
          "I_FOLLOW_ORDER": "2",
          "ACTIVE": "Y",
          "NAME": "79684881034",
          "REDIRECT_NUMBER": "79684881034",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        }
      ]
    ]
  }
}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueDeepComplexWithArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(
		`{
    "redirect_type": 1,
    "followme_struct": [
      3,
      [
        {
          "I_FOLLOW_ORDER": "2",
          "ACTIVE": "Y",
          "NAME": "79684881034",
          "REDIRECT_NUMBER": "79684881034",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        }
      ]
    ]
}`,
		`{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "redirect_type": 1,
    "followme_struct": [
      3,
      [
        {
          "I_FOLLOW_ORDER": "1",
          "ACTIVE": "Y",
          "NAME": "79684881033",
          "REDIRECT_NUMBER": "79684881033",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        },
        {
          "I_FOLLOW_ORDER": "2",
          "ACTIVE": "Y",
          "NAME": "79684881034",
          "REDIRECT_NUMBER": "79684881034",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        }
      ]
    ]
  }
}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesFalseDeepComplexWithArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(
		`{
    "redirect_type": 1,
    "followme_struct": [
      3,
      [
        {
          "I_FOLLOW_ORDER": "2",
          "ACTIVE": "Y",
          "NAME": "WRONG_NAME",
          "REDIRECT_NUMBER": "79684881034",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        }
      ]
    ]
}`,
		`{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "redirect_type": 1,
    "followme_struct": [
      3,
      [
        {
          "I_FOLLOW_ORDER": "1",
          "ACTIVE": "Y",
          "NAME": "79684881033",
          "REDIRECT_NUMBER": "79684881033",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        },
        {
          "I_FOLLOW_ORDER": "2",
          "ACTIVE": "Y",
          "NAME": "79684881034",
          "REDIRECT_NUMBER": "79684881034",
          "PERIOD": "always",
          "PERIOD_DESCRIPTION": "always",
          "TIMEOUT": "15"
        }
      ]
    ]
  }
}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesTrueAgainstJSONRootAsArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	{
        "name": "Object 2",
        "set": false,
        "age": 400
    }`, `[{
	"objects": [
		{
			"name": "Object 1",
			"set": true
		},{
			"name": "Object 2",
			"set": false,
			"age": 400
		}]
	}]`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithJSONRootAsArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 1",
        "set": true
    },{
        "name": "Object 2",
        "set": false,
        "age": 400
	}]`, `
	{
		"objects": [
		{
			"name": "Object 1",
			"set": true
		},{
			"name": "Object 2",
			"set": false,
			"age": 400
		}]
	}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithJSONRootAsPartialArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 1",
        "set": true
    }]`, `
	{
		"objects": [
		{
			"name": "Object 1",
			"set": true
		},{
			"name": "Object 2",
			"set": false,
			"age": 400
		}]
	}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesTrueWithJSONRootAsPartialArrayWithPartialObject(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 2",
        "set": false
    }]`, `
	{
		"objects": [
		{
			"name": "Object 1",
			"set": true
		},{
			"name": "Object 2",
			"set": false,
			"age": 400
		}]
	}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesFalseWithJSONRootAsArrayWithDifferentElement(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 3",
        "set": true
    }]`, `
	{
		"objects": [
		{
			"name": "Object 1",
			"set": true
		},{
			"name": "Object 2",
			"set": false,
			"age": 400
		}]
	}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPartialMatch_MatchesTrueWithJSONRootAsArrayAgainstJSONRootAsArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 1",
        "set": true
    }]`,
		`[
	{
		"name": "Object 1",
		"set": true
	},{
		"name": "Object 2",
		"set": false,
		"age": 400
	}]`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPartialMatch_MatchesFalseWithJSONRootAsArrayAgainstJSONRootAsArrayWithDifferentElement(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPartialMatch(`
	[{
        "name": "Object 3",
        "set": false
    }]`,
		`[
	{
		"name": "Object 1",
		"set": true
	},{
		"name": "Object 2",
		"set": false,
		"age": 400
	}]`, nil)
	Expect(isMatched).To(BeFalse())
}
