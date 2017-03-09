package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type optest struct {
	name     string
	path     string
	expected []int
}

var optests = []optest{
	optest{"single key (period) ", `$.aKey`, []int{opTypeName}},
	optest{"single key (bracket)", `$["aKey"]`, []int{opTypeName}},
	optest{"single key (period) ", `$.*`, []int{opTypeNameWild}},
	optest{"single index", `$[12]`, []int{opTypeIndex}},
	optest{"single key", `$[23:45]`, []int{opTypeIndexRange}},
	optest{"single key", `$[*]`, []int{opTypeIndexWild}},

	optest{"double key", `$["aKey"]["bKey"]`, []int{opTypeName, opTypeName}},
	optest{"double key", `$["aKey"].bKey`, []int{opTypeName, opTypeName}},
}

func TestQueryOperators(t *testing.T) {
	as := assert.New(t)

	for _, t := range optests {
		path, err := parsePath(t.path)
		as.NoError(err)

		as.EqualValues(len(t.expected), len(path.operators))

		for x, op := range t.expected {
			as.EqualValues(pathTokenNames[op], pathTokenNames[path.operators[x].typ])
		}
	}
}
