package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestNumber(t *testing.T) {
	var np *int
	var nnp int = 3

	test(t, "basic number", true, v.Number(), 3)
	test(t, "non-number", false, v.Number(), "")
	test(t, "nil", false, v.Number(), nil)
	test(t, "nil int pointer", false, v.Number(), np)
	test(t, "int pointer", true, v.Number(), &nnp)

	test(t, "numis", true, v.Number(v.NumIs(3)), 3)
	test(t, "!numis", false, v.Number(v.NumIs(3)), 2)

	test(t, "minlen1", false, v.Number(v.NumMin(3)), 2)
	test(t, "minlen2", true, v.Number(v.NumMin(3)), 3)
	test(t, "minlen2", true, v.Number(v.NumMin(3)), 4)

	test(t, "maxlen1", true, v.Number(v.NumMax(3)), 2)
	test(t, "maxlen2", true, v.Number(v.NumMax(3)), 3)
	test(t, "maxlen3", false, v.Number(v.NumMax(3)), 4)

	test(t, "combination1", false, v.Number(v.NumMin(3), v.NumMax(3)), 2)
	test(t, "combination2", true, v.Number(v.NumMin(3), v.NumMax(3)), 3)
	test(t, "combination3", false, v.Number(v.NumMin(3), v.NumMax(3)), 4)
}
