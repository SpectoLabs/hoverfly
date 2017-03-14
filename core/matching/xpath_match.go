package matching

import (
	"bytes"

	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func XpathMatch(matchingString string, toMatch string) bool {
	xpathRule, err := goxpath.Parse(matchingString)
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
