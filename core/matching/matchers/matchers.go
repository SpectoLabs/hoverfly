package matchers

type MatcherFunc func(data interface{}, toMatch string, config map[string]interface{}) (string, bool)

var Matchers = map[string]MatcherFunc{

	"":              ExactMatch,
	Exact:           ExactMatch,
	Glob:            GlobMatch,
	Regex:           RegexMatch,
	JsonPath:        JsonPathMatch,
	Xpath:           XpathMatch,
	Array:           ArrayMatch,
	ContainsExactly: ContainsExactlyMatch,
	Json:            JsonMatch,
	Xml:             XmlMatch,
	JsonPartial:     JsonPartialMatch,
	XmlTemplated:    XmlTemplatedMatch,
}
