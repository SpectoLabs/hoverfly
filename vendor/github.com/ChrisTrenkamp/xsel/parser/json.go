package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/ChrisTrenkamp/xsel/node"
)

type JsonElement struct {
	local string
}

func (j JsonElement) Space() string {
	return ""
}

func (j JsonElement) Local() string {
	return j.local
}

type JsonCharData struct {
	value string
}

func (j JsonCharData) CharDataValue() string {
	return j.value
}

type stateType int

const (
	arrayState stateType = iota
	objectState
	rootState
)

type jsonState struct {
	stateType      stateType
	onField        bool
	emitEndElement bool
}

type jsonParser struct {
	jsonReader *json.Decoder
	stateStack []jsonState
}

func (j *jsonParser) pushState(s stateType) {
	j.stateStack = append(j.stateStack, jsonState{stateType: s})
}

func (j *jsonParser) currentState() stateType {
	if len(j.stateStack) == 0 {
		return rootState
	}

	return j.stateStack[len(j.stateStack)-1].stateType
}

func (j *jsonParser) popState() {
	j.stateStack = j.stateStack[:len(j.stateStack)-1]
}

func (j *jsonParser) setOnField(b bool) {
	if len(j.stateStack) > 0 {
		j.stateStack[len(j.stateStack)-1].onField = b
	}
}

func (j *jsonParser) isOnField() bool {
	if len(j.stateStack) == 0 {
		return false
	}

	return j.stateStack[len(j.stateStack)-1].onField
}

func (j *jsonParser) setEmitEndElement(b bool) {
	if len(j.stateStack) > 0 {
		j.stateStack[len(j.stateStack)-1].emitEndElement = b
	}
}

func (j *jsonParser) isOnEmitEndElement() bool {
	if len(j.stateStack) == 0 {
		return false
	}

	return j.stateStack[len(j.stateStack)-1].emitEndElement
}

func jsonTokenValue(tok json.Token) string {
	switch t := tok.(type) {
	case bool:
		return fmt.Sprintf("%t", t)
	case float64:
		str := strconv.FormatFloat(t, 'g', -1, 64)
		return str
	case json.Number:
		return string(t)
	case string:
		return t
	}

	return "null"
}

func (j *jsonParser) Pull() (node.Node, bool, error) {
	if j.isOnEmitEndElement() {
		j.setEmitEndElement(false)
		return nil, true, nil
	}

	tok, err := j.jsonReader.Token()

	if err != nil {
		return nil, false, err
	}

	switch t := tok.(type) {
	case json.Delim:
		switch t.String() {
		case "{":
			if j.currentState() == objectState {
				j.setOnField(true)
			}

			j.pushState(objectState)
			j.setOnField(true)
			return JsonElement{local: "#obj"}, false, nil
		case "}":
			j.popState()

			if j.isOnField() {
				j.setEmitEndElement(true)
			}

			return nil, true, nil
		case "[":
			if j.currentState() == objectState {
				j.setOnField(true)
			}

			j.pushState(arrayState)
			return JsonElement{local: "#arr"}, false, nil
		case "]":
			j.popState()

			if j.isOnField() {
				j.setEmitEndElement(true)
			}

			return nil, true, nil
		}
	}

	val := jsonTokenValue(tok)

	switch j.currentState() {
	case rootState:
		return JsonCharData{value: val}, false, nil
	case objectState:
		if j.isOnField() {
			j.setOnField(false)
			return JsonElement{local: val}, false, nil
		} else {
			j.setOnField(true)
			j.setEmitEndElement(true)
		}
	}

	return JsonCharData{value: val}, false, nil
}

// Create a Parser that reads the given JSON document.
func ReadJson(in io.Reader) Parser {
	jsonReader := json.NewDecoder(in)

	return &jsonParser{
		jsonReader: jsonReader,
		stateStack: make([]jsonState, 0),
	}
}
