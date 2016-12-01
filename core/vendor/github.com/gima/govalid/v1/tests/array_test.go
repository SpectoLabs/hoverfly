package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

func TestArray(t *testing.T) {
	var np *[]interface{}
	var nnsp = []interface{}{1, 2, 3}
	var nnap = [3]interface{}{1, 2, 3}

	test(t, "basic slice", true, v.Array(), []interface{}{1, 2, 3})
	test(t, "nil", false, v.Array(), nil)
	test(t, "nil slice pointer", false, v.Array(), np)
	test(t, "slice pointer", true, v.Array(), &nnsp)
	test(t, "non-array/slice", false, v.Array(), 3)

	test(t, "basic array ", true, v.Array(), [3]interface{}{1, 2, 3})
	test(t, "nil array pointer", true, v.Array(), [3]interface{}{1, 2, 3})
	test(t, "array pointer", true, v.Array(), &nnap)

	test(t, "int slice", true, v.Array(), []int{1, 2, 3})

	test(t, "minlen1", false, v.Array(v.ArrMin(3)), []interface{}{1, 2})
	test(t, "minlen2", true, v.Array(v.ArrMin(3)), []interface{}{1, 2, 3})
	test(t, "minlen3", true, v.Array(v.ArrMin(3)), []interface{}{1, 2, 3, 4})

	test(t, "maxlen", true, v.Array(v.ArrMax(3)), []interface{}{1, 2})
	test(t, "maxlen2", true, v.Array(v.ArrMax(3)), []interface{}{1, 2, 3})
	test(t, "maxlen3", false, v.Array(v.ArrMax(3)), []interface{}{1, 2, 3, 4})

	test(t, "combination1", false, v.Array(v.ArrMin(3), v.ArrMax(3)), []interface{}{1, 2})
	test(t, "combination2", true, v.Array(v.ArrMin(3), v.ArrMax(3)), []interface{}{1, 2, 3})
	test(t, "combination3", false, v.Array(v.ArrMin(3), v.ArrMax(3)), []interface{}{1, 2, 3, 4})

	test(t, "each1", true, v.Array(v.ArrEach(v.Number(v.NumMin(3)))), []interface{}{})
	test(t, "each2", false, v.Array(v.ArrEach(v.Number(v.NumMin(3)))), []interface{}{2, 3})
	test(t, "each3", true, v.Array(v.ArrEach(v.Number(v.NumMin(3)))), []interface{}{3, 4, 5})
}
