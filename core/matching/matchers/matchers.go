package matchers

type MatcherFunc func(data interface{}, toMatch string) bool

var Matchers = map[string]MatcherFunc{
	// Default matcher
	"": ExactMatch,

	Exact:    ExactMatch,
	Glob:     GlobMatch,
	Json:     JsonMatch,
	JsonPath: JsonPathMatch,
	Regex:    RegexMatch,
	Xml:      XmlMatch,
	Xpath:    XpathMatch,
}
