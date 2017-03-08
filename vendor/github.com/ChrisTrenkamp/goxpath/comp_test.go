package goxpath

import (
	"fmt"
	"testing"
)

func TestAddition(t *testing.T) {
	p := `1 + 1`
	x := `<?xml version="1.0" encoding="UTF-8"?><test></test>`
	exp := "2"
	execVal(p, x, exp, nil, t)
}

func TestParenths(t *testing.T) {
	p := `(1 + 2) * 3`
	x := `<?xml version="1.0" encoding="UTF-8"?><test></test>`
	exp := "9"
	execVal(p, x, exp, nil, t)
}

func TestParenths2(t *testing.T) {
	p := `(1 + 2 * 3)`
	x := `<?xml version="1.0" encoding="UTF-8"?><test></test>`
	exp := "7"
	execVal(p, x, exp, nil, t)
}

func TestEquals(t *testing.T) {
	p := `/test/test2 = 3`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><test2>3</test2></test>`
	exp := "true"
	execVal(p, x, exp, nil, t)
}

func TestEqualsStr(t *testing.T) {
	p := `/test/test2 = 'foobar'`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><test2>foobar</test2></test>`
	exp := "true"
	execVal(p, x, exp, nil, t)
}

func TestEqualsStr2(t *testing.T) {
	p := `/root[@test="foo"]`
	x := `<?xml version="1.0" encoding="UTF-8"?><root test="foo">test</root>`
	exp := "test"
	execVal(p, x, exp, nil, t)
}

func TestNumberOps(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><t><t1>2</t1><t2>3</t2><t3>5</t3><t4>2</t4></t>`
	testFloatMap := make(map[string]float64)
	testFloatMap[`/t/t1 * 3`] = 2 * 3
	testFloatMap[`5 div /t/t1`] = 5.0 / 2.0
	testFloatMap[`/t/t2 + /t/t3`] = 3 + 5
	testFloatMap[`/t/t2 - /t/t3`] = 3 - 5
	testFloatMap[`/t/t3 mod /t/t1`] = 5 % 2
	testFloatMap[`number('5')`] = 5
	testFloatMap[`sum(/t/*)`] = 2 + 3 + 5 + 2
	testFloatMap[`floor(/t/t3 div /t/t1)`] = 2
	testFloatMap[`ceiling(t/t3 div /t/t1)`] = 3
	testFloatMap[`round(-1.5)`] = -2
	testFloatMap[`round(1.5)`] = 2
	testFloatMap[`round(0)`] = 0
	for k, v := range testFloatMap {
		execVal(k, x, fmt.Sprintf("%g", float64(v)), nil, t)
	}

	execVal(`/t/t3 div 0`, x, "Infinity", nil, t)
	execVal(`-1 div 0`, x, "-Infinity", nil, t)
	execVal(`0 div 0`, x, "NaN", nil, t)

	testBoolMap := make(map[string]string)
	testBoolMap[`/t/t1 = 2`] = "true"
	testBoolMap[`2 = /t/t1`] = "true"
	testBoolMap[`/t/t1 != 2`] = "false"
	testBoolMap[`4 = /t/t2`] = "false"
	testBoolMap[`/t/t1 != /t/t2`] = "true"
	testBoolMap[`2 < /t/t4`] = "false"
	testBoolMap[`/t/t1 <= 2`] = "true"
	testBoolMap[`/t/t1 > /t/t4`] = "false"
	testBoolMap[`2 >= /t/t4`] = "true"
	testBoolMap[`/t/t1 >= /t/t4`] = "true"
	testBoolMap[`/t/t1 != /t/t2 and /t/t1 < /t/t4`] = "false"
	testBoolMap[`/t/t1 != /t/t2 or /t/t1 < /t/t4`] = "true"
	testBoolMap[`(/t/t1 != /t/t2 or /t/t1 < /t/t4) = true()`] = "true"
	testBoolMap[`(/t/t1 != /t/t2 and /t/t1 < /t/t4) != true()`] = "true"
	for k, v := range testBoolMap {
		execVal(k, x, v, nil, t)
	}
}
