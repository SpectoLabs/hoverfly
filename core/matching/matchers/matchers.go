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
		MatcherFunction:     ExactMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Exact: {
		MatcherFunction:     ExactMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Glob: {
		MatcherFunction:     GlobMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Json: {
		MatcherFunction:     JsonMatch,
		MatchValueGenerator: IdentityValueGenerator,
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
		MatcherFunction:     XmlMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Xpath: {
		MatcherFunction:     XpathMatch,
		MatchValueGenerator: XPathMatchValueGenerator,
	},
	XmlTemplated: {
		MatcherFunction:     XmlTemplatedMatch,
		MatchValueGenerator: IdentityValueGenerator,
	},
	Array: {
		MatcherFunction:     ArrayMatchWithoutConfig,
		MatchValueGenerator: IdentityValueGenerator,
	},
	JWT: {
		MatcherFunction:     JwtMatcher,
		MatchValueGenerator: JwtMatchValueGenerator,
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
