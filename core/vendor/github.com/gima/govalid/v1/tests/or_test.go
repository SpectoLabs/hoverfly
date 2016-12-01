package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestOr(t *testing.T) {
	test(t, "combination1", true, v.Or(v.String(v.StrIs("a")), v.String(v.StrIs("b"))), "a")
	test(t, "combination2", true, v.Or(v.String(v.StrIs("a")), v.String(v.StrIs("b"))), "b")
	test(t, "combination3", false, v.Or(v.String(v.StrIs("a")), v.String(v.StrIs("b"))), 3)

	test(t, "combination4", true, v.Or(v.String(v.StrIs("a")), v.String(v.StrIs("b")), v.String(v.StrIs("c"))), "b")

	test(t, "combination5", true, v.Or(v.String(v.StrIs("a"))), "a")
	test(t, "combination6", false, v.Or(v.String(v.StrIs("a"))), "b")

	test(t, "combination7", true, v.Or(), nil)
}
