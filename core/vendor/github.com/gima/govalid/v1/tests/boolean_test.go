package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestBoolean(t *testing.T) {
	var np *bool
	var nnp bool = true

	test(t, "type check: true", true, v.Boolean(), true)
	test(t, "type check: false", true, v.Boolean(), false)
	test(t, "type check: non-boolean", false, v.Boolean(), nil)

	test(t, "nil bool pointer", false, v.Boolean(), np)
	test(t, "non-nil bool pointer", true, v.Boolean(), &nnp)

	test(t, "should true value", true, v.Boolean(v.BoolIs(true)), true)
	test(t, "!should true value", false, v.Boolean(v.BoolIs(true)), false)

	test(t, "should false value", true, v.Boolean(v.BoolIs(false)), false)
	test(t, "!should false value", false, v.Boolean(v.BoolIs(false)), true)
}
