package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestOptional(t *testing.T) {
	var np *bool
	test(t, "string", true, v.Optional(v.String(v.StrIs("a"))), "a")
	test(t, "wrong string", false, v.Optional(v.String(v.StrIs("a"))), "b")

	test(t, "nil ptr", true, v.Optional(v.String(v.StrIs("a"))), np)
}
