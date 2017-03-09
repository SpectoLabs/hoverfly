package jsonpath

import (
	"bytes"
	"fmt"
)

const (
	JsonObject = iota
	JsonArray
	JsonString
	JsonNumber
	JsonNull
	JsonBool
)

type Result struct {
	Keys  []interface{}
	Value []byte
	Type  int
}

func (r *Result) Pretty(showPath bool) string {
	b := bytes.NewBufferString("")
	printed := false
	if showPath {
		for _, k := range r.Keys {
			switch v := k.(type) {
			case int:
				b.WriteString(fmt.Sprintf("%d", v))
			default:
				b.WriteString(fmt.Sprintf("%q", v))
			}
			b.WriteRune('\t')
			printed = true
		}
	} else if r.Value == nil {
		if len(r.Keys) > 0 {
			printed = true
			switch v := r.Keys[len(r.Keys)-1].(type) {
			case int:
				b.WriteString(fmt.Sprintf("%d", v))
			default:
				b.WriteString(fmt.Sprintf("%q", v))
			}
		}
	}

	if r.Value != nil {
		printed = true
		b.WriteString(fmt.Sprintf("%s", r.Value))
	}
	if printed {
		b.WriteRune('\n')
	}
	return b.String()
}
