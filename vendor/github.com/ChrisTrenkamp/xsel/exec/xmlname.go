package exec

import (
	"fmt"
	"strings"
)

type XmlName struct {
	Space string
	Local string
}

func (n XmlName) String() string {
	if n.Space == "" {
		return n.Local
	}

	return "{" + n.Space + "}" + n.Local
}

func GetQName(input string, namespaces map[string]string) (XmlName, error) {
	spl := strings.Split(input, ":")
	ret := XmlName{}

	if len(spl) == 1 {
		ret.Local = spl[0]
	} else {
		ns, ok := namespaces[strings.TrimSpace(spl[0])]

		if !ok {
			return XmlName{}, fmt.Errorf("unknown namespace binding '%s'", spl[0])
		}

		ret.Space = ns
		ret.Local = spl[1]
	}

	ret.Local = strings.TrimSpace(ret.Local)
	return ret, nil
}
