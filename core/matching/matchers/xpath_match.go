package matchers

import (
	"bytes"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	log "github.com/sirupsen/logrus"
)

var Xpath = "xpath"

func XpathMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	results, err := XpathExecution(matchString, toMatch)
	if err != nil {
		return false
	}

	return len(results) > 0
}

func XpathExecution(matchString, toMatch string) (tree.NodeSet, error) {
	xpathRule, err := goxpath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse xpath query %s: %s", matchString, err.Error())
		return nil, err
	}

	xTree, err := xmltree.ParseXML(bytes.NewBufferString(toMatch), func(s *xmltree.ParseOptions) {
		s.Strict = false
	})

	if err != nil {
		log.Errorf("Failed to load XML tree: %s", err.Error())
		return nil, err
	}

	results, err := xpathRule.ExecNode(xTree)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return nil, err
	}

	return results, nil
}
