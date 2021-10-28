package matchers

import (
	"bytes"

	"github.com/ChrisTrenkamp/xsel/exec"
	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/parser"
	"github.com/ChrisTrenkamp/xsel/store"
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

	return results.Bool()
}

func XpathExecution(matchString, toMatch string) (exec.Result, error) {
	xpath := grammar.MustBuild(matchString)
	parsedXml := parser.ReadXml(bytes.NewBufferString(toMatch))
	cursor, _ := store.CreateInMemory(parsedXml)

	results, err := exec.Exec(cursor, &xpath)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return nil, err
	}

	return results, nil
}
