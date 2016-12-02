package govalid_test

import (
	v "github.com/gima/govalid/v1"
	"testing"
)

type maep map[interface{}]interface{}

func TestObject(t *testing.T) {
	var np *map[interface{}]interface{}

	test(t, "basic object", true, v.Object(), maep{})
	test(t, "nil", false, v.Object(), nil)
	test(t, "nil map ptr", false, v.Object(), np)
	test(t, "non-object", false, v.Object(), 3)

	testObjKeys(t)
	testObjValues(t)
	testObjKVs(t)
}

func testObjKeys(t *testing.T) {
	counter, countingValidator := createCountingValidator()

	sch := v.Object(
		v.ObjKeys(v.String()),
		v.ObjKeys(v.Function(countingValidator)),
	)
	m := maep{
		"a": nil,
		"b": 1,
		"c": true,
	}
	test(t, "only string objkeys", true, sch, m)
	if *counter != 3 {
		t.Fatalf("key counter should be 3, got %d", *counter)
	}

	m = maep{
		"a": nil,
		1:   1,
	}
	test(t, "!only string objkeys", false, sch, m)
}

func testObjValues(t *testing.T) {
	counter, countingValidator := createCountingValidator()

	sch := v.Object(
		v.ObjValues(v.String()),
		v.ObjValues(v.Function(countingValidator)),
	)
	m := maep{
		nil:  "1",
		1:    "b",
		true: "c",
	}
	test(t, "only string objvalues", true, sch, m)
	if *counter != 3 {
		t.Fatalf("value counter should be 3, got %d", *counter)
	}

	m = maep{
		nil: "1",
		1:   1,
	}
	test(t, "!only string objvalues", false, sch, m)
}

func testObjKVs(t *testing.T) {
	counter, countingValidator := createCountingValidator()

	sch := v.Object(
		v.ObjKV(nil, v.And(v.String(v.StrIs("1")), v.Function(countingValidator))),
		v.ObjKV("1", v.And(v.String(v.StrIs("b")), v.Function(countingValidator))),
		v.ObjKV(true, v.And(v.Number(v.NumIs(3)), v.Function(countingValidator))),
	)
	m := maep{
		nil:  "1",
		"1":  "b",
		true: 3,
	}
	test(t, "mixed objkvs", true, sch, m)
	if *counter != 3 {
		t.Fatalf("value counter should be 3, got %d", *counter)
	}

	m = maep{
		nil:  "1",
		"1":  2,
		true: 3,
	}
	test(t, "!mixed objkvs", false, sch, m)

	m = maep{
		nil:  "1",
		"1":  nil,
		true: 3,
	}
	test(t, "!mixed objkvs (nil)", false, sch, m)
}
