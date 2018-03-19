package matching

import (
	"bytes"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	log "github.com/Sirupsen/logrus"
)

func XpathMatch(matchingString string, toMatch string) bool {
	xpathRule, err := goxpath.Parse(matchingString)
	if err != nil {
		log.Errorf("Failed to parse xpath query %s: %s", matchingString, err.Error())
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
