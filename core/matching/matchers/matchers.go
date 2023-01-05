package matchers

type MatcherFunc func(data interface{}, toMatch string) bool

type MatcherFuncWithConfig func(data interface{}, toMatch string, config map[string]interface{}) bool

var Matchers = map[string]MatcherFunc{
	// Default matcher
	"": ExactMatch,

	Exact:           ExactMatch,
	Glob:            GlobMatch,
	Json:            JsonMatch,
	JsonPath:        JsonPathMatch,
	JsonPartial:     JsonPartialMatch,
	Regex:           RegexMatch,
	Xml:             XmlMatch,
	Xpath:           XpathMatch,
	XmlTemplated:    XmlTemplatedMatch,
	ContainsExactly: ContainsExactlyMatch,
	Array:           ContainsExactlyMatch,
}

var MatchersWithConfig = map[string]MatcherFuncWithConfig{
	Array: ArrayMatch,
}
