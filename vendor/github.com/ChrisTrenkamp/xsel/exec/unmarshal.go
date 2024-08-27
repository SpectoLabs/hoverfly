package exec

import (
	"fmt"
	"reflect"

	"github.com/ChrisTrenkamp/xsel/grammar"
)

// Unmarshal maps a XPath result to a struct or slice.
// When unmarshaling a slice, the result must be a NodeSet. When unmarshaling
// a struct, the result must be a NodeSet with one result. To unmarshal a
// value to a struct field, give it a "xsel" tag name, and a XPath expression
// for its value (e.g. `xsel:"//my-struct[@my-id = 'my-value']"`).
//
// For struct fields, Unmarshal can set fields that are ints and uints, bools,
// strings, slices, and nested structs.
//
// For slice elements, Unmarshal can set ints and uints, bools, strings, and
// structs.  It cannot Unmarshal multidimensional slices.
//
// Arrays, maps, and channels are not supported.
func Unmarshal(result Result, value any, settings ...ContextApply) error {
	return unmarshal(result, value, settings...)
}

func unmarshal(result Result, value any, settings ...ContextApply) error {
	val := reflect.ValueOf(value)
	typ := val.Type()

	for typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	kind := typ.Kind()

	if kind == reflect.Struct {
		return unmarshalStruct(result, val.Addr(), settings...)
	}

	if kind == reflect.Slice {
		return unmarshalSlice(result, val, settings...)
	}

	return fmt.Errorf("unsupported data type")
}

func unmarshalStruct(result Result, val reflect.Value, settings ...ContextApply) error {
	cursor, ok := result.(NodeSet)

	if !ok || len(cursor) != 1 {
		return fmt.Errorf("struct unmarshals must operate on a NodeSet with one result")
	}

	val = val.Elem()

	numField := val.NumField()
	valType := val.Type()

	for i := 0; i < numField; i++ {
		fieldValType := valType.Field(i)
		name := fieldValType.Name
		tag := fieldValType.Tag.Get("xsel")

		if tag == "" {
			continue
		}

		xselExec, err := grammar.Build(tag)
		if err != nil {
			return err
		}

		result, err := Exec(cursor[0], &xselExec, settings...)
		if err != nil {
			return err
		}

		field := val.Field(i)
		fieldType := field.Type()
		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		fieldVal, ok := createValue(fieldType.Kind(), result)

		if ok {
			err = setField(name, field, fieldVal, false)
		} else {
			ptr := reflect.New(fieldType)
			ptr.Elem().Set(reflect.Zero(fieldType))
			err = unmarshal(result, ptr.Interface(), settings...)
			if err != nil {
				return err
			}

			err = setField(name, field, ptr.Elem(), false)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalSlice(result Result, val reflect.Value, settings ...ContextApply) error {
	nodeset, ok := result.(NodeSet)

	if !ok {
		return fmt.Errorf("slice unmarshals must operate on a NodeSet")
	}

	sliceType := reflect.TypeOf(val.Interface())
	sliceElement := sliceType.Elem()
	sliceElementKind := sliceElement.Kind()

	for sliceElementKind == reflect.Pointer {
		sliceElement = sliceElement.Elem()
		sliceElementKind = sliceElement.Kind()
	}

	for _, i := range nodeset {
		var sliceValue reflect.Value

		if sliceElementKind == reflect.Slice {
			return fmt.Errorf("slice unmarshals can only operate on 1-dimensional slices")
		} else if sliceElementKind == reflect.Struct {
			ptr := reflect.New(sliceElement)
			ptr.Elem().Set(reflect.Zero(sliceElement))

			err := unmarshal(NodeSet{i}, ptr.Interface(), settings...)
			if err != nil {
				return err
			}

			sliceValue = ptr.Elem()
		} else {
			val, ok := createValue(sliceElementKind, NodeSet{i})

			if !ok {
				return fmt.Errorf("invalid slice element type")
			}

			sliceValue = val
		}

		err := setField("<slice>", val, sliceValue, true)

		if err != nil {
			return err
		}
	}

	return nil
}

func setField(name string, field reflect.Value, val reflect.Value, checkSlice bool) error {
	isSlice := false
	dereferences := 0
	typ := field.Type()
	assignableType := typ

	if checkSlice && typ.Kind() == reflect.Slice {
		isSlice = true
		typ = typ.Elem()
		assignableType = typ
	}

	elemKind := typ.Kind()

	for elemKind == reflect.Pointer {
		typ = typ.Elem()
		elemKind = typ.Kind()
		dereferences++
	}

	ptrVal := val

	for dereferences != 0 {
		ptr := reflect.New(ptrVal.Type())
		ptr.Elem().Set(ptrVal)
		ptrVal = ptr
		dereferences--
	}

	if !assignableType.AssignableTo(ptrVal.Type()) {
		return fmt.Errorf("could not set field, %s", name)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s is not settable", name)
	}

	if isSlice {
		field.Set(reflect.Append(field, ptrVal))
	} else {
		field.Set(ptrVal)
	}

	return nil
}

func createValue(kind reflect.Kind, result Result) (reflect.Value, bool) {
	switch kind {
	case reflect.String:
		return reflect.ValueOf(result.String()), true

	case reflect.Bool:
		return reflect.ValueOf(result.Bool()), true

	case reflect.Int:
		return reflect.ValueOf(int(result.Number())), true
	case reflect.Uint:
		return reflect.ValueOf(uint(result.Number())), true
	case reflect.Uint8:
		return reflect.ValueOf(uint8(result.Number())), true
	case reflect.Int8:
		return reflect.ValueOf(int8(result.Number())), true
	case reflect.Uint16:
		return reflect.ValueOf(uint16(result.Number())), true
	case reflect.Int16:
		return reflect.ValueOf(int16(result.Number())), true
	case reflect.Uint32:
		return reflect.ValueOf(uint32(result.Number())), true
	case reflect.Int32:
		return reflect.ValueOf(int32(result.Number())), true
	case reflect.Uint64:
		return reflect.ValueOf(uint64(result.Number())), true
	case reflect.Int64:
		return reflect.ValueOf(int64(result.Number())), true
	case reflect.Float32:
		return reflect.ValueOf(float32(result.Number())), true
	case reflect.Float64:
		return reflect.ValueOf(result.Number()), true
	}

	return reflect.ValueOf(0), false
}
