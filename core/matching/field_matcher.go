package matching

import (
	"bytes"
	"regexp"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	"github.com/NodePrime/jsonpath"
	"github.com/SpectoLabs/hoverfly/core/models"
	glob "github.com/ryanuber/go-glob"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) bool {
	if field == nil {
		return true
	}

	if field.ExactMatch != nil {
		return *field.ExactMatch == toMatch
	}

	if field.XpathMatch != nil {
		xpathRule, err := goxpath.Parse(*field.XpathMatch)
		if err != nil {
			return false
		}

		xTree, err := xmltree.ParseXML(bytes.NewBufferString(toMatch))
		if err != nil {
			return false
		}

		results, err := xpathRule.ExecNode(xTree)
		if err != nil {
			return false
		}

		return len(results) > 0
	}

	if field.JsonMatch != nil {
		paths, err := jsonpath.ParsePaths(*field.JsonMatch)
		if err != nil {
			return false
		}

		eval, err := jsonpath.EvalPathsInBytes([]byte(toMatch), paths)
		if err != nil {
			return false
		}

		_, ok := eval.Next()

		return ok
	}

	if field.RegexMatch != nil {
		match, err := regexp.MatchString(*field.RegexMatch, toMatch)
		if err != nil {
			return false
		}

		return match
	}

	if field.GlobMatch != nil {
		return glob.Glob(*field.GlobMatch, toMatch)
	}

	return false
}
