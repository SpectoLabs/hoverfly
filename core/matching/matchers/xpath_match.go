package matchers

import (
	"bytes"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	log "github.com/Sirupsen/logrus"
)

var Xpath = "xpath"

func XpathMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	xpathRule, err := goxpath.Parse(matchString)
	if err != nil {
		log.Errorf("Failed to parse xpath query %s: %s", matchString, err.Error())
		return false
	}

	xTree, err := xmltree.ParseXML(bytes.NewBufferString(toMatch))
	if err != nil {
		log.Errorf("Failed to load XML tree: %s", err.Error())
		return false
	}

	results, err := xpathRule.ExecNode(xTree)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return false
	}

	return len(results) > 0
}
