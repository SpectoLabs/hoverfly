package matchers

type MatcherFunc func(data interface{}, toMatch string) bool

type MatcherFuncWithConfig func(data interface{}, toMatch string, config map[string]interface{}) bool

/*
	 this is called only if matcher returns true and there is chaining to feed value
		this is set as nil for matchers where we are doing complete details match and there is no need of chaining
*/
type MatcherValueGenerator func(data interface{}, toMatch string) string

var Matchers = map[string]MatcherDetails{
	// Default matcher
	"": {
		MatcherFunction: ExactMatch,
	},
	Exact: {
		MatcherFunction: ExactMatch,
	},
	Glob: {
		MatcherFunction:     GlobMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Json: {
		MatcherFunction: JsonMatch,
	},
	JsonPath: {
		MatcherFunction:     JsonPathMatch,
		MatchValueGenerator: JsonPathMatcherValueGenerator,
	},
	JsonPartial: {
		MatcherFunction:     JsonPartialMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Regex: {
		MatcherFunction:     RegexMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Xml: {
		MatcherFunction: XmlMatch,
	},
	Xpath: {
		MatcherFunction:     XpathMatch,
		MatchValueGenerator: XPathMatchValueGenerator,
	},
	XmlTemplated: {
		MatcherFunction:     XmlTemplatedMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	ContainsExactly: {
		MatcherFunction: ContainsExactlyMatch,
	},
	Array: {
		MatcherFunction:     ContainsExactlyMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
}

type MatcherDetails struct {
	MatcherFunction     interface{}
	MatchValueGenerator MatcherValueGenerator
}

var MatchersWithConfig = map[string]MatcherDetails{
	Array: {
		MatcherFunction:     ArrayMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
}
