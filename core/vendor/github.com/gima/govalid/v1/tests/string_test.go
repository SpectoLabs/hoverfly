package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestString(t *testing.T) {
	var np *string
	var nnp string = "a"

	test(t, "non-string", false, v.String(), 3)
	test(t, "basic string", true, v.String(), "")
	test(t, "nil", false, v.String(), nil)
	test(t, "nil string pointer", false, v.String(), np)
	test(t, "string pointer", true, v.String(), &nnp)

	test(t, "equals", true, v.String(v.StrIs("abc")), "abc")
	test(t, "!equals", false, v.String(v.StrIs("abc")), "abd")

	test(t, "minlen1", false, v.String(v.StrMin(3)), "aa")
	test(t, "minlen2", true, v.String(v.StrMin(3)), "aaa")
	test(t, "minlen3", true, v.String(v.StrMin(3)), "aaaa")

	test(t, "maxlen1", true, v.String(v.StrMax(4)), "aaa")
	test(t, "maxlen2", true, v.String(v.StrMax(4)), "aaaa")
	test(t, "maxlen3", false, v.String(v.StrMax(4)), "aaaaa")

	test(t, "regexp1", true, v.String(v.StrRegExp("^.{3}$")), "bbb")
	test(t, "regexp2", false, v.String(v.StrRegExp("^.{3}$")), "bbbb")
	test(t, "regexp3", false, v.String(v.StrRegExp("[")), "c")

	test(t, "combination1", false, v.String(v.StrMin(3), v.StrMax(3)), "cc")
	test(t, "combination2", true, v.String(v.StrMin(3), v.StrMax(3)), "ccc")
	test(t, "combination1", false, v.String(v.StrMin(3), v.StrMax(3)), "cccc")

}
