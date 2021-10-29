package matchers

import (
	"bytes"
	"encoding/xml"
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

	contextSettings := func(c *exec.ContextSettings) {
		xmlns := xmlns{}
		_ = xml.Unmarshal([]byte(toMatch), &xmlns)
		for key, value := range xmlns.Namespaces {
			c.NamespaceDecls[key] = value
		}
	}
	xpath := grammar.MustBuild(matchString)
	parsedXml := parser.ReadXml(bytes.NewBufferString(toMatch))
	cursor, _ := store.CreateInMemory(parsedXml)

	results, err := exec.Exec(cursor, &xpath, contextSettings)
	if err != nil {
		log.Errorf("Failed to execute xpath match: %s", err.Error())
		return nil, err
	}

	return results, nil
}


type xmlns struct {
	Namespaces map[string]string
}

func (a *xmlns) UnmarshalXML(_ *xml.Decoder, start xml.StartElement) error {
	a.Namespaces = map[string]string{}
	for _, attr := range start.Attr {
		if attr.Name.Space == "xmlns" {
			a.Namespaces[attr.Name.Local] = attr.Value
		}
	}
	return nil
}
