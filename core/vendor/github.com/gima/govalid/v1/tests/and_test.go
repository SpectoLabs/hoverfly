package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestAnd(t *testing.T) {
	test(t, "combination1", true, v.And(), nil)

	test(t, "combination2", false, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(3))), "aa")
	test(t, "combination3", true, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(3))), "aaa")
	test(t, "combination4", false, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(3))), "aaaa")

	test(t, "combination5", false, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(4))), "bb")
	test(t, "combination6", true, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(4))), "bbb")
	test(t, "combination7", true, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(4))), "bbbb")
	test(t, "combination8", false, v.And(v.String(v.StrMin(3)), v.String(v.StrMax(4))), "bbbbb")
}
